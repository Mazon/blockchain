package transaction

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
type TxInput struct {
	PrevOutputHash  []byte
	PrevOutputIndex uint32
	Signature       []byte // private key signature of owner.
}

//
// An output of a transaction.  It contains the public key that the next input
// must be able to sign with to claim it.
//
type TxOutput struct {
	Value  uint64 //the number of coins.
	PubKey []byte //pubkey of receiver.
}

//
// The basic transaction that is broadcasted on the network and contained in
// blocks.  A transaction can contain multiple inputs and outputs.
//
type Transaction struct {
	Input  []TxInput
	Output []TxOutput
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

//<Signature from Private Key A> <Public Key A> OP_CHECKSIG
func validateTransaction() {
}
