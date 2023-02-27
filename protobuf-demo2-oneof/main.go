package main

import (
	"fmt"
	"protobuf-demo/api"
)

func oneofDemo() {
	req1 := &api.NoticeReaderRequest{
		Msg: "blog is up to date.",
		NoticeWay: &api.NoticeReaderRequest_Email{
			Email: "shiotya@hotmail.com",
		},
	}

	_ = &api.NoticeReaderRequest{
		Msg: "blog is up to date.",
		NoticeWay: &api.NoticeReaderRequest_Phone{
			Phone: "1222222222",
		},
	}

	switch req1.NoticeWay.(type) {
	case *api.NoticeReaderRequest_Email:
		fmt.Println("send by email.")
	case *api.NoticeReaderRequest_Phone:
		fmt.Println("send by phone.")
	}
}

func main() {
	oneofDemo()
}
