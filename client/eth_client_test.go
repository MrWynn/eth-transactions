package client

import (
	"fmt"
	"testing"
)

func TestGetBlockNumber(t *testing.T) {
	c := NewEthereumClient()
	t.Log(c.GetBlockNumber())
}

func TestBlockTransactions(t *testing.T) {
	c := NewEthereumClient()
	blockNumber, err := c.GetBlockNumber()
	if err != nil {
		t.Fatal(err)
	}
	tx, err := c.GetBlockTransactions(blockNumber)
	if err != nil {
		t.Fatal(err)
	}
	length := len(tx)
	fmt.Printf("%d block has %d transactions\n", blockNumber, length)
	for _, v := range tx {
		fmt.Printf("%+v\n", v)
	}
}
