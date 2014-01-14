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

func Move(x1, y1, x2, y2 float32) error {
	var (
		errbuf bytes.Buffer
		outbuf bytes.Buffer
	)

	if x1 == x2 && y1 == y2 {
		return nil
	}

	cmd := exec.Command(
		`/system/bin/sh`,
		`-c`,
		fmt.Sprintf("input swipe %f %f %f %f", x1, y1, x2, y2),
	)
	cmd.Stderr, cmd.Stdout = &errbuf, &outbuf
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func Back() error {
	var (
		errbuf bytes.Buffer
		outbuf bytes.Buffer
	)
	cmd := exec.Command(
		`/system/bin/sh`,
		`-c`,
		fmt.Sprintf("input keyeven KEYCODE_BACK"),
	)
	cmd.Stderr, cmd.Stdout = &errbuf, &outbuf
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
