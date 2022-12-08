package global

import (
	"github.com/Yosorable/ms-admin/core/config"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

var (
	CONFIG   *config.Config
	DATABASE *gorm.DB

	ConcurrencyControl = &singleflight.Group{}
)
