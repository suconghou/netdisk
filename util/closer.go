package util

import (
	"errors"
	"io"
)

// onCloseReadCloser 是我们的装饰器结构体
type onCloseReadCloser struct {
	source      io.ReadCloser // 被包装的原始对象
	onCloseFunc func() error  // 在 Close 时额外调用的函数
}

// Read 方法直接调用原始 source 的 Read 方法
func (o *onCloseReadCloser) Read(p []byte) (n int, err error) {
	return o.source.Read(p)
}

// Close 方法会先调用我们自定义的函数，然后调用原始 source 的 Close 方法
func (o *onCloseReadCloser) Close() error {
	// 调用自定义的清理函数
	err1 := o.onCloseFunc()

	// 调用原始对象的 Close 方法
	err2 := o.source.Close()

	// 使用 errors.Join 合并两个操作可能返回的错误，这是最健壮的做法
	return errors.Join(err1, err2)
}

// NewOnCloseReadCloser 是一个工厂函数，用于创建我们的装饰器实例
// 注意它返回的是 io.ReadCloser 接口类型，而不是具体的 struct 类型，这是 Go 的惯例
func NewOnCloseReadCloser(rc io.ReadCloser, onClose func() error) io.ReadCloser {
	return &onCloseReadCloser{
		source:      rc,
		onCloseFunc: onClose,
	}
}
