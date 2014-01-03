// +build !android

package testlib

import (
	"os/exec"
	"strconv"
)

func Tap(x, y float32) error {
	// Move the mouse on the given position
	cmd := exec.Command(
		`xdotool`,
		`search`,
		`--name`,
		`Gorgasm Test`,
		`mousemove`,
		`--window`,
		`%1`,
		strconv.FormatFloat(float64(x), 'f', -1, 32),
		strconv.FormatFloat(float64(y), 'f', -1, 32),
	)
	if err := cmd.Run(); err != nil {
		return err
	}

	// Press down left button
	cmd = exec.Command(
		`xdotool`,
		`search`,
		`--name`,
		`Gorgasm Test`,
		`mousedown`,
		`--window`,
		`%1`,
		`1`,
	)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func Move(x, y float32) error {
	// Move the mouse on the given position
	cmd := exec.Command(
		`xdotool`,
		`search`,
		`--name`,
		`Gorgasm Test`,
		`mousemove`,
		`--window`,
		`%1`,
		strconv.FormatFloat(float64(x), 'f', -1, 32),
		strconv.FormatFloat(float64(y), 'f', -1, 32),
	)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
