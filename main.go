package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
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

func getXColor(res string) string {
	buf, err := exec.Command("xrq", res).Output()
	if err != nil {
		log.Println(err)
		return ""
	}
	return strings.TrimSpace(string(buf))
}

func colorsCSSHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, ":root {")

	fmt.Fprintf(w, "--fg-color: %s;\n", getXColor("*foreground"))
	fmt.Fprintf(w, "--bg-color: %s;\n", getXColor("*background"))
	fmt.Fprintf(w, "--cr-color: %s;\n", getXColor("*cursorColor"))

	for i := 0; i < 16; i++ {
		fmt.Fprintf(w, "--color%d: %s;\n", i, getXColor("*color"+strconv.Itoa(i)))
	}

	fmt.Fprintln(w, "}")
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
