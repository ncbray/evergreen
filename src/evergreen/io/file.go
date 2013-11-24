package io

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func WriteFile(filename string, data []byte) error {
	dirname := filepath.Dir(filename)
	err := os.MkdirAll(dirname, 0700)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0600)
}

func WriteDot(data string, outfile string) error {
	cmd := exec.Command("dot", "-Tsvg")
	cmd.Stdin = strings.NewReader(data)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	if stderr.Len() > 0 {
		return fmt.Errorf("dot failed: %s", strings.TrimSpace(stderr.String()))
	}
	return WriteFile(outfile, stdout.Bytes())
}
