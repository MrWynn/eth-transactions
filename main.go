package main

import (
	"encoding/json"
	"eth-transactions/client"
	"eth-transactions/memory"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Initialize Ethereum client
	client := client.NewEthereumClient()

	// Initialize memory storage
	// Set a block number to avoid scanning from block 0
	storage := memory.NewMemoryStorage(19916787)

	// Initialize transaction parser
	parser := NewTransactionParser(client, storage)

	// Example usage
	parser.Subscribe("0x0AFfB0a96FBefAa97dCe488DfD97512346cf3Ab8")

	// Parse new blocks and update storage
	parser.ParseNewBlocks()

	// Expose public interface via REST API
	http.HandleFunc("/currentBlock", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Current Block: %d", parser.GetCurrentBlock())
	})

	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		if parser.Subscribe(address) {
			fmt.Fprintf(w, "Subscribed to address: %s", address)
		} else {
			fmt.Fprintf(w, "Already subscribed to address: %s", address)
		}
	})

	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		transactions := parser.GetTransactions(address)
		json.NewEncoder(w).Encode(transactions)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
