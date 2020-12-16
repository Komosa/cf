package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var usageText = `usage:
cf submit FILE[.EXT] [-lang=LANG] [-prob=PROBLEM]
cf login [HANDLE]
cf con [CONTEST]
cf conf default
cf conf lang=LANGID

Where FILE is file to be submitted
when file has EXTension, it helps determine default language
otherwise file and language will be guessed.

When LANG is specified file extension will be not guessed
ang given lang will be used.

PROBLEM may contain problem letter (A,b,..) and/or contest number (eg. 42)
if part is not specified it will be guessed from extension (letter)
or saved config (contest).

CONTEST will be saved in conf.

HANDLE will be used to login, when not specified last login will be used.

Language config could be set to default for all languages
	or set to LANGID for specific one.
`

func main() {
	err := run()
	show(err)
	if err != nil {
		os.Exit(1)
	}
}

func show(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "err:", err)
	}
}

func run() error {
	cmd, param, prob, lang, err := parseArgs(os.Args[1:])
	if err != nil {
		return err
	}
	cf, err := newCF()
	if err != nil {
		return err
	}
	switch cmd {
	case "login":
		if param == "" {
			if u, ok := cf.config["user"]; ok {
				param = u
			} else {
				param, err = cf.askUser()
				if err != nil {
					return err
				}
			}
		}
		cf.config["user"] = param
		p, err := cf.password()
		if err != nil {
			return err
		}
		err = cf.login(param, p)
		if err != nil {
			return err
		}
		err = cf.save()
		if err != nil {
			return err
		}
		return cf.saveCookie()
	case "submit":
		if param == "" {
			if prob.task != "" {
				param = prob.task
			} else {
				return errors.New("submit: empty problem code")
			}
		}
		if prob.contest == 0 {
			prob.contest, err = cf.contest()
			if err != nil {
				return err
			}
		}
		if prob.task == "" {
			base := filepath.Base(param)
			idx := strings.LastIndexByte(base, '.')
			if idx != -1 {
				base = base[:idx]
			}

			if len(base) == 1 && isletter(rune(base[0])) {
				prob.task = base
			} else if len(base) == 2 && isletter(rune(base[0])) && base[1] >= '0' && base[1] <= '9' {
				prob.task = base
			} else {
				return errors.New("submit: please explicitly pass problem code (-prob=...)")
			}

		}
		prob.task = strings.ToUpper(prob.task)

		err = cf.submit(param, prob, lang)
		if err != nil {
			return err
		}
		return cf.saveCookie()
	case "con":
		cf.config["contest"] = param
		return cf.save()
	case "conf":
		return cf.reconf(param)
	case "help":
		fallthrough
	default: // help
		return errors.New(usageText)
	}

}

func parseArgs(args []string) (cmd, param string, prob probCode, lang int, err error) {
	for _, a := range args {
		if len(a) == 0 {
			continue
		}
		if a[0] == '-' {
			if strings.HasPrefix(a, "-prob=") {
				p := a[6:]
				// contest must be first or skipped
				// then goes problem code
				// eg: 'a', '555a', '670d1', 'd1'
				// 'a670' isn't allowed ('d1670' would be ambigous)
				switch {
				case len(p) == 0:
					err = errors.New("cmdline: empty problem code")
					return
				case isletter(rune(p[0])):
					prob.task = p
				default:
					taskIdx := strings.IndexFunc(p, isletter)
					if taskIdx == -1 {
						prob.contest, err = strconv.Atoi(p)
					} else {
						prob.task = p[taskIdx:]
						prob.contest, err = strconv.Atoi(p[:taskIdx])
					}
					if err != nil {
						err = fmt.Errorf("cmdline: invalid problem code %q", p)
						return
					}
				}
			} else if strings.HasPrefix(a, "-lang=") {
				lang, err = strconv.Atoi(a[6:])
				if err != nil {
					return
				}
			} else {
				err = fmt.Errorf("cmdline: unknown parameter %q", a)
				return
			}
			continue
		}
		if len(cmd) == 0 {
			cmd = a
		} else if len(param) == 0 {
			param = a
		} else {
			err = fmt.Errorf("cmdline: too much positional args")
			return
		}
	}
	return
}

func isletter(ch rune) bool {
	return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z')
}
