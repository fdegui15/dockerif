// ToolAntivirus for "inspectFile"
// base on ClamAV
// http://www.clamav.net/

package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func init() {
	//Integration of the new tool
	tools["av"] = toolsStr{
		Name: "antivirus",
		Cmd:  "clamdscan", //use clamdscan instead of clamscan to speed up the process.
		Args: []string{"--no-summary", "--stdout"},
		Vers: []string{"--version"},
		fn:   launchAntiVirus}
	//Clamdlaunhed()
	//gettoolsVersion("av")
}

func launchAntiVirus(fbyte []byte) []byte {
	toolFlag := "av"

	Clamdlaunhed()

	//Struct of Output
	type AVOut struct {
		Passed bool
		Error  string
	}
	var trait AVOut

	f := make(map[string]interface{})
	json.Unmarshal(fbyte, &f)
	fn := fmt.Sprintf("%s", f["FileName"])

	//str, err := execCmd(tools[toolName].Cmd, tools[toolName].Args, fn)
	str, err := exectoolsCmd(toolFlag, fn)
	if err != nil {
		trait.Error = str
		trait.Passed = false
	} else {
		//traitement.AVError = ""
		if str == fn+": OK" {
			trait.Passed = true
		} else {
			trait.Passed = false
			trait.Error = str
		}
	}
	f[tools[toolFlag].Name] = trait
	//settoolsVersion(toolFlag)
	output, _ := json.MarshalIndent(f, "", "    ")
	return output
}

func Clamdlaunhed() {
	// Launch 'clamd' if is not already daemonised
	if _, err := os.Stat("/tmp/clamd.ctl"); os.IsNotExist(err) {
		// if /tmp/clamd.ctl does not exist
		// then clamd is not launched !!!
		toolCmd := "clamd"
		toolArgs := []string{" "}
		fmt.Println(msgWithDate("clamd launching"))
		str, _ := execCmd(toolCmd, toolArgs, "")
		fmt.Println(msgWithDate("clamd launched : " + str))
	}
}
