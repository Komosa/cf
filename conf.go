package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mattn/go-forlines"
)

var defConf = map[string]string{
	// first five langs here have more than one complier on CF currently
	"cpp":  "1",
	"c":    "10",
	"py":   "31",
	"java": "36",
	"pas":  "4",
	// below are obvious ones
	"d":     "28",
	"pl":    "13",
	"rb":    "8",
	"php":   "6",
	"js":    "34",
	"cs":    "29",
	"ml":    "19",
	"go":    "32",
	"scala": "20",
	"hs":    "12",
	"tcl":   "14",
	// esoteric langs skipped
}

func (cf *cf) contest() (int, error) {
	if s, ok := cf.config["contest"]; ok {
		return strconv.Atoi(s)
	}
	return 0, errors.New("conf: 'contest' must be specified via config file, 'con' subcommand or '-prob' switch")

}

func (cf *cf) save() error {
	sysuser, err := user.Current()
	if err != nil {
		return err
	}

	confFile := filepath.Clean(sysuser.HomeDir + "/.config/cf/config")
	if err = os.MkdirAll(filepath.Dir(confFile), 0700); err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	for k, v := range cf.config {
		fmt.Fprintf(buf, "%s=%s\n", k, v)
	}

	return ioutil.WriteFile(confFile, buf.Bytes(), 0600)
}

func load() (map[string]string, error) {
	sysuser, err := user.Current()
	if err != nil {
		return nil, err
	}

	confFile := filepath.Clean(sysuser.HomeDir + "/.config/cf/config")
	if _, err = os.Stat(confFile); os.IsNotExist(err) {
		return defConf, nil
	}

	f, err := os.Open(confFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m := make(map[string]string)
	return m, forlines.Do(f, func(line string) error {
		p := strings.SplitN(line, "=", 3)
		if len(p) != 2 {
			return fmt.Errorf("load conf: invalid line %q", line)
		}
		m[p[0]] = p[1]
		return nil
	})
}
