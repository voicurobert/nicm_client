package main

import (
	"bufio"
	"github.com/fatih/color"
	"github.com/gofrs/flock"
	"os"
	"strconv"
	"strings"
	"time"
)

func StartApplication() {
	if checkWorkingDirectory() {
		waitUserInput()
	}

	clientVersion := getClientVersion()
	customConfig := GetConfigForName(CustomConfigPath)
	customVersion := getCustomVersion(customConfig)

	var locked bool
	var fileLock *flock.Flock
	var lockErr error

	if clientVersion != customVersion {
		locked, fileLock, lockErr = LockFile()
		FileLock = fileLock
	}

	if lockErr != nil {
		color.Magenta("[# INFO #] Please wait! Sync already running!\n\n")
	} else {
		if locked {
			StopExecution = false
			startSync()
		} else {
			color.Magenta("[# INFO #] Please wait! Sync already running!\n\n")
		}
	}

	waitUserInput()
}

func getClientVersion() int {
	clientVersionStr := GetVersion()
	clientVersion, _ := strconv.Atoi(clientVersionStr)
	return clientVersion
}

func getCustomVersion(cfg ConfigMap) int {
	customVersionStr := cfg["BASE"]["version"]
	customVersion, _ := strconv.Atoi(strings.TrimSpace(customVersionStr))
	return customVersion
}

func startSync() {
	color.Magenta("[# INFO #] Sync started\n\n")

	clientVersion := getClientVersion()

	if clientVersion == 0 {
		defaultConfig := GetConfigForName(DefaultConfigPath)
		SyncArchives(defaultConfig)
		UpdateVersion(defaultConfig["BASE"]["version"])
		startTimer()
		StartNICM()
	} else {
		customConfig := GetConfigForName(CustomConfigPath)
		if len(customConfig) == 0 {
			startTimer()
			StartNICM()
		} else {
			customVersion := getCustomVersion(customConfig)

			if clientVersion != customVersion {
				SyncWithRepo(customConfig)
				startTimer()
				StartNICM()
			} else {
				color.Magenta("[# INFO #] NICM is up to date!\n\n")
				startTimer()
				StartNICM()
			}
		}
	}
}

func checkWorkingDirectory() bool {
	wd, _ := os.Getwd()
	if strings.HasPrefix(wd, "\\") {
		color.Red("Cannot execute NICM Client from this path!\n")
		return true
	}
	return false
}

func waitUserInput() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
	}
}

func startTimer() {
	timer := time.NewTimer(2 * time.Minute)
	go func() {
		<-timer.C
		StopExecution = true
		return
	}()
}
