package cqrs

import "reflect"

type HandlerMap map[reflect.Type]func(interface{}, interface{})

func isHandlerMethod(method reflect.Method) bool {
	return method.Type.NumIn() == 2 && method.Name == "Handle"
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

func buildHandlerMap(entity interface{}) HandlerMap {
	handlers := make(HandlerMap)
	for _, method := range handlerMethods(entity) {
		eventType := method.Type.In(1)
		handlers[eventType] = func(entity interface{}, payload interface{}) {
			method.Func.Call([]reflect.Value{reflect.ValueOf(entity), reflect.ValueOf(payload)})
		}
	}
	return handlers
}
