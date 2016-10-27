package main

import (
	"fmt"
	"reflect"
	"service/impl"
)

func main() {
	si := new(impl.ServiceImpl)

	fmt.Println(si)
	vv := reflect.ValueOf(si)

	nm := vv.NumMethod()
	fmt.Println("NumMethod:", nm)

	for i := 0; i < nm; i++ {
		f := vv.Method(i)
		fmt.Println("index:", i, ",", f.)
	}

}

func init() {

}

/*
	log.Println("serverTest start...")
	log.Println("正在注册rpc服务")

	myServiceTest := new(service.ServiceTest)
	rpc.Register(myServiceTest)

	// 注册信息到注册中心
	register.RegisterInfo()

	// http方式
	rpc.HandleHTTP()
	err := http.ListenAndServe(":9877", nil)

	if err != nil {
		log.Fatalln(err.Error())
	}
*/
