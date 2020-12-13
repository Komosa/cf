package main

import (
	"errors"
	"log"
	"strings"

	"golang.org/x/net/html"
)

func (cf *cf) login(user, pass string) error {
	err := cf.initClient()
	if err != nil {
		return err
	}

	resp, err := cf.client.Get(CFURL + "/enter")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if url := resp.Request.URL.String(); strings.Contains(url, "/profile/") {
		log.Println("login: already logged as", url[strings.LastIndex(url, "/")+1:])
		return nil
	}

	tree, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}
	enterForm := selEnter.MatchFirst(tree)
	if tree == nil {
		return errors.New("login: could not match enter form")
	}

	inputs := selInput.MatchAll(enterForm)
	f := form(inputs)
	f.Set("handleOrEmail", user)
	f.Set("password", pass)
	f.Set("remember", "checked")

	resp, err = cf.client.PostForm(CFURL+"/enter", f)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	log.Println("login:", resp.Status)

	return err
}
