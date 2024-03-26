package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	dac "github.com/Snawoot/go-http-digest-auth-client"
)

func boolToStr(b bool) string {
	if b {
		return "YES"
	}
	return "NO"
}

func ptr(s string) *string {
	return &s
}

func getLog(url string, pwd string) {
	client := &http.Client{
		Transport: dac.NewDigestTransport("root", pwd, http.DefaultTransport),
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Println(err)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	displayLastLines(string(body), 70)
}

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func displayLastLines(logContent string, nline int) {
	fmt.Println("Update..")
	time.Sleep(time.Millisecond * 500)
	clearScreen()
	lines := strings.Split(logContent, "\n")
	startLine := 0
	if len(lines) > nline {
		startLine = len(lines) - nline
	}
	for _, line := range lines[startLine:] {
		fmt.Println(line)
	}
}
