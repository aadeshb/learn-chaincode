package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

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
	ObjectType string `json:"docType"`
}



// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "Register" { //create a new marble
		return t.Register(stub, args)
	} 
	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// initMarble - create a new marble, store into chaincode state
// ============================================================
func (t *SimpleChaincode) Register(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

//	Firstname  string `json:"firstname"`		0
//	Lastname   string `json:"lastname"`			1
//	userID	   string `json:"userid"`			2
//	DOB        string `json:"dob"`				3
//	Email      string `json:"email"`			4
//	Mobile     string `json:"mobile"`			5
//	Class	   string `json:"class"`			6
	
	//   0       1       2     3
	// "asdf", "blue", "35", "bob"
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	// ==== Input sanitation ====
	
	fname := args[0]
	lname := args[1]
	uid := args[2]
	userdob := args[3]
	useremail := args[4]
	usermobile := args[5]
	userclass := args[6]

	

	// ==== Check if user already exists ====
	fnameAsBytes, err := stub.GetState(fname)
	if err != nil {
		return shim.Error("Failed to get marble: " + err.Error())
	} else if fnameAsBytes != nil {
		fmt.Println("This marble already exists: " + fname)
		return shim.Error("This marble already exists: " + fname)
	}

	// ==== Create user object and marshal to JSON ====
	objectType := "user"
	user := &user{fname, lname, uid, userdob, useremail, usermobile, userclass}
	userJSONasBytes, err := json.Marshal(user)
	if err != nil {
		return shim.Error(err.Error())
	}
	//Alternatively, build the marble json string manually if you don't want to use struct marshalling
	//marbleJSONasString := `{"docType":"Marble",  "name": "` + marbleName + `", "color": "` + color + `", "size": ` + strconv.Itoa(size) + `, "owner": "` + owner + `"}`
	//marbleJSONasBytes := []byte(str)

	// === Save user to state ===
	err = stub.PutState(uid, userJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//  ==== Index the marble to enable color-based range queries, e.g. return all blue marbles ====
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on indexName~color~name.
	//  This will enable very efficient state range queries based on composite keys matching indexName~color~*
	indexName := "uid~fname"
	uidIndexKey, err := stub.CreateCompositeKey(indexName, []string{user.userID, user.fname})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(uidIndexKey, value)

	// ==== Marble saved and indexed. Return success ====
	fmt.Println("- end init marble")
	return shim.Success(nil)
}
