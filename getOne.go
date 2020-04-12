package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var sum int64 = 0

var productNum int64 = 10000

var mutex sync.Mutex

var count int64 = 0

func GetOneProduct() bool {
	mutex.Lock()
	defer mutex.Unlock()

	//count += 1

	//if count % 100 == 0 {
		if sum < productNum {
			sum += 1
			return true
		}
	//}


	return false
}

func GetProduct(w http.ResponseWriter, req *http.Request) {
	if GetOneProduct() {
		fmt.Println("success")
		w.Write([]byte("true"))
		return
	}
	w.Write([]byte("false"))
	return
}

func main() {
	http.HandleFunc("/getOne", GetProduct)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("Err:", err)
	}
}