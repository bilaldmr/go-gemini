package pkgFunction

import (
	"sync"
)

type ProcessRegistry struct {
	functions map[string]Function
}

var (
	instance *ProcessRegistry
	once     sync.Once
)

func GetRegistry() *ProcessRegistry {
	once.Do(func() {
		instance = &ProcessRegistry{
			functions: make(map[string]Function),
		}
	})
	return instance
}

func RegisterFunction(function Function) {
	GetRegistry().functions[function.Name] = function
}

func GetFunction(name string) (Function, bool) {
	function, ok := GetRegistry().functions[name]
	return function, ok
}

func GetAllFunctions() []Function {
	functions := make([]Function, 0, len(GetRegistry().functions))
	for _, function := range GetRegistry().functions {
		functions = append(functions, function)
	}
	return functions
}
