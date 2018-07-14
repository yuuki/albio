package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestRun_versionFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("albio --version", " ")

	status := cli.Run(args)
	if status != 0 {
		t.Errorf("expected %d to eq %d", status, exitCodeOK)
	}

	expected := fmt.Sprintf("albio version %s", version)
	if !strings.Contains(errStream.String(), expected) {
		t.Errorf("expected %q to eq %q", errStream.String(), expected)
	}
}

func TestRun_parseError(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("albio --not-exist", " ")

	status := cli.Run(args)
	if status != exitCodeErr {
		t.Errorf("expected %d to eq %d", status, exitCodeErr)
	}

	expected := "Usage: albio"
	if !strings.Contains(errStream.String(), expected) {
		t.Fatalf("expected %q to contain %q", errStream.String(), expected)
	}
}
