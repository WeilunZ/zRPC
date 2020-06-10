package zRPC

import (
	"fmt"
	"reflect"
	"testing"
)

func TestReflect(t *testing.T) {
	f := func(a string, b int) (bool, error) {
		return true, nil
	}
	ft := reflect.TypeOf(f)
	for i := 0; i < ft.NumIn(); i++ {
		fmt.Println(ft.In(i).Name())
	}
	for i := 0; i < ft.NumOut(); i++ {
		fmt.Println(ft.Out(i).Name())
	}
}
