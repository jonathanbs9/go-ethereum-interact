package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
	"github.com/jonathanbs9/go-ethereum-interact/models"
	"github.com/jonathanbs9/go-ethereum-interact/modules"
)

// ClientHandler ethereum client instance
type ClientHandler struct {
	*ethclient.Client
}

func (client ClientHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get param from url
	vars := mux.Vars(r)
	module := vars["module"]

	// Get the query parameters from url request
	address := r.URL.Query().Get("address")
	hash := r.URL.Query().Get("hash")

	// Set our response header
	w.Header().Set("Content-Type", "application/json")

	// Hanlde each request using module parameter
	switch module {
	case "latest-block":
		_block := modules.GetLatestBlock(*client.Client)
		json.NewEncoder(w).Encode(_block)
	case "get-tx":
		if hash == "" {
			json.NewEncoder(w).Encode(&models.Error{
				Code:    400,
				Message: "Malformed request",
			})
			return
		}
		txHash := common.HexToHash(hash)
		_tx := modules.GetTxByHash(*client.Client, txHash)

		if _tx != nil {
			json.NewEncoder(w).Encode(_tx)
			return
		}

		json.NewEncoder(w).Encode(&models.Error{
			Code:    404,
			Message: "Tx Not found!",
		})

	case "send-eth":
		decoder := json.NewDecoder(r.Body)
		var t models.TransferEtheRequest

		err := decoder.Decode(&t)
		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode(&models.Error{
				Code:    400,
				Message: "Malformed request",
			})

			return
		}
		_hash, err := modules.TransferEth(*client.Client, t.PrivKey, t.To, t.Amount)

		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode(&models.Error{
				Code:    500,
				Message: "Internal Server Error",
			})
			return
		}
		json.NewEncoder(w).Encode(&models.HashResponse{
			Hash: _hash,
		})

	case "get-balance":
		if address == "" {
			json.NewEncoder(w).Encode(&models.Error{
				Code:    400,
				Message: "Malformed request",
			})
			return
		}
		balance, err := modules.GetAddressBalance(*client.Client, address)

		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode(&models.Error{
				Code:    500,
				Message: "Internal Server Error",
			})
			return
		}
		json.NewEncoder(w).Encode(&models.BalanceResponse{
			Address: address,
			Balance: balance,
			Symbol:  "Ether",
			Units:   "Wei",
		})

	}

}
