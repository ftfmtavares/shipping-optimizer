// Package product handles services for products management
package product

import (
	"context"

	"github.com/ftfmtavares/shipping-optimizer/internal/domain/product"
)

// Storage provides storage access to products package sizes
type Storage interface {
	PackSizes(int) ([]int, error)
	Store(int, []int)
}

// Configurator provides the products package sizes management service
type Configurator struct {
	storage Storage
}

// NewConfigurator returns an initialized Configurator
func NewConfigurator(storage Storage) Configurator {
	return Configurator{
		storage: storage,
	}
}

// PackSizes method retrieves the package sizes set of a given product
func (c Configurator) PackSizes(ctx context.Context, pid int) (product.Product, error) {
	packsizes, err := c.storage.PackSizes(pid)
	if err != nil {
		return product.Product{}, err
	}

	return product.Product{
		PID:   pid,
		Packs: packsizes,
	}, nil
}

// Update method stores a new package sizes set for a given product
func (c Configurator) Update(ctx context.Context, pid int, packsizes []int) {
	c.storage.Store(pid, packsizes)
}
