package repositories

import (
	"github.com/ftfmtavares/shipping-optimizer/internal/repositories/memory/packsizes"
)

type Repositories struct {
	PackSizes *packsizes.PackSizes
}

func NewAPIRepositories() (r Repositories) {
	return Repositories{
		PackSizes: packsizes.NewPackSizes(),
	}
}
