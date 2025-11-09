package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ftfmtavares/shipping-optimizer/internal/domain/product"
)

type Product interface {
	PackSizes(context.Context, int) (product.Product, error)
	Update(context.Context, int, []int) (product.Product, error)
}

type ProductPackSizesResponse struct {
	PID   int   `json:"pid"`
	Packs []int `json:"packs"`
}

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

type ProductPackSizesRequest struct {
	Packs []int `json:"packs"`
}

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

		prd, err := updater.Update(ctx, productID, packSizes)
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
