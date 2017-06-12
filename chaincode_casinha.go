/*
Copyright IBM Corp 2016 All Rights Reserved.

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
	//"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Pagamento struct {
	Pagador     string `json:"pagador"`
	Recebedor   string `json:"recebedor"`
	DataEntrada string `json:"dataEntrada"`
	DataSaida   string `json:"dataSaida"`
	Valor       string `json:"valor"`
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("init_casinha", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "pagar" {
		return t.pagar(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		// fmt.Println("Retornando pagamento")
		// pagamento, err := GetPagamento(args[0], stub)
		// if err != nil {
		// 	fmt.Println("Error from GetPagamento")
		// 	return nil, err
		// } else {
		// 	pagamentoBytes, err1 := json.Marshal(&pagamento)
		// 	if err1 != nil {
		// 		fmt.Println("Erro convertendo o pagamento")
		// 		return nil, err1
		// 	}
		// 	fmt.Println("Sucesso! Retornando pagamento.")
		// 	return pagamentoBytes, nil
		// }
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) pagar(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var pagador, recebedor, dataEntrada, dataSaida, valor string
	var err error
	//var pagamento Pagamento

	fmt.Println("running pagar()")

	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5.")
	}

	pagador = args[0]
	recebedor = args[1]
	dataEntrada = args[3]
	dataSaida = args[4]
	valor = args[5]

	//pagamento = Pagamento{ Pagador: pagador, Recebedor: recebedor, DataEntrada: dataEntrada, DataSaida: dataSaida, Valor: valor }

	str := `"`+ pagador + `,` + recebedor + `,` + dataEntrada + `,` + dataSaida + `,` + valor + `"`

	fmt.Println("Pagar - Json: " + str)
	//return nil, errors.New("STR: " + str)

	//pagamentoBytes, err := json.Marshal(&pagamento)
	// if err != nil {
	// 	fmt.Println("Erro ao criar objeto pagamento ")
	// 	return nil, errors.New("Erro ao criar objeto pagamento")
	// }

	err = stub.PutState(recebedor, []byte(str)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}

	fmt.Println("Pagamento criado")
	return nil, nil
}

func (t *SimpleChaincode) testePagar(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	err = stub.PutState(args[0], []byte(args[1])) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}

	fmt.Println("Pagamento criado")
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

// func GetPagamento(recebedor string, stub shim.ChaincodeStubInterface) (Pagamento, error) {
// 	var pagamento Pagamento
// 	pagamentoBytes, err := stub.GetState(recebedor)
// 	if err != nil {
// 		fmt.Println("Pagamento não encontrado: " + recebedor)
// 		return pagamento, errors.New("Pagamento não encontrado: " + recebedor)
// 	}

// 	err = json.Unmarshal(pagamentoBytes, &pagamento)
// 	if err != nil {
// 		fmt.Println("Erro ao converter o pagamento " + recebedor + "\n err:" + err.Error())
// 		return pagamento, errors.New("Erro ao converter o pagamento " + recebedor + "\n err:" + err.Error())
// 	}

// 	return pagamento, nil
// }
