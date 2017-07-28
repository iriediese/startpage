package main

import (
	"flag"
	"fmt"
	"go/build"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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
			s = COLOR0
		case "1":
			s = COLOR1
		case "2":
			s = COLOR2
		case "3":
			s = COLOR3
		case "4":
			s = COLOR4
		case "5":
			s = COLOR5
		case "6":
			s = COLOR6
		case "7":
			s = COLOR7
		case "8":
			s = COLOR8
		case "9":
			s = COLOR9
		case "10":
			s = COLOR10
		case "11":
			s = COLOR11
		case "12":
			s = COLOR12
		case "13":
			s = COLOR13
		case "14":
			s = COLOR14
		case "15":
			s = COLOR15

		case "background":
			s = COLORBG
		case "foreground":
			s = COLORFG
		default:
			s = COLORCU
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
		filepath.Join(
			build.Default.GOPATH,
			"src/github.com/tudurom/startpage/static",
		),
		"path to the startpage",
	)
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
