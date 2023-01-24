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

const MINING_DIFICULTY = 3

type Block struct {
	timestamp    int64
	nonce        int
	previousHash [32]byte
	transactions []*Transaction
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	b := new(Block)
	b.nonce = nonce
	b.timestamp = time.Now().UnixNano()
	b.previousHash = previousHash
	b.transactions = transactions
	return b
}

func (b *Block) Print() {
	fmt.Printf("timestamp:\t%d\n", b.timestamp)
	fmt.Printf("nonce:\t\t%d\n", b.nonce)
	fmt.Printf("previous hash:\t%x\n", b.previousHash)
	for _, t := range b.transactions {
		t.Print()
	}

}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256(m)
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			Timestamp    int64          `json:"timestamp"`
			Nonce        int            `json:"nonce"`
			PreviousHash [32]byte       `json:"previous_hash"`
			Transactions []*Transaction `json:"transactions"`
		}{
			Timestamp:    b.timestamp,
			Nonce:        b.nonce,
			PreviousHash: b.previousHash,
			Transactions: b.transactions,
		})
}

// Sample blockchain

type Blockchain struct {
	transactionPool []*Transaction
	chain           []*Block
}

func NewBlockchain() *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

func (bc *Blockchain) AddTransaction(sender string, recepient string, value float32) {
	t := NewTansaction(sender, recepient, value)
	bc.transactionPool = append(bc.transactionPool, t)
}

func (bc *Blockchain) CopyTransaction() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions, NewTansaction(t.senderBlockChainAddress, t.recepientBlockChainAddress, t.value))
	}
	return transactions
}

func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{0, nonce, previousHash, transactions}
	guessHashStr := fmt.Sprintf("%x/n", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeros
}

func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransaction()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFICULTY) {
		nonce += 1
	}
	return nonce
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) Print() {
	fmt.Println()
	for i, block := range bc.chain {
		fmt.Printf("%s Chain: %d %s\n\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n\n\n", strings.Repeat("*", 50))
}

type Transaction struct {
	senderBlockChainAddress    string
	recepientBlockChainAddress string
	value                      float32
}

func NewTansaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{sender, recipient, value}
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf("sender_blockchain_address:\t%s\n", t.senderBlockChainAddress)
	fmt.Printf("recipient_blockchain_address:\t%s\n", t.recepientBlockChainAddress)
	fmt.Printf("value:\t\t\t\t%.1f\n", t.value)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recepient string  `json:"recepient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.senderBlockChainAddress,
		Recepient: t.recepientBlockChainAddress,
		Value:     t.value,
	})
}

func main() {
	blockchain := NewBlockchain()
	blockchain.Print()

	blockchain.AddTransaction("A", "B", 1.0)
	previousHash := blockchain.LastBlock().Hash()
	nonce := blockchain.ProofOfWork()
	blockchain.CreateBlock(nonce, previousHash)
	blockchain.Print()

	blockchain.AddTransaction("C", "D", 2.0)
	blockchain.AddTransaction("X", "Y", 0.5)
	previousHash = blockchain.LastBlock().Hash()
	nonce = blockchain.ProofOfWork()
	blockchain.CreateBlock(nonce, previousHash)
	blockchain.Print()
}
