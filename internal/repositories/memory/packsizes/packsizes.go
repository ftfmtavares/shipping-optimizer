// Package packsizes handles in memory package sizes storage
package packsizes

import (
	"errors"
	"sync"
)

// PackSizes provides in memory storage for products package sizes
type PackSizes struct {
	m         sync.RWMutex
	packsizes map[int][]int
}

// NewPackSizes initializes a new PackSizes
func NewPackSizes() *PackSizes {
	return &PackSizes{
		packsizes: make(map[int][]int),
	}
}

// Store method stores a new package sizes set for a given product
func (p *PackSizes) Store(pid int, packs []int) {
	p.m.Lock()
	defer p.m.Unlock()

	p.packsizes[pid] = packs
}

// PackSizes method retrieves the package sizes set of a given product
func (p *PackSizes) PackSizes(pid int) ([]int, error) {
	p.m.Lock()
	defer p.m.Unlock()

	packs, found := p.packsizes[pid]
	if !found {
		return nil, errors.New("product not found")
	}

	return packs, nil
}
