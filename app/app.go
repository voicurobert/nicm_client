package app

import (
	"github.com/fatih/color"
	"github.com/gofrs/flock"
	"nicm_client/app/consts"
	"nicm_client/app/utils"
	"path"
	"strconv"
	"strings"
)

func StartApplication() {
	fileLock := flock.New(path.Join(consts.ClientRootDir, consts.LockFileName))
	locked, err := fileLock.TryLock()

	if err != nil {
		panic(err.Error())
	}

	if locked {
		color.Magenta("[# INFO #] Sync started\n\n")
		clientVersionStr := utils.GetVersion()
		clientVersion, _ := strconv.Atoi(clientVersionStr)

		if clientVersion == 0 {
			defaultConfig := utils.GetConfigForName(consts.DefaultConfigPath)
			utils.SyncArchives(defaultConfig)
			utils.UpdateVersion(defaultConfig["BASE"]["version"])
			utils.StartNICM()
		} else {

			customConfig := utils.GetConfigForName(consts.CustomConfigPath)

			if len(customConfig) == 0 {
				utils.StartNICM()
			} else {
				customVersionStr := customConfig["BASE"]["version"]
				customVersion, _ := strconv.Atoi(strings.TrimSpace(customVersionStr))

				if clientVersion != customVersion {
					utils.SyncWithRepo(customConfig)
					utils.StartNICM()
				} else {
					color.Magenta("[# INFO #] NICM is up to date!\n\n")
					utils.StartNICM()
				}
			}
		}

		err := fileLock.Unlock()
		if err != nil {
			return
		}
	} else {
		color.Magenta("[# INFO #] Please wait! Sync already running!\n\n")
	}
}
