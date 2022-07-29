package app

import (
	"bufio"
	"github.com/fatih/color"
	"github.com/gofrs/flock"
	"nicm_client/app/consts"
	"nicm_client/app/utils"
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
	customConfig := utils.GetConfigForName(consts.CustomConfigPath)
	customVersion := getCustomVersion(customConfig)

	var locked bool
	var fileLock *flock.Flock

	if clientVersion != customVersion {
		locked, fileLock = utils.LockFile()
		utils.FileLock = fileLock
	}

	if locked {
		utils.StopExecution = false
		startSync()
	} else {
		color.Magenta("[# INFO #] Please wait! Sync already running!\n\n")
	}
	waitUserInput()
}

func getClientVersion() int {
	clientVersionStr := utils.GetVersion()
	clientVersion, _ := strconv.Atoi(clientVersionStr)
	return clientVersion
}

func getCustomVersion(cfg utils.ConfigMap) int {
	customVersionStr := cfg["BASE"]["version"]
	customVersion, _ := strconv.Atoi(strings.TrimSpace(customVersionStr))
	return customVersion
}

func startSync() {
	color.Magenta("[# INFO #] Sync started\n\n")

	clientVersion := getClientVersion()

	if clientVersion == 0 {
		defaultConfig := utils.GetConfigForName(consts.DefaultConfigPath)
		utils.SyncArchives(defaultConfig)
		utils.UpdateVersion(defaultConfig["BASE"]["version"])
		startTimer()
		utils.StartNICM()
	} else {
		customConfig := utils.GetConfigForName(consts.CustomConfigPath)
		if len(customConfig) == 0 {
			startTimer()
			utils.StartNICM()
		} else {
			customVersion := getCustomVersion(customConfig)

			if clientVersion != customVersion {
				utils.SyncWithRepo(customConfig)
				startTimer()
				utils.StartNICM()
			} else {
				color.Magenta("[# INFO #] NICM is up to date!\n\n")
				startTimer()
				utils.StartNICM()
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
		utils.StopExecution = true
		return
	}()
}
