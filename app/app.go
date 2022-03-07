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
		panic(err.Error())
	}

	if locked {
		fmt.Printf("[# INFO #] Sync started\n\n")
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
					fmt.Printf("[# INFO #] NICM is up to date!\n'n")
					utils.StartNICM()
				}
			}
		}

		err := fileLock.Unlock()
		if err != nil {
			return
		}
	} else {
		fmt.Printf("[# INFO #] Please wait! Sync already running!\n\n")
	}
}
