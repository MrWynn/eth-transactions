package memory

import (
	"eth-transactions/def"
	"github.com/ethereum/go-ethereum/common"
	"sync"
)

// MemoryStorage to store data in memory
type MemoryStorage struct {
	currentBlock  int
	subscriptions map[common.Address]bool
	transactions  map[common.Address][]def.Transaction
	mu            sync.RWMutex
}

// NewMemoryStorage creates a new MemoryStorage instance
func NewMemoryStorage(blockNumber int) *MemoryStorage {
	return &MemoryStorage{
		currentBlock:  blockNumber,
		subscriptions: make(map[common.Address]bool),
		transactions:  make(map[common.Address][]def.Transaction),
	}
}

// GetCurrentBlock returns the last parsed block number
func (ms *MemoryStorage) GetCurrentBlock() int {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return ms.currentBlock
}

// Subscribe adds an address to the observer list
func (ms *MemoryStorage) Subscribe(address string) bool {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	addr := common.HexToAddress(address)
	if _, exists := ms.subscriptions[addr]; exists {
		return false
	}
	ms.subscriptions[addr] = true
	return true
}

func (ms *MemoryStorage) IsSubscribe(address string) bool {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	addr := common.HexToAddress(address)
	if _, exists := ms.subscriptions[addr]; exists {
		return true
	}
	return false
}

// GetTransactions returns the list of inbound or outbound transactions for an address
func (ms *MemoryStorage) GetTransactions(address string) []def.Transaction {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return ms.transactions[common.HexToAddress(address)]
}

// UpdateCurrentBlock updates the last parsed block number
func (ms *MemoryStorage) UpdateCurrentBlock(blockNumber int) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.currentBlock = blockNumber
}

// AddTransaction adds a transaction to the storage
func (ms *MemoryStorage) AddTransaction(tx def.Transaction) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	from := common.HexToAddress(tx.From)
	to := common.HexToAddress(tx.To)
	ms.transactions[from] = append(ms.transactions[from], tx)
	ms.transactions[to] = append(ms.transactions[to], tx)
}
