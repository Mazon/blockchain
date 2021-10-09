package main

import (
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
)

//hashGenesisBlock := "0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"

func main() {
	go func() {
		t := time.Now()
		genesisBlock := Block{}
		genesisBlock = Block{Version: 0, Timestamp: t.String(), Difficulty: 1, Nonce: 0}
		spew.Dump(genesisBlock)

		mutex.Lock()
		Blockchain = append(Blockchain, genesisBlock)
		mutex.Unlock()
	}()
	log.Fatal(run())
}
