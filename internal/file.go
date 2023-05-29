package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type datum map[string]interface{}

func (d datum) String(keys []string, fields bool) string {
	var str strings.Builder
	for _, key := range keys {
		if value, ok := d[key]; ok {
			if fields {
				str.WriteString(fmt.Sprintf("%s: %v\n", key, value))
			} else {
				str.WriteString(fmt.Sprintf("%v\n", value))
			}
		}
	}
	return str.String()
}

func isEqual(a datum, b datum) bool {
	delete(a, "ambrosia")
	delete(b, "ambrosia")
	return reflect.DeepEqual(a, b)
}

func contains(s []datum, e datum) bool {
	for _, a := range s {
		if isEqual(a, e) {
			return true
		}
	}
	return false
}

func datumSub(a []datum, b []datum) []datum {
	var ret []datum
	for _, item := range a {
		if !contains(b, item) {
			ret = append(ret, item)
		}
	}
	return ret
}

func load(path string) ([]datum, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []datum

	scanner := bufio.NewScanner(file)
	scanner.Buffer(nil, 10*1024*1024)

	for scanner.Scan() {
		var d datum
		if err := json.Unmarshal(scanner.Bytes(), &d); err != nil {
			return nil, fmt.Errorf("%w: are you using a valid JSONL file?", err)
		}
		data = append(data, d)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func loadWordlist(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var ret []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}
	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func loadResumable(cmd string, infilePath string) ([]datum, error) {
	searchPath := filepath.Dir(infilePath)
	searchPrefix := searchPrefix(infilePath)

	var resumableFiles []string
	files, err := os.ReadDir(searchPath)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if !f.IsDir() && strings.HasPrefix(f.Name(), searchPrefix) {
			resumableFiles = append(resumableFiles, filepath.Join(searchPath, f.Name()))
		}
	}

	var ret []datum

	for _, f := range resumableFiles {
		data, err := load(f)
		if err != nil {
			return nil, err
		}
		ret = append(ret, data...)
	}

	return ret, nil
}

func write(path string, data []datum) (err error) {
	// Error if the file already exists.
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := file.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	writer := bufio.NewWriter(file)
	defer func() {
		flushErr := writer.Flush()
		if flushErr != nil && err == nil {
			err = flushErr
		}
	}()

	for _, d := range data {
		b, err := json.Marshal(d)
		if err != nil {
			return err
		}
		if _, err := writer.Write(append(b, '\n')); err != nil {
			return err
		}
	}

	return nil
}
