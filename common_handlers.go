package tgbotapi

import (
	"fmt"
	"regexp"
)

type HandleUpdateFunc func(tctx TgbotapiContext, update Update) error

type Handler interface {
	CheckUpdate(update Update) (bool, error)
	HandleUpdate(tctx TgbotapiContext, update Update) error
}

type CommonHandler interface {
	Handler
}

type DefalutHander struct {
	CommonHandler
	Callback    func(tctx TgbotapiContext, update Update) error
	Filters     []func(Update) bool
	Middlewares []Middleware
}

func (dh *DefalutHander) CheckUpdate(update Update) (bool, error) {
	if len(dh.Filters) > 0 {
		for _, filterFunc := range dh.Filters {
			if ok := filterFunc(update); !ok {
				return false, nil
			}
		}
	}
	return true, nil
}
func (dh *DefalutHander) HandleUpdate(tctx TgbotapiContext, update Update) error {
	callbackFunc := applyMiddlewares(dh.Callback, dh.Middlewares...)
	return callbackFunc(tctx, update)
}

type MessageHander struct {
	Handler
	Callback    func(tctx TgbotapiContext, update Update) error
	Filters     []func(Update) bool
	Middlewares []Middleware
}

func (mh *MessageHander) CheckUpdate(update Update) (bool, error) {
	if update.Message == nil {
		return false, nil
	}
	if len(mh.Filters) > 0 {
		for _, filterFunc := range mh.Filters {
			if ok := filterFunc(update); !ok {
				return false, nil
			}
		}
	}
	return true, nil
}

func (mh *MessageHander) HandleUpdate(tctx TgbotapiContext, update Update) error {
	callbackFunc := applyMiddlewares(mh.Callback, mh.Middlewares...)
	return callbackFunc(tctx, update)
}

type CommandHander struct {
	DefalutHander
	Command     string
	Callback    func(tctx TgbotapiContext, update Update) error
	Filters     []func(Update) bool
	Middlewares []Middleware
}

func (cmh *CommandHander) CheckUpdate(update Update) (bool, error) {
	if update.Message == nil {
		return false, nil
	}
	if !update.Message.IsCommand() {
		return false, nil
	}
	if update.Message.Command() != cmh.Command {
		return false, nil
	}
	if len(cmh.Filters) > 0 {
		for _, filterFunc := range cmh.Filters {
			if ok := filterFunc(update); !ok {
				return false, nil
			}
		}
	}
	return true, nil

}

func (cmh *CommandHander) HandleUpdate(tctx TgbotapiContext, update Update) error {
	callbackFunc := applyMiddlewares(cmh.Callback, cmh.Middlewares...)
	return callbackFunc(tctx, update)
}

type CallbackHander struct {
	CommonHandler
	Pattern     *regexp.Regexp
	Callback    func(tctx TgbotapiContext, update Update) error
	Middlewares []Middleware
}

func NewCallbackHander(pattern, handlerHame string, Callback func(tctx TgbotapiContext, update Update) error) (*CallbackHander, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		panic(fmt.Sprintf("%v not compiled for CallbackHander %v", pattern, handlerHame))
	}

	return &CallbackHander{Pattern: re, Callback: Callback}, nil
}

func (cbh *CallbackHander) CheckUpdate(update Update) (bool, error) {
	if update.CallbackQuery == nil {
		return false, nil
	}
	matched := cbh.Pattern.Match([]byte(update.CallbackData()))
	if !matched {
		return false, nil
	}
	return true, nil
}

func (ch *CallbackHander) HandleUpdate(tctx TgbotapiContext, update Update) error {
	callbackFunc := applyMiddlewares(ch.Callback, ch.Middlewares...)
	return callbackFunc(tctx, update)
}

type InlineQueryHander struct {
	CommonHandler
	Pattern *regexp.Regexp
	// Type of the chat, from which the inline query was sent. Can be either
	// “sender” for a private chat with the inline query sender, “private”,
	// “group”, “supergroup”, or “channel”. The chat type should be always known
	// for requests sent from official clients and most third-party clients,
	// unless the request was sent from a secret chat
	//
	// optional
	ChatTypes   map[string]struct{}
	Callback    func(tctx TgbotapiContext, update Update) error
	Middlewares []Middleware
}

func NewInlineQueryHander(pattern, handlerHame string, Callback func(tctx TgbotapiContext, update Update) error) (*InlineQueryHander, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		panic(fmt.Sprintf("%v not compiled for InlineQueryHander %v", pattern, handlerHame))
	}
	chatTypes := make(map[string]struct{})
	chatTypes["sender"] = struct{}{}
	return &InlineQueryHander{Pattern: re, Callback: Callback, ChatTypes: chatTypes}, nil
}

func (iqh *InlineQueryHander) CheckUpdate(update Update) (bool, error) {
	if update.InlineQuery == nil {
		return false, nil
	}
	if _, ok := iqh.ChatTypes[update.InlineQuery.ChatType]; len(iqh.ChatTypes) > 0 && !ok {
		return false, nil
	}
	matched := iqh.Pattern.Match([]byte(update.InlineQuery.Query))
	if !matched {
		return false, nil
	}
	return true, nil
}

func (iqh *InlineQueryHander) HandleUpdate(tctx TgbotapiContext, update Update) error {
	callbackFunc := applyMiddlewares(iqh.Callback, iqh.Middlewares...)
	return callbackFunc(tctx, update)
}

func (iqh *InlineQueryHander) Print() {
	fmt.Printf("InlineQueryHander %v", iqh.Pattern)
}
