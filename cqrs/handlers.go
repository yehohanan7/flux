package cqrs

import (
	"reflect"
	"strings"

	"github.com/yehohanan7/flux/utils"
)

type Handlers map[reflect.Type]func(interface{}, interface{})

func isHandlerMethod(method reflect.Method) bool {
	return method.Type.NumIn() == 2 && strings.HasPrefix(method.Name, "Handle")
}

func createHandler(method reflect.Method) func(interface{}, interface{}) {
	return func(entity interface{}, payload interface{}) {
		method.Func.Call([]reflect.Value{reflect.ValueOf(entity), reflect.ValueOf(payload)})
	}
}

func NewHandlers(entity interface{}) Handlers {
	handlers := make(Handlers)
	for _, method := range utils.FindMethods(entity, isHandlerMethod) {
		eventType := method.Type.In(1)
		handlers[eventType] = createHandler(method)
	}
	return handlers
}
