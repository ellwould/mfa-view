package mfaviewresource

import (
	"log"
	"os"
)

func getFile(file string) string {
	var path string
	path = "/usr/local/etc/mfaview/html/"

	var1, err := os.ReadFile(path + file)
	if err != nil {
		log.Fatal(err)

	}
	return string(var1)
}

func StartHTML() string {
	var1 := getFile("mfaview-start.html")
	return var1
}

func EndHTML() string {
	var1 := getFile("mfaview-end.html")
	return var1
}
