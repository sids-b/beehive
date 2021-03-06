package context

import (
	"sync"
	"time"

	"github.com/kubeedge/kubeedge/common/beehive/pkg/common/log"
	"github.com/kubeedge/kubeedge/common/beehive/pkg/core/model"
)

//define channel type
const (
	MsgCtxTypeChannel = "channel"
)

var (
	// singleton
	context *Context
	once    sync.Once
)

// GetContext gets global context instance
func GetContext(contextType string) *Context {
	once.Do(func() {
		context = &Context{}
		switch contextType {
		case MsgCtxTypeChannel:
			channelContext := NewChannelContext()
			context.messageContext = channelContext
			context.moduleContext = channelContext
		default:
			log.LOGGER.Warnf("do not support context type(%s)", contextType)
		}
	})
	return context
}

// AddModule adds module into module context
func (ctx *Context) AddModule(module string) {
	ctx.moduleContext.AddModule(module)
}

// AddModuleGroup adds module into module context group
func (ctx *Context) AddModuleGroup(module, group string) {
	ctx.moduleContext.AddModuleGroup(module, group)
}

// Cleanup cleans up module
func (ctx *Context) Cleanup(module string) {
	ctx.moduleContext.Cleanup(module)
}

// Send the message
func (ctx *Context) Send(module string, message model.Message) {
	ctx.messageContext.Send(module, message)
}

// Receive the message
// module : local module name
func (ctx *Context) Receive(module string) (model.Message, error) {
	message, err := ctx.messageContext.Receive(module)
	if err == nil {
		return message, nil
	}
	log.LOGGER.Warnf("failed to receive message")
	return message, err
}

// SendSync sends message in sync mode
// module: the destination of the message
// timeout: if <= 0 using default value(30s)
func (ctx *Context) SendSync(module string,
	message model.Message, timeout time.Duration) (model.Message, error) {
	resp, err := ctx.messageContext.SendSync(module, message, timeout)
	if err == nil {
		return resp, nil
	}
	return model.Message{}, err
}

// SendResp sends response
// please get resp message using model.NewRespByMessage
func (ctx *Context) SendResp(resp model.Message) {
	ctx.messageContext.SendResp(resp)
}

// Send2Group broadcasts the message to all of group members
func (ctx *Context) Send2Group(moduleType string, message model.Message) {
	ctx.messageContext.Send2Group(moduleType, message)
}

// send2GroupSync broadcasts the message to all of group members in sync mode
func (ctx *Context) send2GroupSync(moduleType string, message model.Message, timeout time.Duration) error {
	return ctx.messageContext.Send2GroupSync(moduleType, message, timeout)
}
