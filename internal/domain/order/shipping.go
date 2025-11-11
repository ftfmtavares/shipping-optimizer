// Package order holds logic and representation of orders data
package order

// Order holds data of a given order
type Order struct {
	PID int
	Qty int
}

// Pack holds data of a given package size quantity

type Pack struct {
	PackSize int
	Quantity int
}

// Shipping holds data of an optimized shipping plan

type Shipping struct {
	PID        int
	Order      int
	Packs      []Pack
	PacksCount int
	Total      int
	Excess     int
}
