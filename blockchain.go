package main

import (
	"blockchain/block"
	"blockchain/web"

	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"sync"
)

//hashGenesisBlock := "0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"

func main() {
	mutex := &sync.Mutex{}

	go func() {
		t := time.Now()
		genesisBlock := block.Block{}
		genesisBlock = block.Block{Version: 0, Timestamp: t.String(), Difficulty: 1, Nonce: 0}
		spew.Dump(genesisBlock)

		mutex.Lock()
		block.Blockchain = append(block.Blockchain, genesisBlock)
		mutex.Unlock()
	}()
	log.Fatal(web.Run())
}
