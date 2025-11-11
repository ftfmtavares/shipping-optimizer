// Package order handles services for orders management
package order

import (
	"context"
	"errors"
	"math"
	"sort"

	"github.com/ftfmtavares/shipping-optimizer/internal/domain/order"
)

// Storage provides storage retrieval access to products package sizes
type Storage interface {
	PackSizes(int) ([]int, error)
}

// Optimizer provides the order packages calculation service
type Optimizer struct {
	storage Storage
}

// NewOptimizer returns an initialized Optimizer
func NewOptimizer(storage Storage) Optimizer {
	return Optimizer{
		storage: storage,
	}
}

// Calculate method calculates the best packages distribution for a given order
func (o Optimizer) Calculate(ctx context.Context, req order.Order) (order.Shipping, error) {
	if req.Qty <= 0 {
		return order.Shipping{}, errors.New("empty order")
	}

	packsizes, err := o.storage.PackSizes(req.PID)
	if err != nil {
		return order.Shipping{}, errors.New("no product found")
	}

	if len(packsizes) == 0 {
		return order.Shipping{}, errors.New("no pack sizes found for product")
	}

	packs, totalCount, packsCount := optimizeShipping(packsizes, req.Qty)

	return order.Shipping{
		PID:        req.PID,
		Order:      req.Qty,
		Packs:      packs,
		PacksCount: packsCount,
		Total:      totalCount,
		Excess:     totalCount - req.Qty,
	}, nil
}

func optimizeShipping(packSizes []int, qty int) ([]order.Pack, int, int) {
	const inf = math.MaxInt

	// a checkpoint holds an intermediate calculation for a given order quantity
	// packsCount stores the least amount of packages that serves that exact quantity
	// packSize indicates the size of the last package so that it can be back tracked
	type checkpoint struct {
		packsCount int
		packSize   int
	}

	// the total items can't exceed the actual order quantity plus the size of the smallest package
	sort.Ints(packSizes)
	limit := qty + packSizes[0]

	// the initial checkpoints slice starts with the highest number of packages for comparison purposes
	cps := make([]checkpoint, limit)
	for i := range cps {
		cps[i].packsCount = inf
	}
	cps[0].packsCount = 0

	// each checkpoint is checked for a possible matching combination
	// all existing packages are added on top of valid checkpoints and the best ones are kept
	for t := range limit {
		if cps[t].packsCount == inf {
			continue
		}

		for _, size := range packSizes {
			next := t + size
			if next < limit && cps[t].packsCount+1 < cps[next].packsCount {
				cps[next].packsCount = cps[t].packsCount + 1
				cps[next].packSize = size
			}
		}
	}

	// finds the best valid combination
	var best int
	for t := qty; t < limit; t++ {
		if cps[t].packsCount < inf {
			best = t
			break
		}
	}

	// backtracks through the checkpoints counting the number of each package size
	packCountsMap := make(map[int]int)
	for t := best; t > 0; {
		packCountsMap[cps[t].packSize]++
		t -= cps[t].packSize
	}

	// prepares the package sizes combination filtering the not used ones
	bestCounts := make([]order.Pack, 0, len(packSizes))
	for _, size := range packSizes {
		count := packCountsMap[size]
		if count == 0 {
			continue
		}

		bestCounts = append(bestCounts, order.Pack{
			PackSize: size,
			Quantity: count,
		})
	}

	return bestCounts, best, cps[best].packsCount
}
