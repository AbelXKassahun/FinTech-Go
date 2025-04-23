package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AbelXKassahun/Digital-Wallet-Platform/internal/fee_engine"
)

type report struct {
	FinalFee float64 `json:"final_fee"`
	FeeBreakdown fee_engine.FeeBreakdown `json:"fee_breakdown"`
}

func CalculateFee(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var details fee_engine.TransactionDetails
	err := json.NewDecoder(r.Body).Decode(&details)
    if err != nil {
        http.Error(w, "Invalid request body: " + err.Error(), http.StatusBadRequest)
        return
    }

	fee, breakdown, err := fee_engine.CalculateFee(details)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	report := report{
		FinalFee: fee,
		FeeBreakdown: breakdown,
	}
	err = json.NewEncoder(w).Encode(report)
	if err != nil {
		log.Println("Couldnt encode tokens: " + err.Error())
		http.Error(w, "Couldnt encode tokens: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
