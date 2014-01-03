// +build android

package testlib

import (
	"bytes"
	"fmt"
	"os/exec"
)

func Tap(x, y float32) error {
	var (
		errbuf bytes.Buffer
		outbuf bytes.Buffer
	)
	cmd := exec.Command(
		`/system/bin/sh`,
		`-c`,
		fmt.Sprintf("input tap %f %f", x, y),
	)
	cmd.Stderr, cmd.Stdout = &errbuf, &outbuf
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func Move(x, y float32) error {
	var (
		errbuf bytes.Buffer
		outbuf bytes.Buffer
	)
	cmd := exec.Command(
		`/system/bin/sh`,
		`-c`,
		fmt.Sprintf("input swipe %f %f %f %f", x, y, x+2, x+2),
	)
	cmd.Stderr, cmd.Stdout = &errbuf, &outbuf
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
