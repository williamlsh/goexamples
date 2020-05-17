// Origin: https://eli.thegreenplace.net/2020/faking-stdin-and-stdout-in-go/

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

type FakeStdio struct {
	origStdout   *os.File
	stdoutReader *os.File
	
	outCh chan []byte
	
	origStdin   *os.File
	stdinWriter *os.File
}

func New(stdinText string) (*FakeStdio, error) {
	stdinReader, stdinWriter, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	
	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	
	origStdin := os.Stdin
	os.Stdin = stdinReader
	
	_, err = stdinWriter.Write([]byte(stdinText))
	if err != nil {
		stdinWriter.Close()
		os.Stdin = origStdin
		return nil, err
	}
	
	origStdout := os.Stdout
	os.Stdout = stdoutWriter
	
	outCh := make(chan []byte)
	
	go func() {
		var b bytes.Buffer
		if _, err := io.Copy(&b, stdoutReader); err != nil {
			log.Println(err)
		}
		outCh <- b.Bytes()
	}()
	
	return &FakeStdio{
		origStdout:   origStdout,
		stdoutReader: stdoutReader,
		outCh:        outCh,
		origStdin:    origStdin,
		stdinWriter:  stdinWriter,
	}, nil
}

func (sf *FakeStdio) ReadAndRestore() ([]byte, error) {
	if sf.stdoutReader == nil {
		return nil, errors.New("ReadAndRestore from closed FakeStdio")
	}
	
	os.Stdout.Close()
	out := <-sf.outCh
	
	os.Stdout = sf.origStdout
	os.Stdin = sf.origStdin
	
	if sf.stdoutReader != nil {
		sf.stdoutReader.Close()
		sf.stdoutReader = nil
	}
	
	if sf.stdinWriter != nil {
		sf.stdinWriter.Close()
		sf.stdinWriter = nil
	}
	
	return out, nil
}

func main() {
	fs, err := New("Input text")
	if err != nil {
		log.Fatal(err)
	}
	
	var scanned string
	fmt.Scanf("%s", &scanned)
	
	fmt.Print("Some output")
	
	b, err := fs.ReadAndRestore()
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Scanned: %q, Captured: %q", scanned, string(b))
}
