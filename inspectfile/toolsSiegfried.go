package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func init() {
	//Integration of the new tool
	tools["sf"] = toolsStr{
		Cmd:  "sf",
		Name: "siegfried",
		Args: []string{"-json"},
		Vers: []string{"-version"},
		fn:   launchSiegfried}
}

func launchSiegfried(fbyte []byte) []byte {
	toolFlag := "sf"

	f := make(map[string]interface{})
	json.Unmarshal(fbyte, &f)
	fn := fmt.Sprintf("%s", f["FileName"])

	sfstring, _ := exectoolsCmd(toolFlag, fn)

	sftbyte := []byte(sfstring)

	//sfm will contain the struct of the output of Siegfried
	sfm := make(map[string]interface{})
	json.Unmarshal(sftbyte, &sfm)

	//sfmf will contain the sfm["files"]
	sfmf := make(map[string]interface{})
	sftbyte, _ = json.MarshalIndent(sfm["files"], "", "    ")
	sfstring = string(sftbyte)               //output can contain multiple files !!
	sfstring = strings.Trim(sfstring, " []") // Trim to remove space, first and last bracket

	//ERROR if contains different files.
	sftbyte = []byte(sfstring)
	//fmt.Printf("%s\n\n", sftbyte)
	json.Unmarshal(sftbyte, &sfmf)

	//Export only the matches section. It can contain 0, 1 or n formats...
	//f["Siegfried"] = sfmf["matches"]
	//matches struct
	var matchesstr []map[string]interface{}
	sftbyte, _ = json.Marshal(sfmf["matches"])
	json.Unmarshal(sftbyte, &matchesstr)
	//matchesstr := []interface{}(sfmf["matches"])
	matchesstr[0]["NbMatches"] = len(matchesstr)
	if len(matchesstr) > 1 {
		//In case there is an other format !
		//We need to try it ;=)
		matchesstr[0]["OtherMatches"] = matchesstr[1:len(matchesstr)]
	}

	f[tools[toolFlag].Name] = matchesstr[0]

	//settoolsVersion(toolFlag)

	output, _ := json.MarshalIndent(f, "", "    ")
	return output
}
