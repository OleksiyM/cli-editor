package ui

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"cli-editor/internal/config"
	"cli-editor/internal/highlight"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

// Editor represents the main editor component
type Editor struct {
	screen        tcell.Screen
	filename      string
	buffer        []string
	cursorX       int
	cursorY       int
	offsetX       int
	offsetY       int
	width         int
	height        int
	statusMsg     string
	statusMsgTime int64
	modified      bool
	activeDialog  *Dialog
	settings      *config.Settings
	highlighter   *highlight.SyntaxHighlighter
	history       []string // List of recently opened files
}

// NewEditor creates a new editor instance
func NewEditor(filename string) (*Editor, error) {
	// Initialize tcell screen with UTF-8 support
	// Set the locale to UTF-8 before creating the screen
	tcell.SetEncodingFallback(tcell.EncodingFallbackUTF8)
	
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	// Set encoding to UTF-8
	screen.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))

	if err := screen.Init(); err != nil {
		return nil, err
	}

	// Load settings
	settingsPath := config.GetSettingsPath()
	settings, err := config.LoadSettings(settingsPath)
	if err != nil {
		// If there's an error, use default settings
		settings = config.DefaultSettings()
	}

	// Create editor
	editor := &Editor{
		screen:      screen,
		filename:    filename,
		buffer:      []string{""},
		settings:    settings,
		highlighter: highlight.NewSyntaxHighlighter(filename),
	}

	// Load file if provided
	if filename != "" {
		if err := editor.loadFile(filename); err != nil {
			// Instead of returning an error, create an error dialog to show on first draw
			errorDialog := NewDialogWithInput(screen, "Error", []string{
				fmt.Sprintf("Cannot open file: %s", filename),
				fmt.Sprintf("%v", err),
				"",
				"Press Enter to continue",
			}, false)
			errorDialog.SetCallback(func(result string) {
				// Just close the dialog
			})
			editor.activeDialog = errorDialog
		}
	}

	// Load history
	editor.loadHistory()

	// Add current file to history if it's not empty
	if filename != "" && editor.settings.SaveHistory {
		found := false
		for _, file := range editor.history {
			if file == filename {
				found = true
				break
			}
		}
		if !found {
			editor.history = append([]string{filename}, editor.history...)
			// Limit history size
			if len(editor.history) > 20 {
				editor.history = editor.history[:20]
			}
			editor.saveHistory()
		}
	}

	// Get screen dimensions
	editor.width, editor.height = screen.Size()

	return editor, nil
}

// loadFile loads a file into the editor buffer
func (e *Editor) loadFile(filename string) error {
	// Check if file exists first
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filename)
	}

	// Get file info for size check
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return fmt.Errorf("error getting file info: %v", err)
	}

	// Check if file is too large based on settings
	maxSizeBytes := int64(e.settings.MaxFileSizeMB) * 1024 * 1024
	fileSize := fileInfo.Size()
	
	// Convert file size to MB for display
	fileSizeMB := float64(fileSize) / (1024 * 1024)
	
	if fileSize > maxSizeBytes {
		return fmt.Errorf("file too large (%.1f MB), max %d MB", 
			fileSizeMB, e.settings.MaxFileSizeMB)
	}

	// Use ioutil.ReadFile which properly handles UTF-8
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Convert to string and handle line endings
	contentStr := string(content)

	// Split by newlines, handling both Unix and Windows line endings
	lines := strings.Split(strings.ReplaceAll(contentStr, "\r\n", "\n"), "\n")

	// Handle case where file doesn't end with newline
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		// Remove the last empty line if it's just from trailing newline
		lines = lines[:len(lines)-1]
	}

	// Update buffer
	e.buffer = lines
	if len(e.buffer) == 0 {
		e.buffer = []string{""}
	}

	// Update highlighter for the new file
	e.highlighter = highlight.NewSyntaxHighlighter(filename)

	// Reset cursor and scroll position
	e.cursorX = 0
	e.cursorY = 0
	e.offsetX = 0
	e.offsetY = 0

	return nil
}

