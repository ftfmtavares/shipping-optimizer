package api

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ftfmtavares/shipping-optimizer/internal/domain/product"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type mockProduct struct {
	calledPackSizes *bool
	calledUpdate    *bool
	pid             *int
	packs           *[]int
	response        product.Product
	err             error
}

func (m mockProduct) PackSizes(ctx context.Context, pid int) (product.Product, error) {
	*m.calledPackSizes = true
	*m.pid = pid
	return m.response, m.err
}

func (m mockProduct) Update(ctx context.Context, pid int, packs []int) {
	*m.calledUpdate = true
	*m.pid = pid
	*m.packs = packs
}

func TestProductPackSizes(t *testing.T) {
	var (
		requestedPackSizes bool
		requestedPID       int
	)
	ctx := context.Background()

	testCases := []struct {
		desc              string
		product           mockProduct
		url               string
		pid               string
		expectedPackSizes bool
		expectedPID       int
		expectedCode      int
		expectedBody      string
	}{
		{
			desc:              "invalid product id",
			product:           mockProduct{},
			url:               "/product/abc/packsizes",
			pid:               "abc",
			expectedPackSizes: false,
			expectedPID:       0,
			expectedCode:      http.StatusBadRequest,
			expectedBody:      "product id not valid\n",
		},
		{
			desc: "pack sizes retrieval error",
			product: mockProduct{
				calledPackSizes: &requestedPackSizes,
				calledUpdate:    nil,
				pid:             &requestedPID,
				packs:           nil,
				response:        product.Product{},
				err:             errors.New("error"),
			},
			url:               "/product/1/packsizes",
			pid:               "1",
			expectedPackSizes: true,
			expectedPID:       1,
			expectedCode:      http.StatusInternalServerError,
			expectedBody:      "internal error\n",
		},
		{
			desc: "pack sizes retrieval success",
			product: mockProduct{
				calledPackSizes: &requestedPackSizes,
				calledUpdate:    nil,
				pid:             &requestedPID,
				packs:           nil,
				response: product.Product{
					PID:   1,
					Packs: []int{5, 10, 12},
				},
				err: nil,
			},
			url:               "/product/1/packsizes",
			pid:               "1",
			expectedPackSizes: true,
			expectedPID:       1,
			expectedCode:      http.StatusOK,
			expectedBody:      "{\"pid\":1,\"packs\":[5,10,12]}\n",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			requestedPackSizes = false
			requestedPID = 0

			req := httptest.NewRequest(http.MethodGet, tC.url, nil)
			req = mux.SetURLVars(req, map[string]string{"pid": tC.pid})
			rec := httptest.NewRecorder()

			ProductPackSizes(ctx, tC.product)(rec, req)

			assert.Equal(t, tC.expectedCode, rec.Code)
			assert.Equal(t, tC.expectedBody, rec.Body.String())

			assert.Equal(t, tC.expectedPackSizes, requestedPackSizes)
			assert.Equal(t, tC.expectedPID, requestedPID)
		})
	}
}

func TestStoreProductPackSizes(t *testing.T) {
	var (
		requestedUpdate bool
		requestedPID    int
		requestedPacks  []int
	)
	ctx := context.Background()

	testCases := []struct {
		desc           string
		product        mockProduct
		url            string
		pid            string
		body           string
		expectedUpdate bool
		expectedPID    int
		expectedPacks  []int
		expectedCode   int
		expectedBody   string
	}{
		{
			desc:           "invalid product id",
			product:        mockProduct{},
			url:            "/product/abc/packsizes",
			pid:            "abc",
			body:           "{\"packs\":[5,10,12]}",
			expectedUpdate: false,
			expectedPID:    0,
			expectedPacks:  nil,
			expectedCode:   http.StatusBadRequest,
			expectedBody:   "product id not valid\n",
		},
		{
			desc:           "invalid request json payload",
			product:        mockProduct{},
			url:            "/product/1/packsizes",
			pid:            "1",
			body:           "invalid",
			expectedUpdate: false,
			expectedPID:    0,
			expectedPacks:  nil,
			expectedCode:   http.StatusBadRequest,
			expectedBody:   "invalid request payload\n",
		},
		{
			desc:           "negative pack sizes request",
			product:        mockProduct{},
			url:            "/product/1/packsizes",
			pid:            "1",
			body:           "{\"packs\":[-5,10,12]}",
			expectedUpdate: false,
			expectedPID:    0,
			expectedPacks:  nil,
			expectedCode:   http.StatusBadRequest,
			expectedBody:   "pack sizes must be positive integers\n",
		},
		{
			desc: "pack sizes update success",
			product: mockProduct{
				calledPackSizes: nil,
				calledUpdate:    &requestedUpdate,
				pid:             &requestedPID,
				packs:           &requestedPacks,
				response: product.Product{
					PID:   1,
					Packs: []int{5, 10, 12},
				},
				err: nil,
			},
			url:            "/product/1/packsizes",
			pid:            "1",
			body:           "{\"packs\":[5,10,12]}",
			expectedUpdate: true,
			expectedPID:    1,
			expectedPacks:  []int{5, 10, 12},
			expectedCode:   http.StatusOK,
			expectedBody:   "{\"pid\":1,\"packs\":[5,10,12]}\n",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			requestedUpdate = false
			requestedPID = 0
			requestedPacks = nil

			req := httptest.NewRequest(http.MethodPost, tC.url, bytes.NewReader([]byte(tC.body)))
			req = mux.SetURLVars(req, map[string]string{"pid": tC.pid})
			rec := httptest.NewRecorder()

			StoreProductPackSizes(ctx, tC.product)(rec, req)

			assert.Equal(t, tC.expectedCode, rec.Code)
			assert.Equal(t, tC.expectedBody, rec.Body.String())

			assert.Equal(t, tC.expectedUpdate, requestedUpdate)
			assert.Equal(t, tC.expectedPID, requestedPID)
			assert.Equal(t, tC.expectedPacks, requestedPacks)
		})
	}
}
