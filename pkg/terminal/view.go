package terminal

type ViewRenderer interface {
	Position(int, int) Position
	HasBorder() bool
	ShouldRender() bool
	Render(Position) []LineAppender
	Commands() map[string]Command
}

type Command func(helper Helper, args ...interface{}) error

type Helper interface {
	ExecuteInTerminal(string) (string, error)
}