// Run starts the editor main loop
func (e *Editor) Run() error {
	defer e.screen.Fini()

	for {
		e.draw()

		// Handle events
		ev := e.screen.PollEvent()

		// Debug key events - commented out for normal use
		// if keyEv, ok := ev.(*tcell.EventKey); ok {
		// 	key := keyEv.Key()
		// 	rune := keyEv.Rune()
		// 	mod := keyEv.Modifiers()
		// 	e.statusMsg = fmt.Sprintf("Key: %v, Rune: %c, Mod: %v", key, rune, mod)
		// }

		// If there's an active dialog, let it handle the event first
		if e.activeDialog != nil && e.activeDialog.IsActive() {
			if e.activeDialog.HandleEvent(ev) {
				continue
			}
		}

		// Otherwise, handle the event in the editor
		switch ev := ev.(type) {
		case *tcell.EventResize:
			e.width, e.height = ev.Size()
			e.screen.Sync()
		case *tcell.EventKey:
			// Debug UTF-8 characters - only show for a short time
			if ev.Key() == tcell.KeyRune {
				r := ev.Rune()
				// Check if it's a non-ASCII character
				if r > 127 {
					// This is a UTF-8 character, make sure we're handling it correctly
					e.statusMsg = fmt.Sprintf("UTF-8 char: %c (U+%04X)", r, r)
					// Reset the status message after a short time
					// In a real implementation, you'd want to use a timer
				}
			}
			if e.handleKeyEvent(ev) {
				return nil // Exit requested
			}
		}
	}
}

// draw renders the editor UI
func (e *Editor) draw() {
	e.screen.Clear()

	// Draw title bar with filename
	titleStyle := tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite)
	for x := 0; x < e.width; x++ {
		e.screen.SetContent(x, 0, ' ', nil, titleStyle)
	}

	// Show filename in title bar
	title := fmt.Sprintf(" %s ", e.filename)
	if e.filename == "" {
		title = " [New File] "
	}
	if e.modified {
		title += "[modified]"
	}

	// Draw title with proper character width handling
	xPos := 0
	for _, r := range title {
		if xPos >= e.width {
			break
		}
		e.screen.SetContent(xPos, 0, r, nil, titleStyle)
		// Use runewidth to get correct width of character
		xPos += runewidth.RuneWidth(r)
	}

	// Define line number width
	lineNumWidth := 0
	if e.settings.LineNumbers {
		lineNumWidth = 4 // Width of line numbers column
	}

	// Draw text buffer (adjusted to account for title bar)
	for y := 0; y < e.height-2; y++ {
		fileY := y + e.offsetY
		if fileY >= len(e.buffer) {
			break
		}

		// Draw line numbers if enabled
		if e.settings.LineNumbers {
			lineNumStyle := tcell.StyleDefault.Foreground(tcell.ColorGray)
			lineNum := fmt.Sprintf("%3d ", fileY+1)
			for i, r := range lineNum {
				e.screen.SetContent(i, y+1, r, nil, lineNumStyle)
			}
		}

		// In the draw function, modify the text rendering to use syntax highlighting
		// Draw text content
		line := e.buffer[fileY]
		var styles []tcell.Style
		if e.settings.SyntaxHighlight {
			styles = e.highlighter.HighlightLine(line)
		} else {
			// Create default styles if syntax highlighting is disabled
			styles = make([]tcell.Style, len(line))
			for i := range styles {
				styles[i] = tcell.StyleDefault
			}
		}

		// Use runes instead of bytes for proper UTF-8 handling
		runes := []rune(line)
		xPos := lineNumWidth
		// Iterate over the runes that are actually visible, starting from the horizontal offset
		for i := e.offsetX; i < len(runes); i++ {
			r := runes[i]

			// Stop if we've gone past the right edge of the screen
			if xPos >= e.width {
				break
			}

			// The highlighter may be byte-indexed, so style could be wrong for unicode.
			// But this loop fixes the character rendering issue.
			style := tcell.StyleDefault
			if e.settings.SyntaxHighlight && i < len(styles) {
				style = styles[i]
			}

			// Handle tab characters
			if r == '\t' {
				// Calculate tab width based on settings
				tabSize := e.settings.TabSize
				tabWidth := tabSize - (xPos-lineNumWidth)%tabSize

				// Draw spaces for the tab
				for j := 0; j < tabWidth; j++ {
					if xPos+j >= e.width {
						break
					}
					e.screen.SetContent(xPos+j, y+1, ' ', nil, style)
				}
				xPos += tabWidth
			} else {
				// Normal character - ensure proper handling of Unicode characters
				e.screen.SetContent(xPos, y+1, r, nil, style)
				// Use runewidth package to get correct width of character
				xPos += runewidth.RuneWidth(r)
			}
		}
	}

	// Draw scrollbar if content is larger than screen
	if len(e.buffer) > e.height-2 {
		scrollbarStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue)
		// Ensure division by zero doesn't happen if buffer is empty (though unlikely here)
		scrollbarPos := 0
		if len(e.buffer) > 0 {
			scrollbarPos = (e.offsetY * (e.height - 2)) / len(e.buffer) // Use offsetY for scrollbar position
		}
		scrollbarHeight := ((e.height - 2) * (e.height - 2)) / len(e.buffer)
		if scrollbarHeight < 1 {
			scrollbarHeight = 1
		}

		for y := 1; y < e.height-1; y++ {
			// Draw the scrollbar track using box drawing character
			e.screen.SetContent(e.width-1, y, '│', nil, scrollbarStyle)
		}
		// Draw the scrollbar thumb
		for y := 0; y < scrollbarHeight; y++ {
			thumbY := scrollbarPos + 1 + y
			if thumbY < e.height-1 { // Ensure thumb stays within bounds
				e.screen.SetContent(e.width-1, thumbY, '█', nil, scrollbarStyle)
			}
		}
	}

	// Draw status bar
	statusStyle := tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite)
	for x := 0; x < e.width; x++ {
		e.screen.SetContent(x, e.height-1, ' ', nil, statusStyle)
	}

	// Show cursor position in status bar
	status := fmt.Sprintf(" %d:%d ", e.cursorY+1, e.cursorX+1)

	for i, r := range status {
		if i >= e.width {
			break
		}
		e.screen.SetContent(i, e.height-1, r, nil, statusStyle)
	}

	// Show temporary message if set
	if e.statusMsg != "" {
		msgStyle := statusStyle.Foreground(tcell.ColorYellow)
		msgX := (e.width - runewidth.StringWidth(e.statusMsg)) / 2
		if msgX < 0 {
			msgX = 0
		}

		// Draw status message with proper character width handling
		xPos := msgX
		for _, r := range e.statusMsg {
			if xPos >= e.width {
				break
			}
			e.screen.SetContent(xPos, e.height-1, r, nil, msgStyle)
			// Use runewidth to get correct width of character
			xPos += runewidth.RuneWidth(r)
		}
	}

	// Draw function key hints
	hints := GetFunctionKeyHints()
	hintsPos := e.width - runewidth.StringWidth(hints)
	if hintsPos >= 0 {
		// Draw hints with proper character width handling
		xPos := hintsPos
		for _, r := range hints {
			if xPos >= e.width {
				break
			}
			e.screen.SetContent(xPos, e.height-1, r, nil, statusStyle)
			// Use runewidth to get correct width of character
			xPos += runewidth.RuneWidth(r)
		}
	}

	// Position cursor (adjusted for title bar and line numbers)
	cursorScreenX := e.cursorX - e.offsetX
	if e.settings.LineNumbers {
		cursorScreenX += lineNumWidth
	}
	// Ensure cursor is within visible bounds on screen
	cursorScreenY := e.cursorY - e.offsetY + 1
	if cursorScreenX >= 0 && cursorScreenX < e.width && cursorScreenY >= 1 && cursorScreenY < e.height-1 {
		e.screen.ShowCursor(cursorScreenX, cursorScreenY)
	} else {
		e.screen.HideCursor() // Hide cursor if it's outside the visible area
	}

	// Draw active dialog if any
	if e.activeDialog != nil && e.activeDialog.IsActive() {
		e.activeDialog.Draw()
	}

	e.screen.Show()
}

