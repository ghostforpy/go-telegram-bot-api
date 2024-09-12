package tgbotapi

import (
	"context"
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
	Filters     []func(Update) (bool, error)
	Middlewares []Middleware
}

func (dh *DefalutHander) CheckUpdate(update Update) (bool, error) {
	return true, nil
}
func (dh *DefalutHander) HandleUpdate(tctx TgbotapiContext, update Update) error {
	callbackFunc := applyMiddlewares(dh.Callback, dh.Middlewares...)
	return callbackFunc(tctx, update)
}

type MessageHander struct {
	Handler
	Callback    func(tctx TgbotapiContext, update Update) error
	Filters     []func(Update) (bool, error)
	Middlewares []Middleware
}

func (mh *MessageHander) CheckUpdate(update Update) (bool, error) {
	if update.Message == nil {
		return false, nil
	}
	if len(mh.Filters) > 0 {
		for _, filterFunc := range mh.Filters {
			if ok, err := filterFunc(update); err == nil && !ok {
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
	Filters     []func(Update) (bool, error)
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
