package mfaviewresource

import (
	"log"
	"os"
)

// Needs work
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

func CSVToArray(rootDirPath string, fileName string) (result [][3]string) {

	// Go introduced OpenRoot in version 1.24, it restricts file operations to a single directory
	rootDir, err := os.OpenRoot(rootDirPath)

	if err != nil {
		panic("Directory path does not exist")
	}

	defer rootDir.Close()

	// Open the file
	file, err := rootDir.Open(fileName)

	if err != nil {
		log.Fatal("File " + fileName + " cannot be opened or does not exist")
	}

	// Close the file
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// For loop to scan lines in the file and add to a slice
	for scanner.Scan() {
		line := scanner.Text()
		// Make a slice and split each line into an index based on the, delimiter
		var splitSlice = []string{}
		splitSlice = strings.Split(line, ",")
		// Convert splitSlice into an array named splitArray
		splitArray := [3]string{}
		copy(splitArray[:], splitSlice)
		// Check if any array elements are empty
		if splitArray[0] == "" || splitArray[1] == "" || splitArray[2] == "" {
		} else {
			result = append(result, splitArray)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return
}