// handleKeyEvent processes keyboard input
func (e *Editor) handleKeyEvent(ev *tcell.EventKey) bool {
	// Check for function keys first
	switch ev.Key() {
	case tcell.KeyF1:
		e.showHelp()
		return false

	case tcell.KeyF2:
		e.saveFile()
		return false

	case tcell.KeyF3:
		e.openFile()
		return false

	case tcell.KeyF4:
		e.showHistory()
		return false

	case tcell.KeyF7:
		e.showSearch()
		return false

	case tcell.KeyF9:
		e.showSettings()
		return false

	case tcell.KeyF10:
		e.showAIAssistant()
		return false
	case tcell.KeyPgDn:
		// Page down - scroll down by screen height
		if e.cursorY < len(e.buffer)-1 {
			scrollAmount := e.height - 3 // Account for title and status bars
			newY := e.cursorY + scrollAmount
			if newY >= len(e.buffer) {
				newY = len(e.buffer) - 1
			}
			e.cursorY = newY
			if e.cursorX > len([]rune(e.buffer[e.cursorY])) {
				e.cursorX = len([]rune(e.buffer[e.cursorY]))
			}
		}
		return false

	case tcell.KeyPgUp:
		// Page up - scroll up by screen height
		if e.cursorY > 0 {
			scrollAmount := e.height - 3 // Account for title and status bars
			newY := e.cursorY - scrollAmount
			if newY < 0 {
				newY = 0
			}
			e.cursorY = newY
			if e.cursorX > len([]rune(e.buffer[e.cursorY])) {
				e.cursorX = len([]rune(e.buffer[e.cursorY]))
			}
		}
		return false
	}

	// Check for Ctrl key combinations (alternatives to function keys)
	if ev.Modifiers() == tcell.ModCtrl {
		switch ev.Rune() {
		case 'h', 'H': // Ctrl+H for Help (alternative to F1)
			e.showHelp()
			return false

		case 's', 'S': // Ctrl+S for Save (alternative to F2)
			e.saveFile()
			return false

		case 'o', 'O': // Ctrl+O for Open (alternative to F3)
			e.openFile()
			return false

		case 'a', 'A': // Ctrl+A for AI Assistant (alternative to F10)
			e.showAIAssistant()
			return false
		case 'p', 'P': // Ctrl+P for Preferences/Settings
			e.showSettings()
			return false

		}
	}

	// Handle other keys
	switch ev.Key() {
	case tcell.KeyEscape, tcell.KeyCtrlQ:
		return true // Exit

	case tcell.KeyUp:
		if e.cursorY > 0 {
			e.cursorY--
			// Scroll up if cursor moves above visible area
			if e.cursorY < e.offsetY {
				e.offsetY--
			}
			// Adjust cursor X if new line is shorter
			if e.cursorX > len([]rune(e.buffer[e.cursorY])) {
				e.cursorX = len([]rune(e.buffer[e.cursorY]))
			}
		}
	case tcell.KeyDown:
		if e.cursorY < len(e.buffer)-1 {
			e.cursorY++
			// Scroll down if cursor moves below visible area
			if e.cursorY >= e.offsetY+e.height-2 { // -2 for title and status bar
				e.offsetY++
			}
			// Adjust cursor X if new line is shorter
			if e.cursorX > len([]rune(e.buffer[e.cursorY])) {
				e.cursorX = len([]rune(e.buffer[e.cursorY]))
			}
		}
	case tcell.KeyLeft:
		if e.cursorX > 0 {
			// Move cursor left by the width of the previous character
			runes := []rune(e.buffer[e.cursorY])
			if e.cursorX <= len(runes) {
				// Get the character to the left of cursor
				prevChar := runes[e.cursorX-1]
				charWidth := runewidth.RuneWidth(prevChar)
				e.cursorX -= charWidth
			} else {
				// Fallback in case of inconsistency
				e.cursorX--
			}
			// Scroll left if cursor moves left of visible area
			if e.cursorX < e.offsetX {
				e.offsetX = e.cursorX
			}
		} else if e.cursorY > 0 { // Move to end of previous line
			e.cursorY--
			e.cursorX = len([]rune(e.buffer[e.cursorY])) // Use rune length for proper UTF-8 handling
			// Scroll up if necessary
			if e.cursorY < e.offsetY {
				e.offsetY--
			}
			// Scroll right if necessary (to bring end of line into view)
			if e.cursorX >= e.offsetX+e.width {
				e.offsetX = e.cursorX - e.width + 1
			}
			if e.cursorX < e.offsetX { // Handle case where line is shorter than offset
				e.offsetX = e.cursorX
			}
		}
	case tcell.KeyRight:
		runes := []rune(e.buffer[e.cursorY])
		if e.cursorX < len(runes) {
			// Move cursor right by the width of the current character
			currentChar := runes[e.cursorX]
			charWidth := runewidth.RuneWidth(currentChar)
			e.cursorX += charWidth
			// Scroll right if cursor moves beyond visible area
			if e.cursorX >= e.offsetX+e.width {
				e.offsetX = e.cursorX - e.width + 1
			}
		} else if e.cursorY < len(e.buffer)-1 { // Move to start of next line
			e.cursorY++
			e.cursorX = 0
			// Scroll down if necessary
			if e.cursorY >= e.offsetY+e.height-2 {
				e.offsetY++
			}
			// Scroll left if necessary (to bring start of line into view)
			if e.cursorX < e.offsetX {
				e.offsetX = 0
			}
		}

	case tcell.KeyF1:
		e.showHelp()

	case tcell.KeyF2:
		e.saveFile()

	case tcell.KeyF10:
		e.showAIAssistant()

	case tcell.KeyTab:
		// Insert tab character or spaces based on settings
		if e.settings.UseSpaces {
			// Insert spaces
			tabSize := e.settings.TabSize
			// Convert line to runes to properly handle UTF-8
			runes := []rune(e.buffer[e.cursorY])
			// Insert the spaces at cursor position
			newRunes := make([]rune, len(runes)+tabSize)
			copy(newRunes, runes[:e.cursorX])
			for i := 0; i < tabSize; i++ {
				newRunes[e.cursorX+i] = ' '
			}
			copy(newRunes[e.cursorX+tabSize:], runes[e.cursorX:])
			e.buffer[e.cursorY] = string(newRunes)
			e.cursorX += tabSize
		} else {
			// Insert tab character
			runes := []rune(e.buffer[e.cursorY])
			// Insert the tab character at cursor position
			newRunes := make([]rune, len(runes)+1)
			copy(newRunes, runes[:e.cursorX])
			newRunes[e.cursorX] = '\t'
			copy(newRunes[e.cursorX+1:], runes[e.cursorX:])
			e.buffer[e.cursorY] = string(newRunes)
			e.cursorX++ // Tab is a single character
		}
		e.modified = true

	case tcell.KeyEnter:
		// Split line at cursor using runes for proper UTF-8 handling
		runes := []rune(e.buffer[e.cursorY])
		leftRunes := runes[:e.cursorX]
		rightRunes := runes[e.cursorX:]
		e.buffer[e.cursorY] = string(leftRunes)
		e.buffer = append(e.buffer[:e.cursorY+1], append([]string{string(rightRunes)}, e.buffer[e.cursorY+1:]...)...)
		e.cursorY++
		e.cursorX = 0
		e.modified = true
		
		// Ensure cursor stays in view when adding new lines
		if e.cursorY >= e.offsetY+e.height-2 { // -2 for title and status bar
			e.offsetY = e.cursorY - (e.height-3) // Keep cursor in view with one line of context
		}

	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if e.cursorX > 0 {
			// Delete character before cursor, handling runes correctly
			runes := []rune(e.buffer[e.cursorY])
			if e.cursorX <= len(runes) {
				// Get the width of the character we're about to delete
				deletedChar := runes[e.cursorX-1]
				charWidth := runewidth.RuneWidth(deletedChar)
				
				// Delete the character
				runes = append(runes[:e.cursorX-1], runes[e.cursorX:]...)
				e.buffer[e.cursorY] = string(runes)
				
				// Move cursor back by the width of the deleted character
				e.cursorX -= charWidth
				e.modified = true
			}
		} else if e.cursorY > 0 {
			// Join with previous line using runes for proper UTF-8 handling
			prevRunes := []rune(e.buffer[e.cursorY-1])
			currentRunes := []rune(e.buffer[e.cursorY])
			combinedRunes := append(prevRunes, currentRunes...)
			e.buffer[e.cursorY-1] = string(combinedRunes)
			e.buffer = append(e.buffer[:e.cursorY], e.buffer[e.cursorY+1:]...)
			e.cursorY--
			e.cursorX = len(prevRunes) // Set cursor to end of previous line
			e.modified = true
		}

	case tcell.KeyDelete:
		runes := []rune(e.buffer[e.cursorY])
		if e.cursorX < len(runes) {
			// Delete character at cursor
			runes = append(runes[:e.cursorX], runes[e.cursorX+1:]...)
			e.buffer[e.cursorY] = string(runes)
			e.modified = true
			// Note: cursor position doesn't change when deleting at cursor
		} else if e.cursorY < len(e.buffer)-1 {
			// Join with next line
			// Convert both lines to runes for proper UTF-8 handling
			currentRunes := []rune(e.buffer[e.cursorY])
			nextRunes := []rune(e.buffer[e.cursorY+1])
			combinedRunes := append(currentRunes, nextRunes...)
			e.buffer[e.cursorY] = string(combinedRunes)
			e.buffer = append(e.buffer[:e.cursorY+1], e.buffer[e.cursorY+2:]...)
			e.modified = true
		}

	default:
		// Insert character
		if ev.Rune() != 0 {
			r := ev.Rune()
			// Convert line to runes to properly handle UTF-8
			runes := []rune(e.buffer[e.cursorY])
			// Insert the new rune at cursor position
			newRunes := make([]rune, len(runes)+1)
			copy(newRunes, runes[:e.cursorX])
			newRunes[e.cursorX] = r
			copy(newRunes[e.cursorX+1:], runes[e.cursorX:])
			e.buffer[e.cursorY] = string(newRunes)
			// Use runewidth to properly handle UTF-8 character width
			e.cursorX += runewidth.RuneWidth(r)
			e.modified = true
		}
	}

	return false
}

