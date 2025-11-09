package shipping

import (
	"context"
	"errors"

	"github.com/ftfmtavares/shipping-optimizer/internal/domain/product"
)

type Storage interface {
	PackSizes(int) []int
	Store(int, []int)
}

type Configurator struct {
	storage Storage
}

func NewConfigurator(storage Storage) Configurator {
	return Configurator{
		storage: storage,
	}
}

func (c Configurator) PackSizes(ctx context.Context, pid int) (product.Product, error) {
	packsizes := c.storage.PackSizes(pid)
	if packsizes == nil {
		return product.Product{}, errors.New("product not found")
	}

	return product.Product{
		PID:   pid,
		Packs: packsizes,
	}, nil
}

func (c Configurator) Update(ctx context.Context, pid int, packsizes []int) (product.Product, error) {
	c.storage.Store(pid, packsizes)
	if packsizes == nil {
		return product.Product{}, errors.New("product not found")
	}

	return product.Product{
		PID:   pid,
		Packs: packsizes,
	}, nil
}
