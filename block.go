package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"
)

const difficulty = 1

type Block struct {
	//header
	Version    uint32
	Timestamp  string
	PrevHash   []byte
	Difficulty uint32
	Nonce      uint32 //The solution to the block.
	//body
	Transactions []Transaction
}

//Metadata about the chain.
type Metadata struct {
	B []string
}

var Blockchain []Block

var mutex = &sync.Mutex{}

func isBlockValid(newBlock, oldBlock Block) bool {
	//if oldBlock.Index+1 != newBlock.Index {
	//	return false
	//}

	//oldBlockHash := calculateHash(oldBlock)
	//if oldBlockHash != newBlock.PrevHash {
	//		return false
	//	}

	//if calculateHash(newBlock) != newBlock.Hash {
	//		return false
	//	}

	return true
}

func isHashValid(hash []byte, difficulty uint32) bool {
	prefix := strings.Repeat("0", int(difficulty))
	//fmt.Println(string(hash[:]))
	return strings.HasPrefix(string(hash[:]), prefix)
}

func generateBlock(oldBlock Block, tx []Transaction) Block {
	var newBlock Block

	t := time.Now()

	newBlock.Timestamp = t.String()
	newBlock.Transactions = tx

	// generate block hash of old block header
	oldBlockHash := calculateHash(oldBlock)
	newBlock.PrevHash = oldBlockHash

	newBlock.Difficulty = difficulty

	for i := 0; ; i++ {
		// increase nonce until hash is valid.
		newBlock.Nonce = uint32(i)
		if !isHashValid(calculateHash(newBlock), newBlock.Difficulty) {
			//fmt.Println(calculateHash(newBlock), " do more work!")
			h := calculateHash(newBlock)
			fmt.Println(hex.EncodeToString(h) + " do more work!")
			time.Sleep(time.Second)
			continue
		} else {
			h := calculateHash(newBlock)
			fmt.Println(hex.EncodeToString(h) + " work done!")
			break
		}

	}
	//fmt.Println(newBlock)
	return newBlock
}

//calculates the block header sha256 hash.
func calculateHash(block Block) []byte {
	bVersion := uinttobyte(block.Version)
	bNonce := uinttobyte(block.Nonce)
	bDifficulty := uinttobyte(block.Difficulty)

	record := []byte{}
	record = append(record, bVersion[:]...)
	record = append(record, block.PrevHash[:]...)
	record = append(record, bNonce[:]...)
	record = append(record, []byte(block.Timestamp)[:]...)
	record = append(record, bDifficulty[:]...)

	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	//fmt.Println(hex.EncodeToString(hashed))
	return hashed
}
