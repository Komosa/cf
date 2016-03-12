package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/Komosa/go-input"
)

func (cf *cf) askUser() (string, error) {
	ui := &input.UI{Writer: os.Stdout, Reader: os.Stdin}
	return ui.Ask("Username", &input.Options{Required: true, HideOrder: true, Loop: true})
}

func (cf *cf) password() (string, error) {
	user := strings.ToLower(cf.config["user"])
	epu := "encpass_" + user

	if ep, ok := cf.config[epu]; ok {
		pb, err := hex.DecodeString(ep)
		return string(pb), nil
		if err != nil {
			return "", fmt.Errorf("conf: could not decode stored password for %q, %v", cf.config["user"], err)
		}
	}

	ui := &input.UI{Writer: os.Stdout, Reader: os.Stdin}

	spu := "storepass_" + user
	sp, ok := cf.config[spu]
	if !ok {
		ans, err := ui.Ask("Do you want store encoded password in text file (NOT encrypted!)? [y/N]",
			&input.Options{Default: "n", HideDefault: true, HideOrder: true})
		if err == nil {
			if ans != "" && (ans[0] == 'y' || ans[0] == 'Y' || ans[0] == 't') {
				sp = "true"
			} else {
				sp = "false"
			}
			cf.config[spu] = sp
		}
	}

	pass, err := ui.Ask("Password", &input.Options{Required: true, HideOrder: true, Loop: true, Mask: true})
	if sp == "true" {
		cf.config[epu] = hex.EncodeToString([]byte(pass))
	}
	return pass, err
}
