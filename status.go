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
	Result []MsgSub
}

type MsgSub struct {
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

func (cf *cf) status(prob probCode) {
	urlBase := CFURL + "/api/user.status?handle=" + cf.config["user"] + "&count="
	url := urlBase + "1"
	subCnt := 1
	var err error
	var resp *http.Response
	var body []byte
	spin := `|/-\`
	empty := strings.Repeat(" ", 80)
	var subID int
	sleepTimeInc := time.Second / 5
	sleepTime := sleepTimeInc
	for i := 0; ; i++ {
		time.Sleep(sleepTime)
		sleepTime += sleepTimeInc
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
		var m *MsgSub
		for _, msg := range Message.Result {
			if (subID == 0 || subID == msg.Id) && msg.ContestId == prob.contest && msg.Problem.Index == prob.task {
				m = &msg
				subID = m.Id
				break
			}
		}
		if m == nil {
			if subCnt > 100 {
				if subID != 0 {
					err = fmt.Errorf("could not find submission %d in last %d submissions", subID, subCnt)
				} else {
					err = fmt.Errorf("could not find submission for problem %d%s", prob.contest, prob.task)
				}
				break
			}
			subCnt *= 2
			url = fmt.Sprintf("%s%d", urlBase, subCnt)
			continue
		}

		buf := &bytes.Buffer{}
		fmt.Fprintf(buf, "\rid=%v problem=%v%v ... %v   %s%s", m.Id, m.ContestId, prob.task, string(spin[i%len(spin)]), formatVerdict(m), empty)
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

func formatVerdict(m *MsgSub) string {
	switch m.Verdict {
	case "TESTING":
		return m.Verdict
	case "OK":
		if m.Testset == "PRETESTS" {
			return greenColor("PRETESTS_PASSED")
		}
		return greenColor("ACCEPTED")
	case "RUNTIME_ERROR":
		fallthrough
	case "WRONG_ANSWER":
		fallthrough
	case "PRESENTATION_ERROR":
		fallthrough
	case "TIME_LIMIT_EXCEEDED":
		fallthrough
	case "MEMORY_LIMIT_EXCEEDED":
		fallthrough
	case "IDLENESS_LIMIT_EXCEEDED":
		return fmt.Sprintf("%s %s", redColor(m.Verdict), yellowColor(fmt.Sprintf("on test %d", m.PassedTestCount+1)))
	default:
		return redColor(m.Verdict)
	}
}

func redColor(s string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", 31, s)
}

func greenColor(s string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", 32, s)
}

func yellowColor(s string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", 33, s)
}
