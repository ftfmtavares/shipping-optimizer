package order

import (
	"context"
	"errors"
	"testing"

	"github.com/ftfmtavares/shipping-optimizer/internal/domain/order"
	"github.com/stretchr/testify/assert"
)

type mockStorage struct{}

func (m mockStorage) PackSizes(pid int) ([]int, error) {
	switch pid {
	case 0:
		return []int{}, nil
	case 1:
		return []int{5, 10, 12}, nil
	case 2:
		return []int{23, 31, 53}, nil
	case 3:
		return []int{23, 31, 53, 79, 97, 113, 137}, nil
	}
	return nil, errors.New("error")
}

func TestShippingCalculateShipping(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		desc          string
		pid           int
		order         order.Order
		expected      order.Shipping
		expectedError assert.ErrorAssertionFunc
	}{
		{
			desc: "no order",
			pid:  1,
			order: order.Order{
				PID: 1,
				Qty: 0,
			},
			expected:      order.Shipping{},
			expectedError: assert.Error,
		},
		{
			desc: "product not found",
			pid:  -1,
			order: order.Order{
				PID: -1,
				Qty: 21,
			},
			expected:      order.Shipping{},
			expectedError: assert.Error,
		},
		{
			desc: "no pack sizes defined",
			pid:  0,
			order: order.Order{
				PID: 0,
				Qty: 21,
			},
			expected:      order.Shipping{},
			expectedError: assert.Error,
		},
		{
			desc: "simple case",
			pid:  1,
			order: order.Order{
				PID: 1,
				Qty: 21,
			},
			expected: order.Shipping{
				PID:   1,
				Order: 21,
				Packs: []order.Pack{
					{
						PackSize: 10,
						Quantity: 1,
					},
					{
						PackSize: 12,
						Quantity: 1,
					},
				},
				PacksCount: 2,
				Total:      22,
				Excess:     1,
			},
			expectedError: assert.NoError,
		},
		{
			desc: "target case",
			pid:  2,
			order: order.Order{
				PID: 2,
				Qty: 500000,
			},
			expected: order.Shipping{
				PID:   2,
				Order: 500000,
				Packs: []order.Pack{
					{
						PackSize: 23,
						Quantity: 2,
					},
					{
						PackSize: 31,
						Quantity: 7,
					},
					{
						PackSize: 53,
						Quantity: 9429,
					},
				},
				PacksCount: 9438,
				Total:      500000,
				Excess:     0,
			},
			expectedError: assert.NoError,
		},
		{
			desc: "load case",
			pid:  3,
			order: order.Order{
				PID: 3,
				Qty: 100000000,
			},
			expected: order.Shipping{
				PID:   3,
				Order: 100000000,
				Packs: []order.Pack{
					{
						PackSize: 97,
						Quantity: 1,
					},
					{
						PackSize: 113,
						Quantity: 4,
					},
					{
						PackSize: 137,
						Quantity: 729923,
					},
				},
				PacksCount: 729928,
				Total:      100000000,
				Excess:     0,
			},
			expectedError: assert.NoError,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			optimizer := NewOptimizer(mockStorage{})

			res, err := optimizer.Calculate(ctx, tC.order)
			tC.expectedError(t, err)
			assert.Equal(t, tC.expected, res)
		})
	}
}