// showHelp displays help information
func (e *Editor) showHelp() {
	content := []string{
		"Function Keys:",
		"F1: Show this help",
		"F2: Save file",
		"F3: Open file",
		"F4: File history",
		"F7: Search",
		"F9: Settings",
		"F10: AI Assistant",
		"",
		"Navigation:",
		"Arrow keys: Move cursor",
		"Ctrl+Q: Quit",
	}

	dialog := NewDialogWithInput(e.screen, "Help", content, false)
	dialog.SetCallback(func(result string) {
		// Just close the dialog
	})

	// Add dialog to active dialogs
	e.activeDialog = dialog
}

// openFile displays a dialog to open a file
func (e *Editor) openFile() {
	dialog := NewDialog(e.screen, "Open File", []string{"Enter file path:"})
	dialog.SetCallback(func(result string) {
		if result != "" {
			if err := e.loadFile(result); err != nil {
				// Show error dialog instead of status message
				errorDialog := NewDialogWithInput(e.screen, "Error", []string{
					fmt.Sprintf("Cannot open file:"),
					fmt.Sprintf("%v", err),
					"",
					"Press Enter to continue",
				}, false)
				errorDialog.SetCallback(func(result string) {
					// Just close the dialog
				})
				e.activeDialog = errorDialog
			} else {
				e.filename = result
				e.modified = false
				e.statusMsg = fmt.Sprintf("Opened %s", result)

				// Add to history if enabled and not already present
				if e.settings.SaveHistory && e.filename != "" {
					found := false
					for _, file := range e.history {
						if file == e.filename {
							found = true
							break
						}
					}
					if !found {
						e.history = append([]string{e.filename}, e.history...)
						// Limit history size
						if len(e.history) > 20 { // Example limit
							e.history = e.history[:20]
						}
						e.saveHistory() // Persist history
					}
				}
			}
		}
	})

	// Add dialog to active dialogs
	e.activeDialog = dialog
}

