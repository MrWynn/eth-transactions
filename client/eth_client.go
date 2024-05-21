package client

import (
	"bytes"
	"encoding/json"
	"eth-transactions/def"
	"fmt"
	"math/big"
	"net/http"
)

var httpURL = "https://mainnet.infura.io/v3/3a02ea0cb0bb4d39beda168dcc993b04"

// EthereumClient represents an Ethereum JSON-RPC client
type EthereumClient struct {
	url string
}

// NewEthereumClient creates a new EthereumClient instance
func NewEthereumClient() *EthereumClient {
	return &EthereumClient{url: httpURL}
}

// JSONRPCRequest represents a JSON-RPC request
type JSONRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// JSONRPCResponse represents a JSON-RPC response
type JSONRPCResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	Result  json.RawMessage  `json:"result"`
	Error   *json.RawMessage `json:"error"`
	ID      int              `json:"id"`
}

// Call makes a JSON-RPC call to the Ethereum node
func (ec *EthereumClient) Call(method string, params []interface{}, result interface{}) error {
	reqBody, err := json.Marshal(JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(ec.url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var rpcResp JSONRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return err
	}

	if rpcResp.Error != nil {
		return fmt.Errorf("JSON-RPC error: %s", string(*rpcResp.Error))
	}

	return json.Unmarshal(rpcResp.Result, result)
}

// GetBlockNumber returns the latest block number
func (ec *EthereumClient) GetBlockNumber() (int, error) {
	var result string
	if err := ec.Call("eth_blockNumber", []interface{}{}, &result); err != nil {
		return 0, err
	}

	var blockNumber int
	if _, err := fmt.Sscanf(result, "0x%x", &blockNumber); err != nil {
		return 0, err
	}

	return blockNumber, nil
}

// GetBlockTransactions returns the transactions of a given block
func (ec *EthereumClient) GetBlockTransactions(blockNumber int) ([]def.Transaction, error) {
	var block struct {
		Transactions []struct {
			From  string `json:"from"`
			To    string `json:"to"`
			Value string `json:"value"`
			Hash  string `json:"hash"`
		} `json:"transactions"`
	}
	param := fmt.Sprintf("0x%x", blockNumber)
	if err := ec.Call("eth_getBlockByNumber", []interface{}{param, true}, &block); err != nil {
		return nil, err
	}

	var transactions []def.Transaction
	for _, tx := range block.Transactions {
		val, ok := new(big.Int).SetString(tx.Value, 0)
		if !ok {
			return nil, fmt.Errorf("invalid value: %s", tx.Value)
		}
		transactions = append(transactions, def.Transaction{
			From:        tx.From,
			To:          tx.To,
			Amount:      val.String(),
			Hash:        tx.Hash,
			BlockNumber: blockNumber,
		})
	}

	return transactions, nil
}
