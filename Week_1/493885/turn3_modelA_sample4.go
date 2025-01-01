package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

func readFileWithTimeout(ctx context.Context, file *os.File) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		buf := make([]byte, 1024)
		n, err := io.ReadFull(file, buf)
		if err != nil && err != io.EOF {
			return err
		}
		fmt.Println("Read:", string(buf[:n]))
	}
	return nil
}

func main() {
	file1, err := os.Open("file1.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file1.Close()

	file2, err := os.Open("file2.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file2.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = readFileWithTimeout(ctx, file1)
	if err != nil {
		fmt.Println("Error reading file1:", err)
	}

	err = readFileWithTimeout(ctx, file2)
	if err != nil {
		fmt.Println("Error reading file2:", err)
	}
}