// showAIAssistant displays the AI assistant dialog
func (e *Editor) showAIAssistant() {
	// Remove or comment out this line if not used
	// currentLine := e.buffer[e.cursorY]

	dialog := NewDialog(e.screen, "AI Assistant", []string{
		"What would you like to do with this code?",
		"1. Explain",
		"2. Optimize",
		"3. Document",
		"",
		"Enter your choice or question:",
	})

	dialog.SetCallback(func(result string) {
		if result != "" {
			// Here you would call your AI service
			// For now, just show a placeholder message
			e.statusMsg = "AI Assistant: Feature coming soon!"
		}
	})

	// Add dialog to active dialogs
	e.activeDialog = dialog
}

// saveFile saves the current buffer to a file
func (e *Editor) saveFile() {
	if e.filename == "" {
		// Prompt for filename if it's empty
		dialog := NewDialog(e.screen, "Save File As", []string{"Enter file path:"})
		dialog.SetCallback(func(result string) {
			if result != "" {
				e.filename = result
				// Update highlighter for new filename/extension
				e.highlighter = highlight.NewSyntaxHighlighter(e.filename)
				e.saveFile() // Call saveFile again now that filename is set
			} else {
				e.statusMsg = "Save cancelled"
			}
		})
		e.activeDialog = dialog
		return
	}

	// Join buffer lines into a single string
	content := strings.Join(e.buffer, "\n")
	// Write content to the file
	// Use os.WriteFile which is simpler than ioutil.WriteFile (deprecated)
	err := os.WriteFile(e.filename, []byte(content), 0644) // Use standard file permissions
	if err != nil {
		e.statusMsg = fmt.Sprintf("Error saving file: %v", err)
		return
	}

	e.modified = false // File is no longer modified
	e.statusMsg = fmt.Sprintf("Saved %s", e.filename)

	// Add to history if enabled and not already present
	if e.settings.SaveHistory && e.filename != "" {
		found := false
		for _, file := range e.history {
			if file == e.filename {
				found = true
				break
			}
		}
		if !found {
			e.history = append([]string{e.filename}, e.history...)
			// Limit history size
			if len(e.history) > 20 { // Example limit
				e.history = e.history[:20]
			}
			e.saveHistory() // Persist history
		}
	}
}

