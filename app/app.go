package app

import (
	"bufio"
	"github.com/fatih/color"
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

	locked, fileLock := utils.LockFile()
	utils.FileLock = fileLock

	if locked {
		utils.StopExecution = false
		startSync()
	} else {
		color.Magenta("[# INFO #] Please wait! Sync already running!\n\n")
	}
	waitUserInput()
}

func startSync() {
	color.Magenta("[# INFO #] Sync started\n\n")
	clientVersionStr := utils.GetVersion()
	clientVersion, _ := strconv.Atoi(clientVersionStr)

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
			customVersionStr := customConfig["BASE"]["version"]
			customVersion, _ := strconv.Atoi(strings.TrimSpace(customVersionStr))

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
