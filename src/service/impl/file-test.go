package impl

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

func (si *ServiceImpl) DownloadFileTest() (fByte []byte, other map[string]interface{}, err error) {
	fileName := si.InParams["file_name"].(string)
	t1 := time.Now().UnixNano()
	fByte, err = ioutil.ReadFile("./config/" + fileName)
	if err != nil {
		log.Panicln("读取文件失败", err)
		return
	}
	t2 := time.Now().UnixNano()

	other = make(map[string]interface{})
	other["READ_FILE_BEGIN_TIME"] = t1
	other["READ_FILE_END_TIME"] = t2
	fmt.Println("读取文件成功")
	return
}
