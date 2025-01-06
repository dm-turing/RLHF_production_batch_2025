package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

func sample1(filenames []string) {
	for _, filename := range filenames {
		data, err := readFile(filename)
		if err != nil {
			continue
		}
		processData(data)
	}
}

func sample2() {
	conn, err := net.Dial("tcp", "example.com:80")
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}
	defer conn.Close()

	resp, err := http.Get("http://example.com/")
	if err != nil {
		log.Fatalf("Error requesting: %v", err)
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}
}

func processFile(filename string) error {
	data, err := readFile(filename)
	if err != nil {
		return err
	}

	err = processData(data)
	if err != nil {
		return err
	}

	return nil
}

func processData(data string) error {
	// DO nothing
	return nil
}

func readFile(data string) (string, error) {
	return "", nil
}

func main() {
	sample1([]string{"file1.txt", "file2.txt", "file3.txt"})
	sample2()
	processFile("example1.txt")
}
