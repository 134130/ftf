package terminal

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

var reportRegex *regexp.Regexp

func init() {
	reportRegex = regexp.MustCompile("\x1b\\[(\\d+);(\\d+)R")
}

const (
	csi = "\x1b["

	enableAltBuf  = csi + "?1049h"
	disableAltBuf = csi + "?1049l"
	showCursor    = csi + "?25h"
	hideCursor    = csi + "?25l"
	enableWrap    = csi + "?7h"
	disableWrap   = csi + "?7l"

	deviceStatusReport = csi + "6n"
	saveCursor         = csi + "s"
	restoreCursor      = csi + "u"

	eraseDisplayEnd = csi + "0J"
	eraseDisplayAll = csi + "2J"
	eraseLineAll    = csi + "2K"

	resetGraphics = csi + "m"

	bold      = "1"
	faint     = "2" // Not widely supported.
	italic    = "3" // Not widely supported. Sometimes treated as inverse.
	underline = "4"
	reverse   = "7"
	nobold    = "21"
	nofaint   = "22"
	noreverse = "27"
)

func cursorUp(args ...int) string {
	i := 1
	if len(args) > 0 {
		i = args[0]
	}
	if i == 0 {
		return ""
	}
	return fmt.Sprint(csi, i, "A")
}

func cursorDown(args ...int) string {
	i := 1
	if len(args) > 0 {
		i = args[0]
	}
	if i == 0 {
		return ""
	}
	return fmt.Sprint(csi, i, "B")
}

func cursorForward(args ...int) string {
	i := 1
	if len(args) > 0 {
		i = args[0]
	}
	if i == 0 {
		return ""
	}
	return fmt.Sprint(csi, i, "C")
}

func cursorBack(args ...int) string {
	i := 1
	if len(args) > 0 {
		i = args[0]
	}
	if i == 0 {
		return ""
	}
	return fmt.Sprint(csi, i, "D")
}

func cursorPosition(row int, column int) string {
	return fmt.Sprint(csi, row, ";", column, "H")
}

func readReport(in io.Reader) (int, int, error) {
	input := make([]byte, 128)
	n, err := in.Read(input)
	if err != nil {
		return 0, 0, err
	}
	match := reportRegex.FindStringSubmatch(string(input[:n]))
	if match == nil {
		return 0, 0, errors.New("could not parse report")
	}
	row, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, 0, err
	}
	col, err := strconv.Atoi(match[2])
	if err != nil {
		return 0, 0, err
	}
	return row, col, nil
}
