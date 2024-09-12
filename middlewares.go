package tgbotapi

import (
	"fmt"
	"time"
)

type Middleware func(next HandleUpdateFunc) HandleUpdateFunc

func applyMiddlewares(h HandleUpdateFunc, m ...Middleware) HandleUpdateFunc {
	if len(m) < 1 {
		return h
	}
	wrapped := h
	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped)
	}
	return wrapped
}

// AutoAnswerCallbackQueryMiddleware answer every callback query
func AutoAnswerCallbackQueryMiddleware(next HandleUpdateFunc) HandleUpdateFunc {
	return func(tctx TgbotapiContext, update Update) error {
		if update.CallbackQuery != nil {
			update.Bot.AnswerCallbackQuery(update)
		}
		return next(tctx, update)
	}
}

// ExampleMiddleware proceed update
func ExampleMiddleware(next HandleUpdateFunc) HandleUpdateFunc {
	return func(tctx TgbotapiContext, update Update) error {
		// do something
		return next(tctx, update)
	}
}

// ExampleStopMiddleware stop proceed update
func ExampleStopMiddleware(next HandleUpdateFunc) HandleUpdateFunc {
	return func(tctx TgbotapiContext, update Update) error {
		// do something
		return nil
	}
}

// ExampleNotHanledMiddleware stop proceed update in handlers block
// ConvHandlers -> CommonHandlers -> DefaultHandler
func ExampleNotHanledMiddleware(next HandleUpdateFunc) HandleUpdateFunc {
	return func(tctx TgbotapiContext, update Update) error {
		// do something
		return fmt.Errorf("some error")
	}
}

// ExampleTimerMiddleware calc time
func ExampleTimerMiddleware(next HandleUpdateFunc) HandleUpdateFunc {
	return func(tctx TgbotapiContext, update Update) error {
		t := time.Now().Unix()
		r := next(tctx, update)
		fmt.Printf("\nTime for handle: %v\n", time.Now().Unix()-t)
		return r
	}
}
