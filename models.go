package main

import (
	"fmt"
	"github.com/emersion/go-imap"
	"log"
	"time"
)

type Item struct {
	CreatedAt time.Time
	SaveAt    time.Time
	Save      bool
	CopyAt    time.Time
	Copy      bool
	SaveFile  string
	Folder    string        `gorm:"primaryKey"`
	MailAt    time.Time     `gorm:"primaryKey"`
	MailID    string        `gorm:"primaryKey"`
	cfg       *Job          `gorm:"-"`
	msg       *imap.Message `gorm:"-"`
	FileSave  string
}

func ItemStart(folder string, msg *imap.Message, job *Job) {
	var i Item
	i.Folder = folder
	i.MailAt = msg.InternalDate
	i.MailID = msg.Envelope.MessageId
	i.Pull(job)
	i.cfg = job
	i.msg = msg
	if job.SaveMail {
		i.Save = true
	}
	err := i.SaveMail()
	if err != nil {
		job.Log.Error(err)
		return
	}
}

func (i *Item) Pull(job *Job) {
	job.Db.Take(i)
}
func (i *Item) Push(job *Job) {
	job.Db.Save(i)
}

func (i *Item) SaveMail() error {
	if !i.Save {
		return nil
	}
	if !i.SaveAt.IsZero() {
		return nil
	}
	if i.FileSave != "" {
		folder := fmt.Sprintf("%s/mail/%s", i.cfg.Path, i.Folder)
		fileMsg := fmt.Sprintf("%s/%s", folder, i.msg.Envelope.MessageId)
		i.FileSave = fileMsg
	}
	if !i.SaveAt.IsZero() {
		return nil
	}
	literal := i.msg.GetBody(&imap.BodySectionName{})
	buf := make([]byte, literal.Len())
	read, err := literal.Read(buf)
	if err != nil {
		return err
	}
	log.Println(read)
	log.Println(buf)
	return nil
}
