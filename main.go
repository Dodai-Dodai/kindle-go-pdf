package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	pageNumber := 1

	// imgディレクトリが存在しない場合は作成
	if _, err := os.Stat("img"); os.IsNotExist(err) {
		os.Mkdir("img", 0777)
	}

	cmd := exec.Command("adb", "shell", "wm", "size")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing adb command:", err)
		return
	}

	screenResolution := string(output)
	screenResolution = strings.Replace(screenResolution, "Physical size: ", "", 1)
	resolutionParts := strings.Split(screenResolution, "x")
	if len(resolutionParts) != 2 {
		fmt.Println("Unexpected screen resolution format:", screenResolution)
		return
	}

	screenWidth, err := strconv.Atoi(strings.TrimSpace(resolutionParts[0]))
	if err != nil {
		fmt.Println("Error converting screen width:", err)
		return
	}

	screenHeight, err := strconv.Atoi(strings.TrimSpace(resolutionParts[1]))
	if err != nil {
		fmt.Println("Error converting screen height:", err)
		return
	}

	tapX := 0
	if len(os.Args) > 1 && os.Args[1] == "1" {
		tapX = screenWidth - (screenWidth - 10)
	} else {
		tapX = screenWidth - 10
	}

	tapY := screenHeight / 2

	for {
		formattedPageNumber := fmt.Sprintf("%04d", pageNumber)
		screenshotCmd := exec.Command("adb", "exec-out", "screencap", "-p")
		screenshotFile, err := os.Create(fmt.Sprintf("img/%s.png", formattedPageNumber))
		if err != nil {
			fmt.Println("Error creating screenshot file:", err)
			return
		}
		defer screenshotFile.Close()

		screenshotCmd.Stdout = screenshotFile
		if err := screenshotCmd.Run(); err != nil {
			fmt.Println("Error taking screenshot:", err)
			return
		}

		tapCmd := exec.Command("adb", "shell", "input", "touchscreen", "tap", strconv.Itoa(tapX), strconv.Itoa(tapY))
		if err := tapCmd.Run(); err != nil {
			fmt.Println("Error executing tap command:", err)
			return
		}

		time.Sleep(300 * time.Millisecond)
		pageNumber++
	}
}
