package prompt

import (
	"github.com/Ak-Army/git-prompt/color"
	"strings"
	"sync"
	"time"
	"strconv"
)

type GitStatus interface {
	Prompt(color color.Color) string
}

func GetCurrentStatus() *gitStatus {
	gS := &gitStatus{}
	gS.waitGroup = &sync.WaitGroup{}
	err := gS.getStatus()
	if err != nil {
		return gS
	}
	gS.waitGroup.Add(2)
	go gS.getStashCount()
	go gS.getTag()
	gS.waitGroup.Wait()
	return gS
}

// Status ...
type gitStatus struct {
	GitStatus
	waitGroup *sync.WaitGroup
	branch    string
	tag       string
	detached  bool
	hasremote bool
	ahead     int
	behind    int
	staged    int
	conflicts int
	changed   int
	untracked int
	stashs    int
}

// Format ...
func (gS *gitStatus) Prompt(color color.Color) string {
	if gS.branch == "" {
		return ""
	}
	ret := make([]string, 2)
	ret[0] = color.Yellow("(")

	// branch
	if gS.detached {
		ret[0] += color.Red("%s", gS.branch)
	} else {
		ret[0] += color.Cyan("%s", gS.branch)
	}

	if gS.tag != "" {
		ret[0] += color.Cyan(" %s", gS.tag)
	}

	if gS.hasremote {
		if gS.ahead > 0 {
			ret[1] += color.Green("↑%d", gS.ahead)
		}
		if gS.behind > 0 {
			ret[1] += color.Red("↓%d", gS.behind)
		}
	}

	if gS.staged > 0 {
		ret[1] += color.Blue("+%d", gS.staged)
	}

	if gS.changed > 0 {
		ret[1] += color.Blue("△%d", gS.changed)
	}

	if gS.untracked > 0 {
		ret[1] += color.Blue("?%d", gS.untracked)
	}

	if gS.conflicts > 0 {
		ret[1] += color.Red("!%d", gS.conflicts)
	}

	if gS.stashs > 0 {
		ret[1] += color.Green("|s%d", gS.stashs)
	}

	r := strings.TrimSpace(strings.Join(ret, " "))
	r += color.Yellow(")")

	return r
}

func (gS *gitStatus) getStatus() error {
	gitDir, _, err := communicate("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return err
	}
	lastFetch, _, err := communicate("stat", "-c", `%Y`, strings.TrimSpace(gitDir)+`/.git/FETCH_HEAD`)
	if err == nil {
		i, err := strconv.ParseInt(strings.TrimSpace(lastFetch), 10, 64)
		if err == nil && time.Since(time.Unix(i, 0)).Hours() > 1 {
			communicate("git", "fetch")
		}
	}
	lines, err := getLines("git", "status", "--porcelain", "--branch")
	if err != nil {
		return err
	}

	gS.waitGroup.Add(2)
	go gS.parseBranchLine(lines[0])
	go gS.collectChanges(lines[1:])
	return nil
}

func (gS *gitStatus) getStashCount() {
	defer gS.waitGroup.Done()

	stdout, stderr, err := communicate("git", "log", "--format=\"%%gd: %%gs\"", "-g", "--first-parent", "-m", "refs/stash", "--")
	if err != nil {
		return
	} else if strings.Contains(stderr, "fatal") {
		return
	}

	gS.stashs = len(strings.Split(stdout, "\n")) - 1
}

func (gS *gitStatus) parseBranchLine(line string) {
	defer gS.waitGroup.Done()

	if strings.Contains(line, "no branch") {
		gS.detached = true
		gS.branch = "no branch"
		hash, _, _ := communicate("git", "rev-parse", "--short", "HEAD")
		gS.branch += ":" + strings.TrimSpace(hash[0:len(hash)-1])
	} else if strings.Contains(line, "...") {
		gS.hasremote = true

		splitted := strings.Split(line, " ")
		gS.branch = strings.Split(splitted[1], "...")[0]

		if len(splitted) >= 3 {
			joined := strings.Join(splitted[2:], " ")

			gS.ahead = parsePattern(`ahead (\d+)`, joined)
			gS.behind = parsePattern(`behind (\d+)`, joined)
		}
	} else {
		gS.branch = strings.TrimSpace(strings.Split(line, " ")[1])
	}
}

func (gS *gitStatus) collectChanges(lines []string) {
	defer gS.waitGroup.Done()

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		idxStatus := line[0]
		wtStatus := line[1]

		// " M  hoge.txt" , "AM ahoo.png" , ...
		if wtStatus != ' ' && wtStatus != '?' {
			gS.changed++
		}

		// "MT hoge.cpp" , "A  fuga.txt" , ...
		if idxStatus != ' ' && idxStatus != '?' {
			gS.staged++
		}

		// "?? hoge.txt", ...
		if idxStatus == '?' && wtStatus == '?' {
			gS.untracked++
		}

		// "UU hogehoge.txt" ...
		if idxStatus == 'U' && wtStatus == 'U' {
			gS.conflicts++
		}
	}
}

func (gS *gitStatus) getTag() {
	defer gS.waitGroup.Done()
	tag, _, _ := communicate("git", "describe", "--tags", "--exact")
	if tag != "" {
		gS.tag = strings.TrimSpace(tag[0 : len(tag)-1])
	}
}
