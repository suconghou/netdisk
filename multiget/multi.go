package multiget

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"netdisk/request"
)

const part = 1024 * 512

// 一个并发的多DNS记录同时请求然后合并
type Multi struct {
	ctx     context.Context
	cancel  context.CancelFunc
	target  string
	tasks   chan *task
	res     chan *taskres
	dataMap map[int]*taskres
	current int
	readed  int64
	total   int64
	parts   int
}

type task struct {
	no    int
	start int64
	end   int64
}

type taskres struct {
	no   int
	data *bytes.Buffer
	err  error
}

func (m *Multi) Read(p []byte) (int, error) {
	if m.readed == m.total {
		return 0, io.EOF
	}
	if res, ok := m.dataMap[m.current]; ok {
		if res.data != nil && res.data.Len() > 0 {
			n, err := res.data.Read(p)
			m.readed += int64(n)
			if m.current > 0 && m.current == m.parts && res.err == nil { // 本次读取的是最后一块数据,并且都读取完了,块也没有错误
				delete(m.dataMap, m.current)
				return n, io.EOF
			}
			return n, err
		} else if res.err != nil {
			m.cancel()
			return 0, res.err
		} else {
			delete(m.dataMap, m.current)
			m.current++
			return 0, nil
		}
	}
	// 没找到想read的块,就在这里等它;而不去read,会确保res被塞满暂停新线程下载
	for {
		select {
		case <-m.ctx.Done():
			return 0, m.ctx.Err()
		case task := <-m.res:
			m.dataMap[task.no] = task
			// 有新的块来了,更新存储桶,如果来的块是要找的块,我们中断,下次调用就会读这个块了
			if task.no == m.current {
				return 0, nil
			}
		}
	}
}
func (m *Multi) Close() error {
	m.cancel()
	return nil
}

func (m *Multi) thread(addr []string) {
	for _, a := range addr {
		go func(host string) {
			for {
				select {
				case <-m.ctx.Done():
					return
				case t := <-m.tasks:
					hh := http.Header{
						"Range": []string{fmt.Sprintf("bytes=%d-%d", t.start, t.end)},
					}
					res, ok, err := request.Request(m.ctx, m.target, http.MethodGet, hh, 60, nil, host)
					if err == nil && ok {
						var buf = &bytes.Buffer{}
						_, err = buf.ReadFrom(io.LimitReader(res.Body, part))
						_ = res.Body.Close()
						if err == nil {
							m.res <- &taskres{data: buf, no: t.no}
						} else {
							m.res <- &taskres{err: err, no: t.no}
						}
					} else {
						if err == nil {
							var msg = "failed"
							if res != nil {
								msg = fmt.Sprintf("%s %s", res.Request.Method, res.Status)
							}
							err = fmt.Errorf("%s %s %d [%d-%d] ", msg, m.target, t.no, t.start, t.end)
						}
						m.res <- &taskres{err: err, no: t.no}
					}
				}
			}
		}(a)
	}
}

func (m *Multi) run(parts []*task) {
	for _, item := range parts {
		select {
		case <-m.ctx.Done():
			return
		case m.tasks <- item:
		}
	}
}

// addr 传入DNS解析IP列表，不包含端口
func Get(ctx context.Context, target string, start int64, end int64, addr []string) (io.ReadCloser, int64, error) {
	u, total, res, err := resURI(target)
	if err != nil {
		return nil, 0, err
	}
	if total < part {
		return res, total, nil
	}
	if err := res.Close(); err != nil {
		return nil, 0, err
	}
	target = u.String()
	parts, err := itemParts(total, start, end)
	if err != nil {
		return nil, total, err
	}
	context, cancel := context.WithCancel(ctx)
	m := &Multi{ctx: context, cancel: cancel, target: target, tasks: make(chan *task, 20), res: make(chan *taskres, 20), dataMap: make(map[int]*taskres), parts: len(parts), total: total}
	go m.thread(addr)
	go m.run(parts)
	return m, total, nil
}

func itemParts(filesize int64, start int64, end int64) ([]*task, error) {
	if end <= 0 || end >= filesize {
		end = filesize - 1
	}
	if start >= filesize || start > end {
		return nil, fmt.Errorf("error start-end")
	}
	var (
		items = []*task{}
		last  bool
		index = 0
	)
	for {
		offset := start + part - 1
		if offset >= end-1 {
			offset = end - 1
			last = true
		}
		items = append(items, &task{no: index, start: start, end: offset})
		index++
		start = offset + 1
		if last {
			break
		}
	}
	return items, nil
}
