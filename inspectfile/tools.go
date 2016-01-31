// Tools module from "inspectFile"
// This file contains all the common structures, variables and functions
// used by the different tools implemented

package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

//list all the tools to use
var listtools []string

type toolsStr struct {
	Cmd    string              //Cmd line to launch the tool
	Name   string              //Name of the tool
	Args   []string            //Arguments used with the Cmd Line
	Vers   []string            //Arguments used to return the Version Txt
	VerTxt string              //The Version Text returned by the tool
	fn     func([]byte) []byte //function to call to use the tool
}

//mapping Flag => Tool
var tools = make(map[string]toolsStr)

//ToolsVersion will contain the versions of all the tools used
var ToolsVersion = make(map[string]interface{})

func settoolsVersion(toolFlag string) {
	if ToolsVersion[toolFlag] == nil {
		var outVers string
		if tools[toolFlag].Cmd != "" {
			var err error
			outVers, err = execCmd(tools[toolFlag].Cmd, tools[toolFlag].Vers, "")
			if err != nil {
				fmt.Println(msgWithDate(tools[toolFlag].Name + " executed with ERROR : " + err.Error()))
			} else {
				fmt.Println(msgWithDate(tools[toolFlag].Name + " executed without issue"))
			}
		} else {
			//No Command line
			outVers = tools[toolFlag].VerTxt
		}
		AddNewToolsVersion(tools[toolFlag].Name, outVers)
		//fmt.Sprint(msgWithDate(outVers))
		fmt.Println(msgWithDate(tools[toolFlag].Name + " version is " + outVers))
	}
}

//Launch a command with args and the filename (the filename is optional for the version)
func exectoolsCmd(toolname string, filename string) (string, error) {
	fmt.Println(msgWithDate(toolname + " is launched..."))

	//cmdargs := append(tools[toolname].Args, filename)
	//cmdexec := exec.Command(tools[toolname].Cmd, cmdargs...)
	//output, err := cmdexec.CombinedOutput()
	output, err := execCmd(tools[toolname].Cmd, tools[toolname].Args, filename)
	if err != nil {
		fmt.Println(msgWithDate(toolname + " executed with ERROR : " + err.Error()))
	} else {
		fmt.Println(msgWithDate(toolname + " executed without issue."))
	}
	return output, err
}

//Launch a command with args and the filename (the filename is optional for the version)
func execCmd(cmd string, args []string, filename string) (string, error) {
	var cmdargs []string
	if filename != "" {
		cmdargs = append(args, filename)
	} else {
		cmdargs = args
	}
	cmdexec := exec.Command(cmd, cmdargs...)
	output, err := cmdexec.CombinedOutput()
	return strings.Trim(string(output), " \n"), err
}

func msgWithDate(msg string) string {
	//Display the msg with the date in front (Logs type)
	t := time.Now()
	return t.Local().String() + ": " + msg
} //

func AddNewToolsVersion(name string, ntv string) {
	ToolsVersion[name] = ntv
}

func initToolsVersion() {
	for key, _ := range ToolsVersion {
		delete(ToolsVersion, key)
	}
}

func ExportToolsVersion() []byte {
	fmt.Println(msgWithDate("ExportToolsVersion launched"))
	initToolsVersion()
	for toolFlag := range listtools {
		settoolsVersion(listtools[toolFlag])
	}
	fmt.Println(msgWithDate("ExportToolsVersion done"))

	out, _ := json.MarshalIndent(ToolsVersion, "", "    ")
	return out
}
