package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type reqParams struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type resParams struct {
	Code int `json:"code"`
	Res  int `json:"res"`
}

func main() {
	url := `http://localhost:18090/add`
	p := reqParams{
		X: 1,
		Y: 2,
	}
	b, _ := json.Marshal(p)
	resp, _ := http.Post(url, "application/json", bytes.NewReader(b))
	resB, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	var res resParams
	json.Unmarshal(resB, &res)
	fmt.Println(res.Res)
}
