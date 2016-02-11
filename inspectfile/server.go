// Server module for "inspectFile"

package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

//Constant MAXProcess can be implemented in Dockerfile !!
var MAXProcess int = 10

//NbActivProcess
var NbActivProcess int = 0

// Management of Process
type ProcUnit struct {
	Name      string
	FileName  string
	BeginTime time.Time
	EndTime   time.Time
	Activ     bool
}

//List of Processes
var ListProc []ProcUnit

func msgWithDateProc(NProc int, msg string) string {
	//Display the msg with the date in front (Logs type)
	t := time.Now()
	return t.Local().String() + " Proc#" + strconv.Itoa(NProc) + ": " + msg
}

func AddNewProc(procname, fn string) (int, error) {
	// Try to add a new process
	// if NbActivProcess < MAXProcess create a new Procunit AND return its id AND NbActivProcess ++
	// if NbAdtivProcess >= MAXProcess return ERROR
	if NbActivProcess < MAXProcess {
		/*
			if !(LastProcess > 0) {
				LastProcess = 0
			}
			if !(NbActivProcess > 0) {
				NbActivProcess = 0
			}
		*/
		NbActivProcess++
		var cpu ProcUnit
		cpu.Activ = true
		cpu.BeginTime = time.Now()
		cpu.Name = procname
		cpu.FileName = fn
		ListProc = append(ListProc, cpu)
		curProc := len(ListProc) - 1
		fmt.Println(msgWithDate("AddNewProcess: " + procname + " create process #" + strconv.Itoa(curProc) + ". NbActivProcess = " + strconv.Itoa(NbActivProcess)))
		return curProc, nil
	} else {
		fmt.Println(msgWithDate("AddNewProcess: Nb Maximum of Process reached : " + procname + "can't create new process!"))
		return -1, errors.New("Nb Maximum of Process reached!")
	}
}

func EndProc(id int) {
	if ListProc[id].Activ {
		ListProc[id].EndTime = time.Now()
		ListProc[id].Activ = false
		//ListProc[id].EndTime = time.Now()
		//ListProc[id].Activ = false
		NbActivProcess--
		fmt.Println(msgWithDateProc(id, "EndProc: The Process "+strconv.Itoa(id)+" is ended. NbActivProcess = "+strconv.Itoa(NbActivProcess)))
	} else {
		fmt.Println(msgWithDateProc(id, "EndProc: The Process "+strconv.Itoa(id)+" is already ended !"))
	}
}

func handleErr(w http.ResponseWriter, status int, e error) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, e.Error())
}

func httpinspect() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//Upload a file, Copy in the /tmp directory, Inspect the file and Remove the file
		mime := "application/json"
		NProc, err := AddNewProc("httpinspect", "Which one ?")
		if err != nil {
			fmt.Println(msgWithDate(err.Error()))
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		defer EndProc(NProc) //We must End the process at the end of the function!
		tools := r.URL.Query().Get("tools")
		fmt.Println("*****************************************")
		fmt.Println(msgWithDateProc(NProc, "httpinspect launched."))
		fmt.Println(msgWithDateProc(NProc, "URL String: "+r.URL.String()+" & tools"+tools))
		initTools(tools)
		initToolsVersion()

		var fn string //fn is the filename

		if r.Header.Get("Content-Type") == "application/octet-stream" {
			//used by filedrop.js
			//the size of the file is limited to 2 GB (not tested yet !)
			fn = r.Header.Get("X-File-Name") //getbootstrap the filename
			ListProc[NProc].FileName = fn
			fmt.Printf("Size of data to read = %d (%s)\n", r.ContentLength, r.Header.Get("X-File-Size"))
			data := make([]byte, r.ContentLength) //get the data upload in memory
			fmt.Println(msgWithDateProc(NProc, "Trying to read the data !!"))
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "ERROR in reading the data: "+err.Error(), http.StatusBadRequest)
				fmt.Println(msgWithDateProc(NProc, "ReadAll Error = "+err.Error()))
				//EndProc(NProc)
				return
			} else {
				fmt.Println(msgWithDateProc(NProc, "The data is correctly read"))
			}
			err = ioutil.WriteFile("/tmp/"+fn, data, 0644) //
			if err != nil {
				http.Error(w, "ERROR creating file in /tmp/: "+err.Error(), http.StatusBadRequest)
				fmt.Println(msgWithDateProc(NProc, "ERROR creating file /tmp/"+fn))
				//EndProc(NProc)
				return
			}
			output = inspectfile("/tmp/"+fn, nil)

		} else {
			//used by the curl command line
			// Content-Type = multipart which is standard!!
			infile, header, err := r.FormFile("file")
			fn = header.Filename
			if err != nil {
				http.Error(w, "ERROR parsing uploaded file: "+err.Error(), http.StatusBadRequest)
				fmt.Println(msgWithDateProc(NProc, "ERROR parsing uploaded file: "+err.Error()))
				//EndProc(NProc)
				return
			}
			outfile, err := os.Create("/tmp/" + fn)
			if err != nil {
				http.Error(w, "ERROR creating file: "+err.Error(), http.StatusBadRequest)
				fmt.Println(msgWithDateProc(NProc, "ERROR creating file: "+err.Error()))
				//EndProc(NProc)
				return
			}
			_, err = io.Copy(outfile, infile)
			if err != nil {
				http.Error(w, "ERROR saving file: "+err.Error(), http.StatusBadRequest)
				fmt.Println(msgWithDateProc(NProc, "ERROR saving file: "+err.Error()))
				//EndProc(NProc)
				return
			}
			output = inspectfile(outfile.Name(), nil)
		}

		w.Header().Set("Content-Type", mime)
		w.Write(output)

		err2 := os.Remove("/tmp/" + fn)
		if err2 != nil {
			http.Error(w, "ERROR removing file: "+err.Error(), http.StatusBadRequest)
			fmt.Println(msgWithDateProc(NProc, "ERROR removing file: "+err.Error()))
			//EndProc(NProc)
			return
		}
		fmt.Println(msgWithDateProc(NProc, "httpinspect done."))
		//EndProc(NProc)
		fmt.Println("*****************************************")
		return
	}
}

