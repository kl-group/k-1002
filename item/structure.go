package item

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type handler struct {
	Db  *gorm.DB
	Log *logrus.Logger
}
