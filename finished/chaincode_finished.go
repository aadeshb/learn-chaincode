package main

import (
	//"bytes"
	"encoding/json"
	"fmt"
	//"strconv"
	//"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//=============================================================================================================================================================================

//																					Structs

//=============================================================================================================================================================================



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
	Item			string `json:"item"`
	Creator  		string `json:"creator"`
	Current_Owner   string `json:"currentowner"`
	ClaimTags       string `json:"claimtags"`
	Location      	string `json:"location"`
	Date     		string `json:"date"`
	CertID	   		string `json:"certid"`
	ObjectType      string `json:"docType"`
}


type FinishedGood struct {
	FPID			string `json:"fpid"`
	Name 			string `json:"name"`
	Creator  		string `json:"creator"`
	Current_Owner   string `json:"currentowner"`
	Ingredients 	string `json:"ingredients"`
	//Previous_Owner  string `json:"previousowner"`
	Certificates	string `json:"certificates"`
	ClaimTags       string `json:"claimtags"`
	Location      	string `json:"location"`
	Date     		string `json:"date"`
	CertID	   		string `json:"certid"`
	ObjectType 		string `json:"docType"`
}

type PurchaseOrder struct{

	PurchaseOrderID	string `json:"purchaseorderid"`
	Customer  		string `json:"customer"`
	Vendor   		string `json:"vendor"`
	ProductID   	string `json:"productid"`
	Price       	string `json:"price"`
	Date        	string `json:"date"`
	ObjectType 		string `json:"docType"`
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
		return t.RegisterRM(stub, args) 
	} else if function == "RegisterFP" { 
		return t.RegisterFP(stub, args) 
	} else if function == "makePurchaseOrder" { 
		return t.makePurchaseOrder(stub, args) 
	} else if function == "replyPurchaseOrder" { 
		return t.replyPurchaseOrder(stub, args) 
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
//  Item 			string `json:"item"`					1
//	Creator  		string `json:"creator"`					2
//	Current_Owner   string `json:"currentowner"`			3
//	ClaimTags       string `json:"claimtags"`				4
//	Location      	string `json:"location"`				5
//	Date     		string `json:"date"`					6
//	CertID	   		string `json:"certid"`					7
//	ObjectType      string `json:"docType"`					8
	
	
	

	// ==== Input sanitation ====
	
	rawid := args[0]
	item := args[1]
	originalcreator := args[2]
	cowner := args[3]
	claimtags := args[4]
	loc := args[5]
	dates := args[6]
	userclass := args[7]
	

	
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
	RawMaterial := &RawMaterial{rawid, item, originalcreator, cowner, claimtags, loc, dates, userclass, objectType}
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


func (t *SimpleChaincode) RegisterFP(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

//	ObjectType 		string `json:"docType"`					0
//	FPID			string `json:"fpid"`					1	
//	Name 			string `json:"name"`					2
//	Creator  		string `json:"creator"`					3
//	Current_Owner   string `json:"currentowner"`			4
//	Ingredients 	string `json:"ingredients"`				5
//	//Previous_Owner  string `json:"previousowner"`			6
//	Certificates	string `json:"certificates"`			7
//	ClaimTags       string `json:"claimtags"`				8
//	Location      	string `json:"location"`				9	
//	Date     		string `json:"date"`					10
//	CertID	   		string `json:"certid"`					11
	
	
	

	// ==== Input sanitation ====
	
	fpid_i := args[0]
	name_i := args[1]
	originalcreator_i := args[2]
	cowner_i := args[3]
	ingredients_i := args[4]
	certificates_i := args[5]
	claimtags_i := args[6]
	loc_i := args[7]
	dates_i := args[8]
	certid_i := args[9]
	

	// This wont matter once we implement UUID 
	// ==== Check if user already exists ====
	fpid_iAsBytes, err := stub.GetState(fpid_i)		
	if err != nil {
		return shim.Error("Failed to get user: " + err.Error())
	} else if fpid_iAsBytes != nil {
		fmt.Println("This user already exists: " + fpid_i)
		return shim.Error("This user already exists: " + fpid_i)
	}

	// ==== Create user object and marshal to JSON ====
	objectType := "FinishedGood"
	FinishedGood := &FinishedGood{fpid_i, name_i, originalcreator_i, cowner_i, ingredients_i, certificates_i, claimtags_i, loc_i, dates_i, certid_i, objectType}
	FinishedGoodJSONasBytes, err := json.Marshal(FinishedGood)
	if err != nil {
		return shim.Error(err.Error())
	}
	

	// === Save user to state ===
	err = stub.PutState(fpid_i, FinishedGoodJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//  ==== Index the user to enable color-based range queries, e.g. return all blue users ====
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on indexName~color~name.
	//  This will enable very efficient state range queries based on composite keys matching indexName~color~*
	indexName := "fpid_i~cowner"
	fpiIndexKey, err := stub.CreateCompositeKey(indexName, []string{FinishedGood.FPID, FinishedGood.Current_Owner})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the user.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(fpiIndexKey, value)

	// ==== user saved and indexed. Return success ====
	fmt.Println("- end init user")
	return shim.Success(nil)
}

func (t *SimpleChaincode) makePurchaseOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

//	PurchaseOrderID	string `json:"purchaseorderid"`		0
//	Customer  		string `json:"customer"`			1
//	Vendor   		string `json:"vendor"`				2
//	ProductID   	string `json:"productid"`			3
//	Price       	string `json:"price"`				4
//	Date        	string `json:"date"`				5
//	ObjectType 		string `json:"docType"`				6
	
	// ==== Input sanitation ====
	
	purchid := args[0]
	cust := args[1]
	vend := args[2]
	prodid := args[3]
	price:= args[4]
	dat := args[5]
	

	
	// This wont matter once we implement UUID 
	// ==== Check if product already exists ====
	purchAsBytes, err := stub.GetState(purchid)		
	if err != nil {
		return shim.Error("Failed to get product: " + err.Error())
	} else if purchAsBytes != nil {
		fmt.Println("This product already exists: " + purchid)
		return shim.Error("This product already exists: " + purchid)
	}

	// ==== Create user object and marshal to JSON ====
	objectType := "PurchaseOrder"
	PurchaseOrder := &PurchaseOrder{purchid, cust, vend, prodid, price, dat, objectType}
	prodJSONasBytes, err := json.Marshal(PurchaseOrder)
	if err != nil {
		return shim.Error(err.Error())
	}
	

	// === Save user to state ===
	err = stub.PutState(purchid, prodJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== user saved and indexed. Return success ====
	fmt.Println("- end init user")
	return shim.Success(nil)
}

func (t *SimpleChaincode) replyPurchaseOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	var key, value string
	
	
	var a = time.Now()
	var b = a.Format("20060102150405") 
	key = args[0] 
	var body = args[2] //this will be the yes or no
	value = args[1] + "-" + b +"-"+  key + " " + body


	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return shim.Error(err.Error())
	}
	
	return shim.Success(nil)
}






//===========================================================================================================================================================================

//																				Reading

//===========================================================================================================================================================================


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
