/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"strconv"
	"strings"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ExchangePlaceChaincode example simple Chaincode implementation
type ExchangePlaceChaincode struct {

}

//exo 1 : flexible data model
type Wallet struct {
	CompanyID        string   `json:"companyID"`
	HoldingEUR       float64  `json:"holdingEUR"`
	HoldingDOL       float64  `json:"holdingDOL"`
}

func (t *ExchangePlaceChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}


func (t *ExchangePlaceChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	raiseError := func(errorMessage string) pb.Response {
		return shim.Error(errorMessage)
	}

	functionTypeName, args := stub.GetFunctionAndParameters()
	functionType := strings.Split(functionTypeName, ":")[0]
	if functionType == "invoke" {
		functionName := strings.Split(functionTypeName, ":")[1]
		switch functionName {
		case "createWallet":
			var Company string
			var Eur, Dol float64
			var WalletCompany Wallet
			var err error

			if len(args) != 3 {
				return shim.Error("Incorrect number of arguments. Expecting 3")
			}
			Company = args[0]
			Eur, err = strconv.ParseFloat(args[1], 64)
			if err != nil {
				return shim.Error("Expecting integer value for asset holding")
			}
			Dol, err = strconv.ParseFloat(args[2], 64)
			if err != nil {
				return shim.Error("Expecting integer value for asset holding")
			}
			fmt.Printf("Company = %s, Eur = %d, Dol = %d\n", Company, Eur, Dol)

			// Write the state to the ledger
			WalletCompany.CompanyID = Company

			WalletCompany.HoldingEUR = Eur
			WalletCompany.HoldingDOL = Dol
			walletAsBytes, err := json.Marshal(WalletCompany)
			if err != nil {
				return shim.Error(err.Error())
			}
			err = stub.PutState(Company, walletAsBytes)
			if err != nil {
				return shim.Error(err.Error())
			}
			bytes, err := json.Marshal(WalletCompany)
			if err == nil {
				return shim.Success(bytes)
			}
			return shim.Success(nil)
		default:
			return raiseError("Received unknown invoke function name")
		}
	} else if functionType == "query" {
		return t.query(stub, args)
	}
	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}
// Deletes an entity from state
func (t *ExchangePlaceChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

// query callback representing the query of a chaincode
func (t *ExchangePlaceChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

func main() {
	err := shim.Start(new(ExchangePlaceChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
