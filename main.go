package main

import (
	"errors"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	msg, err := getMessages()
	if err != nil {
		println("Can't read commit messages or version tag is not present")
		return
	}
	msgs := strings.Split(msg, "\n")
	for i, m := range msgs {
		msgs[i] = strings.Trim(m, `"`)
	}
	println(strings.Join(msgs, "\n"))
}

func getMessages() (string, error) {
	currentVersion, err := getCurrentTagVersion()
	if err != nil {
		return "", err
	}
	var msg string
	for i := 0; i < 50; i++ {
		prevVersion := getPrevTagVersion(currentVersion)
		msg, err = getMessagesBetweenTags(prevVersion, currentVersion)
		if err != nil {
			continue
		}
		return msg, nil
	}
	return "", err
}

func getMessagesBetweenTags(tag1, tag2 string) (string, error) {
	cmd := exec.Command("git", `log`, `--pretty=format:"%s"`, tag1+".."+tag2)
	res, err := cmd.Output()
	return string(res), err
}

func getCurrentTagVersion() (string, error) {
	cmd := exec.Command("git", "log", "-n1", "--pretty='%h'")
	res, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	commit := strings.Trim(strings.TrimRight(string(res), "\n"), "'")
	cmd = exec.Command("git", "describe", "--exact-match", "--tags", commit)
	res, err = cmd.Output()
	rgxNum := regexp.MustCompile(`v\d+\.\d+\.\d+`)
	currentVersion := rgxNum.Find(res)
	if currentVersion == nil {
		return "", errors.New("no tag version")
	}
	return string(currentVersion), nil
}

func getPrevTagVersion(v string) string {
	segments := strings.Split(v, ".")
	patchVersion, _ := strconv.Atoi(segments[2])
	segments[2] = strconv.Itoa(patchVersion - 1)
	return strings.Join(segments, ".")
}
