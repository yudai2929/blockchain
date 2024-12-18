package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// 直前のブロックのハッシュ値
type Hash [32]byte

type Transactions []string

type Block struct {
	previousHash Hash
	transactions Transactions
	timestamp    int64
}

type Blocks []*Block

func (b Blocks) Last() *Block {
	return b[len(b)-1]
}

func NewBlock(previousHash Hash, transaction Transactions) *Block {
	return &Block{
		previousHash: previousHash,
		transactions: transaction,
		timestamp:    time.Now().UnixNano(),
	}
}

func (b *Block) Hash() (Hash, error) {
	m, err := json.Marshal(b)
	if err != nil {
		return Hash{}, fmt.Errorf("failed to marshal block: %w", err)
	}
	return sha256.Sum256(m), nil
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64    `json:"timestamp"`
		PreviousHash [32]byte `json:"previous_hash"`
		Transaction  []string `json:"transaction"`
	}{
		Timestamp:    b.timestamp,
		PreviousHash: b.previousHash,
		Transaction:  b.transactions,
	})
}

func (b *Block) Print() {
	fmt.Printf("Timestamp       %d\n", b.timestamp)
	fmt.Printf("PreviousHash    %x\n", b.previousHash)
	fmt.Printf("Transaction     %s\n", b.transactions)
}

type Blockchain struct {
	transactionPool Transactions // トランザクションを一時的にプールするフィールド
	blocks          Blocks
}

func New() *Blockchain {
	b := NewBlock([32]byte{}, []string{})
	return &Blockchain{
		transactionPool: []string{},
		blocks:          Blocks{b},
	}
}

func (bc *Blockchain) CreateBlock() error {
	ph, err := bc.blocks.Last().Hash()
	if err != nil {
		return fmt.Errorf("failed to get hash: %w", err)
	}
	b := NewBlock(ph, bc.transactionPool)
	bc.blocks = append(bc.blocks, b)
	bc.transactionPool = nil
	return nil
}

func (bc *Blockchain) AddTransaction(transaction string) {
	bc.transactionPool = append(bc.transactionPool, transaction)
}

func (bc *Blockchain) Print() {
	for i, b := range bc.blocks {
		fmt.Printf("%s Block %d %s\n", strings.Repeat("=", 15), i, strings.Repeat("=", 15))
		b.Print()
	}
}
