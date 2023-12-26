package config

import term "github.com/134130/ftf/pkg/terminal"

type KeyBindings map[string][]string
type Graphics map[string]*term.Graphic

var (
	DefaultKeyBindings = KeyBindings{
		//(&term.Event{Symbol: term.Rune, Value: 'j'}).HashKey(): []string{"tree:next"},
		//(&term.Event{Symbol: term.Rune, Value: 'k'}).HashKey(): []string{"tree:prev"},
		//(&term.Event{Symbol: term.Rune, Value: 'h'}).HashKey(): []string{"tree:parent", "tree:close"},
		//(&term.Event{Symbol: term.Rune, Value: 'l'}).HashKey(): []string{"tree:open", "tree:next"},
		(&term.Event{Symbol: term.Down}).HashKey():             []string{"tree:next"},
		(&term.Event{Symbol: term.Up}).HashKey():               []string{"tree:prev"},
		(&term.Event{Symbol: term.Left}).HashKey():             []string{"tree:parent", "tree:close"},
		(&term.Event{Symbol: term.Right}).HashKey():            []string{"tree:open", "tree:next"},
		(&term.Event{Symbol: term.Rune, Value: ' '}).HashKey(): []string{"tree:selectPath"},
		(&term.Event{Symbol: term.Enter}).HashKey():            []string{"print"},
		(&term.Event{Symbol: term.Rune, Value: 'q'}).HashKey(): []string{"quit"},
		(&term.Event{Symbol: term.Escape}).HashKey():           []string{"quit"},
		(&term.Event{Symbol: term.CtrlC}).HashKey():            []string{"quit"},
	}
	DefaultGraphics = Graphics{
		"tree:cursor": &term.Graphic{
			Reverse: true,
		},
		"tree:selected": &term.Graphic{
			Bold: true,
		},
	}
)
