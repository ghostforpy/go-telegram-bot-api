package tgbotapi

import (
	"fmt"
	"regexp"
	"strings"
)

type ConvHandler struct {
	EntryPoints        []CommonHandler
	States             map[string][]CommonHandler
	Fallbacks          []CommonHandler
	Name               string
	stateCheckerRegexp *regexp.Regexp
	AllowReentry       bool
}

func NewConvHandler(name string, entryPoints []CommonHandler, states map[string][]CommonHandler) *ConvHandler {
	keys := make([]string, 0, len(states))
	if len(entryPoints) == 0 {
		panic(fmt.Sprintf("ConvHandler %v coudn't be empty entry points", name))
	}
	if len(states) == 0 {
		panic(fmt.Sprintf("ConvHandler %v coudn't be empty states", name))
	}
	internaleStates := make(map[string][]CommonHandler)
	for key, val := range states {
		keys = append(keys, key)
		internaleStates[fmt.Sprintf("%v:%v", name, key)] = val
	}
	r, err := regexp.Compile(fmt.Sprintf(
		"^%v:(%v)$",
		name,
		strings.Join(keys, "|"),
	))
	if err != nil {
		panic(fmt.Sprintf("ConvHandler %v not compile", name))
	}
	return &ConvHandler{
		Name:               name,
		EntryPoints:        entryPoints,
		States:             internaleStates,
		stateCheckerRegexp: r,
	}
}

func (cvh *ConvHandler) HandleUpdate(tctx TgbotapiContext, update Update) (HandleUpdateFunc, error) {
	state := tctx.UserState
	if state == "" || state == "main" || cvh.AllowReentry {
		for _, handler := range cvh.EntryPoints {
			if ok, err := handler.CheckUpdate(update); err == nil && ok {
				return handler.HandleUpdate, nil
			}
		}
	} else if cvh.stateCheckerRegexp.Match([]byte(state)) {
		// states
		for _, handler := range cvh.States[state] {
			if ok, err := handler.CheckUpdate(update); err == nil && ok {
				return handler.HandleUpdate, nil
			}
		}
		// fallbacks
		for _, handler := range cvh.Fallbacks {
			if ok, err := handler.CheckUpdate(update); err == nil && ok {
				return handler.HandleUpdate, nil
			}
		}
	}
	return nil, nil
}