// Add history functions
func (e *Editor) loadHistory() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	historyFile := filepath.Join(homeDir, ".cli-editor-history")
	data, err := ioutil.ReadFile(historyFile)
	if err != nil {
		return
	}

	e.history = strings.Split(string(data), "\n")
	// Remove empty entries
	var cleanHistory []string
	for _, file := range e.history {
		if file != "" {
			cleanHistory = append(cleanHistory, file)
		}
	}
	e.history = cleanHistory
}

func (e *Editor) saveHistory() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	historyFile := filepath.Join(homeDir, ".cli-editor-history")
	content := strings.Join(e.history, "\n")
	ioutil.WriteFile(historyFile, []byte(content), 0644)
}

func (e *Editor) showHistory() {
	if len(e.history) == 0 {
		e.statusMsg = "No history available"
		return
	}

	content := []string{"Recent Files:"}
	for i, file := range e.history {
		content = append(content, fmt.Sprintf("%d. %s", i+1, file))
	}
	content = append(content, "", "Enter number to open:")

	dialog := NewDialog(e.screen, "File History", content)
	dialog.SetCallback(func(result string) {
		if result == "" {
			return
		}

		index, err := strconv.Atoi(result)
		if err != nil || index < 1 || index > len(e.history) {
			e.statusMsg = "Invalid selection"
			return
		}

		filename := e.history[index-1]
		// Check if file exists
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			e.statusMsg = fmt.Sprintf("File not found: %s", filename)
			return
		}

		// Save current file if modified
		if e.modified {
			e.saveFile()
		}

		// Load selected file
		e.loadFile(filename)
		e.filename = filename
		e.cursorX = 0
		e.cursorY = 0
		e.offsetX = 0
		e.offsetY = 0
		e.modified = false
	})

	e.activeDialog = dialog
}

