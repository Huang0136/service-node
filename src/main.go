package main

import (
	_ "service"

	_ "monitor"

	_ "github.com/coreos/etcd/clientv3"
	_ "golang.org/x/net/context"
)

var Count int = 0

func main() {

	for {

	}
}

func init() {

}
