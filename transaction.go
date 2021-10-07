package main

//
// An input of a transaction.  It contains the location of the previous
// transaction's output that it claims and a signature that matches the
// output's public key.
//
type CTXIn struct {
  COutPoint prevout
}

//
// An output of a transaction.  It contains the public key that the next input
// must be able to sign with to claim it.
//
type CTXOut struct {
  nValue int64
}

//
// The basic transaction that is broadcasted on the network and contained in
// blocks.  A transaction can contain multiple inputs and outputs.
//
type Transaction struct {
  version int
  vin CTXIn
  vout CTXOut
}
