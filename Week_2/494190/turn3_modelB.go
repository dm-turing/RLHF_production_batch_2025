package main

import (
	"fmt"
	"sync"
)

func processFile(filename string, errc chan error) {
	// ... process the file
	err := fmt.Errorf("error in %s", filename)
	if err != nil {
		errc <- err
		return
	}
	errc <- nil
}

func main() {
	filenames := []string{"file1.txt", "file2.txt", "file3.txt"}
	errc := make(chan error, len(filenames))

	var wg sync.WaitGroup
	wg.Add(len(filenames))

	for _, filename := range filenames {
		go func(f string) {
			defer wg.Done()
			processFile(f, errc)
		}(filename)
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}
