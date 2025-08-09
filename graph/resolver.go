package graph

import (
	"trade_company/internal/config"

	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB  *gorm.DB
	Cfg *config.Config
}
