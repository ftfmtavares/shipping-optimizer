// Package repositories handles general repositories logic
package repositories

import (
	"github.com/ftfmtavares/shipping-optimizer/internal/repositories/memory/packsizes"
)

// Repositories holds all repositories
type Repositories struct {
	PackSizes *packsizes.PackSizes
}

// NewAPIRepositories initializes a Repositories for the api application
func NewAPIRepositories() (r Repositories) {
	return Repositories{
		PackSizes: packsizes.NewPackSizes(),
	}
}
