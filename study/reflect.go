package main

import (
	"fmt"
	"reflect"
)

func main() {
	 x := 3.2
	fmt.Println("type:", reflect.TypeOf(x))
	fmt.Println("value:", reflect.ValueOf(x))

	v := reflect.ValueOf(x)
	fmt.Println("type:", v.Type())
	fmt.Println("kind is float64:", v.Kind() == reflect.Float64)
	fmt.Println("value:", v.Float())

	g := reflect.ValueOf(&x)
	k := g.Elem()
	fmt.Println("can set:" ,k.CanSet())
	k.SetFloat(3.1)
	fmt.Println(k.Interface())
	fmt.Println(x)

	type T struct {
		A int
		B string
	}
	t := T{23, "skidoo"}
	s := reflect.ValueOf(&t).Elem()
	typeOfT := s.Type()//把s.Type()返回的Type对象复制给typeofT，typeofT也是一个反射。
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)//迭代s的各个域，注意每个域仍然是反射。
		fmt.Printf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())//提取了每个域的名字
	}
}


