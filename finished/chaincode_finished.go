package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/nu7hatch/gouuid"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}


// The main function is used to bootstrap the code, however we don't have any functionality for it right now
// it only reports if an error occurs, which never should
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}


}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	//If we are given more than one argument for the init function, then it errors
	//Init shouldn't need any arguments at all actually
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	//We are writing some key under the name "hello_world" to the ledger, the hello world is simply a key to hold whatever
	//argument we pass, if we query "hello_world" with the query function, then it would return whatever argument is placed in args[0]
	err := stub.PutState("hello_world", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke is the entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	//If the init function is called, then we send the args to the init command to be stored under the "hello_world" key
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "sOne" {
		return t.sOne(stub, args)
	} else if function == "purchaseOrder" {
		return t.makePurchaseOrder(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	// If the read function is called the read function activiates
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

//=============================================================================================================================================
// write - invoke function to write key/value pair
func (t *SimpleChaincode) RegisterSupplies(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")


	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	 //rename 
	value = args[1]
	var id, err := uuid.NewV4()
	err = stub.PutState(id, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}


// Purchase order code and "write" code are the exact same, because in essence, both should do the same job, which is to write
// data to the ledger which can be read later on 

func (t *SimpleChaincode) makePurchaseOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){

	var key, value string
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename 
	value = args[1]
	var id, err := uuid.NewV4()

	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil







}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
