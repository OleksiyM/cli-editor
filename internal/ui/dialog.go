package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

// Dialog represents a modal dialog
type Dialog struct {
	screen     tcell.Screen
	title      string
	content    []string
	width      int
	height     int
	x          int
	y          int
	active     bool
	hasInput   bool
	inputField string
	inputPos   int
	callback   func(result string)
}

// NewDialog creates a new dialog
func NewDialog(screen tcell.Screen, title string, content []string) *Dialog {
	return NewDialogWithInput(screen, title, content, true)
}

// NewDialogWithInput creates a new dialog with optional input field
func NewDialogWithInput(screen tcell.Screen, title string, content []string, hasInput bool) *Dialog {
	screenWidth, screenHeight := screen.Size()

	// Calculate dialog dimensions
	width := 60
	if width > screenWidth-4 {
		width = screenWidth - 4
	}

	// Calculate height with the new structure:
	// Title + content + empty line + input field (if applicable) + empty line + buttons + empty line + borders
	height := len(content) + 7 // Base structure: content lines + 6 fixed lines (title, empty, empty, empty, buttons, empty)
	if hasInput {
		height += 1 // Add space for input field
	}
	if height > screenHeight-4 {
		height = screenHeight - 4
	}

	// Center dialog
	x := (screenWidth - width) / 2
	y := (screenHeight - height) / 2

	return &Dialog{
		screen:   screen,
		title:    title,
		content:  content,
		width:    width,
		height:   height,
		x:        x,
		y:        y,
		active:   true,
		hasInput: hasInput,
	}
}

// Draw renders the dialog
func (d *Dialog) Draw() {
	if !d.active {
		return
	}

	style := tcell.StyleDefault
	borderStyle := style.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue)
	titleStyle := style.Foreground(tcell.ColorYellow).Background(tcell.ColorBlue)
	contentStyle := style.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue)

	// Draw border with box-drawing characters
	for y := d.y; y < d.y+d.height; y++ {
		for x := d.x; x < d.x+d.width; x++ {
			if y == d.y || y == d.y+d.height-1 {
				if x == d.x && y == d.y {
					d.screen.SetContent(x, y, '┌', nil, borderStyle) // Top-left corner
				} else if x == d.x+d.width-1 && y == d.y {
					d.screen.SetContent(x, y, '┐', nil, borderStyle) // Top-right corner
				} else if x == d.x && y == d.y+d.height-1 {
					d.screen.SetContent(x, y, '└', nil, borderStyle) // Bottom-left corner
				} else if x == d.x+d.width-1 && y == d.y+d.height-1 {
					d.screen.SetContent(x, y, '┘', nil, borderStyle) // Bottom-right corner
				} else {
					d.screen.SetContent(x, y, '─', nil, borderStyle) // Horizontal border
				}
			} else if x == d.x || x == d.x+d.width-1 {
				d.screen.SetContent(x, y, '│', nil, borderStyle) // Vertical border
			} else {
				d.screen.SetContent(x, y, ' ', nil, contentStyle) // Dialog background
			}
		}
	}

	// Draw title
	titleX := d.x + (d.width-runewidth.StringWidth(d.title))/2
	for i, r := range d.title {
		d.screen.SetContent(titleX+i, d.y, r, nil, titleStyle)
	}

	// Calculate content area with new structure:
	// Text content -> Empty line -> Input field (if applicable) -> Empty line -> Buttons -> Empty line
	contentStartY := d.y + 2
	buttonsY := d.y + d.height - 4  // Position buttons correctly
	inputY := buttonsY - 2         // Input field above buttons with empty line
	
	// Draw content (text content)
	for i, line := range d.content {
		// Draw line with proper character width handling
		xPos := d.x + 2
		for _, r := range line {
			if xPos >= d.x+d.width-2 {
				break
			}
			d.screen.SetContent(xPos, contentStartY+i, r, nil, contentStyle)
			// Use runewidth to get correct width of character
			xPos += runewidth.RuneWidth(r)
		}
	}

	// Draw input field if needed
	if d.hasInput {
		inputFieldStyle := style.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
		
		for x := d.x + 2; x < d.x+d.width-2; x++ {
			d.screen.SetContent(x, inputY, ' ', nil, inputFieldStyle)
		}

		// Draw input text with proper character width handling
		xPos := d.x + 2
		for _, r := range d.inputField {
			if xPos >= d.x+d.width-2 {
				break
			}
			d.screen.SetContent(xPos, inputY, r, nil, inputFieldStyle)
			// Use runewidth to get correct width of character
			xPos += runewidth.RuneWidth(r)
		}

		// Position cursor in input field
		// Make sure cursor doesn't go beyond the input field bounds
		cursorX := d.x + 2
		// Calculate cursor position based on input text
		textRunes := []rune(d.inputField)
		for i := 0; i < d.inputPos && i < len(textRunes); i++ {
			cursorX += runewidth.RuneWidth(textRunes[i])
		}
		maxCursorX := d.x + d.width - 3
		if cursorX > maxCursorX {
			cursorX = maxCursorX
		}
		d.screen.ShowCursor(cursorX, inputY)
	}

	// Draw buttons with better styling
	// Create better-looking buttons
	okButtonText := " OK "
	cancelButtonText := " Cancel "
	
	// Calculate button positions
	okButtonWidth := runewidth.StringWidth(okButtonText) + 2  // Add 2 for left/right borders
	cancelButtonWidth := runewidth.StringWidth(cancelButtonText) + 2
	
	spacing := 4
	totalWidth := okButtonWidth + cancelButtonWidth + spacing
	startX := d.x + (d.width - totalWidth) / 2
	
	okButtonX := startX
	cancelButtonX := okButtonX + okButtonWidth + spacing
	
	// Draw OK button
	okButtonStyle := style.Background(tcell.ColorGreen).Foreground(tcell.ColorBlack)
	
	// Left border
	d.screen.SetContent(okButtonX, buttonsY, '[', nil, okButtonStyle)
	
	// Button text
	textX := okButtonX + 1
	for _, r := range okButtonText {
		d.screen.SetContent(textX, buttonsY, r, nil, okButtonStyle)
		textX += runewidth.RuneWidth(r)
	}
	
	// Right border
	d.screen.SetContent(okButtonX+okButtonWidth-1, buttonsY, ']', nil, okButtonStyle)
	
	// Draw Cancel button
	cancelButtonStyle := style.Background(tcell.ColorRed).Foreground(tcell.ColorWhite)
	
	// Left border
	d.screen.SetContent(cancelButtonX, buttonsY, '[', nil, cancelButtonStyle)
	
	// Button text
	textX = cancelButtonX + 1
	for _, r := range cancelButtonText {
		d.screen.SetContent(textX, buttonsY, r, nil, cancelButtonStyle)
		textX += runewidth.RuneWidth(r)
	}
	
	// Right border
	d.screen.SetContent(cancelButtonX+cancelButtonWidth-1, buttonsY, ']', nil, cancelButtonStyle)
}

