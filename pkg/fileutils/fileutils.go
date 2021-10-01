package fileutils

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

// FileExists checks if the file exists in the provided path
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// FolderExists checks if the folder exists
func FolderExists(foldername string) bool {
	info, err := os.Stat(foldername)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	return info.IsDir()
}

// CreateFolders in the list
func CreateFolders(paths ...string) error {
	for _, path := range paths {
		if err := CreateFolder(path); err != nil {
			return err
		}
	}

	return nil
}

// CreateFolder path
func CreateFolder(path string) error {
	return os.MkdirAll(path, 0700)
}

// ReadFileWithReader and stream on a channel
func ReadFileWithReader(r io.Reader) (chan string, error) {
	out := make(chan string)
	go func() {
		defer close(out)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			out <- scanner.Text()
		}
	}()

	return out, nil
}

// ReadFileWithReader with specific buffer size and stream on a channel
func ReadFileWithReaderAndBufferSize(r io.Reader, maxCapacity int) (chan string, error) {
	out := make(chan string)
	go func() {
		defer close(out)
		scanner := bufio.NewScanner(r)
		buf := make([]byte, maxCapacity)
		scanner.Buffer(buf, maxCapacity)
		for scanner.Scan() {
			out <- scanner.Text()
		}
	}()

	return out, nil
}

// ReadFile with filename
func ReadFile(filename string) (chan string, error) {
	if !FileExists(filename) {
		return nil, errors.New("file doesn't exist")
	}
	out := make(chan string)
	go func() {
		defer close(out)
		f, err := os.Open(filename)
		if err != nil {
			return
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			out <- scanner.Text()
		}
	}()

	return out, nil
}

// ReadFile with filename and specific buffer size
func ReadFileWithBufferSize(filename string, maxCapacity int) (chan string, error) {
	if !FileExists(filename) {
		return nil, errors.New("file doesn't exist")
	}
	out := make(chan string)
	go func() {
		defer close(out)
		f, err := os.Open(filename)
		if err != nil {
			return
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		buf := make([]byte, maxCapacity)
		scanner.Buffer(buf, maxCapacity)
		for scanner.Scan() {
			out <- scanner.Text()
		}
	}()

	return out, nil
}

func ReadFileAsStringSlice(filename string) ([]string, error) {
	if !FileExists(filename) {
		return nil, errors.New("file doesn't exist")
	}

	linesSlice := make([]string, 0)
	file, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		linesSlice = append(linesSlice, line)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
	}

	return linesSlice, nil
}

// GetTempFileName generate a temporary file name
func GetTempFileName() (string, error) {
	tmpfile, err := os.CreateTemp("", "")
	if err != nil {
		return "", err
	}
	tmpFileName := tmpfile.Name()
	if err := tmpfile.Close(); err != nil {
		return tmpFileName, err
	}
	err = os.RemoveAll(tmpFileName)
	return tmpFileName, err
}

// CopyFile from source to destination
func CopyFile(src, dst string) error {
	if !FileExists(src) {
		return errors.New("source file doesn't exist")
	}
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return dstFile.Sync()
}
