package app

import (
	"fmt"
	"github.com/gofrs/flock"
	"nicm_client/app/consts"
	"nicm_client/app/utils"
	"path"
)

func StartApplication() {
	fileLock := flock.New(path.Join(consts.ClientRootDir, consts.LockFileName))
	locked, err := fileLock.TryLock()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if locked {
		fmt.Println("[# INFO #] Sync started")
		fmt.Println()
		clientVersion := utils.GetVersion()

		if clientVersion == "0" {
			defaultConfig := utils.GetConfigForName(consts.DefaultConfigPath)
			utils.SyncArchives(defaultConfig)
			utils.UpdateVersion(defaultConfig["BASE"]["version"])
			utils.SyncWithRepo(utils.GetConfigForName(consts.CustomConfigPath))
			utils.StartNICM()
		} else {
			customConfig := utils.GetConfigForName(consts.CustomConfigPath)

			if len(customConfig) == 0 {
				utils.StartNICM()
			} else {
				customVersion := customConfig["BASE"]["version"]
				if clientVersion != customVersion {
					utils.SyncWithRepo(customConfig)
					utils.StartNICM()
				} else {
					fmt.Println("[# INFO #] NICM is up to date!")
					fmt.Println()
					utils.StartNICM()
				}
			}
		}

		err := fileLock.Unlock()

		if err != nil {
			return
		}

	} else {
		fmt.Println("[# INFO #] Please wait! Sync already running!")
		fmt.Println()
	}
}
