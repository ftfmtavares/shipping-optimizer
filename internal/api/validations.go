// Package api handles the api requests and definitions
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const maxOrder = 10000000

func validatePidVar(w http.ResponseWriter, r *http.Request) (int, bool) {
	pidVar := mux.Vars(r)["pid"]
	convertedPid, err := strconv.Atoi(pidVar)
	if err != nil || convertedPid <= 0 {
		http.Error(w, "product id not valid", http.StatusBadRequest)
		return 0, false
	}

	return convertedPid, true
}

func validateOrderQuery(w http.ResponseWriter, r *http.Request) (int, bool) {
	orders := r.URL.Query()["order"]
	if len(orders) == 0 {
		http.Error(w, "order query parameter must be specified", http.StatusBadRequest)
		return 0, false
	}

	convertedOrder, err := strconv.Atoi(orders[0])
	if err != nil || convertedOrder <= 0 {
		http.Error(w, "order query parameter not valid", http.StatusBadRequest)
		return 0, false
	}
	if convertedOrder > maxOrder {
		http.Error(w, fmt.Sprintf("order too large: maximum %d", maxOrder), http.StatusBadRequest)
		return 0, false
	}

	return convertedOrder, true
}

func validatePackSizesRequest(w http.ResponseWriter, r *http.Request) ([]int, bool) {
	var req *ProductPackSizesRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return nil, false
	}

	for _, size := range req.Packs {
		if size <= 0 {
			http.Error(w, "pack sizes must be positive integers", http.StatusBadRequest)
			return nil, false
		}
	}

	return req.Packs, true
}
