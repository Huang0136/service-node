package service

/*

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// 加法
func (s *ServiceTest) Addition(req Req, resp *Resp) error {
	//	val1 := int(req.Params["val1"])
	//	val2 := int(req.Params["val2"])
	//	fmt.Println("入参,val1:", val1, " val2:", val2)

	//	resp.Params["ret"] = val1 + val2
	fmt.Println("Addition...", time.Now())

	for key, val := range req.Params {
		fmt.Println("key:", key, " value:", val)
	}

	fmt.Println("req:", req)

	resp.Params = make(map[string]interface{})
	resp.Params["status"] = int(0)

	str, err := json.Marshal(req.Params)
	if err != nil {
		log.Fatalln(err)
	}

	resp.Params["resutl"] = str
	resp.Params["time"] = time.Now().Format("2006-01-02 15:04:05.9999")

	fmt.Println("resp:", resp)

	return nil
}
*/
