package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines the key bindings for the application
type KeyMap struct {
	// Navigation
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	Home     key.Binding
	End      key.Binding

	// Actions
	Enter    key.Binding
	Back     key.Binding
	Quit     key.Binding
	Help     key.Binding
	Refresh  key.Binding
	Search   key.Binding
	Filter   key.Binding
	Open     key.Binding
	Copy     key.Binding

	// View switching
	Dashboard key.Binding
	Problems  key.Binding
	Practice  key.Binding
	Stats     key.Binding

	// Tab navigation
	NextTab key.Binding
	PrevTab key.Binding

	// Problem-specific
	Random    key.Binding
	MarkSolved key.Binding
	Skip      key.Binding
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		// Navigation
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "right"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "ctrl+u"),
			key.WithHelp("pgup", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown", "ctrl+d"),
			key.WithHelp("pgdn", "page down"),
		),
		Home: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("home/g", "go to start"),
		),
		End: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("end/G", "go to end"),
		),

		// Actions
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc", "backspace"),
			key.WithHelp("esc", "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r", "ctrl+r"),
			key.WithHelp("r", "refresh"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		Filter: key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "filter"),
		),
		Open: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "open in browser"),
		),
		Copy: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "copy link"),
		),

		// View switching
		Dashboard: key.NewBinding(
			key.WithKeys("1", "d"),
			key.WithHelp("1/d", "dashboard"),
		),
		Problems: key.NewBinding(
			key.WithKeys("2", "p"),
			key.WithHelp("2/p", "problems"),
		),
		Practice: key.NewBinding(
			key.WithKeys("3", "P"),
			key.WithHelp("3/P", "practice"),
		),
		Stats: key.NewBinding(
			key.WithKeys("4", "s"),
			key.WithHelp("4/s", "stats"),
		),

		// Tab navigation
		NextTab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next tab"),
		),
		PrevTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev tab"),
		),

		// Problem-specific
		Random: key.NewBinding(
			key.WithKeys("R"),
			key.WithHelp("R", "random problem"),
		),
		MarkSolved: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "mark solved"),
		),
		Skip: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "skip/next"),
		),
	}
}

// ShortHelp returns a short help view
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Help,
		k.Quit,
		k.Dashboard,
		k.Problems,
		k.Practice,
		k.Stats,
	}
}

// FullHelp returns the full help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		// Navigation
		{k.Up, k.Down, k.Left, k.Right, k.PageUp, k.PageDown},
		// Actions
		{k.Enter, k.Back, k.Refresh, k.Search, k.Filter},
		// Views
		{k.Dashboard, k.Problems, k.Practice, k.Stats},
		// Problem actions
		{k.Open, k.Copy, k.Random, k.MarkSolved, k.Skip},
		// General
		{k.Help, k.Quit},
	}
}

// NavigationHelp returns navigation-focused help bindings
func (k KeyMap) NavigationHelp() []key.Binding {
	return []key.Binding{
		k.Up, k.Down, k.Enter, k.Back,
	}
}

// ViewHelp returns view switching help bindings
func (k KeyMap) ViewHelp() []key.Binding {
	return []key.Binding{
		k.Dashboard, k.Problems, k.Practice, k.Stats,
	}
}
