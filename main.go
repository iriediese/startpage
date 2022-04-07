package main

import (
	"flag"
	"fmt"
//	"go/build"
	"log"
	"net/http"
	"os"
	"os/exec"
//	"path/filepath"
	"strconv"
	"strings"
)

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

func getXRes(col string) string {
	var q, s string
	if len(col) <= 2 {
		q = "*color" + col
	} else {
		q = "*" + col
	}

	b, err := exec.Command("xrq", q).Output()
	if err == nil {
		s = strings.TrimSpace(string(b))
	} else {
		switch col {
		case "0":
			s = "#050D18"
		case "1":
			s = "#345268"
		case "2":
			s = "#486176"
		case "3":
			s = "#3B6B84"
		case "4":
			s = "#53738A"
		case "5":
			s = "#647B90"
		case "6":
			s = "#6C8EA6"
		case "7":
			s = "#b4cad9"
		case "8":
			s = "#7d8d97"
		case "9":
			s = "#345268"
		case "10":
			s = "#486176"
		case "11":
			s = "#3B6B84"
		case "12":
			s = "#53738A"
		case "13":
			s = "#647B90"
		case "14":
			s = "#6C8EA6"
		case "15":
			s = "#b4cad9"

		case "background":
			s = "#050D18"
		case "foreground":
			s = "#b4cad9"
		default: // cursor color
			s = "#c5c8c6"
		}
	}
	return s
}

func colorsCSSHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	fmt.Fprintln(w, ":root {")

	fmt.Fprintf(w, "--fg-color: %s;\n", getXRes("foreground"))
	fmt.Fprintf(w, "--bg-color: %s;\n", getXRes("background"))
	fmt.Fprintf(w, "--cr-color: %s;\n", getXRes("cursorColor"))

	for i := 0; i < 16; i++ {
		fmt.Fprintf(w, "--color%d: %s;\n", i, getXRes(strconv.Itoa(i)))
	}

	fmt.Fprintln(w, "}")
}

func main() {
	var startpagePath, listenAddress string

	flag.StringVar(
		&startpagePath,
		"path",
		"/home/tedi/sources/startpage/static",
		"path to the startpage",
	)
	flag.StringVar(&listenAddress, "listen", ":443", "address to listen for http connections")
	flag.Parse()

	// checking flags
	if !exists(startpagePath, true) {
		log.Fatal("`" + startpagePath + "` is an invalid startpage path")
	}

	fs := http.FileServer(http.Dir(startpagePath))
	http.Handle("/", fs)
	http.HandleFunc("/colors.css", colorsCSSHandler)
	log.Println("Listening on " + listenAddress)
	err := http.ListenAndServeTLS(listenAddress, "/home/tedi/ssl/server.crt", "/home/tedi/ssl/server.key", nil)

	if err != nil {
		log.Fatal(err)
	}
}
