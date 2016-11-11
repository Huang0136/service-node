// 读取配置文件的值
package constants

import (
	"bufio"
	"io"
	"logs"
	"net"
	"os"
	"strings"
)

// 节点配置信息
var Configs map[string]string = make(map[string]string)

func init() {
	readServerNodeProperties()
	getIpv4()

	/*
		for k, v := range Configs {
			fmt.Printf("%s = %s\n", k, v)
		}*/
}

// 读取server-node配置文件
func readServerNodeProperties() {
	logs.MyInfoLog.Println("正在读取配置文件：server-node.properties")

	file, err := os.Open("./config/server-node.properties")
	logs.MyErrorLog.CheckPaniclnError("读取server-node配置文件出错:", err)

	read := bufio.NewReader(file)
	for {
		b, _, err := read.ReadLine()
		if err != nil || err == io.EOF {
			break
		}

		str := string(b)
		if !strings.HasPrefix(str, "#") && strings.Contains(str, "=") {
			strArray := strings.Split(str, "=")
			Configs[strings.TrimSpace(strArray[0])] = strings.TrimSpace(strArray[1])
		}
	}

	logs.MyInfoLog.Println("server-node.properties配置文件读取完成!")
}

// 获取本机IP
func getIpv4() {
	/*
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			logs.MyErrorLog.CheckFatallnError("获取本地ip失败", err)
		}

		for _, addr := range addrs {
			fmt.Println(addr)
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					Configs["serverNode.ip"] = ipnet.IP.String()
				}
			}
		}
	*/

	conn, err := net.Dial("udp", "www.hao123.com:80")
	logs.MyErrorLog.CheckPaniclnError("连接到www.hao123.com失败", err)

	Configs["serverNode.ip"] = strings.Split(conn.LocalAddr().String(), ":")[0] // 本机的IP
}
