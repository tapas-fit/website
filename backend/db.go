package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

func fname(path string, day uint32) string {
	return fmt.Sprintf("%s/waitlist-%d.csv", path, day)
}

type DB struct {
	path string
	fd   *os.File
	day  uint32
	sync.Mutex
}

func OpenDB(path string) (*DB, error) {
	day := uint32(time.Now().Unix() / 86400)

	fd, err := os.OpenFile(fname(path, day), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	d := &DB{path: path, fd: fd, day: day}

	go d.rotate()
	return d, nil
}

func (d *DB) rotate() {
	for {
		// sleep until it's 00:00:00 and then rotate the file
		t := time.Now()
		tomorrow := t.Add(time.Hour * 24)
		tomorrow = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
		time.Sleep(tomorrow.Sub(t))
		d.Lock()
		d.fd.Close()
		d.day = uint32(time.Now().Unix() / 86400)
		var err error
		d.fd, err = os.OpenFile(fname(d.path, d.day), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		d.Unlock()
	}
}

func (d *DB) add(f url.Values) error {
	d.Lock()
	defer d.Unlock()

	_, err := fmt.Fprintf(d.fd, "ts=%s;", time.Now().Format(time.RFC3339))
	if err != nil {
		return err
	}

	for k, v := range f {
		if k == "cf-turnstile-response" || len(v) == 0 {
			continue
		}
		_, err := fmt.Fprintf(d.fd, "%s=%s;", k, strings.TrimSuffix(strings.TrimPrefix(v[0], "["), "]"))
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintf(d.fd, "\n")
	if err != nil {
		return err
	}
	return d.fd.Sync()
}
