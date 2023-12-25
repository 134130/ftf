package config

import term "ftf/pkg/terminal"

type KeyBindings map[string][]string
type Graphics map[string]*term.Graphic

var (
	DefaultKeyBindings = KeyBindings{
		(&term.Event{Symbol: term.Rune, Value: 'j'}).HashKey():         []string{"tree:next"},
		(&term.Event{Symbol: term.Rune, Value: 'k'}).HashKey():         []string{"tree:prev"},
		(&term.Event{Symbol: term.Rune, Value: 'h'}).HashKey():         []string{"tree:close"},
		(&term.Event{Symbol: term.Rune, Value: 'l'}).HashKey():         []string{"tree:open"},
		(&term.Event{Symbol: term.Rune, Value: term.Up}).HashKey():     []string{"tree:next"},
		(&term.Event{Symbol: term.Rune, Value: term.Down}).HashKey():   []string{"tree:prev"},
		(&term.Event{Symbol: term.Rune, Value: term.Left}).HashKey():   []string{"tree:close"},
		(&term.Event{Symbol: term.Rune, Value: term.Right}).HashKey():  []string{"tree:open"},
		(&term.Event{Symbol: term.Rune, Value: 'q'}).HashKey():         []string{"quit"},
		(&term.Event{Symbol: term.Rune, Value: term.Escape}).HashKey(): []string{"quit"},
		(&term.Event{Symbol: term.Rune, Value: term.CtrlC}).HashKey():  []string{"quit"},
	}
	DefaultGraphics = Graphics{
		"tree:cursor": &term.Graphic{
			Reverse: true,
		},
	}
)
