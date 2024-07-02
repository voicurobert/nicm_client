package main

import (
	"os"
	"path"
)

var (
	ClientRootDir, _ = os.Getwd()

	RepoRootDir       = "\\\\dakar.office.orange.intra\\NICM_Repository\\nicm_sync_repos2"
	DefaultConfigPath = path.Join(RepoRootDir, "sync_nicm_default.config")
	CustomConfigPath  = path.Join(RepoRootDir, "sync_nicm.config")
	VersionFilePath   = path.Join(ClientRootDir, "version")

	ArchiveExtension = ".zip"

	// NicmPathToBat nicm test
	// NicmPathToBat = "nicm\\run5\\nicm_test\\"

	// NicmPathToBat nicm disaster recovery
	// NicmPathToBat = "nicm\\run5\\nicm_dr\\"

	// NicmPathToBat prod path
	NicmPathToBat = "nicm\\run5\\nicm\\"

	NicmBatName = "start_nicm_client.bat"

	LockFileName = "sync_nicm_local_env.lock"

	NicmPathToFile = "nicm\\run5\\"
	NicmFileName   = "nicm_"

	LockFullPath = path.Join(ClientRootDir, LockFileName)
)
