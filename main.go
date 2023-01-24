package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

func init() {
	log.SetPrefix("Blockchain")
}

type Block struct {
	nonce        int
	previousHash [32]byte
	timestamp    int64
	transactions []string
}

func NewBlock(nonce int, previousHash [32]byte) *Block {
	b := new(Block)
	b.nonce = nonce
	b.timestamp = time.Now().UnixNano()
	b.previousHash = previousHash
	return b
}

func (b *Block) Print() {
	fmt.Printf("timestamp:\t%d\n", b.timestamp)
	fmt.Printf("nonce:\t\t%d\n", b.nonce)
	fmt.Printf("previous hash:\t%x\n", b.previousHash)
	fmt.Printf("transactions:\t%s\n", b.transactions)

}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256(m)
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			Timestamp    int64    `json:"timestamp"`
			Nonce        int      `json:"nonce"`
			PreviousHash [32]byte `json:"previous_hash"`
			Transactions []string `json:"transactions"`
		}{
			Timestamp:    b.timestamp,
			Nonce:        b.nonce,
			PreviousHash: b.previousHash,
			Transactions: b.transactions,
		})
}

// Sample blockchain

type Blockchain struct {
	transactionPool []string
	chain           []*Block
}

func NewBlockchain() *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash)
	bc.chain = append(bc.chain, b)
	return b
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain: %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

func main() {
	blockchain := NewBlockchain()
	blockchain.Print()

	previousHash := blockchain.LastBlock().Hash()
	blockchain.CreateBlock(5, previousHash)
	blockchain.Print()

	previousHash = blockchain.LastBlock().Hash()
	blockchain.CreateBlock(1, previousHash)
	blockchain.Print()
}
