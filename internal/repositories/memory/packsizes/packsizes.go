package packsizes

import (
	"sync"
)

type PackSizes struct {
	m         sync.RWMutex
	packsizes map[int][]int
}

func NewPackSizes() *PackSizes {
	return &PackSizes{
		packsizes: make(map[int][]int),
	}
}

func (p *PackSizes) Store(pid int, packs []int) {
	p.m.Lock()
	defer p.m.Unlock()

	p.packsizes[pid] = packs
}

func (p *PackSizes) PackSizes(pid int) []int {
	p.m.Lock()
	defer p.m.Unlock()

	return p.packsizes[pid]
}
