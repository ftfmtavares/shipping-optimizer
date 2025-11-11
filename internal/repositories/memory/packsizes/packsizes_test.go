package packsizes

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackSizesStore(t *testing.T) {
	ps := NewPackSizes()

	testCases := []struct {
		desc  string
		pid   int
		packs []int
	}{
		{
			desc:  "new product store",
			pid:   1,
			packs: []int{5, 10, 12},
		},
		{
			desc:  "existing product update",
			pid:   1,
			packs: []int{5, 10, 15, 20},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ps.Store(tC.pid, tC.packs)
			res, err := ps.PackSizes(tC.pid)
			assert.NoError(t, err)
			assert.Equal(t, tC.packs, res)
		})
	}
}

func TestPackSizesPackSizes(t *testing.T) {
	ps := NewPackSizes()

	testCases := []struct {
		desc          string
		pid           int
		expected      []int
		expectedError assert.ErrorAssertionFunc
	}{
		{
			desc:          "non existant product",
			pid:           1,
			expected:      nil,
			expectedError: assert.Error,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			res, err := ps.PackSizes(tC.pid)
			tC.expectedError(t, err)
			assert.Equal(t, tC.expected, res)
		})
	}
}

func TestPackSizesConcurrentAccess(t *testing.T) {
	ps := NewPackSizes()
	wg := sync.WaitGroup{}

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(pid int) {
			defer wg.Done()
			ps.Store(pid, []int{pid})
		}(i)
	}

	wg.Wait()

	wg = sync.WaitGroup{}
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(pid int) {
			defer wg.Done()
			val, err := ps.PackSizes(pid)
			assert.NoError(t, err)
			assert.Equal(t, []int{pid}, val)
		}(i)
	}

	wg.Wait()
}
