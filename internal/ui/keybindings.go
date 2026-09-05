package ui

import (
	"github.com/gdamore/tcell/v2"
)

// KeyAction represents an action to be performed when a key is pressed
type KeyAction struct {
	Description string
	Handler     func() bool
}

// KeyMap maps key events to actions
type KeyMap map[tcell.Key]KeyAction

// DefaultKeyMap returns the default key bindings for the editor
func DefaultKeyMap() KeyMap {
	return KeyMap{
		tcell.KeyF1: {
			Description: "Help",
			Handler:     nil,
		},
		tcell.KeyF2: {
			Description: "Save",
			Handler:     nil,
		},
		tcell.KeyF3: {
			Description: "Open",
			Handler:     nil,
		},
		tcell.KeyF10: {
			Description: "AI",
			Handler:     nil,
		},
		tcell.KeyCtrlQ: {
			Description: "Quit",
			Handler:     nil,
		},
	}
}

// GetKeyDescription returns a description for a key
func GetKeyDescription(key tcell.Key) string {
	switch key {
	case tcell.KeyF1:
		return "F1:Help"
	case tcell.KeyF2:
		return "F2:Save"
	case tcell.KeyF3:
		return "F3:Open"
	case tcell.KeyF10:
		return "F10:AI"
	case tcell.KeyCtrlQ:
		return "^Q:Quit"
	default:
		return ""
	}
}

// GetFunctionKeyHints returns a string with function key hints
func GetFunctionKeyHints() string {
    return " F1:Help F2:Save F3:Open F4:History F7:Search F9:Settings F10:AI ^Q:Quit "
}