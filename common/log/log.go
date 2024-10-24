package log

import (
	"sync"

	"github.com/gzjjjfree/gzv2ray/common/serial"
)

// Message is the interface for all log messages.
// Message 是所有日志消息的接口
type Message interface {
	String() string
}

// Handler is the interface for log handler.
// Handler 是日志处理程序的接口。
type Handler interface {
	Handle(msg Message)
}

// GeneralMessage is a general log message that can contain all kind of content.
// GeneralMessage 是通用日志消息，可以包含所有类型的内容。
type GeneralMessage struct {
	Severity Severity
	Content  interface{}
}

// String implements Message.
// 字符串实现消息。
func (m *GeneralMessage) String() string {
	return serial.Concat("[", m.Severity, "] ", m.Content)
}

// Record writes a message into log stream.
// 记录将消息写入日志流。
func Record(msg Message) {
	logHandler.Handle(msg)
}

var (
	logHandler syncHandler
)

// RegisterHandler register a new handler as current log handler. Previous registered handler will be discarded.
// RegisterHandler 注册一个新的处理程序作为当前日志处理程序。先前注册的处理程序将被丢弃。
func RegisterHandler(handler Handler) {
	if handler == nil {
		panic("Log handler is nil")
	}
	logHandler.Set(handler)
}

type syncHandler struct {
	sync.RWMutex
	Handler
}

func (h *syncHandler) Handle(msg Message) {
	h.RLock()
	defer h.RUnlock()

	if h.Handler != nil {
		h.Handler.Handle(msg)
	}
}

func (h *syncHandler) Set(handler Handler) {
	h.Lock()
	defer h.Unlock()

	h.Handler = handler
}
