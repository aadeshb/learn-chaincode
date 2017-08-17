package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"time"
	"encoding/json"

	//"strconv"
)

//==============================================================================================================================================================

//																			Structs

//=============================================================================================================================================================

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}



type user struct {
	//ObjectType string `json:"docType"`
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	userID	   string `json:"userid"`
	DOB        string `json:"dob"`
	Email      string `json:"email"`
	Mobile     string `json:"mobile"`
	Class	   string `json:"class"`
}
// transfer function will append new owner and move current owner to previous owner, 
type RawMaterial struct {
	//ObjectType string `json:"docType"`
	Creator  		string `json:"creator"`
	Current_Owner   string `json:"currentowner"`
	//Previous_Onwer  string `json:"previousowner"`
	//State 			string `json:"state"`
	ClaimTags       string `json:"claimtags"`
	Location      	string `json:"location"`
	Date     		string `json:"date"`
	CertID	   		string `json:"certid"`
	Referencer		string `json:"referencer"`
}


type PurchaseOrder struct{

	Customer  		string `json:"customer"`
	Vendor   		string `json:"vendor"`
	ProductID   	string `json:"productid"`
	Price       	string `json:"price"`
	Date        	string `json:"date"`
	PurchaseOrderID	string `json:"purchaseorderid"`
}


//================================================================================================================================================================

//																	Main, Init, Invoke, Query

//================================================================================================================================================================




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
	} else if function == "RegisterRM" {
		return t.RegisterRM(stub, args)
	} else if function == "makePurchaseOrder" {
		return t.makePurchaseOrder(stub, args)
	} else if function == "replyPurchaseOrder" {
		return t.replyPurchaseOrder(stub, args)
	} else if function == "TransferAsset" {
		return t.TransferAsset(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)
    
	if function == "read" { 
		return t.read(stub, args)
	} 



	fmt.Println("query did not find func: " + function)
	return nil, errors.New("Received unknown function query: " + function)
}



//===============================================================================================================================================================

//															REGISTRATION CODE

//===============================================================================================================================================================
// write - invoke function to write key/value pair
func (t *SimpleChaincode) Register(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var v5User user
	var username = args[0]
	v5User.Firstname = args[1]
    v5User.Lastname = args[2]
	v5User.DOB = args[3]
	v5User.Email = args[4]
	v5User.Mobile = args[5]
	v5User.Class = args[6]

	bytes, err := json.Marshal(v5User)
    
    if err != nil { return nil, errors.New("Error creating v5User record") }

	err = stub.PutState(username, bytes)
	return nil, nil
}



func (t *SimpleChaincode) RegisterRM(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var r RawMaterial
	var prodid = args[0]

	var a = time.Now()
	var b = a.Format("20060102150405") 

	r.Creator = args[1]
    r.Current_Owner = args[2]
	r.ClaimTags = args[3]
	r.Location = args[4]
	r.Date = args[5]
	r.CertID = args[6]
	r.Referencer = prodid + "-" + b

	bytes, err := json.Marshal(r)
    
    if err != nil { return nil, errors.New("Error creating raw material") }

	err = stub.PutState(prodid, bytes)


	return nil, nil
}


//====================================================================================================================================================================

//																PURCHASE ORDERS

//====================================================================================================================================================================

// Purchase order code and "write" code are the exact same, because in essence, both should do the same job, which is to write
// data to the ledger which can be read later on 
// makePurchaseOrder has two user given inputs, 1 - supplier id, 2- manufacturer id
func (t *SimpleChaincode) makePurchaseOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){

	var p PurchaseOrder
	var poid = args[0]

	var a = time.Now()
	var b = a.Format("20060102150405") 

	//var userkeycombo = username + "-" + b

	p.Customer = args[1]
    p.Vendor = args[2]
	p.ProductID = args[3]
	p.Price = args[4]
	p.Date = args[5]
	p.PurchaseOrderID = poid + "-" + b
	//r.Referencer = userkeycombo

	bytes, err := json.Marshal(p)
    
    if err != nil { return nil, errors.New("Error creating raw material") }

	err = stub.PutState(poid, bytes)
	return nil, nil
}



func (t *SimpleChaincode) replyPurchaseOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){

	var key, value string
	var err error
	
	var a = time.Now()
	var b = a.Format("20060102150405") 
	key = args[0] 
	var body = args[2] //this will be the yes or no
	value = args[1] + "-" + b +"-"+  key + " " + body



	//here will be the automatic transfer functions calling









	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}



// ==================================================================================================================================================================

// 																	SENDING GOODS/TRANSFERING ASSETS

// ==================================================================================================================================================================


func (t *SimpleChaincode) TransferAsset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	
	//      0       		 1
	// "ProductID", "new owner name"
	
	
	RMAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get thing")
	}

	res := RawMaterial{}
	json.Unmarshal(RMAsBytes, &res)										//un stringify it aka JSON.parse()
	res.Current_Owner = args[1]														//change the user
	
	jsonAsBytes, _ := json.Marshal(res)
	err = stub.PutState(args[0], jsonAsBytes)								//rewrite the marble with id as key
	if err != nil {
		return nil, err
	}
	
	fmt.Println("- end set user")
	return nil, nil
}


//========================================================================================================================================================================

//															Certificate Stuff

//========================================================================================================================================================================


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



//===========================================================================================================================================================================

//																			Read

//===========================================================================================================================================================================




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


