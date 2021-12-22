package utils

import (
	"archive/zip"
	"fmt"
	"github.com/bigkevmcd/go-configparser"
	"io"
	"nicm_client/app/consts"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetVersion() string {
	nr, err := os.ReadFile(consts.VersionFilePath)

	if err != nil {
		panic(err)
	}
	return string(nr)
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
	err := os.WriteFile(consts.VersionFilePath, []byte(newVersion+"\n"), 0644)
	if err != nil {
		panic(err)
	}
}

func SyncArchives(config ConfigMap) {
	archives, exists := config["ARCHIVES"]
	if exists == false {
		return
	}
	for archiveName, value := range archives {
		if value == "0" {
			continue
		}

		ext := ""
		if strings.Contains(archiveName, "nicm") || strings.Contains(archiveName, "run") {
			ext = ""
		} else {
			ext = "externals"
		}

		clientArchivePath := consts.ClientRootDir + "\\" + archiveName + consts.ArchiveExtension

		clientPath := ""
		fullClientPath := ""
		sourceArchivePath := ""
		if ext != "" {
			fullClientPath = consts.ClientRootDir + "\\" + ext + "\\" + archiveName
			clientPath = consts.ClientRootDir + "\\" + ext + "\\"
			sourceArchivePath = consts.RepoRootDir + "\\" + ext + "\\" + archiveName + consts.ArchiveExtension
		} else {
			fullClientPath = consts.ClientRootDir + "\\" + archiveName
			clientPath = consts.ClientRootDir + "\\"
			sourceArchivePath = consts.RepoRootDir + "\\" + archiveName + consts.ArchiveExtension
		}
		fmt.Printf("Copying %s, please wait...\n", archiveName)
		copyCommand(sourceArchivePath, consts.ClientRootDir)

		if _, err := os.Stat(fullClientPath); err == nil {
			_ = os.RemoveAll(fullClientPath)
		}
		_ = os.MkdirAll(fullClientPath, os.ModePerm)

		unzip(clientArchivePath, clientPath, archiveName)
		_ = os.Remove(clientArchivePath)
	}
}

func unzip(archivePath, path, zipName string) {
	fmt.Printf("Unzipping %s, please wait...\n", zipName)
	archive, err := zip.OpenReader(archivePath)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(path, f.Name)

		if !strings.HasPrefix(filePath, filepath.Clean(path)+string(os.PathSeparator)) {
			fmt.Println("invalid file path")
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
}

func copyCommand(src, dest string) {
	executeCommand("/C", "copy", src, dest)
}

func StartNICM() {
	fmt.Println("Starting NICM Application...")
	startPath := consts.ClientRootDir + "\\" + consts.NicmPathToBat + consts.NicmBatName
	_ = os.Chdir(startPath)
	executeCommand("/C", startPath)

}

func executeCommand(args ...string) {
	cmd := exec.Command("cmd.exe", args...)
	_ = cmd.Start()
	//cmdReader, _ := cmd.StdoutPipe()
	//scanner := bufio.NewScanner(cmdReader)
	//go func() {
	//	for scanner.Scan() {
	//		//fmt.Printf("\t > %s\n", scanner.Text())
	//	}
	//}()
	//_ = cmd.Start()
	//_ = cmd.Wait()
}

func SyncWithRepo(config ConfigMap) {
	SyncArchives(config)
	UpdateVersion(config["BASE"]["version"])
}
