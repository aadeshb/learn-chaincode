package main

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}




func main() {
	err := shim.Start(new(SampleChaincode))
	if err != nil {
	fmt.Println("Could not start SampleChaincode")
	} else {
	fmt.Println("SampleChaincode successfully started")
	}
}




func (t *SampleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
if len(args) != 1 {
       return nil, errors.New("Incorrect number of arguments. Expecting 1")
   }
   err := stub.PutState("emp_name", []byte(args[0]))
   if err != nil {
       return nil, err
   }
   return nil, nil
}


func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
   fmt.Println("invoke is running " + function)
   // Handle different functions
   if function == "init" {
       return t.Init(stub, "init", args)
   } else if function == "write" {
       return t.write(stub, args)
   }
   fmt.Println("invoke did not find func: " + function)
   return nil, errors.New("Received unknown function invocation")
}



func (t *SampleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)
	// Handle different functions
	if function == "dummy_query" { //read a variable
	fmt.Println("hi there " + function) //error
	return nil, nil;
	}

	fmt.Println("query did not find func: " + function) //error
	return nil, errors.New("Received unknown function query: " + function)
}
