// ToolFido for "inspectFile"
// Use FIDO
// http://openpreservation.org/technology/products/fido/

package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func init() {
	//Integration of the new tool
	tools["fido"] = toolsStr{
		Cmd:  "/home/fido/fido.py",
		Name: "fido",
		Args: []string{"-q", "-matchprintf", "OK§%(info.time)s§%(info.puid)s§%(info.formatname)s§%(info.signaturename)s§%(info.mimetype)s§%(info.matchtype)s\n", "-nomatchprintf", "KO§§§§§§%(info.matchtype)s"},
		Vers: []string{"-v"},
		fn:   launchFido}
}

func launchFido(fbyte []byte) []byte {
	toolFlag := "fido"

	f := make(map[string]interface{})
	json.Unmarshal(fbyte, &f)
	fn := fmt.Sprintf("%s", f["FileName"])

	//Struct of Output
	type FidoOut struct {
		Passed        bool   `json:"passed"`
		Puid          string `json:"puid,omitempty"`
		Formatname    string `json:"formatname,omitempty"`
		Signaturename string `json:"signaturename,omitempty"`
		Mimetype      string `json:"mimetype,omitempty"`
		Matchtype     string `json:"matchtype,omitempty"`
	}
	var fo FidoOut

	sfstring, _ := exectoolsCmd(toolFlag, fn)

	sft := strings.Split(sfstring, "§")

	fo.Passed = (sft[0] == "OK")
	fo.Puid = sft[2]
	fo.Formatname = sft[3]
	fo.Signaturename = sft[4]
	fo.Mimetype = sft[5]
	fo.Matchtype = sft[6]

	f[tools[toolFlag].Name] = fo

	//settoolsVersion(toolFlag)

	output, _ := json.MarshalIndent(f, "", "    ")
	return output
}
