package product

import (
	"context"
	"errors"
	"testing"

	"github.com/ftfmtavares/shipping-optimizer/internal/domain/product"
	"github.com/stretchr/testify/assert"
)

type mockStorage struct {
	calledPackSizes *bool
	calledStore     *bool
	pid             *int
	packs           *[]int
	response        []int
	err             error
}

func (m mockStorage) PackSizes(pid int) ([]int, error) {
	*m.calledPackSizes = true
	*m.pid = pid
	return m.response, m.err
}

func (m mockStorage) Store(pid int, packs []int) {
	*m.calledStore = true
	*m.pid = pid
	*m.packs = packs
}

func TestPackSizes(t *testing.T) {
	var (
		requestedPackSizes bool
		requestedPID       int
	)
	ctx := context.Background()

	testCases := []struct {
		desc              string
		storage           mockStorage
		pid               int
		expectedPackSizes bool
		expectedPID       int
		expected          product.Product
		expectedError     assert.ErrorAssertionFunc
	}{
		{
			desc: "product not found",
			storage: mockStorage{
				calledPackSizes: &requestedPackSizes,
				calledStore:     nil,
				pid:             &requestedPID,
				packs:           nil,
				response:        nil,
				err:             errors.New("error"),
			},
			pid:               1,
			expectedPackSizes: true,
			expectedPID:       1,
			expected:          product.Product{},
			expectedError:     assert.Error,
		},
		{
			desc: "product found",
			storage: mockStorage{
				calledPackSizes: &requestedPackSizes,
				calledStore:     nil,
				pid:             &requestedPID,
				packs:           nil,
				response:        []int{5, 10, 12},
				err:             nil,
			},
			pid:               1,
			expectedPackSizes: true,
			expectedPID:       1,
			expected: product.Product{
				PID:   1,
				Packs: []int{5, 10, 12},
			},
			expectedError: assert.NoError,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			requestedPackSizes = false
			requestedPID = 0

			cfg := NewConfigurator(tC.storage)
			res, err := cfg.PackSizes(ctx, tC.pid)
			tC.expectedError(t, err)
			assert.Equal(t, tC.expected, res)

			assert.Equal(t, tC.expectedPackSizes, requestedPackSizes)
			assert.Equal(t, tC.expectedPID, requestedPID)
		})
	}
}

func TestUpdate(t *testing.T) {
	var (
		requestedUpdate bool
		requestedPID    int
		requestedPacks  []int
	)
	ctx := context.Background()

	testCases := []struct {
		desc           string
		storage        mockStorage
		pid            int
		packSizes      []int
		expectedUpdate bool
		expectedPID    int
		expectedPacks  []int
	}{
		{
			desc: "update success",
			storage: mockStorage{
				calledPackSizes: nil,
				calledStore:     &requestedUpdate,
				pid:             &requestedPID,
				packs:           &requestedPacks,
				response:        nil,
			},
			pid:            1,
			packSizes:      []int{5, 10, 12},
			expectedUpdate: true,
			expectedPID:    1,
			expectedPacks:  []int{5, 10, 12},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			requestedUpdate = false
			requestedPID = 0
			requestedPacks = nil

			cfg := NewConfigurator(tC.storage)
			cfg.Update(ctx, tC.pid, tC.packSizes)

			assert.Equal(t, tC.expectedUpdate, requestedUpdate)
			assert.Equal(t, tC.expectedPID, requestedPID)
			assert.Equal(t, tC.expectedPacks, requestedPacks)
		})
	}
}
