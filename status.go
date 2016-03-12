package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var Message struct {
	Status string
	Result []struct {
		Id                  int
		ContestId           int
		CreationTimeSeconds int
		RelativeTimeSeconds int
		Problem             struct {
			ContestId int
			Index     string
			Name      string
			Points    float64
			Tags      []string
		}
		Author struct {
			ContestId int
			Members   []struct {
				Handle string
			}
			ParticipantType  string
			Ghost            bool
			StartTimeSeconds int
		}
		ProgrammingLanguage string
		Verdict             string
		Testset             string
		PassedTestCount     int
		TimeConsumedMillis  int
		MemoryConsumedBytes int
	}
}

func (cf *cf) status() {
	url := CFURL + "/api/user.status?handle=" + cf.config["user"] + "&count=1"
	var err error
	var resp *http.Response
	var body []byte
	spin := `|/-\`
	empty := strings.Repeat(" ", 80)
	for i := 0; ; i++ {
		time.Sleep(time.Second / 5)
		resp, err = http.Get(url)
		if err != nil {
			break
		}

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			break
		}
		resp.Body.Close()

		err = json.Unmarshal(body, &Message)
		if len(Message.Result) == 0 {
			break
		}
		m := Message.Result[0]
		buf := &bytes.Buffer{}
		fmt.Fprintf(buf, "\rid=%v problem=%v%v ... %v   %s%s", m.Id, m.ContestId, m.Problem.Index, string(spin[i%len(spin)]), m.Verdict, empty)
		fmt.Print(buf.String()[:80])
		if m.Verdict != "TESTING" && m.Verdict != "" {
			break
		}
	}
	fmt.Println()
	if err != nil {
		log.Println("status:", err)
	}
}
