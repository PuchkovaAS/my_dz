package main

import (
	"fmt"
	"math/rand"
	"net/http"
)

func randomDice(w http.ResponseWriter, req *http.Request) {
	randomNum := rand.Intn(6) + 1
	buf := fmt.Append(nil, randomNum)
	w.Write(buf)
}

func main() {
	http.HandleFunc("/", randomDice)
	fmt.Println("Server started at :8081")
	http.ListenAndServe(":8081", nil)
}
