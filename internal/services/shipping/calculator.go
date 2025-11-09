package shipping

import (
	"context"
	"errors"
	"math"
	"sort"

	"github.com/ftfmtavares/shipping-optimizer/internal/domain/shipping"
)

type Storage interface {
	PackSizes(int) []int
}

type Optimizer struct {
	storage Storage
}

func NewOptimizer(storage Storage) Optimizer {
	return Optimizer{
		storage: storage,
	}
}

func (o Optimizer) Calculate(ctx context.Context, order shipping.Order) (shipping.Shipping, error) {
	if order.Qty <= 0 {
		return shipping.Shipping{}, errors.New("empty order")
	}

	packsizes := o.storage.PackSizes(order.PID)
	if len(packsizes) == 0 {
		return shipping.Shipping{}, errors.New("no pack sizes found for product")
	}

	packs, totalCount, packsCount := optimizeShipping(packsizes, order.Qty)

	return shipping.Shipping{
		PID:        order.PID,
		Order:      order.Qty,
		Packs:      packs,
		PacksCount: packsCount,
		Total:      totalCount,
		Excess:     totalCount - order.Qty,
	}, nil
}

func optimizeShipping(packSizes []int, order int) ([]shipping.Pack, int, int) {
	const inf = math.MaxInt

	type checkpoint struct {
		packsCount int
		packsize   int
	}

	sort.Ints(packSizes)
	limit := order + packSizes[len(packSizes)-1]

	cps := make([]checkpoint, limit)
	for i := range cps {
		cps[i].packsCount = inf
	}
	cps[0].packsCount = 0

	for t := 0; t < limit; t++ {
		if cps[t].packsCount == inf {
			continue
		}

		for _, size := range packSizes {
			next := t + size
			if next < limit && cps[t].packsCount+1 < cps[next].packsCount {
				cps[next].packsCount = cps[t].packsCount + 1
				cps[next].packsize = size
			}
		}
	}

	var best int
	for t := order; t < limit; t++ {
		if cps[t].packsCount < inf {
			best = t
			break
		}
	}

	packCountsMap := make(map[int]int)
	for t := best; t > 0; {
		packCountsMap[cps[t].packsize]++
		t -= cps[t].packsize
	}

	bestCounts := make([]shipping.Pack, 0, len(packSizes))
	for _, size := range packSizes {
		count := packCountsMap[size]
		if count == 0 {
			continue
		}

		bestCounts = append(bestCounts, shipping.Pack{
			PackSize: size,
			Quantity: count,
		})
	}

	return bestCounts, best, cps[best].packsCount
}
