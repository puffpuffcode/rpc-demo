package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type addParams struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type resParams struct {
	Code int `json:"code"`
	Res  int `json:"res"`
}

func add(x, y int) int {
	return x + y
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	p, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll(r.Body) failed...\n")
		return
	}
	defer r.Body.Close()
	params := new(addParams)
	json.Unmarshal(p, params)
	res := add(params.X, params.Y)
	resBytes, _ := json.Marshal(resParams{
		Code: 200,
		Res:  res,
	})
	w.Write(resBytes)
}

func main() {
	http.HandleFunc("/add", addHandler)
	http.ListenAndServe(":18090", nil)
}
