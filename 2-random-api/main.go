package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
)

func randomDice(w http.ResponseWriter, req *http.Request) {
	randomNum := rand.Intn(6) + 1
	w.Write([]byte(strconv.Itoa(randomNum)))
}

func main() {
	http.HandleFunc("/", randomDice)
	fmt.Println("Server started at :8081")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
