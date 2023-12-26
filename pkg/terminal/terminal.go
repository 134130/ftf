package terminal

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"golang.org/x/term"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"reflect"
	"runtime/debug"
	"strings"
	"syscall"
)

type Flag int

const (
	FlagNone Flag = iota
	FlagPrint
)

type Terminal struct {
	config          *Config
	originalState   term.State
	insertedNewline bool
	previousRender  map[string]bool
	rows            int
	cols            int
	in              *os.File
	out             *os.File
	loop            bool
	currentRow      int
	flag            Flag
}

type Config struct {
	Height float64
}

func OpenTerm(config *Config) (*Terminal, error) {
	inFd, err := syscall.Open("/dev/tty", syscall.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	outFd, err := syscall.Open("/dev/tty", syscall.O_WRONLY, 0)
	t := Terminal{
		config:     config,
		in:         os.NewFile(uintptr(inFd), "/dev/tty"),
		out:        os.NewFile(uintptr(outFd), "/dev/tty"),
		currentRow: 1,
	}

	return &t, t.initTerm()
}

func (t *Terminal) initTerm() error {

	state, err := term.MakeRaw(int(t.out.Fd()))
	if err != nil {
		return err
	}
	t.originalState = *state

	if t.config.Height == 1.0 {
		t.out.WriteString(enableAltBuf)
		t.out.WriteString(cursorPosition(1, 1))
	} else {
		t.out.WriteString(deviceStatusReport)
		_, col, err := readReport(t.in)
		if err == nil && col > 1 {
			t.out.WriteString("\n")
			t.insertedNewline = true
		}
	}

	t.out.WriteString(disableWrap)
	t.out.WriteString(hideCursor)
	return nil
}

func (t *Terminal) revertTerm() {
	if t.config.Height == 1.0 {
		t.out.WriteString(enableAltBuf)
		t.out.WriteString(disableAltBuf)
	} else if t.insertedNewline {
		t.out.WriteString(cursorUp())
	}
	t.previousRender = map[string]bool{}
	t.out.WriteString(enableWrap)
	t.out.WriteString(showCursor)
	term.Restore(int(t.out.Fd()), &t.originalState)
}

func (t *Terminal) Close() {
	t.out.WriteString(eraseDisplayEnd)
	t.revertTerm()
	t.in.Close()
	t.out.Close()
}

func (t *Terminal) position(row int, col int) string {
	out := &strings.Builder{}
	out.WriteString(cursorBack(t.cols))
	vertical := row - t.currentRow
	if vertical > 0 {
		out.WriteString(cursorDown(vertical))
	} else if vertical < 0 {
		out.WriteString(cursorUp(-vertical))
	}
	out.WriteString(cursorForward(col - 1))
	t.currentRow = row
	return out.String()
}

func (t *Terminal) border(p Position) string {
	if p.Rows < 2 || p.Cols < 2 {
		return ""
	}
	out := &strings.Builder{}
	out.WriteString(t.position(p.Top, p.Left))
	out.WriteString("┌" + strings.Repeat("─", p.Cols-2) + "┐")
	for i := 1; i < p.Rows-1; i++ {
		out.WriteString(t.position(p.Top+i, p.Left+p.Cols-1))
		out.WriteString("│")
	}
	for i := 1; i < p.Rows-1; i++ {
		out.WriteString(t.position(p.Top+i, p.Left))
		out.WriteString("│")
	}
	out.WriteString(t.position(p.Top+p.Rows-1, p.Left))
	out.WriteString("└" + strings.Repeat("─", p.Cols-2) + "┘")
	out.WriteString(t.position(1, 1))
	return out.String()
}

func (t *Terminal) render(views []ViewRenderer) {
	out := &strings.Builder{}
	newRender := map[string]bool{}
	for _, view := range views {
		if !view.ShouldRender() {
			continue
		}

		p := view.Position(t.rows, t.cols)
		if view.HasBorder() {
			s := t.border(p)
			if _, ok := t.previousRender[s]; !ok {
				out.WriteString(s)
			}
			newRender[s] = true
			p = p.Shrink(1)
		}

		_ = reflect.TypeOf(view).Name()
		lines := view.Render(p)
		for row := 0; row < p.Rows; row++ {
			lineRender := &strings.Builder{}
			lineRender.WriteString(t.position(p.Top+row, p.Left))
			if row < len(lines) {
				lineRender.WriteString(lines[row].Text())
				if p.Cols > lines[row].Length() {
					lineRender.WriteString(strings.Repeat(" ", p.Cols-lines[row].Length()))
				}
			} else {
				lineRender.WriteString(strings.Repeat(" ", p.Cols))
			}
			if p.Top+row < t.rows {
				lineRender.WriteString("\n")
				t.currentRow += 1
			}
			lineRender.WriteString(t.position(1, 1))
			if _, ok := t.previousRender[lineRender.String()]; !ok {
				out.WriteString(lineRender.String())
			}
			newRender[lineRender.String()] = true
		}
	}
	t.out.WriteString(out.String())
	t.previousRender = newRender
}

func (t *Terminal) fetchWinSize() error {
	width, height, err := term.GetSize(int(t.out.Fd()))
	if err != nil {
		return err
	}
	t.rows = int(float64(height) * t.config.Height)
	t.cols = width
	return nil
}

func (t *Terminal) StartLoop(bindings map[string][]string, views []ViewRenderer) (flag Flag, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("terminal error: %v, stacktrace: %s", r, string(debug.Stack()))
		}
	}()

	intSigs := make(chan os.Signal, 1)
	signal.Notify(intSigs, syscall.SIGINT, syscall.SIGTERM)

	winChSig := make(chan os.Signal, 1)
	signal.Notify(winChSig, syscall.SIGWINCH)

	events := make(chan Event)
	nextEvents := make(chan bool)
	go readEvents(t.in, events, nextEvents)

	err = t.fetchWinSize()
	if err != nil {
		return t.flag, err
	}
	t.render(views)

	t.loop = true
	for {
		select {
		case <-intSigs:
			log.Debug().Msg("Received interrupt.")
			t.loop = false
		case <-winChSig:
			log.Debug().Msg("Received window change.")
			t.fetchWinSize()
			t.render(views)
			log.Debug().Msg("Rerendered.")
		case event := <-events:
			cmdKeys, ok := bindings[event.HashKey()]
			if !ok {
				continue
			}
			for _, cmdKey := range cmdKeys {
				if cmd, ok := t.getCommands()[cmdKey]; ok {
					err := cmd(t)
					if err != nil {
						return t.flag, err
					}
				} else {
					for _, view := range views {
						if cmd, ok := view.Commands()[cmdKey]; ok {
							err := cmd(t)
							if err != nil {
								return t.flag, err
							}
							break
						}
					}
				}
			}
			t.render(views)
		case nextEvents <- true:
		}
		if !t.loop {
			break
		}
	}
	return t.flag, nil
}

func (t *Terminal) getCommands() map[string]Command {
	return map[string]Command{
		"print": func(_ Helper, args ...interface{}) error {
			t.flag = FlagPrint
			t.loop = false
			return nil
		},
		"quit": func(_ Helper, args ...interface{}) error {
			t.loop = false
			return nil
		},
	}
}

func (t *Terminal) ExecuteInTerminal(cmd string) (string, error) {
	tempF, err := os.CreateTemp("", "ftf_")
	if err != nil {
		return "", err
	}
	defer os.Remove(tempF.Name())
	defer tempF.Close()

	fzf := exec.Command("bash", "-c", cmd+" > "+tempF.Name())
	fzf.Stdin = t.in
	fzf.Stdout = t.out
	fzf.Stderr = t.out
	t.revertTerm()
	defer t.initTerm()
	err = fzf.Run()
	if err != nil {
		return "", err
	}
	out, err := io.ReadAll(tempF)
	return string(out), err
}