func httpinspectpath() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//inspect a path
		//THe path used the Environment variable %MOUNTDIR%
		mime := "application/json"
		NProc, err := AddNewProc("httpinspectpath", "")
		if err != nil {
			fmt.Println(msgWithDate(err.Error()))
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		defer EndProc(NProc) //We must End the process when leaving this function!
		tools := r.URL.Query().Get("tools")
		fmt.Println("*****************************************")
		fmt.Println(msgWithDateProc(NProc, "httpinspectpath launched."))
		fmt.Println(msgWithDateProc(NProc, "URL String: "+r.URL.String()+" & tools"+tools))
		initTools(tools)
		initToolsVersion()

		path := r.URL.Path
		fmt.Println(msgWithDateProc(NProc, "GET Path: "+path))
		if len(path) < 10 {
			http.Error(w, "You need to pass the filename or directory path after /inspect !!!", http.StatusBadRequest)
			fmt.Println(msgWithDateProc(NProc, "httpinspectpath: ERROR no directory path after /inspect"))
			//EndProc(NProc)
			return
		} else {
			path = os.Getenv("MOUNTDIR") + path[len("/inspectpath"):]
			ListProc[NProc].FileName = "Path= " + path
		}
		info, err := os.Stat(path)
		if err != nil {
			handleErr(w, http.StatusNotFound, err)
			fmt.Println(msgWithDateProc(NProc, "httpinspectpath: ERROR in getting inforamtion from path"+err.Error()))
			//EndProc(NProc)
			return
		}
		w.Header().Set("Content-Type", mime)
		if info.IsDir() {
			output = nil
			output = inspectdir(path)
		} else {
			output = inspectfile(path, nil)
		}

		w.Write(output)
		fmt.Println(msgWithDateProc(NProc, "httpinspect PATH done."))
		fmt.Println("*****************************************")
		//EndProc(NProc)
		return

	}
}

func httplocalinspect() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//var output []byte
		//output := []byte("AV update - TO DO...")
		fmt.Println(msgWithDate("msg"))
		fmt.Println("URL String: ", r.URL.String())
		tools := r.URL.Query().Get("tools")
		fmt.Println("Tools: ", tools)
		/*if tools == "" {
			tools = "hash-sf"
		}*/
		initTools(tools)
		//initToolsVersion()
		filename := os.Getenv("MOUNTDIR") + "/" + r.URL.Query().Get("file")
		fmt.Println("*****************************************")
		fmt.Println(msgWithDate("httplocalinspect : Inspect the Local file name:" + filename))

		w.Header().Set("Content-Type", "application/json")
		output := inspectfile(filename, nil)
		//output := []byte(filename + " - " + tools + "\n")
		//t := time.Now()
		fmt.Println(msgWithDate("httplocalinspect done"))
		fmt.Println("*****************************************")
		w.Write(output)
		return
	}
}

