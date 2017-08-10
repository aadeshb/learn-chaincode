package main

import (
	//"bytes"
	"encoding/json"
	"fmt"
	//"strconv"
	//"strings"
	//"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
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
	if function == "initUser" { //create a new user
		return t.initUser(stub, args)
	} 

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// inituser - create a new user, store into chaincode state
// ============================================================
func (t *SimpleChaincode) initUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
	if err != nil {
		return shim.Error("Failed to get user: " + err.Error())
	} else if userAsBytes != nil {
		fmt.Println("This user already exists: " + firstName)
		return shim.Error("This user already exists: " + firstName)
	}

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
