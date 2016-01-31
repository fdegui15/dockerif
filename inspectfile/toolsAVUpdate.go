// ToolAVUpdate for "inspectFile"
// Use freshclam for the update and clamd to have the date of the signature file
// http://www.clamav.net/

package main

import (
	"encoding/json"
	"fmt"
)

var AVUpdateLaunched string

func launchAVUpdate() []byte {
	//launch an update of ClamAV

	toolName := "AV update" // => THis is not a tool to inspect a file.
	toolCmd := "freshclam"
	toolArgs := []string{" "}
	//toolVers := []string{""}

	//Struct of FileAVStr
	type AVOut struct {
		Passed bool
		Msg    string
	}
	var trait AVOut

	f := make(map[string]interface{})

	if AVUpdateLaunched == "true" {
		trait.Passed = false
		trait.Msg = "An AntiVirus Update is pending. Be patient..."
		fmt.Println(msgWithDate(trait.Msg))
		f[toolName] = trait
		output, _ := json.MarshalIndent(f, "", "    ")
		return output
	}
	AVUpdateLaunched = "true"

	fmt.Println(msgWithDate(toolName + " is launched..."))

	str, err := execCmd(toolCmd, toolArgs, "")
	if err != nil {
		trait.Msg = "Error in AntiVirus update " + str
		fmt.Println(msgWithDate("Error in AntiVirus update " + err.Error()))
		trait.Passed = false
		f[toolName] = trait
		output, _ := json.MarshalIndent(f, "", "    ")
		AVUpdateLaunched = ""
		return output
	} else {
		trait.Passed = true
		trait.Msg = str

		fmt.Println("AV update passed = Result from freschclam ")
		fmt.Println("===========================================")
		fmt.Println(str)
		fmt.Println("===========================================")

		// If the clamd daemon is not launch, we must launch clamdscan
		Clamdlaunhed()

		AVUpdateLaunched = ""
		toolCmd := "clamdscan"
		toolArgs := []string{"--version"}
		str, err := execCmd(toolCmd, toolArgs, "")
		if err != nil {
			trait.Msg = "Error in clamdscan " + str
			trait.Passed = false
			f[toolName] = trait
			output, _ := json.MarshalIndent(f, "", "    ")
			return output
		} else {
			trait.Msg = str
			trait.Passed = true
			f[toolName] = trait
			output, _ := json.MarshalIndent(f, "", "    ")
			return output
		}
	}
}
