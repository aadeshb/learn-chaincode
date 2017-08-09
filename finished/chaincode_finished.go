package main

import (
	"errors"
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type LocationReadable struct {
	District string
	City     string
	State    string
}

type Locale struct {
	Location         string
	CountryCode      string
	CurrencyId       string
	CurrencySymbol   string
	LocationReadable LocationReadable
}

type Media struct {
	Image []string
	Video []string
}

type Variations struct {
	FixedPrice float64
	Media      Media
	Quantity   int
}

type PaymentData struct {
	PaymentName    string
	PaymentService string
}

type Payment struct {
	Online  []PaymentData
	Offline []PaymentData
}

type Shipping struct {
	ShippingService        string
	ShippingName           string
	ShippingCost           float64
	HandlingTimeMax        int
	DispatchTimeMin        int
	DispatchTimeMax        int
	ShippingAdditionalCost string
}
type Item struct {
	_version    string
	CategoryId  []string
	Title       string
	Media       Media
	SellerId    string
	Locale      Locale
	ListingType string
	Payment     Payment
	StartTime   string
	EndTime     string
	Shipping    []Shipping
	TitleSlug   string
	Variations  []Variations
	_fpaiStatus string
}

// The main function is used to bootstrap the code, however we don't have any functionality for it right now
// it only reports if an error occurs, which never should
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}



	itemInfoR := `{"locale":{"location":"51.51121389999999,-0.11982439999997041","countryCode":"GB","currencyId":"GBP","currencySymbol":"Â£","locationReadable":{"district":"City of Westminster","city":"London","state":"Greater London"}},"_version":"serving","categoryId":["Root","Cameras \u0026 Photo","Digital Cameras"],"title":"many pictures","media":{"image":["//lh5.ggpht.com/O_o_N6CFkClY5AV0-LqntpyFjor7Of4u23ZcK7lYwc2uY1ea7GWi61VDJZCB7UCb79svkjKPHIenqwEUhjHi0jdIQnnl6z_p03yktPUB1FBHezIQ","//lh6.ggpht.com/ih3q2d7CenGLPyupH9FpfsoJQWQpw1i8wWA2Kd26bFnSF2fbnKyGU9WePIhCgEeqw5p6YMVmFi1c9oS0Ag93aF_oZ3ZiwK7fQuSYIrZ9VhgXbrTHkw","//lh6.ggpht.com/7RJRsapsnwWL3_KiLIjMz4QojDzUvsztXtvKTFvIfde_AHccDnOibAvXRN73tTB4SeHzlj8S1LWxbYwwWFGn9elfCKdSb8BUIU5QJY1LO791HutQ","//lh6.ggpht.com/qAtjgyHAB734Ox_4NC_fa-ZRqrCjCmJu0Tp8bo-HMO88duv8l4hhuv2REBkB--yneFzOL7annecVlGty-YsKouondiOFVnAZWzjpdrfsGfbL6wh2","//lh3.ggpht.com/dWUbASepwHF4lHaXIPnpv4BNm2pCml9MlJt7s86s1cpu-PsYNmS0yQmKFKTM38q_oMLW_YJMJ19civ2gVViKAGYcZylRW7jN3w77AJvhzS6JE2g","//lh6.ggpht.com/9aXLmPRVeZnxkwvNb3mWTF8kvfEY_lho_lOVVc9AbNqLb8GQmiS_XXVZ3OKqMv2pxgYSayMYPPRh6ACYyh0H8KtS8mPD6MKUkEajwxkTtp5Q4Lo","//lh3.ggpht.com/FG_QXZPHJ2tTYwI_t5Fg1KqivglVg9RlJn0JRsu9Ox8vJ7IcBirb2IV_I1LL_WVOMxfTuBBSDLMlrw9v0MCAdmnPCR29sCbRGjhm6zEfIH-3q2QSdw","//lh4.ggpht.com/Y23DqORrVkM2m55f-rq5_BBrlkvQg4uX7AsAt-ixhMobjK_SFgFaDfktgLhkNsyKwSr9HcF8iiGY3Nw0xOKXG1sn6wyAWg_qsolmKjVOrM5V5mIR","//lh6.ggpht.com/mQ62Ly-DjMKPMzU1OcSPJ7SLBqym0uBjawlkTHfmb-HOKaD56dnitk1duwPFJVdbi0GUpd63RQvr2VMpHp6S1OQ3di-hq4-JPeRoS5FJzksXSvW_","//lh3.ggpht.com/dqWjWPcNsvlR1tMC_agizX19f9MDiNGWFYTYVn4kjJxzIIkEe0mLzNcvS62zVJxAOaitT-IgaUfZ-Ze23BgzbqYY-l600i_LbVe35Uinz6sXIyoB","//lh6.ggpht.com/xhSdFc9uHgghs_6gf3seUWYM-PG2oLmjTrpF7ptEEMqaIrQIa8VPfC6tXE7f3M13eZvDXYqMW_k0AHO5vwCEPNp-iObixskd_lBaKNfz3MH3SNQ","//lh5.ggpht.com/kYeoKPoZGJCow-G1FhnD8kzVjNjbQA8-Kyj8eAh0HL-fMZX9tTeFPQikTZdSU0kks4-5Ui54cZF2CjGut9vfMJAVDKIq3T-bAQewCxvfl2120tH5zQ","//lh5.ggpht.com/4qUl3d-G9EPBzcYKrimNsWhQw7CmONV0jgfVhxFgB9mEU_QLRCyNJTWs2A3xf6wc7AUF2DXrKEkoX-SNLMZ6s-O4aXXV9WOjOPcWdAYreMRBld0E","//lh5.ggpht.com/z-0C4G6EWYkelAF1LjPfl_UQcsp92H4joIPt8NfsOl0nPJ2VpzZYahWadKqTLfl6kq3C6aDBcwfGQyMWSozYoZIAOAW0yRvZrwxia321PlsKTxbZ","//lh4.ggpht.com/U7I12JrDYmMC_pUXpw8DVBjBilU67BvbM8qT8gJE0bQfkhHo7FOdMttiz3syP5IR-LyO4J1WBlfmZjvMjRr4GIBt4o3Vqp-hKz7q2_OGwGtsN5s","//lh3.ggpht.com/fF2XWEtqG23ybhzClhC_p8gvKJalf1vg7k3H7UkuAaIVubil7EgOvJUCwAZk2KiCtlPYp1E5Ep2xaxZjJRmg5EFSEAjqlMHJS_Wd1Bcje6xre4s","//lh3.ggpht.com/jgOebMihBoIZvHE4EOklJvZ_k-9egjNIlUKfKFcLkvXJs8g2FXjPvdFUbwqGrkHrMtyis8uOvgt-E51Vm11hq4bieh7h0cegca0VI4vFtFaAemU","//lh3.ggpht.com/MOrI-zKNMNrQE_aHj5hzbojP3T0hEMJKK6K8UO3e1NBC-nkcQeIM1QnvtJdT_G-W4e7-qv4BiqwdWcNHBpZXOmmX3tcuYEV8u_ANEoa9_aUIfeyg","//lh6.ggpht.com/SyIS5sGOkTG7k_jFF14wzH9Evrblv6o4pHBI6z6X070-xhAeyut_kRO6xHtDID4KLcWFvItjQy-plPcJ6K1T9tlFOrtaryEPvuAYdMVx8e0TTw","//lh6.ggpht.com/2Pp9kLYFhDT3USwHinU5OxnzcWWOLI0nOWe29gOD5KMzyEcXoHkTN-AutJV9M8F_9eqAP379XB9O1d0BWPanhr-MguzKxfHeUvYTs6yHzDkxyfe0NA","//lh4.ggpht.com/7aofqklSkF3AMDfF19yqsA9J3EfEiKy1NdOelEGKNnW0Cv5tGEpq2PF_jZO1MVoBbrrmVVRv0Tdq7I8KyZbIlyHdbTs1jMl7dEFqVMvsPcyaORyHlQ","//lh4.ggpht.com/anYJHqkMCkuhmIHQTBspLtWcDTyx1ZRe84_q5pAgVEOVmsKkaKhS725N4YFoj2zpJrBP7iTC2vf1GUtrp6H7kkm8c1k6zkW6I_Gf5f9A3re_I8Ex","//lh3.ggpht.com/OtSw0rU-DvfoXgoWrQdkln6Kz7O14TF9qrPNJSGJnZLeDqUEctOn1DT09pdwwVpNQV-cXmVYQL-PX4XPhpZLWH1ciSkVT6WHNmTz1D9pHphBwJUv","//lh3.ggpht.com/cTCZnXPIjI-EO2bvQdLgeoSLOSlMFcv805n347Zyci9XDYUdcVDC_5H7SFVYDr4pC5HtQDYnrOHL6AinLW7hWtfSCLlvVhVUNQ-DlDn0NwZ-1iCO-g","//lh4.ggpht.com/i-mL_JcF9rwjQq6HnuKzuAHU41_UGxQ62IOPZvaDrATXaPFbhe-EbT7ZIpboyNA5PXRCsxNsZ9hu58edRvNs5ScgKN8Lg-00J2LhlwMAbdEsv7b0nw","//lh6.ggpht.com/D_YV2BG1WWwl67xNloP3sxzRkqhcVTgJi58L-A8nLrOcMR_tBqLz4fHEGQ-qiNcG_-32MNy3dlSPWrTBKzBcweJxgMnRVet5yuGfelUlwehDtXX_3w"],"video":[]},"sellerId":"mihai","listingType":"fixedPrice","payment":{"online":[{"paymentName":"PayPal","paymentService":"paypal"}],"offline":[{"paymentName":"Pay on Pick-up","paymentService":"payOnPickup"}]},"startTime":"2014-01-04T10:02:18+00:00","endTime":"2014-04-04T10:02:18+00:00","shipping":[{"shippingService":"economy","shippingName":"Economy","shippingCost":1.0,"handlingTimeMax":4,"dispatchTimeMin":1,"dispatchTimeMax":10,"shippingAdditionalCost":"2"},{"shippingService":"localPickup","shippingName":"Local Pick-Up","shippingCost":0.0,"handlingTimeMax":2,"dispatchTimeMin":0,"dispatchTimeMax":0,"shippingAdditionalCost":"0"}],"titleSlug":"many-pictures","variations":[{"fixedPrice":222999.0,"media":{"image":["//lh6.ggpht.com/ih3q2d7CenGLPyupH9FpfsoJQWQpw1i8wWA2Kd26bFnSF2fbnKyGU9WePIhCgEeqw5p6YMVmFi1c9oS0Ag93aF_oZ3ZiwK7fQuSYIrZ9VhgXbrTHkw","//lh6.ggpht.com/9aXLmPRVeZnxkwvNb3mWTF8kvfEY_lho_lOVVc9AbNqLb8GQmiS_XXVZ3OKqMv2pxgYSayMYPPRh6ACYyh0H8KtS8mPD6MKUkEajwxkTtp5Q4Lo","//lh3.ggpht.com/FG_QXZPHJ2tTYwI_t5Fg1KqivglVg9RlJn0JRsu9Ox8vJ7IcBirb2IV_I1LL_WVOMxfTuBBSDLMlrw9v0MCAdmnPCR29sCbRGjhm6zEfIH-3q2QSdw"],"video":[]},"quantity":1121,"Brand":"Bell \u0026 Howell"},{"fixedPrice":211.0,"media":{"image":["//lh6.ggpht.com/qAtjgyHAB734Ox_4NC_fa-ZRqrCjCmJu0Tp8bo-HMO88duv8l4hhuv2REBkB--yneFzOL7annecVlGty-YsKouondiOFVnAZWzjpdrfsGfbL6wh2","//lh3.ggpht.com/FG_QXZPHJ2tTYwI_t5Fg1KqivglVg9RlJn0JRsu9Ox8vJ7IcBirb2IV_I1LL_WVOMxfTuBBSDLMlrw9v0MCAdmnPCR29sCbRGjhm6zEfIH-3q2QSdw","//lh6.ggpht.com/9aXLmPRVeZnxkwvNb3mWTF8kvfEY_lho_lOVVc9AbNqLb8GQmiS_XXVZ3OKqMv2pxgYSayMYPPRh6ACYyh0H8KtS8mPD6MKUkEajwxkTtp5Q4Lo","//lh3.ggpht.com/MOrI-zKNMNrQE_aHj5hzbojP3T0hEMJKK6K8UO3e1NBC-nkcQeIM1QnvtJdT_G-W4e7-qv4BiqwdWcNHBpZXOmmX3tcuYEV8u_ANEoa9_aUIfeyg"],"video":[]},"quantity":2,"Brand":"Fujifilm"},{"fixedPrice":22.0,"media":{"image":["//lh3.ggpht.com/jgOebMihBoIZvHE4EOklJvZ_k-9egjNIlUKfKFcLkvXJs8g2FXjPvdFUbwqGrkHrMtyis8uOvgt-E51Vm11hq4bieh7h0cegca0VI4vFtFaAemU","//lh3.ggpht.com/MOrI-zKNMNrQE_aHj5hzbojP3T0hEMJKK6K8UO3e1NBC-nkcQeIM1QnvtJdT_G-W4e7-qv4BiqwdWcNHBpZXOmmX3tcuYEV8u_ANEoa9_aUIfeyg","//lh4.ggpht.com/anYJHqkMCkuhmIHQTBspLtWcDTyx1ZRe84_q5pAgVEOVmsKkaKhS725N4YFoj2zpJrBP7iTC2vf1GUtrp6H7kkm8c1k6zkW6I_Gf5f9A3re_I8Ex"],"video":[]},"quantity":12,"Brand":"Gateway"}],"_fpaiStatus":"published"}`
	itemInfoBytes := []byte(itemInfoR)
	var ItemInfo Item
	er := json.Unmarshal(itemInfoBytes, &ItemInfo)
	if er != nil {
		panic(er)
	} else {
		fmt.Println(ItemInfo)
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
	} else if function == "supplierInfo" {
		return t.supplierInfo(stub, args)
	} else if function == "purchaseOrder" {
		return t.purchaseOrder(stub, args)
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
func (t *SimpleChaincode) supplierInfo(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}


// Purchase order code and "write" code are the exact same, because in essence, both should do the same job, which is to write
// data to the ledger which can be read later on 

func (t *SimpleChaincode) purchaseOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){

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
