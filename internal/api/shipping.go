package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ftfmtavares/shipping-optimizer/internal/domain/shipping"
)

type ShippingOptimizer interface {
	Calculate(context.Context, shipping.Order) (shipping.Shipping, error)
}

type PackResponse struct {
	PackSize int `json:"packsize"`
	Quantity int `json:"quantity"`
}

type ShippingCalculationResponse struct {
	Order      int            `json:"order"`
	Packs      []PackResponse `json:"packs"`
	PacksCount int            `json:"packscount"`
	Total      int            `json:"total"`
	Excess     int            `json:"excess"`
}

func ShippingCalculation(ctx context.Context, calculator ShippingOptimizer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productID, valid := validatePidVar(w, r)
		if !valid {
			return
		}

		order, valid := validateOrderQuery(w, r)
		if !valid {
			return
		}

		sd, err := calculator.Calculate(ctx, shipping.Order{
			PID: productID,
			Qty: order,
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
