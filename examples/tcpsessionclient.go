package main

import (
	"fmt"
	"github.com/hprose/hprose-go/hprose"
)

type Stub struct {
	Inc func() int
}

type MyClientFilter map[interface{}]int

func (filter MyClientFilter) InputFilter(data []byte, context interface{}) []byte {
	if len(data) > 7 && data[0] == 's' && data[1] == 'i' && data[2] == 'd' {
		sessionid := int(data[3])<<24 | int(data[4])<<16 | int(data[5])<<8 | int(data[6])
		filter[context] = sessionid
		data = data[7:]
	}
	return data
}

func (filter MyClientFilter) OutputFilter(data []byte, context interface{}) []byte {
	if sessionid, ok := filter[context]; ok {
		buf := make([]byte, 7+len(data))
		buf[0] = 's'
		buf[1] = 'i'
		buf[2] = 'd'
		buf[3] = byte(sessionid >> 24 & 0xff)
		buf[4] = byte(sessionid >> 16 & 0xff)
		buf[5] = byte(sessionid >> 8 & 0xff)
		buf[6] = byte(sessionid & 0xff)
		copy(buf[7:], data)
		return buf
	}
	return data
}

func main() {
	client := hprose.NewClient("tcp4://127.0.0.1:4321/").(*hprose.TcpClient)
	client.AddFilter(MyClientFilter{})
	var stub Stub
	client.UseService(&stub)
	go fmt.Println(stub.Inc())
	go fmt.Println(stub.Inc())
	go fmt.Println(stub.Inc())
	go fmt.Println(stub.Inc())
	fmt.Println(stub.Inc())
}
