package impl

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

func (si *ServiceImpl) DownloadFileTest() (fByte []byte, other map[string]interface{}, err error) {
	// 入参
	fileName := si.InParams["file_name"].(string)

	// 读取文件
	tReadFile1 := time.Now()
	fByte, err = ioutil.ReadFile("./config/" + fileName)
	if err != nil {
		log.Panicln("读取文件失败", err)
		return
	}
	tReadFile2 := time.Now()

	// 返回结果(性能分析)
	other = make(map[string]interface{})
	other["TIME_READ_FILE"] = tReadFile2.Sub(tReadFile1).Seconds()

	fmt.Println("读取文件成功")
	return
}