// Add search function
func (e *Editor) showSearch() {
	dialog := NewDialog(e.screen, "Search", []string{"Enter search term:"})
	dialog.SetCallback(func(result string) {
		if result == "" {
			return
		}

		// Start search from current position
		startY := e.cursorY
		startX := e.cursorX + 1 // Start after current position

		found := false
		// Search current line from cursor position
		runes := []rune(e.buffer[startY])
		if startX < len(runes) {
			// Convert search term to runes for proper comparison
			searchRunes := []rune(result)
			// Search in the rest of current line
			for i := startX; i <= len(runes)-len(searchRunes); i++ {
				match := true
				for j, r := range searchRunes {
					if runes[i+j] != r {
						match = false
						break
					}
				}
				if match {
					e.cursorX = i
					found = true
					break
				}
			}
		}

		// If not found in current line, search subsequent lines
		if !found {
			searchRunes := []rune(result)
			for y := startY + 1; y < len(e.buffer); y++ {
				lineRunes := []rune(e.buffer[y])
				// Search entire line
				for i := 0; i <= len(lineRunes)-len(searchRunes); i++ {
					match := true
					for j, r := range searchRunes {
						if lineRunes[i+j] != r {
							match = false
							break
						}
					}
					if match {
						e.cursorY = y
						e.cursorX = i
						found = true
						break
					}
				}
				if found {
					break
				}
			}
		}

		// If not found yet, search from beginning
		if !found {
			searchRunes := []rune(result)
			for y := 0; y <= startY; y++ {
				lineRunes := []rune(e.buffer[y])
				startIdx := 0
				if y == startY {
					startIdx = startX
				}

				// Search from startIdx to end of line
				for i := startIdx; i <= len(lineRunes)-len(searchRunes); i++ {
					match := true
					for j, r := range searchRunes {
						if lineRunes[i+j] != r {
							match = false
							break
						}
					}
					if match {
						e.cursorY = y
						e.cursorX = i
						found = true
						break
					}
				}
				if found {
					break
				}
			}
		}

		if !found {
			// First close the search dialog
			e.activeDialog = nil
			
			// Then show error dialog without input field
			errorDialog := NewDialogWithInput(e.screen, "Search", []string{
				fmt.Sprintf("String not found: %s", result),
			}, false)
			errorDialog.SetCallback(func(result string) {
				// Just close the dialog
			})
			e.activeDialog = errorDialog
		} else {
			// Ensure cursor is visible
			if e.cursorY < e.offsetY {
				e.offsetY = e.cursorY
			} else if e.cursorY >= e.offsetY+e.height-2 {
				e.offsetY = e.cursorY - (e.height - 3)
			}

			if e.cursorX < e.offsetX {
				e.offsetX = e.cursorX
			} else if e.cursorX >= e.offsetX+e.width {
				e.offsetX = e.cursorX - (e.width - 1)
			}
		}
	})

	e.activeDialog = dialog
}