func httpavupdate() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//output := []byte("AV update - TO DO...")
		NProc, err := AddNewProc("httpavupdate", "")
		if err != nil {
			fmt.Println(msgWithDate(err.Error()))
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		defer EndProc(NProc) //We must End the process when leaving this function!
		fmt.Println("*****************************************")
		fmt.Println(msgWithDateProc(NProc, "httpAVUpdate launched."))
		w.Header().Set("Content-Type", "application/json")
		w.Write(launchAVUpdate())
		fmt.Println(msgWithDateProc(NProc, "httpAVUpdate done."))
		fmt.Println("*****************************************")
		//EndProc(NProc)
		return
	}
}

func httpcbavupdate() func(w http.ResponseWriter, r *http.Request) {
	//callback version of httpavupdate
	return func(w http.ResponseWriter, r *http.Request) {
		//output := []byte("AV update - TO DO...")
		NProc, err := AddNewProc("httpcbavupdate", "")
		if err != nil {
			fmt.Println(msgWithDate(err.Error()))
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		defer EndProc(NProc) //We must End the process when leaving this function!
		fmt.Println("*****************************************")
		fmt.Println(msgWithDateProc(NProc, "httpcbAVUpdate launched."))
		cbfunc := r.URL.Query().Get("callback")
		fmt.Println(msgWithDateProc(NProc, "httpcbAVUpdate = callback: "+cbfunc))
		w.Header().Set("Content-Type", "application/json")
		str := cbfunc + "(" + string(launchAVUpdate()) + ");"
		arr := []byte(str)
		w.Write(arr)
		//w.Write(launchAVUpdate())
		fmt.Println(msgWithDateProc(NProc, "httpcbAVUpdate done."))
		fmt.Println("*****************************************")
		//EndProc(NProc)
		return
	}
}

func httpgettoolsversion() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//output := []byte("Tools version - TO DO...")
		NProc, err := AddNewProc("httpgettoolsversion", "")
		if err != nil {
			fmt.Println(msgWithDate(err.Error()))
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		defer EndProc(NProc) //We must End the process when leaving this function!
		fmt.Println("*****************************************")
		fmt.Println(msgWithDateProc(NProc, "httpgettoolsversion launched"))
		tools := r.URL.Query().Get("tools")
		fmt.Println(msgWithDateProc(NProc, "httpgettoolsversion = tools:"+tools))
		initTools(tools)
		//initToolsVersion()
		w.Header().Set("Content-Type", "application/json")
		w.Write(ExportToolsVersion())
		fmt.Println(msgWithDateProc(NProc, "httpgettoolsversion done"))
		fmt.Println("*****************************************")
		//EndProc(NProc)
		return
	}
}

func httpcbgettoolsversion() func(w http.ResponseWriter, r *http.Request) {
	//callback version of httpgettoolsversion
	return func(w http.ResponseWriter, r *http.Request) {
		NProc, err := AddNewProc("httpcbgettoolsversion", "")
		if err != nil {
			fmt.Println(msgWithDate(err.Error()))
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		defer EndProc(NProc) //We must End the process when leaving this function!
		fmt.Println("*****************************************")
		fmt.Println(msgWithDateProc(NProc, "httpcbgettoolsversion launched"))
		tools := r.URL.Query().Get("tools")
		cbfunc := r.URL.Query().Get("callback")
		fmt.Println(msgWithDateProc(NProc, "httpcbgettoolsversion = tools:"+tools+" callback: "+cbfunc))
		initTools(tools)
		//initToolsVersion()
		w.Header().Set("Content-Type", "application/json")
		str := cbfunc + "(" + string(ExportToolsVersion()) + ");"
		arr := []byte(str)
		w.Write(arr)
		fmt.Println(msgWithDateProc(NProc, "httpcbgettoolsversion done"))
		fmt.Println("*****************************************")
		//EndProc(NProc)
		return
	}
}

const usage = `
<html>
<head>
<title>inspectfile</title>
<meta http-equiv="refresh" content="5; URL=/demo/demo.html">
</head>
<body>
</body>
</html>
`

func handleMain(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" || r.URL.Path != "/" {
		handleErr(w, http.StatusNotFound, fmt.Errorf("Not a valid path\n"))
		return
	}
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, usage)
}

func listen(port string) {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/inspect", httpinspect())
	http.HandleFunc("/localinspect", httplocalinspect())
	http.HandleFunc("/inspectpath/", httpinspectpath())
	http.HandleFunc("/avupdate", httpavupdate())
	http.HandleFunc("/cbavupdate", httpcbavupdate())
	http.HandleFunc("/gettoolsversion", httpgettoolsversion())
	http.HandleFunc("/cbgettoolsversion", httpcbgettoolsversion())
	http.Handle("/demo/", http.StripPrefix("/demo/", http.FileServer(http.Dir("/home/demo"))))
	http.ListenAndServe(port, nil)
}
