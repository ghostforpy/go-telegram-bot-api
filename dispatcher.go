package tgbotapi

import (
	"context"
	"encoding/json"
	"fmt"
)

type Dispatcher struct {
	ConvHandlers              []ConvHandler
	CommonHandlers            []CommonHandler
	DefaultHandler            Handler
	Bot                       BotAPI
	StateStorage              StateStorage
	GlobalMiddlewares         []Middleware
	ConvHandlersMiddlewares   []Middleware
	CommonHandlersMiddlewares []Middleware
	DefaultHandlerMiddlewares []Middleware
}

func NewDispatcher(bot BotAPI) *Dispatcher {
	return &Dispatcher{Bot: bot}
}

type TgbotapiContext struct {
	Ctx          context.Context
	UserState    string
	StateStorage StateStorage
}

func NewTgbotapiContext() *TgbotapiContext {
	return &TgbotapiContext{
		Ctx: context.TODO(),
	}
}

func (dispatcher *Dispatcher) HandleUpdate(update Update) error {
	var updateHandled bool
	defer func() {
		if r := recover(); r != nil {
			if JsonUpdate, err := json.Marshal(update); err == nil {
				fmt.Printf("Recovered HandleUpdate %v, err:%v", JsonUpdate, r)
			} else {
				fmt.Printf("Recovered HandleUpdate %v, err:%v", JsonUpdate, r)
			}
		}
	}()
	tctx := NewTgbotapiContext()
	tctx.StateStorage = dispatcher.StateStorage
	state, err := dispatcher.StateStorage.GetState(tctx.Ctx, update)
	if err == nil && state != "" {
		tctx.UserState = state
	} else {
		tctx.UserState = "main"
	}

	for _, convHandler := range dispatcher.ConvHandlers {
		if handlerFunc, err := convHandler.HandleUpdate(*tctx, update); err == nil && handlerFunc != nil {
			handlerFunc := applyMiddlewares(handlerFunc, dispatcher.GlobalMiddlewares...)
			err := handlerFunc(*tctx, update)
			if err == nil {
				updateHandled = true
			}
		}
		if updateHandled {
			break
		}
	}
	if updateHandled {
		return nil
	}
	for _, commonHandler := range dispatcher.CommonHandlers {
		if ok, err := commonHandler.CheckUpdate(update); err == nil && ok {
			handlerFunc := applyMiddlewares(commonHandler.HandleUpdate, dispatcher.GlobalMiddlewares...)
			err := handlerFunc(*tctx, update)
			if err == nil {
				updateHandled = true
			}
		}
		if updateHandled {
			break
		}
	}
	if updateHandled {
		return nil
	}
	if dispatcher.DefaultHandler != nil {
		handlerFunc := applyMiddlewares(dispatcher.DefaultHandler.HandleUpdate, dispatcher.GlobalMiddlewares...)
		err := handlerFunc(*tctx, update)
		if err != nil {
			return err
		}
		updateHandled = true
	}
	if updateHandled {
		return nil
	}
	return fmt.Errorf("%v not handled", update)

}

func (dispatcher *Dispatcher) AddConvHandler(handler ConvHandler) {
	dispatcher.ConvHandlers = append(dispatcher.ConvHandlers, handler)
}

func (dispatcher *Dispatcher) AddCommonHandler(handler CommonHandler) {
	dispatcher.CommonHandlers = append(dispatcher.CommonHandlers, handler)
}

func (dispatcher *Dispatcher) AddGlobalMiddlewares(middlewares ...Middleware) {
	dispatcher.GlobalMiddlewares = append(dispatcher.GlobalMiddlewares, middlewares...)
}
func (dispatcher *Dispatcher) AddConvHandlersMiddlewares(middlewares ...Middleware) {
	dispatcher.ConvHandlersMiddlewares = append(dispatcher.ConvHandlersMiddlewares, middlewares...)
}
func (dispatcher *Dispatcher) AddCommonHandlersMiddlewares(middlewares ...Middleware) {
	dispatcher.CommonHandlersMiddlewares = append(dispatcher.CommonHandlersMiddlewares, middlewares...)
}
func (dispatcher *Dispatcher) AddDefaultHandlerMiddlewares(middlewares ...Middleware) {
	dispatcher.DefaultHandlerMiddlewares = append(dispatcher.DefaultHandlerMiddlewares, middlewares...)
}
