package ui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

// StatusBar represents the status bar at the bottom of the editor
type StatusBar struct {
	screen     tcell.Screen
	width      int
	message    string
	messageTime time.Time
	messageDuration time.Duration
}

// NewStatusBar creates a new status bar
func NewStatusBar(screen tcell.Screen) *StatusBar {
	width, _ := screen.Size()
	return &StatusBar{
		screen:     screen,
		width:      width,
		messageDuration: 5 * time.Second,
	}
}

// Draw renders the status bar
func (s *StatusBar) Draw(filename string, modified bool, cursorX, cursorY int) {
	width, _ := s.screen.Size()
	s.width = width
	
	// Status bar style
	style := tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite)
	
	// Clear status bar
	for x := 0; x < s.width; x++ {
		s.screen.SetContent(x, s.width-1, ' ', nil, style)
	}
	
	// Show filename and position
	status := fmt.Sprintf(" %s - %d:%d ", filename, cursorY+1, cursorX+1)
	if modified {
		status += "[modified]"
	}
	
	for i, r := range status {
		if i >= s.width {
			break
		}
		s.screen.SetContent(i, s.width-1, r, nil, style)
	}
	
	// Show temporary message if set
	if s.message != "" && time.Since(s.messageTime) < s.messageDuration {
		msgStyle := style.Foreground(tcell.ColorYellow)
		msgX := (s.width - runewidth.StringWidth(s.message)) / 2
		if msgX < 0 {
			msgX = 0
		}
		
		for i, r := range s.message {
			if msgX+i >= s.width {
				break
			}
			s.screen.SetContent(msgX+i, s.width-1, r, nil, msgStyle)
		}
	} else {
		s.message = ""
	}
	
	// Draw function key hints
	hints := GetFunctionKeyHints()
	hintsPos := s.width - runewidth.StringWidth(hints)
	if hintsPos >= 0 {
		for i, r := range hints {
			s.screen.SetContent(hintsPos+i, s.width-1, r, nil, style)
		}
	}
}

// SetMessage sets a temporary message to display in the status bar
func (s *StatusBar) SetMessage(message string) {
	s.message = message
	s.messageTime = time.Now()
}

// Resize updates the status bar dimensions
func (s *StatusBar) Resize() {
	width, _ := s.screen.Size()
	s.width = width
}

// Remove this function as it's already defined in keybindings.go
// GetFunctionKeyHints returns a string with function key hints
// func GetFunctionKeyHints() string {
//     return " F1/^H:Help F2/^S:Save F3/^O:Open F10/^A:AI ^Q:Quit "
// }