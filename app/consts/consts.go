package consts

import (
	"os"
	"path"
)

var (
	ClientRootDir, _ = os.Getwd()

	//RepoRootDir       = "\\dakar\\NICM_Repository\\nicm_sync_repos"
	//DefaultConfigPath = RepoRootDir + "sync_nicm_default.config"
	//CustomConfigPath  = RepoRootDir + "sync_nicm.config"
	//VersionFilePath   = RepoRootDir + "version"

	RepoRootDir       = "C:\\sw\\nicm\\nicm_529"
	DefaultConfigPath = path.Join(RepoRootDir, "sync_nicm_default.config")
	CustomConfigPath  = path.Join(RepoRootDir, "sync_nicm.config")
	VersionFilePath   = path.Join(ClientRootDir, "version")

	ArchiveExtension = ".zip"

	NicmPathToBat = "nicm\\run5\\nicm\\"
	NicmBatName   = "start_nicm_client.bat"

	LockFileName = "sync_nicm_local_env.lock"
)