// Add settings functions
func (e *Editor) showSettings() {
	content := []string{
		"Editor Settings:",
		"",
		fmt.Sprintf("1. Tab Size: %d", e.settings.TabSize),
		fmt.Sprintf("2. Use Spaces: %t", e.settings.UseSpaces),
		fmt.Sprintf("3. Line Numbers: %t", e.settings.LineNumbers),
		fmt.Sprintf("4. Syntax Highlighting: %t", e.settings.SyntaxHighlight),
		fmt.Sprintf("5. Save History: %t", e.settings.SaveHistory),
		fmt.Sprintf("6. Max File Size (MB): %d", e.settings.MaxFileSizeMB),
		"",
		"Enter option number to toggle/change:",
	}

	dialog := NewDialog(e.screen, "Settings", content)
	dialog.SetCallback(func(result string) {
		if result == "" {
			return // User cancelled
		}

		option, err := strconv.Atoi(result)
		if err != nil {
			e.statusMsg = "Invalid option number"
			return
		}

		settingsChanged := false
		switch option {
		case 1:
			e.showTabSizeDialog() // This will handle saving itself
			return                // Don't re-save here
		case 2:
			e.settings.UseSpaces = !e.settings.UseSpaces
			e.statusMsg = fmt.Sprintf("Use spaces set to: %t", e.settings.UseSpaces)
			settingsChanged = true
		case 3:
			e.settings.LineNumbers = !e.settings.LineNumbers
			e.statusMsg = fmt.Sprintf("Line numbers set to: %t", e.settings.LineNumbers)
			settingsChanged = true
		case 4:
			e.settings.SyntaxHighlight = !e.settings.SyntaxHighlight
			e.statusMsg = fmt.Sprintf("Syntax highlighting set to: %t", e.settings.SyntaxHighlight)
			settingsChanged = true
		case 5:
			e.settings.SaveHistory = !e.settings.SaveHistory
			e.statusMsg = fmt.Sprintf("Save history set to: %t", e.settings.SaveHistory)
			settingsChanged = true
		case 6:
			e.showMaxFileSizeDialog() // Handle max file size setting
			return                    // Don't re-save here
		default:
			e.statusMsg = "Invalid setting option"
		}

		// Save settings if they were changed by options 2-5
		if settingsChanged {
			settingsPath := config.GetSettingsPath()
			if err := config.SaveSettings(e.settings, settingsPath); err != nil {
				e.statusMsg = fmt.Sprintf("Error saving settings: %v", err)
			}
		}
	})

	e.activeDialog = dialog
}

func (e *Editor) showTabSizeDialog() {
	content := []string{
		fmt.Sprintf("Current tab size: %d", e.settings.TabSize),
		"Enter new tab size (e.g., 2, 4, 8):",
	}

	dialog := NewDialog(e.screen, "Set Tab Size", content)
	dialog.SetCallback(func(result string) {
		if result == "" {
			return // User cancelled
		}

		tabSize, err := strconv.Atoi(result)
		// Add some validation for reasonable tab sizes
		if err != nil || tabSize < 1 || tabSize > 16 {
			e.statusMsg = "Invalid tab size (must be 1-16)"
			return
		}

		e.settings.TabSize = tabSize
		e.statusMsg = fmt.Sprintf("Tab size set to %d", tabSize)

		// Save settings immediately after changing tab size
		settingsPath := config.GetSettingsPath()
		if err := config.SaveSettings(e.settings, settingsPath); err != nil {
			e.statusMsg = fmt.Sprintf("Error saving settings: %v", err)
		}
	})

	e.activeDialog = dialog
}

func (e *Editor) showMaxFileSizeDialog() {
	content := []string{
		fmt.Sprintf("Current max file size: %d MB", e.settings.MaxFileSizeMB),
		"Enter new max file size in MB (e.g., 10, 50, 100):",
	}

	dialog := NewDialog(e.screen, "Set Max File Size", content)
	dialog.SetCallback(func(result string) {
		if result == "" {
			return // User cancelled
		}

		maxSize, err := strconv.Atoi(result)
		// Add validation for reasonable file sizes
		if err != nil || maxSize < 1 || maxSize > 10000 {
			e.statusMsg = "Invalid file size (must be 1-10000 MB)"
			return
		}

		e.settings.MaxFileSizeMB = maxSize
		e.statusMsg = fmt.Sprintf("Max file size set to %d MB", maxSize)

		// Save settings immediately after changing max file size
		settingsPath := config.GetSettingsPath()
		if err := config.SaveSettings(e.settings, settingsPath); err != nil {
			e.statusMsg = fmt.Sprintf("Error saving settings: %v", err)
		}
	})

	e.activeDialog = dialog
}

// ActiveDialog returns the currently active dialog
func (e *Editor) ActiveDialog() *Dialog {
	return e.activeDialog
}
