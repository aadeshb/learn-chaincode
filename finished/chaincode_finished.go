package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"time"
	"encoding/json"

	//"strconv"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}


type user struct {
	//ObjectType string `json:"docType"`
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	DOB        string `json:"dob"`
	Email      string `json:"email"`
	Mobile     string `json:"mobile"`
	Class	   string `json:"class"`
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
	} else if function == "Register" {
		return t.Register(stub, args)
	} else if function == "makePurchaseOrder" {
		return t.makePurchaseOrder(stub, args)
	} else if function == "replyPurchaseOrder" {
		return t.makePurchaseOrder(stub, args)
	} 
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)
    //var u user

	// Handle different functions
	// If the read function is called the read function activiates
	if function == "read" { //read a variable

		return t.read(stub, args)
	} //else if function == "retrieve" {
		//return t.retrieve(stub, args)
	//}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}


func (t *SimpleChaincode) retrieve(stub shim.ChaincodeStubInterface, args []string) (user, error) {
	
	var v user

	bytes, err := stub.GetState(args[0]);

	//if err != nil {	fmt.Printf("RETRIEVE_V5C: Failed to invoke vehicle_code: %s", err); return v, errors.New("RETRIEVE_V5C: Error retrieving vehicle with v5cID = " + args) }

	err = json.Unmarshal(bytes, &v);

    if err != nil {	fmt.Printf("RETRIEVE_V5C: Corrupt vehicle record "+string(bytes)+": %s", err); return v, errors.New("RETRIEVE_V5C: Corrupt vehicle record"+string(bytes))	}

	return v, nil


}



//============================================================================================================================================

//															REGISTRATION CODE

//=============================================================================================================================================
// write - invoke function to write key/value pair
func (t *SimpleChaincode) Register(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var v5User user
	v5User.Firstname = args[0]
    v5User.Lastname = args[1]
	v5User.DOB = args[2]
	v5User.Email = args[3]
	v5User.Mobile = args[4]
	v5User.Class = args[5]

	bytes, err := json.Marshal(v5User)
    
    if err != nil { return nil, errors.New("Error creating v5User record") }

	err = stub.PutState("v5User", bytes)
	return nil, nil
}

// Purchase order code and "write" code are the exact same, because in essence, both should do the same job, which is to write
// data to the ledger which can be read later on 
// makePurchaseOrder has two user given inputs, 1 - supplier id, 2- manufacturer id
func (t *SimpleChaincode) makePurchaseOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){

	var key, value string
	var err error
	
	var a = time.Now()
	var b = a.Format("20060102150405") 
	key = args[0] //the key is simply the suppliers id
	var body = args[2]
	value = args[1] + "-" + b +"-"+  key + " " + body
	//var comm string = value + b + key
	
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}



func (t *SimpleChaincode) replyPurchaseOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){

	var key, value string
	var err error
	
	var a = time.Now()
	var b = a.Format("20060102150405") 
	key = args[0] //the key is simply the suppliers id
	var body = args[2]
	value = args[1] + "-" + b +"-"+  key + " " + body
	//var comm string = value + b + key
	
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}


func (t *SimpleChaincode) awardCertificate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){

	var key, value string
	var err error
	fmt.Println("running write()")


	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename 
	value = args[1]
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



func (t *SimpleChaincode) viewPurchaseOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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


