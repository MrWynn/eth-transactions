package main

import (
	"eth-transactions/client"
	"eth-transactions/def"
	"eth-transactions/memory"
	"log"
	"time"
)

// Parser interface
type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []def.Transaction
}

// TransactionParser implements the Parser interface
type TransactionParser struct {
	storage *memory.MemoryStorage
	client  *client.EthereumClient
}

// NewTransactionParser creates a new TransactionParser instance
func NewTransactionParser(client *client.EthereumClient, storage *memory.MemoryStorage) *TransactionParser {
	return &TransactionParser{
		storage: storage,
		client:  client,
	}
}

// GetCurrentBlock returns the last parsed block number
func (tp *TransactionParser) GetCurrentBlock() int {
	return tp.storage.GetCurrentBlock()
}

// Subscribe adds an address to the observer list
func (tp *TransactionParser) Subscribe(address string) bool {
	return tp.storage.Subscribe(address)
}

// GetTransactions returns the list of inbound or outbound transactions for an address
func (tp *TransactionParser) GetTransactions(address string) []def.Transaction {
	return tp.storage.GetTransactions(address)
}

// ParseNewBlocks fetches and parses new blocks, storing relevant transactions
func (tp *TransactionParser) ParseNewBlocks() {
	currentBlock, err := tp.client.GetBlockNumber()
	if err != nil {
		log.Println("Error fetching block number:", err)
		return
	}

	lastParsedBlock := tp.storage.GetCurrentBlock()
	go func() {
		for {
			if lastParsedBlock < currentBlock {
				log.Printf("Start parser transactions from %d block to %d block\n", lastParsedBlock+1, currentBlock)
			}
			for blockNumber := lastParsedBlock + 1; blockNumber <= currentBlock; blockNumber++ {
				transactions, err := tp.client.GetBlockTransactions(blockNumber)
				if err != nil {
					log.Println("Error fetching block transactions:", err)
					continue
				}

				for _, tx := range transactions {
					if tp.storage.IsSubscribe(tx.From) || tp.storage.IsSubscribe(tx.To) {
						tp.storage.AddTransaction(tx)
					}
				}
				tp.storage.UpdateCurrentBlock(blockNumber)
			}

			time.Sleep(10 * time.Second)

			lastParsedBlock = tp.storage.GetCurrentBlock()
			currentBlock, err = tp.client.GetBlockNumber()
			if err != nil {
				log.Println("Error fetching block number:", err)
			}
		}
	}()

}
