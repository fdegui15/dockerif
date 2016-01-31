// ToolMediainfo for "inspectFile"
// use Mediainfo
// https://mediaarea.net/fr/MediaInfo

package main

import (
	"encoding/json"
	"fmt"
)

func init() {
	//Integration of the new tool
	tools["mi"] = toolsStr{
		Cmd:  "mediainfo",
		Name: "mediainfo",
		Args: []string{"--Inform=file:///home/go/src/inspectfile/MediaInfoTemplate.txt"},
		Vers: []string{"--version"},
		fn:   launchMediaInfo}
}

func launchMediaInfo(fbyte []byte) []byte {
	toolFlag := "mi"

	f := make(map[string]interface{})
	json.Unmarshal(fbyte, &f)
	fn := fmt.Sprintf("%s", f["FileName"])

	sfstring, _ := exectoolsCmd(toolFlag, fn)
	sftbyte := []byte(sfstring)

	sfm := make(map[string]interface{})
	json.Unmarshal(sftbyte, &sfm)

	//Define this ouput here
	f[tools[toolFlag].Name] = sfm

	//settoolsVersion(toolFlag)

	output, _ := json.MarshalIndent(f, "", "    ")
	return output
}
