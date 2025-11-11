// Package api handles the api requests and definitions
package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ftfmtavares/shipping-optimizer/internal/domain/order"
)

// ShippingOptimizer provides the order packages calculation service
type ShippingOptimizer interface {
	Calculate(context.Context, order.Order) (order.Shipping, error)
}

// PackResponse holds information of a package size quantity
type PackResponse struct {
	PackSize int `json:"packsize"`
	Quantity int `json:"quantity"`
}

// ShippingCalculationResponse holds the orders calculation response
type ShippingCalculationResponse struct {
	Order      int            `json:"order"`
	Packs      []PackResponse `json:"packs"`
	PacksCount int            `json:"packscount"`
	Total      int            `json:"total"`
	Excess     int            `json:"excess"`
}

// OrderCalculation handles the orders calculation requests
func OrderCalculation(ctx context.Context, calculator ShippingOptimizer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productID, valid := validatePidVar(w, r)
		if !valid {
			return
		}

		orderQty, valid := validateOrderQuery(w, r)
		if !valid {
			return
		}

		sd, err := calculator.Calculate(ctx, order.Order{
			PID: productID,
			Qty: orderQty,
		})
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		packs := make([]PackResponse, 0, len(sd.Packs))
		for _, pack := range sd.Packs {
			if pack.Quantity > 0 {
				packs = append(packs, PackResponse{
					PackSize: pack.PackSize,
					Quantity: pack.Quantity,
				})
			}
		}

		err = json.NewEncoder(w).Encode(ShippingCalculationResponse{
			Order:      sd.Order,
			Packs:      packs,
			PacksCount: sd.PacksCount,
			Total:      sd.Total,
			Excess:     sd.Excess,
		})
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
	}
}
