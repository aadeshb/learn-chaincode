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
	
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	userID	   string `json:"userid"`
	DOB        string `json:"dob"`
	Email      string `json:"email"`
	Mobile     string `json:"mobile"`
	Class	   string `json:"class"`
	ObjectType string `json:"docType"`
}

type RawMaterial struct {
	
	RMID 			string `json:"rmid"`
	Creator  		string `json:"creator"`
	Current_Owner   string `json:"currentowner"`
	ClaimTags       string `json:"claimtags"`
	Location      	string `json:"location"`
	Date     		string `json:"date"`
	CertID	   		string `json:"certid"`
	Referencer		string `json:"referencer"`
	ObjectType      string `json:"docType"`
}


// =============================================================================================================================================================

// 																					MAIN FUNCTIONS

// ==============================================================================================================================================================
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
	if function == "Register" { //create a new user
		return t.Register(stub, args)
	} else if function == "RegisterRM" { 
		return t.read(stub, args) 
	} else if function == "read" { 
		return t.read(stub, args) 
	}
	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}


//============================================================================================================================================================================

//																				REGISTRATION CODE BELOW

//=============================================================================================================================================================================



func (t *SimpleChaincode) Register(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

//	Firstname  string `json:"firstname"`		0
//	Lastname   string `json:"lastname"`			1
//	userID	   string `json:"userid"`			2
//	DOB        string `json:"dob"`				3
//	Email      string `json:"email"`			4
//	Mobile     string `json:"mobile"`			5
//	Class	   string `json:"class"`			6
	
	
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

	
	// This wont matter once we implement UUID 
	// ==== Check if user already exists ====
	fnameAsBytes, err := stub.GetState(uid)		//Change this to uid not fname
	if err != nil {
		return shim.Error("Failed to get user: " + err.Error())
	} else if fnameAsBytes != nil {
		fmt.Println("This user already exists: " + fname)
		return shim.Error("This user already exists: " + fname)
	}

	// ==== Create user object and marshal to JSON ====
	objectType := "user"
	user := &user{fname, lname, uid, userdob, useremail, usermobile, userclass, objectType}
	userJSONasBytes, err := json.Marshal(user)
	if err != nil {
		return shim.Error(err.Error())
	}
	

	// === Save user to state ===
	err = stub.PutState(uid, userJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//  ==== Index the user to enable color-based range queries, e.g. return all blue users ====
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on indexName~color~name.
	//  This will enable very efficient state range queries based on composite keys matching indexName~color~*
	indexName := "uid~fname"
	uidIndexKey, err := stub.CreateCompositeKey(indexName, []string{user.userID, user.Firstname})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the user.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(uidIndexKey, value)

	// ==== user saved and indexed. Return success ====
	fmt.Println("- end init user")
	return shim.Success(nil)
}


func (t *SimpleChaincode) RegisterRM(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

//	RMID 			string `json:"rmid"`					0
//	Creator  		string `json:"creator"`					1
//	Current_Owner   string `json:"currentowner"`			2
//	ClaimTags       string `json:"claimtags"`				3
//	Location      	string `json:"location"`				4
//	Date     		string `json:"date"`					5
//	CertID	   		string `json:"certid"`					6
//	Referencer		string `json:"referencer"`				7
//	ObjectType      string `json:"docType"`					8
	
	
	

	// ==== Input sanitation ====
	
	rawid := args[0]
	originalcreator := args[1]
	cowner := args[2]
	claimtags := args[3]
	loc := args[4]
	dates := args[5]
	userclass := args[6]
	ref := args[7]

	
	// This wont matter once we implement UUID 
	// ==== Check if user already exists ====
	rawidAsBytes, err := stub.GetState(rawid)		
	if err != nil {
		return shim.Error("Failed to get user: " + err.Error())
	} else if rawidAsBytes != nil {
		fmt.Println("This user already exists: " + rawid)
		return shim.Error("This user already exists: " + rawid)
	}

	// ==== Create user object and marshal to JSON ====
	objectType := "RawMaterial"
	RawMaterial := &RawMaterial{rawid, originalcreator, cowner, claimtags, loc, dates, userclass, ref, objectType}
	RawMaterialJSONasBytes, err := json.Marshal(RawMaterial)
	if err != nil {
		return shim.Error(err.Error())
	}
	

	// === Save user to state ===
	err = stub.PutState(rawid, RawMaterialJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//  ==== Index the user to enable color-based range queries, e.g. return all blue users ====
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on indexName~color~name.
	//  This will enable very efficient state range queries based on composite keys matching indexName~color~*
	indexName := "rawid~cowner"
	rawidIndexKey, err := stub.CreateCompositeKey(indexName, []string{RawMaterial.RMID, RawMaterial.Current_Owner})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the user.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(rawidIndexKey, value)

	// ==== user saved and indexed. Return success ====
	fmt.Println("- end init user")
	return shim.Success(nil)
}

func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	var err error

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}



	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}
