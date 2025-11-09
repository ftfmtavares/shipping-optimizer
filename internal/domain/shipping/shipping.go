package shipping

type Order struct {
	PID int
	Qty int
}

type Pack struct {
	PackSize int
	Quantity int
}

type Shipping struct {
	PID        int
	Order      int
	Packs      []Pack
	PacksCount int
	Total      int
	Excess     int
}
