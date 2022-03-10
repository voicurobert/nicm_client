package main

import (
	"fmt"
	"nicm_client/app"
	"nicm_client/app/consts"
	"time"
)

func main() {
	fmt.Println(fmt.Sprintf("%s\\%s%s%d", consts.ClientRootDir, consts.NicmPathToFile, consts.NicmFileName, time.Now().UnixNano()))
	app.StartApplication()
}
