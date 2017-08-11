package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"time"
	//"strconv"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type user struct {
	ObjectType string `json:"docType"`
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
	} else if function == "Register" {
		return t.Register(stub, args)
	} else if function == "makePurchaseOrder" {
		return t.makePurchaseOrder(stub, args)
	} else if function == "initUser"{
		return t.initUser(stub, args)
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
func (t *SimpleChaincode) Register(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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

func (t *SimpleChaincode) initUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error)  {
	var err error

	//   	0       	1         2       3        4		 5
	// "Firstname", "Lastname", "DOB", "Email", "Mobile", "Class"
	
	firstName := args[0]
	lastName := args[1]
	dob := args[2]
	email := args[3]
	class := args[5]
	mobile := args[4]
	

	// ==== Check if user already exists ====
	userAsBytes, err := stub.GetState(mobile)
	//if err != nil {
//		return shim.Error("Failed to get user: " + err.Error())
//	} else if userAsBytes != nil {
//		fmt.Println("This user already exists: " + firstName)
//		return shim.Error("This user already exists: " + firstName)
//	}

	// ==== Create user object and marshal to JSON ====
	objectType := "user"
	user := &user{objectType, firstName, lastName, dob, email, mobile, class}
	userJSONasBytes, err := json.Marshal(user)
	if err != nil {
		return shim.Error(err.Error())
	}
	//Alternatively, build the user json string manually if you don't want to use struct marshalling
	//userJSONasString := `{"docType":"user",  "name": "` + userName + `", "color": "` + color + `", "size": ` + strconv.Itoa(size) + `, "owner": "` + owner + `"}`
	//userJSONasBytes := []byte(str)

	// === Save user to state ===
	err = stub.PutState(firstName, userJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//  ==== Index the user to enable color-based range queries, e.g. return all blue users ====
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on indexName~color~name.
	//  This will enable very efficient state range queries based on composite keys matching indexName~color~*
	indexName := "class~name"
	classNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{user.Class, user.Firstname})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the user.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(classNameIndexKey, value)

	// ==== user saved and indexed. Return success ====
	fmt.Println("- end init user")
	return shim.Success(nil)
}

