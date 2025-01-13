package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

// FilePool provides a synchronized way to manage file handles.
type FilePool struct {
	mu      sync.Mutex
	release chan *os.File
	files   []*os.File
	cap     int
	count   int
	newFile func() (*os.File, error)
}

func (f *FilePool) NewFile(path string) (*os.File, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.count < f.cap {
		_file, err := f.newFile()
		if err != nil {
			return nil, err
		}
		f.files = append(f.files, _file)
		f.count++
		return _file, nil
	}

	// Wait for a file handle to be returned
	select {
	case _file := <-f.release:
		return _file, nil
	default:
		log.Fatalf("file pool exhausted: can't open file %q", path)
	}
	return nil, nil
}

func (f *FilePool) Release(_file *os.File) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.count > 0 {
		f.files[f.count-1].Close()
		f.files[f.count-1] = _file
		f.count--
	}

	select {
	case f.release <- f.files[f.count]:
	default:
		// If no one is waiting, close the file
		_file.Close()
		f.release <- _file
	}
}

func (f *FilePool) Close() {
	close(f.release)
	for _, _file := range f.files {
		_file.Close()
	}
}

var (
	pool = FilePool{
		cap: 100, // Set the maximum number of open files
		newFile: func() (*os.File, error) {
			return os.Open("example.txt") // Replace with your file handling logic
		},
		release: make(chan *os.File, 1),
	}
)

func main() {
	defer pool.Close()

	// Concurrent processing using a limited number of file handles
	for i := 0; i < 1000; i++ {
		go func(idx int) {
			_file, err := pool.NewFile("example.txt")
			if err != nil {
				log.Fatal(err)
			}
			defer _file.Close()

			data, err := ioutil.ReadAll(_file)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(fmt.Sprintf("Read file %d: %s", idx, string(data)))

			pool.Release(_file)
		}(i)
	}
}
