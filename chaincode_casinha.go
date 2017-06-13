/*
	Autor: Rodrigo Sclosa 
	E-mail: rodrigogs@ciandt.com
	Data: 13/06/2017

	Chaincode paga recebimento de pagamentos da casinha via rede Blockchain.
*/

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

type Pagamento struct {
	Pagador     string `json:"pagador"`
	Recebedor   string `json:"recebedor"`
	DataEntrada string `json:"dataEntrada"`
	DataSaida   string `json:"dataSaida"`
	Valor       string `json:"valor"`
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var blank []string
	blankBytes, _ := json.Marshal(&blank)
	err := stub.PutState("init_casinha", blankBytes)

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
	fmt.Println("Função não reconhecida: " + function)

	return nil, errors.New("Função não reconhecida: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		fmt.Println("Retornando pagamento")
		pagamento, err := GetPagamento(args[0], stub)
		if err != nil {
			fmt.Println("Erro de GetPagamento")
			return nil, err
		} else {
			pagamentoBytes, err1 := json.Marshal(&pagamento)
			if err1 != nil {
				fmt.Println("Erro convertendo o pagamento")
				return nil, err1
			}
			fmt.Println("Sucesso! Retornando pagamento.")
			return pagamentoBytes, nil
		}
		//eturn t.read(stub, args)
	}
	fmt.Println("Função não reconhecida: " + function)

	return nil, errors.New("Função não reconhecida: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) pagar(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//var pagador, recebedor, dataEntrada, dataSaida, valor string
	var err error
	var recebedor string

	fmt.Println("running pagar()")

	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5.")
	}

	recebedor = args[1]

	var pagamento Pagamento
	pagamento.Pagador = args[0]
	pagamento.Recebedor = recebedor
	pagamento.DataEntrada = args[2]
	pagamento.DataSaida = args[3]
	pagamento.Valor = args[4]

	// pagador = args[0]
	// recebedor = args[1]
	// dataEntrada = args[2]
	// dataSaida = args[3]
	// valor = args[4]
	
	// pagamento := Pagamento{ Pagador: pagador, Recebedor: recebedor, DataEntrada: dataEntrada, DataSaida: dataSaida, Valor: valor }
	// fmt.Println(pagamento)

	pagamentoBytes, err := json.Marshal(&pagamento)
	if err != nil {
		fmt.Println("Erro ao criar objeto pagamento ")
		return nil, errors.New("Erro ao criar objeto pagamento")
	}

	err = stub.PutState(recebedor, []byte(pagamentoBytes)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}

	fmt.Println("Pagamento criado")
	return pagamentoBytes, nil
}

func GetPagamento(recebedor string, stub shim.ChaincodeStubInterface) (Pagamento, error) {
	var pagamento Pagamento
	pagamento = Pagamento{}
	pagamentoBytes, err := stub.GetState(recebedor)
	if err != nil {
		fmt.Println("Pagamento não encontrado: " + recebedor)
		return pagamento, errors.New("Pagamento não encontrado: " + recebedor)
	}

	err = json.Unmarshal(pagamentoBytes, &pagamento)
	if err != nil {
		fmt.Println("Erro ao converter o pagamento " + recebedor + "\n err:" + err.Error())
		return pagamento, errors.New("Erro ao converter o pagamento " + recebedor + "\n err:" + err.Error())
	}

	return pagamento, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}