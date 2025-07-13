package assessments

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Enter  key.Binding
	Back   key.Binding
	Quit   key.Binding
	New    key.Binding
	Edit   key.Binding
	Delete key.Binding
	Tab    key.Binding
	ShiftTab key.Binding
	Confirm key.Binding
	Cancel key.Binding
}

var DefaultKeyMap = KeyMap{
	Up:     key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "move up")),
	Down:   key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "move down")),
	Enter:  key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select/submit")),
	Back:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	Quit:   key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "quit")),
	New:    key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new")),
	Edit:   key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
	Delete: key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
	Tab: key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next field")),
	ShiftTab: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "previous field")),
	Confirm: key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "confirm")),
	Cancel: key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "cancel")),
}
