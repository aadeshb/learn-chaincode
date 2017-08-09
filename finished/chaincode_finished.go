package main

import (
	"errors"
	"fmt"
	"encoding/json"
	"net/http"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)


type SimpleChaincode struct {
}

type Input struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	DOB       string `json:"DOB"`
	Email     string `json:"email"`
	Mobile    string `json:"mobile"`
}


func test(rw http.ResponseWriter, req *http.Request) {
    var t Input
    if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
  	fmt.Println(err)
	}	
    fmt.Fprintf(rw, "%s\n", req.Body)
}

func main() {
    http.HandleFunc("/test", test)
    http.ListenAndServe(":8080", nil)
}



