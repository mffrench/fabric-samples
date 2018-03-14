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
	"time"
	"strconv"
	"strings"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

/*

//exo 1 : flexible data model
//
//type Wallet struct {
//	CompanyID        string   `json:"companyID"`
//	HoldingEUR       float64  `json:"holdingEUR"`
//	HoldingDOL       float64  `json:"holdingDOL"`
//}
//
//type Transaction struct {
//	TransactionID    string    `json:"transactionID`
//	CompanyIDSrc     string    `json:"companyIDSrc"`
//	CompanyIDDst     string    `json:"companyIDDst"`
//	HoldingType      string    `json:"holdingType"`
//	HoldingValue     float64   `json:"holdingValue"`
//}

-> createWallet(CompanyName)

exo 2 : several invoke definitions are possible

type HoldingType struct {
	HoldingName     string     `json:"holdingName"`
	HoldingCode     string     `json:"holdingCode"`
}

type Wallet struct {
	CompanyID        string              `json:"companyID"`
	Holdings         map[string]float64  `json:"holdings"`
}

type Exchange struct {
	ExchangeID          string    `json:"exchangeID"`
	FirstTransactionID  string    `json:"firstTransactionID"`
	SecondTransactionID string    `json:"secondTransactionID"`
}

-> transferFromTo(CompanyIDFrom, CompanyIDTo, Type, Value)
-> exchangeBetween(CompanyIDA, CompanyIDB, TypeAB, ValueAB, TypeBA, ValueBA)

exo 3 : parcours blocks & transac + several query

.... COUCHDB + MAP/REDUCE GIVE US MONEY ....

 */

// ExchangePlaceChaincode example simple Chaincode implementation
type ExchangePlaceChaincode struct {

}

//exo 1 : flexible data model

type Wallet struct {
	CompanyID        string   `json:"companyID"`
	HoldingEUR       float64  `json:"holdingEUR"`
	HoldingDOL       float64  `json:"holdingDOL"`
}

type Transaction struct {
	TransactionID    string    `json:"transactionID"`
	CompanyIDSrc     string    `json:"companyIDSrc"`
	CompanyIDDst     string    `json:"companyIDDst"`
	HoldingType      string    `json:"holdingType"`
	HoldingValue     float64   `json:"holdingValue"`
}

