// Package api handles the api requests and definitions
package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ftfmtavares/shipping-optimizer/internal/domain/product"
)

// Product provides the product package sizes management service
type Product interface {
	PackSizes(context.Context, int) (product.Product, error)
	Update(context.Context, int, []int)
}

// ProductPackSizesResponse holds the product package sizes response
type ProductPackSizesResponse struct {
	PID   int   `json:"pid"`
	Packs []int `json:"packs"`
}

// ProductPackSizes handles the product packages sizes retrieval requests
func ProductPackSizes(ctx context.Context, retriever Product) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productID, valid := validatePidVar(w, r)
		if !valid {
			return
		}

		prd, err := retriever.PackSizes(ctx, productID)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(ProductPackSizesResponse{
			PID:   productID,
			Packs: prd.Packs,
		})
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
	}
}

// ProductPackSizesRequest holds the product package sizes update request
type ProductPackSizesRequest struct {
	Packs []int `json:"packs"`
}

// StoreProductPackSizes handles the product packages sizes update requests
func StoreProductPackSizes(ctx context.Context, updater Product) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productID, valid := validatePidVar(w, r)
		if !valid {
			return
		}

		packSizes, valid := validatePackSizesRequest(w, r)
		if !valid {
			return
		}

		updater.Update(ctx, productID, packSizes)

		err := json.NewEncoder(w).Encode(ProductPackSizesResponse{
			PID:   productID,
			Packs: packSizes,
		})
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
	}
}
