package multiget

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"netdisk/request"
	"sync/atomic"
)

const part = 1024 * 256

// 一个并发的多DNS记录同时请求然后合并
type Multi struct {
	target  string
	readers []*lazyReader // 保存引用，用于统一关闭释放内存
	r       io.Reader
}

type partItem struct {
	start int64
	end   int64
}

type lazyReader struct {
	r    *partItem
	res  io.ReadCloser
	err  error
	next *lazyReader
	fn   func(int64, int64) (io.ReadCloser, error)
}

func (l *lazyReader) Read(p []byte) (int, error) {
	if l.res == nil {
		l.res, l.err = makeRequest(l)
	}
	if l.err != nil {
		return 0, l.err
	}
	return l.res.Read(p)
}
func (l *lazyReader) Close() error {
	if l.res == nil {
		return nil
	}
	return l.res.Close()
}

func (m *Multi) Read(p []byte) (int, error) {
	return m.r.Read(p)
}

func (m *Multi) Close() error {
	var errs = []error{}
	for _, r := range m.readers {
		if err := r.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// addr 传入DNS解析IP列表，不包含端口
func Get(target string, start int64, end int64, addr []string) (io.ReadCloser, int64, error) {
	u, filesize, res, err := resURI(target)
	if err != nil {
		return nil, 0, err
	}
	if filesize < part {
		return res, filesize, nil
	}
	if err := res.Close(); err != nil {
		return nil, 0, err
	}
	target = u.String()
	var idx int64 = 0
	fn := func(start int64, end int64) (io.ReadCloser, error) {
		i := atomic.AddInt64(&idx, 1) - 1
		host := addr[i%int64(len(addr))]
		hh := http.Header{
			"range":      []string{fmt.Sprintf("bytes=%d-%d", start, end)},
			"user-agent": ua,
		}
		log.Printf("%s", host)
		res, ok, err := request.Request(context.Background(), target, http.MethodGet, hh, 60, nil, host)
		if err != nil {
			return nil, err
		}
		if !ok {
			_ = res.Body.Close()
			return nil, fmt.Errorf("%s %s:%s", res.Request.Method, res.Request.URL, res.Status)
		}
		return res.Body, nil
	}
	parts, err := itemParts(filesize, start, end)
	if err != nil {
		return nil, filesize, err
	}
	r, readers := newReadConcater(parts, fn)
	return &Multi{target: target, readers: readers, r: r}, filesize, nil
}

func newReadConcater(items []*partItem, fn func(int64, int64) (io.ReadCloser, error)) (io.Reader, []*lazyReader) {
	buffers := make([]*lazyReader, 0, len(items))
	for _, t := range items {
		buffers = append(buffers, &lazyReader{
			r:  t,
			fn: fn,
		})
	}
	max := len(buffers) - 1
	for i := range max {
		buffers[i].next = buffers[i+1]
	}
	rr := make([]io.Reader, len(buffers))
	for i, t := range buffers {
		rr[i] = t
	}
	return io.MultiReader(rr...), buffers
}

func makeRequest(l *lazyReader) (io.ReadCloser, error) {
	n := l.next
	for range 2 {
		if n == nil {
			break
		}
		if n.res == nil && n.err == nil {
			go func(targetLR *lazyReader) {
				targetLR.res, targetLR.err = prefetchToMemory(targetLR)
			}(n)
		}
		n = n.next
	}
	return prefetchToMemory(l)
}

func prefetchToMemory(l *lazyReader) (io.ReadCloser, error) {
	res, err := l.fn(l.r.start, l.r.end)
	if err != nil {
		return nil, err
	}
	bs, err := io.ReadAll(io.LimitReader(res, part))
	_ = res.Close()
	log.Printf("got %d-%d %d %v", l.r.start, l.r.end, len(bs), err)
	if err != nil {
		return nil, err
	}
	return io.NopCloser(bytes.NewBuffer(bs)), nil
}

func itemParts(filesize int64, start int64, end int64) ([]*partItem, error) {
	if end <= 0 || end >= filesize {
		end = filesize - 1
	}
	if start >= filesize || start > end {
		return nil, fmt.Errorf("error start-end")
	}
	var (
		items = []*partItem{}
		last  bool
	)
	for {
		offset := start + part - 1
		if offset >= end-1 {
			offset = end - 1
			last = true
		}
		items = append(items, &partItem{start, offset})
		start = offset + 1
		if last {
			break
		}
	}
	return items, nil
}
