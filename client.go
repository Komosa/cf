package main

import (
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"github.com/Komosa/persistent-cookiejar"

	"golang.org/x/net/publicsuffix"
)

const CFURL = "https://codeforces.com"

func (cf *cf) initClient() error {
	sysuser, err := user.Current()
	if err != nil {
		return err
	}
	fp := filepath.Clean(sysuser.HomeDir + "/.config/cf/")
	if err = os.MkdirAll(fp, 0700); err != nil {
		return err
	}

	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
		Filename:         filepath.Clean(fp + "/" + cf.config["user"] + ".cookie"),
	})
	cf.client = &http.Client{Jar: jar}
	return err
}

func (cf *cf) saveCookie() error {
	if c := cf.client; c != nil {
		j, ok := c.Jar.(*cookiejar.Jar)
		if ok {
			return j.Save()
		}
	}
	return nil
}
