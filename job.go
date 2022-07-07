package main

import (
	"fmt"
	"github.com/emersion/go-imap"
	"log"
	"os"
	"time"
)

func (j Job) start() {
	err := j.DialSrc()
	if err != nil {
		j.Log.Info("Error connect src")
		j.Log.Error(err)
		return
	}
	err = j.DialDst()
	if err != nil {
		j.Log.Info("Error connect dst")
		j.Log.Error(err)
		return
	}

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- j.SrcClient.List("", "*", mailboxes)
	}()
	for m := range mailboxes {
		if err = j.ImapCheckFolder(m.Name); err != nil {
			j.Log.Info("Error Check folder imap")
			j.Log.Error(err)
			return
		}
	}
	if err = <-done; err != nil {
		j.Log.Error("Error mailbox")
		j.Log.Error(err)
		return
	}
	for _, v := range j.ImapFolders {
		if err = j.loadFolder(v); err != nil {
			j.Log.Error("Error get loadFolder")
			j.Log.Error(err)
			continue
		}
	}
	// Select INBOX
	mbox, err := j.SrcClient.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Flags for INBOX:", mbox.Flags)
}
func (j *Job) loadFolder(folder ImapFolder) error {
	mbox, err := j.SrcClient.Select(folder.Name, true)
	if err != nil {
		return err
	}
	countAll := mbox.Messages
	page := 10
	j.Log.Infof("In Folder %s: %d messages", folder, countAll)
	for i := 1; i < 100; i = i + page {
		time.Sleep(1 * time.Second)
		if err = j.LoadMessages(i, i+page, folder.Name); err != nil {
			return err
		}
	}

	log.Println("Done!")
	return nil
}

func (j *Job) DialSrc() error {
	j.Log.Debug("Connecting to server...")
	cc, err := j.Src.dial()
	if err != nil {
		return err
	}
	j.SrcClient = cc

	return nil
}

func (j *Job) DialDst() error {
	j.Log.Debug("Connecting to server...")
	cc, err := j.Dst.dial()
	if err != nil {
		return err
	}
	j.DstClient = cc

	return nil
}

func (j *Job) ImapCheckFolder(folder string) error {
	var x ImapFolder
	x.Name = folder
	fl := fmt.Sprintf("%s/mail/%s", j.Path, folder)
	err := os.MkdirAll(fl, os.ModePerm)
	if err != nil {
		return err
	}
	//todo проверить папку в ящике назначения
	j.ImapFolders = append(j.ImapFolders, x)
	return nil
}

func (j *Job) LoadMessages(start, count int, folder string) error {
	seqset := new(imap.SeqSet)
	seqset.AddRange(uint32(start), uint32(count))
	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- j.SrcClient.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()
	for msg := range messages {
		ItemStart(folder, msg, j)
	}
	if err := <-done; err != nil {
		return err
	}
	return nil
}
