package main

// This app takes a .json file that is the exported / saved contents of a mongodb collection
// and then writes out a new .js file with commands wrapped around the documents in the .json
// file ... so that this new .js script can be used to write the collection into another
// mongodb.
// Thus facilitating the transfer of mongodb collections from develop onto local MacBook mongodb.

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func main() {

	fmt.Printf("processing collections\n")

	processCollection("datasets")
	processCollection("dimension-options")
	processCollection("editions")
	processCollection("instances")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func processCollection(inputFileName string) {
	collectionName := inputFileName
	fmt.Printf("Processing collection: %s\n", collectionName)
	inputFileName += ".json"
	inputJsonFile, err := os.Open(inputFileName)
	check(err)
	defer inputJsonFile.Close()

	outputFileName := "insert-" + collectionName + ".js"

	outputJsFile, err := os.Create(outputFileName)
	check(err)
	defer outputJsFile.Close()

	// read first lines until an opening curly brace is found
	r := bufio.NewReader(inputJsonFile)
	for {
		s, err := Readln(r)
		check(err)
		if s[0] == '{' {
			break
		}
	}
	fmt.Printf("Found first opening curly brace\n")

	_, err = fmt.Fprint(outputJsFile, fmt.Sprintf("// init datasets database with collection: %s\n", collectionName))
	check(err)

	_, err = fmt.Fprint(outputJsFile, "\ndb = db.getSiblingDB('datasets')\n")
	check(err)

	if collectionName == "dimension-options" {
		collectionName = "dimension.options" // force to be the same name as in original collection on develop
	}

	_, err = fmt.Fprint(outputJsFile, fmt.Sprintf("\ndb.%s.remove({})\n\n", collectionName))
	check(err)

	for {
		// prefix command
		_, err = fmt.Fprint(outputJsFile, fmt.Sprintf("db.%s.insertOne({\n", collectionName))
		check(err)

		for {
			s, err := Readln(r)
			check(err)
			if s[0] == '}' {
				break
			}
			_, err = fmt.Fprint(outputJsFile, fmt.Sprintf("%s\n", s))
			check(err)
		}

		// postfix the close
		_, err = fmt.Fprint(outputJsFile, "})\n")
		check(err)

		s, err := Readln(r)
		if err != nil {
			// end of file found
			break
		}

		if s[0] != '{' {
			fmt.Printf("BAD line: %s\n", s)
			panic(errors.New("Opening curly brace was expected, but its not there"))
		}
	}
}

// Readln returns a single line (without the ending \n) from the input buffered reader.
// An error is returned if there is an error with the buffered reader.
func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}
