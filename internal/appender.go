package internal

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"unicode"
)

var (
	errNilResponse = fmt.Errorf("nil response")
)

type prefixAppender struct {
	appenders map[rune]*fileAppender
	mutex     *sync.Mutex
	pathTmpl  string
}

func newPrefixAppender(pathTmpl string) *prefixAppender {
	return &prefixAppender{
		appenders: make(map[rune]*fileAppender),
		mutex:     &sync.Mutex{},
		pathTmpl:  pathTmpl,
	}
}

func (a *prefixAppender) append(resp string, d datum) error {
	if len(resp) == 0 {
		return errNilResponse
	}

	var prefix rune
	for _, char := range resp {
		if !unicode.IsSpace(char) && !unicode.IsPunct(char) {
			prefix = char
			break
		}
	}

	var appender *fileAppender
	var ok bool

	a.mutex.Lock()
	defer a.mutex.Unlock()

	appender, ok = a.appenders[prefix]
	if !ok {
		path := fmt.Sprintf(a.pathTmpl, prefix)
		newA, err := newFileAppender(path)
		if err != nil {
			return err
		}
		a.appenders[prefix] = newA
		appender = newA
	}

	err := appender.append(d)
	if err != nil {
		return err
	}

	return nil
}

func (a *prefixAppender) appendWithResponse(resp string, d datum) error {
	d["ambrosia"] = resp
	return a.append(resp, d)
}

func (a *prefixAppender) close() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	for _, appender := range a.appenders {
		err := appender.close()
		if err != nil {
			return err
		}
	}
	return nil
}

type fileAppender struct {
	path  string
	file  *os.File
	buf   *bufio.Writer
	mutex *sync.Mutex
}

func newFileAppender(path string) (*fileAppender, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening file for append: %w", err)
	}

	w := bufio.NewWriter(f)

	return &fileAppender{
		path:  path,
		file:  f,
		buf:   w,
		mutex: &sync.Mutex{},
	}, nil
}

func (a *fileAppender) append(d datum) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	err := enc.Encode(d)
	if err != nil {
		return err
	}

	a.mutex.Lock()
	defer a.mutex.Unlock()

	_, err = a.buf.Write(buf.Bytes())
	if err != nil {
		return err
	}

	err = a.buf.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (a *fileAppender) close() error {
	err := a.buf.Flush()
	if err != nil {
		return fmt.Errorf("error flushing buffer: %w", err)
	}

	err = a.file.Close()
	if err != nil {
		return fmt.Errorf("error closing file: %w", err)
	}

	return nil
}
