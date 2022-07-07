package main

import "github.com/emersion/go-imap/client"

func (m Mail) dial() (*client.Client, error) {
	var mc *client.Client
	if m.ImapType == "TLS" {
		mcTLS, err := client.DialTLS(m.DialString, nil)
		if err != nil {
			return nil, err
		}
		mc = mcTLS
	}

	if err := mc.Login(m.User, m.Password); err != nil {
		return nil, err
	}
	return mc, nil
}
