// +build !android

package testlib

import (
	"math"
	"os/exec"
	"strconv"
)

func move(x, y float32) error {
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

func Move(x1, y1, x2, y2 float32) error {
	dx := math.Abs(float64(x1) - float64(x2))
	dy := math.Abs(float64(y1) - float64(y2))

	inc := 1.0
	if dx != 0 {
		inc = dy / dx
	}

	y := y1 + 1

	for x := x1 + 1; x <= x2; x++ {
		move(x, y)
		y += float32(inc)
	}

	return nil
}
