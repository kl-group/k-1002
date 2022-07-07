package main

import (
	"fmt"
	"github.com/emersion/go-imap/client"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"strings"
)

type Config struct {
	Path Path  `yaml:"path"`
	Jobs []Job `yaml:"jobs"`
}

type Path struct {
	Jobs string `yaml:"jobs"`
	Logs string `yaml:"logs"`
}
type Job struct {
	Purge       bool           `yaml:"purge" gorm:"-"`
	SaveMail    bool           `yaml:"savemail" gorm:"-"`
	Threads     uint           `yaml:"threads" gorm:"-"`
	Src         Mail           `yaml:"src" gorm:"-"`
	Dst         Mail           `yaml:"dst" gorm:"-"`
	Name        string         `yaml:"-" gorm:"-"`
	Log         *logrus.Logger `yaml:"-" gorm:"-"`
	SrcClient   *client.Client `yaml:"-" gorm:"-"`
	DstClient   *client.Client `yaml:"-" gorm:"-"`
	ImapFolders []ImapFolder   `yaml:"-" gorm:"-"`
	Path        string         `yaml:"-" gorm:"-"`
	Db          *gorm.DB       `yaml:"-" gorm:"-"`
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
type ImapFolder struct {
	Name   string
	Folder string
}

func loadConfig() error {
	fl := "./config.yaml"
	_, err := os.Stat(fl)
	if err != nil {
		return err
	}
	read, err := os.ReadFile(fl)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(read, &cfg)
	if err != nil {
		return err
	}
	cfg.prepareDefault()
	cfg.checkRequirements()
	return nil
}

func (c *Config) prepareDefault() {
	c.Path.prepareDefault()

	for k, _ := range c.Jobs {
		c.Jobs[k].prepareDefault()
	}
}
func (c *Path) prepareDefault() {
	if c.Logs == "" {
		c.Logs = "./logs"
	}
	if c.Jobs == "" {
		c.Logs = "./jobs"
	}
}
func (j *Job) prepareDefault() {
	if j.Threads == 0 {
		GlobalLog.Info("Threads set default 1")
		j.Threads = 1
	}
	j.Src.prepareDefault()
	j.Dst.prepareDefault()
	j.Name = fmt.Sprintf("%s@%s", strings.Replace(j.Src.Server, ".", "", len(j.Src.Server)), j.Src.User)
	j.Path = fmt.Sprintf("%s/%s", cfg.Path.Jobs, j.Name)
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s/histroy.db", j.Path)), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	db.AutoMigrate(&Item{})
	j.Db = db
	j.Log = logrus.New()
	j.Log.SetFormatter(&logrus.TextFormatter{})
	j.Log.SetOutput(
		&lumberjack.Logger{
			Filename:   fmt.Sprintf("%s/%s/logs/jobs.log", cfg.Path.Jobs, j.Name),
			MaxSize:    1,  // megabytes after which new file is created
			MaxBackups: 3,  // number of backups
			MaxAge:     28, //days
		})
	j.Log.SetLevel(logrus.InfoLevel)
}

func (c *Mail) prepareDefault() {
	if c.ImapType == "" {
		c.ImapType = "TLS"
	}
	if c.ImapType == "TLS" && c.Port == 0 {
		c.Port = 993
	}
	if c.ImapType == "STARTTLS" && c.Port == 0 {
		c.Port = 143
	}
	c.DialString = fmt.Sprintf("%s:%d", c.Server, c.Port)
}

func (c *Mail) checkRequirements() {
	if c.User == "" {
		GlobalLog.Fatal("User is null")
		return
	}
	if c.Password == "" {
		GlobalLog.Fatal("Password is null")
		return
	}
	if c.Server == "" {
		GlobalLog.Fatal("Server is null")
		return
	}
}
func (c *Config) checkRequirements() {
	for k, _ := range c.Jobs {
		c.Jobs[k].Src.checkRequirements()
		c.Jobs[k].Dst.checkRequirements()
	}
}
