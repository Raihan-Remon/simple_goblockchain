package main

import (
	"fmt"
	"log"
	"time"
)

func init() {
	log.SetPrefix("Blockchain")
}

type Block struct {
	nonce        int
	previousHash string
	timestamp    int64
	transactions []string
}

func NewBlock(nonce int, previousHash string) *Block {
	b := new(Block)
	b.nonce = nonce
	b.timestamp = time.Now().UnixNano()
	b.previousHash = previousHash
	return b
}

func (b *Block) Print() {
	fmt.Printf("timestamp:\t%d\n", b.timestamp)
	fmt.Printf("nonce:\t\t%d\n", b.nonce)
	fmt.Printf("previous hash:\t%s\n", b.previousHash)
	fmt.Printf("transactions:\t%s\n", b.transactions)

}

func main() {
	b := NewBlock(0, "init hash")
	b.Print()
}