// HandleEvent processes events for the dialog
func (d *Dialog) HandleEvent(ev tcell.Event) bool {
	if !d.active {
		return false
	}

	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape:
			d.active = false
			if d.callback != nil {
				d.callback("") // Cancel returns empty string
			}
			return true

		case tcell.KeyEnter:
			d.active = false
			if d.callback != nil {
				d.callback(d.inputField)
			}
			return true

		case tcell.KeyTab:
			// In dialogs without input, Tab could be used to navigate buttons
			// For now, we'll keep the simple behavior
			return true

		case tcell.KeyBackspace, tcell.KeyBackspace2:
			if d.hasInput && d.inputPos > 0 {
				// Convert to runes for proper UTF-8 handling
				runes := []rune(d.inputField)
				if d.inputPos <= len(runes) {
					// Get the width of the character we're about to delete
					deletedChar := runes[d.inputPos-1]
					charWidth := runewidth.RuneWidth(deletedChar)
					
					// Delete the character
					runes = append(runes[:d.inputPos-1], runes[d.inputPos:]...)
					d.inputField = string(runes)
					
					// Move cursor back by the width of the deleted character
					d.inputPos -= charWidth
				}
			}
			return true

		case tcell.KeyDelete:
			if d.hasInput {
				// Convert to runes for proper UTF-8 handling
				runes := []rune(d.inputField)
				if d.inputPos < len(runes) {
					// Delete the character
					runes = append(runes[:d.inputPos], runes[d.inputPos+1:]...)
					d.inputField = string(runes)
				}
			}
			return true

		case tcell.KeyLeft:
			if d.hasInput && d.inputPos > 0 {
				// Convert to runes for proper UTF-8 handling
				runes := []rune(d.inputField)
				if d.inputPos <= len(runes) {
					// Move cursor left by the width of the previous character
					prevChar := runes[d.inputPos-1]
					charWidth := runewidth.RuneWidth(prevChar)
					d.inputPos -= charWidth
				} else {
					// Fallback in case of inconsistency
					d.inputPos--
				}
			}
			return true

		case tcell.KeyRight:
			if d.hasInput {
				// Convert to runes for proper UTF-8 handling
				runes := []rune(d.inputField)
				if d.inputPos < len(runes) {
					// Move cursor right by the width of the current character
					currentChar := runes[d.inputPos]
					charWidth := runewidth.RuneWidth(currentChar)
					d.inputPos += charWidth
				}
			}
			return true

		default:
			if d.hasInput && ev.Rune() != 0 {
				r := ev.Rune()
				// Convert to runes for proper UTF-8 handling
				runes := []rune(d.inputField)
				// Insert the new rune at cursor position
				newRunes := make([]rune, len(runes)+1)
				copy(newRunes, runes[:d.inputPos])
				newRunes[d.inputPos] = r
				copy(newRunes[d.inputPos+1:], runes[d.inputPos:])
				d.inputField = string(newRunes)
				// Use runewidth to properly handle UTF-8 character width
				d.inputPos += runewidth.RuneWidth(r)
				return true
			}
		}
	}

	return false
}

// SetCallback sets the callback function for when the dialog is closed
func (d *Dialog) SetCallback(callback func(result string)) {
	d.callback = callback
}

// IsActive returns whether the dialog is active
func (d *Dialog) IsActive() bool {
	return d.active
}
