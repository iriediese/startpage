package main

// #cgo LDFLAGS: -lX11
// #include <stdlib.h>
// #include <X11/Xresource.h>
// char * getAddr(XrmValue val) {
//		return val.addr;
// }
import "C"

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

var dpy *C.Display
var xrm *C.char
var xrdb C.XrmDatabase

func exists(dp string, isDir bool) bool {
	fi, err := os.Stat(dp)
	if err == nil {
		return isDir == fi.IsDir()
	} else if os.IsNotExist(err) {
		return false
	} else {
		log.Fatal(err)
		return false
	}
}

func getXColor(res string) string {
	var t *C.char
	var xvalue C.XrmValue

	C.XrmGetResource(xrdb, C.CString(res), C.CString("*"), &t, &xvalue)

	if C.getAddr(xvalue) != nil {
		return C.GoString(C.getAddr(xvalue))
	} else {
		return ""
	}
}

func getXrdb() error {
	dpy = C.XOpenDisplay(nil)
	if dpy == nil {
		return errors.New("could not open display")
	}

	C.XrmInitialize()
	xrm = C.XResourceManagerString(dpy)

	if xrm == nil {
		return errors.New("could not get resource properties")
	}

	xrdb = C.XrmGetStringDatabase(xrm)

	return nil
}

func freeXrdb() {
	C.XrmDestroyDatabase(xrdb)
	C.XFlush(dpy)
	C.XCloseDisplay(dpy)
}

func colorsCSSHandler(w http.ResponseWriter, r *http.Request) {
	err := getXrdb()
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "text/css")
	fmt.Fprintln(w, ":root {")

	fmt.Fprintf(w, "--fg-color: %s;\n", getXColor("*foreground"))
	fmt.Fprintf(w, "--bg-color: %s;\n", getXColor("*background"))
	fmt.Fprintf(w, "--cr-color: %s;\n", getXColor("*cursorColor"))

	for i := 0; i < 16; i++ {
		fmt.Fprintf(w, "--color%d: %s;\n", i, getXColor("*color"+strconv.Itoa(i)))
	}

	fmt.Fprintln(w, "}")
	freeXrdb()
}

func main() {
	var startpagePath, listenAddress string

	flag.StringVar(&startpagePath, "path", "", "path to the startpage")
	flag.StringVar(&listenAddress, "listen", ":8081", "address to listen for http connections")
	flag.Parse()

	// checking flags
	if !exists(startpagePath, true) {
		log.Fatal("`" + startpagePath + "` is an invalid startpage path")
	}

	fs := http.FileServer(http.Dir(startpagePath))
	http.Handle("/", fs)
	http.HandleFunc("/colors.css", colorsCSSHandler)
	log.Println("Listening on " + listenAddress)
	err := http.ListenAndServe(listenAddress, nil)

	if err != nil {
		log.Fatal(err)
	}
}
