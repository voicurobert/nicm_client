package utils

import (
	"archive/zip"
	"fmt"
	"github.com/bigkevmcd/go-configparser"
	"github.com/fatih/color"
	"github.com/gofrs/flock"
	"io"
	"log"
	"nicm_client/app/consts"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var FileLock *flock.Flock
var StopExecution bool

func GetVersion() string {
	nr, err := os.ReadFile(consts.VersionFilePath)

	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(string(nr))
}

type ConfigMap map[string]map[string]string

func GetConfigForName(configFilePath string) ConfigMap {
	config, err := configparser.NewConfigParserFromFile(configFilePath)
	if err != nil {
		panic(fmt.Sprintf("cannot read config file %s, error: %s", configFilePath, err.Error()))
	}

	configMap := make(ConfigMap)

	for _, section := range config.Sections() {
		keyValue, _ := config.Items(section)
		configMap[section] = keyValue
	}

	return configMap
}

func UpdateVersion(newVersion string) {
	err := os.WriteFile(consts.VersionFilePath, []byte(newVersion+""), 0644)
	if err != nil {
		panic(err)
	}
}

var wg sync.WaitGroup

func SyncArchives(config ConfigMap) {
	archives, exists := config["ARCHIVES"]
	if exists == false {
		return
	}
	for archiveName, value := range archives {
		if value == "0" {
			continue
		}
		wg.Add(1)
		go syncArchive(archiveName)
	}
	wg.Wait()
}

func syncArchive(archiveName string) {
	ext := getMainDir(archiveName)

	clientArchivePath := consts.ClientRootDir + "\\" + archiveName + consts.ArchiveExtension

	clientPath, fullClientPath, sourceArchivePath := getPaths(ext, archiveName)

	color.Yellow("Copying %s, please wait...\n", archiveName)
	copyCommand(sourceArchivePath, consts.ClientRootDir)

	if _, err := os.Stat(fullClientPath); err == nil {
		_ = os.RemoveAll(fullClientPath)
	}
	_ = os.MkdirAll(fullClientPath, os.ModePerm)

	unzip(clientArchivePath, clientPath, archiveName)
	_ = os.Remove(clientArchivePath)
	//color.Yellow("removed archive: %s \n", archiveName)

	time.Sleep(500 * time.Millisecond)
	wg.Done()
}

func getMainDir(archiveName string) string {
	ext := ""
	if strings.Contains(archiveName, "nicm") || strings.Contains(archiveName, "run") {
		ext = ""
	} else {
		ext = "externals"
	}
	return ext
}

func getPaths(mainDir, archiveName string) (string, string, string) {
	clientPath := ""
	fullClientPath := ""
	sourceArchivePath := ""
	if mainDir != "" {
		fullClientPath = consts.ClientRootDir + "\\" + mainDir + "\\" + archiveName
		clientPath = consts.ClientRootDir + "\\" + mainDir + "\\"
		sourceArchivePath = consts.RepoRootDir + "\\" + mainDir + "\\" + archiveName + consts.ArchiveExtension
	} else {
		fullClientPath = consts.ClientRootDir + "\\" + archiveName
		clientPath = consts.ClientRootDir + "\\"
		sourceArchivePath = consts.RepoRootDir + "\\" + archiveName + consts.ArchiveExtension
	}

	return clientPath, fullClientPath, sourceArchivePath
}

func unzip(archivePath, path, zipName string) {
	color.Yellow("Unzipping %s, please wait...\n", zipName)
	archive, err := zip.OpenReader(archivePath)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(path, f.Name)

		if !strings.HasPrefix(filePath, filepath.Clean(path)+string(os.PathSeparator)) {
			color.Red("invalid file path")
			return
		}
		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(filePath, os.ModePerm)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}
		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}
		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}
		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		_ = dstFile.Close()
		_ = fileInArchive.Close()
	}
	color.Green("Done unzipping %s\n", zipName)
}

func copyCommand(src, dest string) {
	cmd := exec.Command("cmd.exe", "/C", "copy", src, dest)
	_ = cmd.Run()
}

func StartNICM() {
	fullPath := fmt.Sprintf(
		"%s\\%s%s%d%s",
		consts.ClientRootDir, consts.NicmPathToFile,
		consts.NicmFileName,
		time.Now().UnixNano(),
		".txt")

	_ = os.WriteFile(fullPath, nil, 0644)

	color.Green("Starting NICM Application... this can take a while, please wait. \nDO NOT CLOSE this window, it will close automatically!")
	startPath := consts.ClientRootDir + "\\" + consts.NicmPathToBat + consts.NicmBatName
	_ = os.Chdir(startPath)

	executeCommand("/C", startPath, fullPath)
	checkNicmFile(fullPath)
}

func checkNicmFile(filePath string) {
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for range ticker.C {
			if StopExecution {
				color.HiRed("NICM Application could not start...program will stop.")
				UnlockFile(FileLock)
				ticker.Stop()
				return
			}
			value, err := os.ReadFile(filePath)
			if err != nil {
				panic(err)
			}
			txt := string(value)
			if txt == "done" {
				_ = os.Remove(filePath)
				color.Green("Started NICM application!")
				UnlockFile(FileLock)
				ticker.Stop()
				return
			}
		}
	}()
}

func executeCommand(args ...string) {
	cmd := exec.Command("cmd.exe", args...)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

func SyncWithRepo(config ConfigMap) {
	SyncArchives(config)
	UpdateVersion(config["BASE"]["version"])
}

func LockFile() (bool, *flock.Flock) {
	fileLock := flock.New(path.Join(consts.ClientRootDir, consts.LockFileName))
	locked, err := fileLock.TryLock()
	if err != nil {
		panic(err.Error())
	}
	return locked, fileLock
}

func UnlockFile(fileLock *flock.Flock) {
	err := fileLock.Unlock()
	if err != nil {
		panic(err.Error())
	}

	err = os.Remove(fileLock.Path())
	if err != nil {
		fmt.Println(err.Error())
	}
}
