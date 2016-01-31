// ToolHash for "inspectFile"
// Use only go libraies

// You can change the size of the buffer in this file. See the variable "bufferSize"

package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

//Size of buffer to compute hashes
const bufferSize = 500000

func init() {
	//Integration of the new tool
	tools["hash"] = toolsStr{
		Name:   "hash",
		Cmd:    "", // Using go library
		Args:   []string{""},
		Vers:   []string{"nocmd"},
		VerTxt: "go library : crypto/md5 & crypto/sha512",
		fn:     launchHash}
}

func launchHash(fbyte []byte) []byte {

	toolName := "hash"

	//Struct FileHash
	type FileHashStr struct {
		Md5    string
		Sha512 string
	}
	var fhash FileHashStr

	f := make(map[string]interface{})
	json.Unmarshal(fbyte, &f)
	fn := fmt.Sprintf("%s", f["FileName"])

	str5, str512, _ := md5sha512sum(fn)
	fhash.Sha512 = str5
	fhash.Md5 = str512

	f[toolName] = fhash

	if ToolsVersion[toolName] == nil {
		AddNewToolsVersion(toolName, "go library : crypto/md5 & crypto/sha512")
	}
	output, _ := json.MarshalIndent(f, "", "    ")
	return output
}

// md5sha512sum returns MD5 and sha512 checksum of filename
func md5sha512sum(filename string) (string, string, error) {
	if info, err := os.Stat(filename); err != nil {
		return "", "", err
	} else if info.IsDir() {
		return "", "", nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	hash5 := md5.New()
	hash512 := sha512.New()
	for buf, reader := make([]byte, bufferSize), bufio.NewReader(file); ; {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", "", err
		}

		hash5.Write(buf[:n])
		hash512.Write(buf[:n])
	}

	checksum5 := fmt.Sprintf("%x", hash5.Sum(nil))
	checksum512 := fmt.Sprintf("%x", hash512.Sum(nil))
	fmt.Println(msgWithDate("hash computed from " + filename + " = " + checksum5))
	return checksum5, checksum512, nil
}
