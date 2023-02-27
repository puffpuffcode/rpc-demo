package main

import (
	"fmt"
	"protobuf_demo/api"

	"github.com/iancoleman/strcase"
	fieldmask_utils "github.com/mennanov/fieldmask-utils"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func wrapValueDemo() {
	// cli
	book := &api.Book{
		Title:  "A Book",
		Author: "abc",
		// optional field
		Color:     proto.String("yellow"),
		Price:     &wrapperspb.Int64Value{Value: 9999},
		SalePrice: &wrapperspb.DoubleValue{Value: 9999.99},
		Memo:      &wrapperspb.StringValue{Value: "高い！"},
	}

	// server
	if book.Price == nil { // No value is assigned to it
		fmt.Println("没有设置值！")
	} else {
		fmt.Println("有值了:", book.GetPrice().GetValue())
	}

	// optional f check
	if book.Color != nil {
		fmt.Println("Color:", book.GetColor())
	}
}

// 实现部分更新
func fieldMaskDemo() {
	// cli
	req := &api.UpdateBookRequest{
		Op: "makito",
		Book: &api.Book{
			Price: &wrapperspb.Int64Value{Value: 22},
			Info: &api.Book_Info{
				B: "2",
			},
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"price","info.b"}},
	}

	// server
	mask, _ := fieldmask_utils.MaskFromProtoFieldMask(req.UpdateMask, strcase.ToCamel)
	resMap := make(map[string]interface{})
	// 将部分更新的字段数据读取到 map 中
	err := fieldmask_utils.StructToMap(mask, req.Book, resMap)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n",resMap)
}

func main() {
	// wrapValueDemo()
	fieldMaskDemo()
}
