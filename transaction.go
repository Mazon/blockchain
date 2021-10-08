package main

//
// A pool of unconfirmed transactions waiting to be included in a block.
//
type MemPool struct {
	Transactions []Transaction
}

//
// An input of a transaction.  It contains the location of the previous
// transaction's output that it claims and a signature that matches the
// output's public key.
//
type CTXIn struct {
	CPrevTx   string
	COutIndex int
	ScriptSig Script
}

//
// An output of a transaction.  It contains the public key that the next input
// must be able to sign with to claim it.
//
type CTXOut struct {
	Value        int64 //the number of coins.
	ScriptPubKey Script
}

//
// A list of instructions recorded with each transaction that describe how the next person wanting to spend can gain access.
//
type Script struct {
	Signature    string
	ScriptPubKey string
}

//
// The basic transaction that is broadcasted on the network and contained in
// blocks.  A transaction can contain multiple inputs and outputs.
//
type Transaction struct {
	input  []CTXIn
	output []CTXOut
}

// Verifies incoming transactions and adds them to memPool.
func addTransactionToPool(t Transaction) {
	// data structure is ok.
	// input and output have values.
	// the transaction is less then 1MB
	// the values must be more then 0 and less then max coins.
	//  None of the inputs have hash = 0
	// locktime?
	// transaction size > 100byte
	// sig larger then sig limit
	// For each input -> output must excist and not have been spend.
	// Reject if input less then output.
	// Reject if if tx value is too low to get into empty block.

	// A matching transaction must excist?

}

func CreateTransaction() {
}

func CommitTransactionSpend() {
}

func SendMoney() {
}

// First transaction in a block.
func coinBaseTransaction() {
}
