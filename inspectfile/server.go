// Server module for "inspectFile"

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func handleErr(w http.ResponseWriter, status int, e error) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, e.Error())
}

func httpinspect() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//Upload a file, Copy in the /tmp directory, Inspect the file and Remove the file
		mime := "application/json"
		tools := r.URL.Query().Get("tools")
		fmt.Println("*****************************************")
		fmt.Println(msgWithDate("httpinspect launched."))
		fmt.Println(msgWithDate("URL String: " + r.URL.String() + " & tools" + tools))
		initTools(tools)
		initToolsVersion()

		var fn string //fn is the filename

		if r.Header.Get("Content-Type") == "application/octet-stream" {
			//used by filedrop.js
			//the size of the file is limited to 2 GB (not tested yet !)
			fn = r.Header.Get("X-File-Name") //get the filename
			fmt.Printf("Size of data to read = %d (%s)\n", r.ContentLength, r.Header.Get("X-File-Size"))
			data := make([]byte, r.ContentLength) //get the data upload in memory
			fmt.Println(msgWithDate("Trying to read the data !!"))
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "ERROR in reading the data: "+err.Error(), http.StatusBadRequest)
				fmt.Println(msgWithDate("ReadAll Error = " + err.Error()))
				return
			} else {
				fmt.Println(msgWithDate("The data is correctly read"))
			}
			err = ioutil.WriteFile("/tmp/"+fn, data, 0644) //
			if err != nil {
				http.Error(w, "ERROR creating file in /tmp/: "+err.Error(), http.StatusBadRequest)
				fmt.Println(msgWithDate("ERROR creating file /tmp/" + fn))
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
				fmt.Println(msgWithDate("ERROR parsing uploaded file: " + err.Error()))
				return
			}
			outfile, err := os.Create("/tmp/" + fn)
			if err != nil {
				http.Error(w, "ERROR creating file: "+err.Error(), http.StatusBadRequest)
				fmt.Println(msgWithDate("ERROR creating file: " + err.Error()))
				return
			}
			_, err = io.Copy(outfile, infile)
			if err != nil {
				http.Error(w, "ERROR saving file: "+err.Error(), http.StatusBadRequest)
				fmt.Println(msgWithDate("ERROR saving file: " + err.Error()))
				return
			}
			output = inspectfile(outfile.Name(), nil)
		}

		w.Header().Set("Content-Type", mime)
		w.Write(output)

		err := os.Remove("/tmp/" + fn)
		if err != nil {
			http.Error(w, "ERROR removing file: "+err.Error(), http.StatusBadRequest)
			fmt.Println(msgWithDate("ERROR removing file: " + err.Error()))
			return
		}
		fmt.Println(msgWithDate("httpinspect done."))
		fmt.Println("*****************************************")
		return
	}
}

func httpinspectpath() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//inspect a path
		//THe path used the Environment variable %MOUNTDIR%
		mime := "application/json"
		tools := r.URL.Query().Get("tools")
		fmt.Println("*****************************************")
		fmt.Println(msgWithDate("httpinspectpath launched."))
		fmt.Println(msgWithDate("URL String: " + r.URL.String() + " & tools" + tools))
		initTools(tools)
		initToolsVersion()

		path := r.URL.Path
		fmt.Println(msgWithDate("GET Path: " + path))
		if len(path) < 10 {
			http.Error(w, "You need to pass the filename or directory path after /inspect !!!", http.StatusBadRequest)
			fmt.Println(msgWithDate("httpinspectpath: ERROR no directory path after /inspect"))
			return
		} else {
			path = os.Getenv("MOUNTDIR") + path[len("/inspectpath"):]
		}
		info, err := os.Stat(path)
		if err != nil {
			handleErr(w, http.StatusNotFound, err)
			fmt.Println(msgWithDate("httpinspectpath: ERROR in getting inforamtion from path" + err.Error()))
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
		fmt.Println(msgWithDate("httpinspect PATH done."))
		fmt.Println("*****************************************")
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
		fmt.Println("*****************************************")
		fmt.Println(msgWithDate("httpAVUpdate launched."))
		w.Header().Set("Content-Type", "application/json")
		w.Write(launchAVUpdate())
		fmt.Println(msgWithDate("httpAVUpdate done."))
		fmt.Println("*****************************************")
		return
	}
}

func httpcbavupdate() func(w http.ResponseWriter, r *http.Request) {
	//callback version of httpavupdate
	return func(w http.ResponseWriter, r *http.Request) {
		//output := []byte("AV update - TO DO...")
		fmt.Println("*****************************************")
		fmt.Println(msgWithDate("httpcbAVUpdate launched."))
		cbfunc := r.URL.Query().Get("callback")
		fmt.Println(msgWithDate("httpcbAVUpdate = callback: " + cbfunc))
		w.Header().Set("Content-Type", "application/json")
		str := cbfunc + "(" + string(launchAVUpdate()) + ");"
		arr := []byte(str)
		w.Write(arr)
		//w.Write(launchAVUpdate())
		fmt.Println(msgWithDate("httpcbAVUpdate done."))
		fmt.Println("*****************************************")
		return
	}
}

func httpgettoolsversion() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//output := []byte("Tools version - TO DO...")
		fmt.Println("*****************************************")
		fmt.Println(msgWithDate("httpgettoolsversion launched"))
		tools := r.URL.Query().Get("tools")
		fmt.Println(msgWithDate("httpgettoolsversion = tools:" + tools))
		initTools(tools)
		//initToolsVersion()
		w.Header().Set("Content-Type", "application/json")
		w.Write(ExportToolsVersion())
		fmt.Println(msgWithDate("httpgettoolsversion done"))
		fmt.Println("*****************************************")
		return
	}
}

func httpcbgettoolsversion() func(w http.ResponseWriter, r *http.Request) {
	//callback version of httpgettoolsversion
	return func(w http.ResponseWriter, r *http.Request) {
		//output := []byte("Tools version - TO DO...")
		fmt.Println("*****************************************")
		fmt.Println(msgWithDate("httpcbgettoolsversion launched"))
		tools := r.URL.Query().Get("tools")
		cbfunc := r.URL.Query().Get("callback")
		fmt.Println(msgWithDate("httpcbgettoolsversion = tools:" + tools + " callback: " + cbfunc))
		initTools(tools)
		//initToolsVersion()
		w.Header().Set("Content-Type", "application/json")
		str := cbfunc + "(" + string(ExportToolsVersion()) + ");"
		arr := []byte(str)
		w.Write(arr)
		fmt.Println(msgWithDate("httpcbgettoolsversion done"))
		fmt.Println("*****************************************")
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
