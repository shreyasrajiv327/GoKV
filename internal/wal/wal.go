package wal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type WAL struct{
	file *os.File
	path string
}

func New(path string) (*WAL, error){
	file, err := os.OpenFile(
		path,
		os.O_CREATE|os.O_RDWR|os.O_APPEND,
		0644,
	)

	if err != nil{
		return nil, err
	}

	return &WAL{
		file: file,
		path: path,
	}, nil
}

func (w *WAL) LogPut(key, value string) error{
	entry := fmt.Sprintf("PUT %s %s\n", key, value)

	_,err := w.file.WriteString(entry)
	if err != nil{
		return err
	}

	return w.file.Sync()
}

func ( w*WAL) LogDelete(key string) error{
	entry := fmt.Sprintf("DELETE %s\n", key)
		_,err := w.file.WriteString(entry)
	if err != nil{
		return err
	}

	return w.file.Sync()

}

func (w *WAL) Replay(
	putFunc func(string, string),
	deleteFunc func(string),
) error {

	file, err := os.Open(w.path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, " ")

		if len(parts) < 2 {
			continue
		}

		switch parts[0] {

		case "PUT":
			if len(parts) < 3 {
				continue
			}

			key := parts[1]
			value := strings.Join(parts[2:], " ")

			putFunc(key, value)

		case "DELETE":
			key := parts[1]
			deleteFunc(key)
		}
	}

	return scanner.Err()
}

func (w *WAL) Clear() error {

	err := w.file.Close()
	if err != nil {
		return err
	}

	file, err := os.Create(w.path)
	if err != nil {
		return err
	}

	w.file = file

	return nil
}