package prompt

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"regexp"
	"strconv"
)

// Communicate ...
func communicate(name string, arg ...string) (string, string, error) {
	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command(name, arg...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err := cmd.Run()
	if err != nil {
		return "", "", err
	}
	stdout := outbuf.String()
	stderr := errbuf.String()

	return stdout, stderr, nil
}

// GetLines ...
func getLines(name string, arg ...string) ([]string, error) {
	var lines []string
	stdout, stderr, err := communicate(name, arg...)
	if err != nil {
		return []string{}, err
	} else if strings.Contains(stderr, "fatal") {
		return []string{}, errors.New(stderr)
	}
	if len(stdout) > 0 {
		lines = strings.Split(stdout[0:len(stdout)-1], "\n")
	}

	return lines, nil
}

func parsePattern(pattern string, s string) int {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(s)
	if len(matches) > 1 {
		ahead, err := strconv.Atoi(matches[1])
		if err == nil {
			return ahead
		}
	}
	return 0
}