type Exchange struct {
	ExchangeID          string    `json:"exchangeID"`
	FirstTransactionID  string    `json:"firstTransactionID"`
	SecondTransactionID string    `json:"secondTransactionID"`
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

		case "transaction":
			var walletCompanySrc Wallet
			var walletCompanyDst Wallet
			var transaction Transaction

			transaction.CompanyIDSrc = args[0]
			transaction.CompanyIDDst = args[1]
			transaction.HoldingType = args[3]
			transaction.TransactionID = "TR" + strconv.FormatInt(time.Now().Unix(),  10)


			X, err := strconv.ParseFloat(args[2], 64)
			transaction.HoldingValue = X

			WalletCompanySrcbytes, err := stub.GetState(transaction.CompanyIDSrc)
			if err != nil {
				return shim.Error("Failed to get state")
			}
			if WalletCompanySrcbytes == nil {
				return shim.Error("Entity not found")
			}
			err = json.Unmarshal(WalletCompanySrcbytes, &walletCompanySrc)
			if err != nil {
				return shim.Error("Error while unmarshalling")
			}
			WalletCompanyDstbytes, err := stub.GetState(transaction.CompanyIDDst)
			if err != nil {
				return shim.Error("Failed to get state")
			}
			if WalletCompanyDstbytes == nil {
				return shim.Error("Entity not found")
			}
			err = json.Unmarshal(WalletCompanyDstbytes, &walletCompanyDst)
			if err != nil {
				return shim.Error("Error while unmarshallingg")
			}
			//test := float64(float64(walletCompanySrc.HoldingEUR) - float64(walletCompanyDst.HoldingEUR))
			if transaction.HoldingType == "EUR" {
				walletCompanySrc.HoldingEUR = float64(walletCompanySrc.HoldingEUR) - float64(transaction.HoldingValue)
				walletCompanyDst.HoldingEUR = float64(walletCompanyDst.HoldingEUR) + float64(transaction.HoldingValue)
				fmt.Printf("Src = %d EUR, Dst = %d EUR\n", walletCompanySrc.HoldingEUR, walletCompanyDst.HoldingEUR)
			} else if transaction.HoldingType == "DOL" {
				walletCompanySrc.HoldingDOL = float64(walletCompanySrc.HoldingDOL) - float64(transaction.HoldingValue)
				walletCompanyDst.HoldingDOL = float64(walletCompanyDst.HoldingDOL) + float64(transaction.HoldingValue)
				fmt.Printf("Src = %d DOL, Dst = %d DOL\n", walletCompanySrc.HoldingDOL, walletCompanyDst.HoldingDOL)
			} else {
				return shim.Error("Error wrong currency")
			}

			walletCompanySrcAsBytes, err := json.Marshal(walletCompanySrc)
			if err != nil {
				return shim.Error(err.Error())
			}
			err = stub.PutState(walletCompanySrc.CompanyID, walletCompanySrcAsBytes)
			if err != nil {
				return shim.Error(err.Error())
			}
			walletCompanyDstAsBytes, err := json.Marshal(walletCompanyDst)
			if err != nil {
				return shim.Error(err.Error())
			}
			err = stub.PutState(walletCompanyDst.CompanyID, walletCompanyDstAsBytes)
			if err != nil {
				return shim.Error(err.Error())
			}
			transactionAsBytes, err := json.Marshal(transaction)
			if err != nil {
				return shim.Error(err.Error())
			}
			err = stub.PutState(transaction.TransactionID, transactionAsBytes)
			if err != nil {
				return shim.Error(err.Error())
			}
			return shim.Success(nil)
			//-> exchangeBetween(CompanyIDA, CompanyIDB, ValueAB,TypeAB, ValueBA, TypeBA)
			//-> exchangeBetween(CompanyA,	 CompanyB, 		100,	EUR, 	70, 	DOL)
		case "exchange":
			var walletCompanyA Wallet
			var walletCompanyB Wallet
			var transactionAB Transaction
			transactionAB.CompanyIDSrc = args[0]
			transactionAB.CompanyIDDst = args[1]
			transactionAB.HoldingType = args[3]
			transactionAB.TransactionID = "TR_AB" + strconv.FormatInt(time.Now().Unix(),  10)

			XA, err := strconv.ParseFloat(args[2], 64)
			if err != nil {
				return shim.Error("Invalid transaction amount")
			}
			transactionAB.HoldingValue = XA

			var transactionBA Transaction
			transactionBA.CompanyIDSrc = args[1]
			transactionBA.CompanyIDDst = args[0]
			transactionBA.HoldingType = args[5]
			transactionBA.TransactionID = "TR_BA" + strconv.FormatInt(time.Now().Unix(),  10)

			XB, err := strconv.ParseFloat(args[4], 64)
			if err != nil {
				return shim.Error("Invalid transaction amount")
			}
			transactionBA.HoldingValue = XB

			walletCompanyAbytes, err := stub.GetState(transactionAB.CompanyIDSrc)
			if err != nil {
				return shim.Error("Failed to get state")
			}
			if walletCompanyAbytes == nil {
				return shim.Error("Entity not found")
			}
			err = json.Unmarshal(walletCompanyAbytes, &walletCompanyA)
			if err != nil {
				return shim.Error("Error while unmarshalling")
			}

			walletCompanyBbytes, err := stub.GetState(transactionBA.CompanyIDSrc)
			if err != nil {
				return shim.Error("Failed to get state")
			}
			if walletCompanyBbytes == nil {
				return shim.Error("Entity not found")
			}
			err = json.Unmarshal(walletCompanyBbytes, &walletCompanyB)
			if err != nil {
				return shim.Error("Error while unmarshalling")
			}
			if transactionAB.HoldingType == "EUR" {
				walletCompanyA.HoldingEUR = float64(walletCompanyA.HoldingEUR) - float64(transactionAB.HoldingValue)
				walletCompanyB.HoldingEUR = float64(walletCompanyB.HoldingEUR) + float64(transactionAB.HoldingValue)
			} else if transactionAB.HoldingType == "DOL" {
				walletCompanyA.HoldingDOL = float64(walletCompanyA.HoldingDOL) - float64(transactionAB.HoldingValue)
				walletCompanyB.HoldingDOL = float64(walletCompanyB.HoldingDOL) + float64(transactionAB.HoldingValue)
			} else {
				return shim.Error("Error wrong currency")
			}

			if transactionBA.HoldingType == "EUR" {
				walletCompanyB.HoldingEUR = float64(walletCompanyB.HoldingEUR) - float64(transactionBA.HoldingValue)
				walletCompanyA.HoldingEUR = float64(walletCompanyA.HoldingEUR) + float64(transactionBA.HoldingValue)
			} else if transactionBA.HoldingType == "DOL" {
				walletCompanyB.HoldingDOL = float64(walletCompanyB.HoldingDOL) - float64(transactionBA.HoldingValue)
				walletCompanyA.HoldingDOL = float64(walletCompanyA.HoldingDOL) + float64(transactionBA.HoldingValue)
			}

			walletCompanyAAsBytes, err := json.Marshal(walletCompanyA)
			if err != nil {
				return shim.Error(err.Error())
			}
			err = stub.PutState(walletCompanyA.CompanyID, walletCompanyAAsBytes)
			if err != nil {
				return shim.Error(err.Error())
			}
			walletCompanyBAsBytes, err := json.Marshal(walletCompanyB)
			if err != nil {
				return shim.Error(err.Error())
			}
			err = stub.PutState(walletCompanyB.CompanyID, walletCompanyBAsBytes)
			if err != nil {
				return shim.Error(err.Error())
			}
			transactionABAsBytes, err := json.Marshal(transactionAB)
			if err != nil {
				return shim.Error(err.Error())
			}
			err = stub.PutState(transactionAB.TransactionID, transactionABAsBytes)
			if err != nil {
				return shim.Error(err.Error())
			}
			transactionBAAsBytes, err := json.Marshal(transactionBA)
			if err != nil {
				return shim.Error(err.Error())
			}
			err = stub.PutState(transactionBA.TransactionID, transactionBAAsBytes)
			if err != nil {
				return shim.Error(err.Error())
			}
			return shim.Success(nil)

			var myExchange Exchange
			myExchange.FirstTransactionID = transactionAB.TransactionID
			myExchange.SecondTransactionID = transactionBA.TransactionID
			myExchange.ExchangeID = "EX" + strconv.FormatInt(time.Now().Unix(),  10)
			myExchangeAsBytes, err := json.Marshal(myExchange)
			if err != nil {
				return shim.Error(err.Error())
			}
			err = stub.PutState(myExchange.ExchangeID, myExchangeAsBytes)
			if err != nil {
				return shim.Error(err.Error())
			}

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
