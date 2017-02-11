package cqrs

import (
	"reflect"
	"strings"
)

type handlermap map[reflect.Type]func(interface{}, interface{})

func isHandlerMethod(method reflect.Method) bool {
	return method.Type.NumIn() == 2 && strings.HasPrefix(method.Name, "Handle")
}

func handlerMethods(entity interface{}) []reflect.Method {
	entityType := reflect.TypeOf(entity)
	numberOfMethods := entityType.NumMethod()
	methods := make([]reflect.Method, 0, numberOfMethods)
	for i := 0; i < numberOfMethods; i++ {
		method := entityType.Method(i)
		if isHandlerMethod(method) {
			methods = append(methods, method)
		}
	}
	return methods
}

func createHandler(method reflect.Method) func(interface{}, interface{}) {
	return func(entity interface{}, payload interface{}) {
		method.Func.Call([]reflect.Value{reflect.ValueOf(entity), reflect.ValueOf(payload)})
	}
}

func buildHandlerMap(entity interface{}) handlermap {
	handlers := make(handlermap)
	for _, method := range handlerMethods(entity) {
		eventType := method.Type.In(1)
		handlers[eventType] = createHandler(method)
	}
	return handlers
}
