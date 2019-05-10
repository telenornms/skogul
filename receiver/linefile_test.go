package receiver_test

import (
	"encoding/json"
	"fmt"
	"github.com/KristianLyng/skogul"
	"github.com/KristianLyng/skogul/parser"
	"github.com/KristianLyng/skogul/receiver"
	"github.com/KristianLyng/skogul/sender"
	"math/rand"
	"os"
	"syscall"
	"testing"
	"time"
)

func deleteFile(t *testing.T, file string) {
	err := os.Remove(file)
	if err != nil {
		t.Errorf("Failed to remove old test file %s: %v", file, err)
	}
}
func TestLinefile(t *testing.T) {
	rand.Seed(int64(time.Now().Nanosecond()))
	one := &(sender.Test{})

	file := fmt.Sprintf("%s/skogul-linefiletest-%d-%d", os.TempDir(), os.Getpid(), rand.Int())

	_, err := os.Stat(file)

	if err != nil && !os.IsNotExist(err) {
		t.Errorf("Error statting tmp file %s: %v", file, err)
		return
	}

	if !os.IsNotExist(err) {
		t.Errorf("File possibly exists already: %s", file)
		return
	}

	err = syscall.Mkfifo(file, 0600)

	if err != nil {
		t.Errorf("Unable to make fifo %s: %v", file, err)
		return
	}
	defer deleteFile(t, file)

	h := skogul.Handler{Sender: one, Parser: parser.JSON{}}
	rcv := receiver.LineFile{File: file, Handler: h}
	go rcv.Start()
	b, err := json.Marshal(validContainer)
	if err != nil {
		t.Errorf("Failed to marshal container: %v", err)
		return
	}
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		t.Errorf("Unable to open file/fifo for writing: %v", err)
		return
	}
	defer func() {
		f.Close()
	}()
	f.WriteString(fmt.Sprintf("%s\n", b))
	time.Sleep(time.Duration(10 * time.Millisecond))
	if one.Received != 1 {
		t.Errorf("Didn't receive thing on other end!")
	}
	f.WriteString(fmt.Sprintf("%s\n", b))
	time.Sleep(time.Duration(10 * time.Millisecond))
	if one.Received != 2 {
		t.Errorf("Didn't receive thing on other end!")
	}
}
