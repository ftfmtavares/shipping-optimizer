package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ftfmtavares/shipping-optimizer/internal/domain/order"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type mockShippingCalculator struct {
	called   *bool
	order    *order.Order
	response order.Shipping
	err      error
}

func (m mockShippingCalculator) Calculate(ctx context.Context, order order.Order) (order.Shipping, error) {
	*m.called = true
	*m.order = order
	return m.response, m.err
}

func TestShippingCalculation(t *testing.T) {
	var (
		requestedCalculation bool
		requestedOrder       order.Order
	)
	ctx := context.Background()

	testCases := []struct {
		desc                string
		calculator          mockShippingCalculator
		url                 string
		pid                 string
		expectedCalculation bool
		expectedOrder       order.Order
		expectedCode        int
		expectedBody        string
	}{
		{
			desc:                "invalid product id",
			calculator:          mockShippingCalculator{},
			url:                 "/product/abc/shipping-calculation?order=10",
			pid:                 "abc",
			expectedCalculation: false,
			expectedOrder:       order.Order{},
			expectedCode:        http.StatusBadRequest,
			expectedBody:        "product id not valid\n",
		},
		{
			desc:                "missing order quantity",
			calculator:          mockShippingCalculator{},
			url:                 "/product/1/shipping-calculation",
			pid:                 "1",
			expectedCalculation: false,
			expectedOrder:       order.Order{},
			expectedCode:        http.StatusBadRequest,
			expectedBody:        "order query parameter must be specified\n",
		},
		{
			desc:                "invalid order quantity",
			calculator:          mockShippingCalculator{},
			url:                 "/product/1/shipping-calculation?order=abc",
			pid:                 "1",
			expectedCalculation: false,
			expectedOrder:       order.Order{},
			expectedCode:        http.StatusBadRequest,
			expectedBody:        "order query parameter not valid\n",
		},
		{
			desc:                "maximum order quantity exceeded",
			calculator:          mockShippingCalculator{},
			url:                 "/product/1/shipping-calculation?order=10000001",
			pid:                 "1",
			expectedCalculation: false,
			expectedOrder:       order.Order{},
			expectedCode:        http.StatusBadRequest,
			expectedBody:        "order too large: maximum 10000000\n",
		},
		{
			desc: "calculation error",
			calculator: mockShippingCalculator{
				called:   &requestedCalculation,
				order:    &requestedOrder,
				response: order.Shipping{},
				err:      errors.New("error"),
			},
			url:                 "/product/1/shipping-calculation?order=21",
			pid:                 "1",
			expectedCalculation: true,
			expectedOrder: order.Order{
				PID: 1,
				Qty: 21,
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "internal error\n",
		},
		{
			desc: "calculation success",
			calculator: mockShippingCalculator{
				called: &requestedCalculation,
				order:  &requestedOrder,
				response: order.Shipping{
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
				err: nil,
			},
			url:                 "/product/1/shipping-calculation?order=21",
			pid:                 "1",
			expectedCalculation: true,
			expectedOrder: order.Order{
				PID: 1,
				Qty: 21,
			},
			expectedCode: http.StatusOK,
			expectedBody: "{\"order\":21,\"packs\":[{\"packsize\":10,\"quantity\":1},{\"packsize\":12,\"quantity\":1}],\"packscount\":2,\"total\":22,\"excess\":1}\n",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			requestedCalculation = false
			requestedOrder = order.Order{}

			req := httptest.NewRequest(http.MethodGet, tC.url, nil)
			req = mux.SetURLVars(req, map[string]string{"pid": tC.pid})
			rec := httptest.NewRecorder()

			OrderCalculation(ctx, tC.calculator)(rec, req)

			assert.Equal(t, tC.expectedCode, rec.Code)
			assert.Equal(t, tC.expectedBody, rec.Body.String())

			assert.Equal(t, tC.expectedCalculation, requestedCalculation)
			assert.Equal(t, tC.expectedOrder, requestedOrder)
		})
	}
}
