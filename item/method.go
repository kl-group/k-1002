package item

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Load(db *gorm.DB, l *logrus.Logger) {
	var h handler
	h.Db = db
	h.Log = l
}
