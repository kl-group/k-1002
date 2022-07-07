package job

import (
	"github.com/emersion/go-imap/client"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type handler struct {
	Db  *gorm.DB
	Log *logrus.Logger
}

type Job struct {
	Purge      bool           `yaml:"purge" gorm:"-"`
	SaveMail   bool           `yaml:"savemail" gorm:"-"`
	Threads    uint           `yaml:"threads" gorm:"-"`
	Src        Mail           `yaml:"src" gorm:"-"`
	Dst        Mail           `yaml:"dst" gorm:"-"`
	Name       string         `yaml:"-" gorm:"-"`
	Log        *logrus.Logger `yaml:"-" gorm:"-"`
	SrcClient  *client.Client `yaml:"-" gorm:"-"`
	DstClient  *client.Client `yaml:"-" gorm:"-"`
	Path       string         `yaml:"-" gorm:"-"`
	sqliteFile string         `yaml:"-" gorm:"-"`
}
type Mail struct {
	ImapType   string `yaml:"imaptype"`
	Server     string `yaml:"server"`
	Port       uint16 `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	SslCheck   bool   `yaml:"sslcheck"`
	DialString string `yaml:"-"`
}
