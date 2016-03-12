package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/net/html"
)

type probCode struct {
	contest int
	task    byte
}

func (pc probCode) String() string {
	return fmt.Sprintf("%d%c", pc.contest, pc.task)
}

// guess not defined args
func (cf *cf) guessArgs(file *string, prob probCode, lang *int) error {
	f, l := *file, *lang
	needExt := len(f) == 1 && toupper(f[0]) == prob.task

	if l == 0 && needExt {
		// search over disk
		f += "."
		var match []string
		for k := range cf.config {
			if st, err := os.Stat(f + k); err == nil && !st.IsDir() {
				match = append(match, k)
			}
		}
		if len(match) == 0 {
			return fmt.Errorf("submit: could not found solution file for %q", prob.task)
		}
		if len(match) > 1 {
			buf := bytes.NewBufferString("submit: more than one file looks like solution for ")
			fmt.Fprintf(buf, "%v, candidates (with lang IDs):", prob.task)
			for _, k := range match {
				fmt.Fprintf(buf, "\n%s%v (%s)", f, k, cf.config[k])
			}
			return errors.New(buf.String() + "\n\nremove all but one of those files, or\n" +
				" remove all but one of those extensions in conf, or\n specify lang at command line")
		}

		*file = f + match[0]
		l, err := strconv.Atoi(cf.config[match[0]])
		*lang = l
		return err
	}

	if l != 0 && needExt {
		s := strconv.Itoa(l)
		for k, v := range cf.config {
			if v == s {
				*file += "." + k
				return nil
			}
		}
		return fmt.Errorf("submit: extension must be specified for lang %q", l)
	}

	if l == 0 {
		ext := filepath.Ext(f)
		if ext == "" {
			return errors.New("submit: lang must be specified")
		}
		s, ok := cf.config[ext]
		if !ok {
			return fmt.Errorf("submit: unknown file extension %q, (you may add it in conf)", ext)
		}
		l, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*lang = l
	}
	return nil
}

func (cf *cf) submit(file string, prob probCode, lang int) error {
	if _, ok := cf.config["user"]; !ok {
		return errors.New("submit: user must be configured (and logged in)")
	}

	err := cf.guessArgs(&file, prob, &lang)
	if err != nil {
		return err
	}

	err = cf.initClient()
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	resp, err := cf.client.Get(CFURL + "/problemset/submit")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.Request.URL.String() == CFURL+"/enter" {
		return errors.New("submit: not logged in")
	}

	tree, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}
	submitForm := selSubmit.MatchFirst(tree)
	if tree == nil {
		return errors.New("login: could not match enter form")
	}
	action := formAction(submitForm)
	fields := form(selInput.MatchAll(submitForm))
	fields.Del("sourceFile")
	fields.Set("submittedProblemCode", prob.String())
	fields.Set("programTypeId", fmt.Sprintf("%d", lang))

	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	for k, vs := range fields {
		for _, v := range vs {
			w.WriteField(k, v)
		}
	}
	fw, err := w.CreateFormFile("sourceFile", filepath.Base(file))
	if err != nil {
		return err
	}
	if _, err = io.Copy(fw, f); err != nil {
		return err
	}
	w.Close()

	req, err := http.NewRequest("POST", CFURL+"/problemset/submit"+action, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err = cf.client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	log.Println("submit:", resp.Status)

	cf.status()
	return nil
}