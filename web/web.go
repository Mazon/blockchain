package web

import (
	"blockchain/transaction"
	"encoding/json"

	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func Run() error {
	mux := makeMuxRouter()
	//	httpAddr := os.Getenv("ADDR")
	httpAddr := "8080"
	//log.Println("Listening on ", os.Getenv("ADDR"))
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	/*bytes, err := json.MarshalIndent(blockchain.Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
	*/
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// create tx
	data := []byte{102, 97, 108, 99, 111, 110}
	i1 := transaction.TxInput{}
	o1 := transaction.TxOutput{Value: 100, PubKey: data}

	tx := transaction.Transaction{[]transaction.TxInput{i1}, []transaction.TxOutput{o1}}
	_ = tx

	/*	decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&tx); err != nil {
			respondWithJSON(w, r, http.StatusBadRequest, r.Body)
			return
		}*/
	defer r.Body.Close()

	//ensure atomicity when creating new block
	//mutex.Lock()
	//	newBlock := blockchain.GenerateBlock(blockchain.Blockchain[len(blockchain.Blockchain)-1], []transaction.Transaction{tx})
	//mutex.Unlock()

	//if blockchain.IsBlockValid(newBlock, blockchain.Blockchain[len(blockchain.Blockchain)-1]) {
	//	blockchain.Blockchain = append(blockchain.Blockchain, newBlock)
	//	spew.Dump(blockchain.Blockchain)
	//}

	//respondWithJSON(w, r, http.StatusCreated, newBlock)
	respondWithJSON(w, r, http.StatusCreated, 200)

}
