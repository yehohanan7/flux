package utils

import "reflect"

func FindMethods(object interface{}, predicate func(reflect.Method) bool) []reflect.Method {
	entityType := reflect.TypeOf(object)
	numberOfMethods := entityType.NumMethod()
	methods := make([]reflect.Method, 0, numberOfMethods)
	for i := 0; i < numberOfMethods; i++ {
		method := entityType.Method(i)
		if predicate(method) {
			methods = append(methods, method)
		}
	}
	return methods
}
