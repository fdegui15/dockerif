// ToolExiftool for "inspectFile"
// Use Exiftool
// http://www.sno.phy.queensu.ca/~phil/exiftool/

package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func init() {
	//Integration of the new tool
	tools["et"] = toolsStr{
		Cmd:  "exiftool",
		Name: "exiftool",
		Args: []string{"-j"},
		Vers: []string{"-ver"},
		fn:   launchExiftool}
}

func launchExiftool(fbyte []byte) []byte {
	toolFlag := "et"

	f := make(map[string]interface{})
	json.Unmarshal(fbyte, &f)
	fn := fmt.Sprintf("%s", f["FileName"])

	sfstring, _ := exectoolsCmd(toolFlag, fn)
	//Trim the bracket [] Only one file !!!
	sfstring = strings.Trim(sfstring, " []")
	sftbyte := []byte(sfstring)

	//sfm will contain the struct of the output of Siegfried
	sfm := make(map[string]interface{})
	json.Unmarshal(sftbyte, &sfm)

	f[tools[toolFlag].Name] = sfm

	//settoolsVersion(toolFlag)

	output, _ := json.MarshalIndent(f, "", "    ")
	return output
}
