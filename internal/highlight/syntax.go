package highlight

import (
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
)

// SyntaxHighlighter provides syntax highlighting for code
type SyntaxHighlighter struct {
	language string
	rules    map[string]tcell.Style
}

// Define defaultStyle at package level
var defaultStyle = tcell.StyleDefault

// NewSyntaxHighlighter creates a new syntax highlighter for the given filename
func NewSyntaxHighlighter(filename string) *SyntaxHighlighter {
    // Default to text if no filename or extension
    language := "text"
    
    if filename != "" {
        ext := strings.ToLower(filepath.Ext(filename))
        
        // Map file extensions to language types
        switch ext {
        case ".go":
            language = "go"
        case ".js", ".jsx", ".ts", ".tsx":
            language = "javascript"
        case ".py":
            language = "python"
        case ".md", ".markdown":
            language = "markdown"
        case ".sh", ".bash":
            language = "shell"
        }
    }
    
    return &SyntaxHighlighter{
        language: language,
        rules:    getHighlightRules(language),
    }
}

// HighlightLine applies syntax highlighting to a line of text
func (h *SyntaxHighlighter) HighlightLine(line string) []tcell.Style {
    // Create a slice of styles for each character in the line
    styles := make([]tcell.Style, len(line))
    
    // Fill with default style
    for i := range styles {
        styles[i] = defaultStyle
    }
    
    // Skip highlighting for empty lines or if language is text
    if len(line) == 0 || h.language == "text" {
        return styles
    }
    
    // Add a safety check to prevent hanging on very long lines
    if len(line) > 10000 {
        return styles // Skip highlighting for extremely long lines
    }
    
    // Apply language-specific highlighting with error handling
    defer func() {
        if r := recover(); r != nil {
            // If highlighting panics, return default styles
            for i := range styles {
                styles[i] = defaultStyle
            }
        }
    }()
    
    // Apply language-specific highlighting
    switch h.language {
		/*
    case "go":
        h.highlightGo(line, styles)
    case "javascript":
        h.highlightJavaScript(line, styles)
    case "python":
        h.highlightPython(line, styles)
	case "shell":
        h.highlightShell(line, styles)
	*/		
    case "markdown":
        h.highlightMarkdown(line, styles)
    default:
		h.language = "text"
        
    }
    
    return styles
}

// highlightGo applies Go syntax highlighting
func (h *SyntaxHighlighter) highlightGo(line string, styles []tcell.Style) {
	// Keywords
	keywords := []string{
		"package", "import", "func", "return", "var", "const", "type",
		"struct", "interface", "map", "chan", "go", "defer", "if",
		"else", "switch", "case", "default", "for", "range", "break",
		"continue", "select",
	}

	// Apply keyword highlighting
	for _, keyword := range keywords {
		h.highlightKeyword(line, styles, keyword, h.rules["keyword"])
	}

	// Highlight strings
	inString := false
	stringStart := 0

	for i, c := range line {
		if c == '"' && (i == 0 || line[i-1] != '\\') {
			if inString {
				// End of string
				for j := stringStart; j <= i; j++ {
					styles[j] = h.rules["string"]
				}
				inString = false
			} else {
				// Start of string
				inString = true
				stringStart = i
			}
		}
	}

	// Highlight comments
	commentIdx := strings.Index(line, "//")
	if commentIdx >= 0 {
		for i := commentIdx; i < len(line); i++ {
			styles[i] = h.rules["comment"]
		}
	}
}

// highlightJavaScript applies JavaScript syntax highlighting
func (h *SyntaxHighlighter) highlightJavaScript(line string, styles []tcell.Style) {
	// Similar to Go highlighting but with JavaScript keywords
	keywords := []string{
		"var", "let", "const", "function", "return", "if", "else",
		"switch", "case", "default", "for", "while", "do", "break",
		"continue", "try", "catch", "finally", "throw", "class",
		"extends", "new", "this", "super", "import", "export",
	}

	for _, keyword := range keywords {
		h.highlightKeyword(line, styles, keyword, h.rules["keyword"])
	}

	// Highlight strings and comments (similar to Go)
	// ... implementation similar to Go ...
}

// highlightPython applies Python syntax highlighting
func (h *SyntaxHighlighter) highlightPython(line string, styles []tcell.Style) {
	// Python keywords
	keywords := []string{
		"def", "class", "import", "from", "as", "return", "if", "elif",
		"else", "for", "while", "break", "continue", "try", "except",
		"finally", "with", "in", "is", "not", "and", "or", "lambda",
		"None", "True", "False",
	}

	for _, keyword := range keywords {
		h.highlightKeyword(line, styles, keyword, h.rules["keyword"])
	}

	// Highlight strings and comments
	// ... implementation similar to Go ...
}

// highlightKeyword highlights a specific keyword in the line
func (h *SyntaxHighlighter) highlightKeyword(line string, styles []tcell.Style, keyword string, style tcell.Style) {
	lower := strings.ToLower(line)
	keywordLower := strings.ToLower(keyword)

	idx := 0
	for {
		idx = strings.Index(lower[idx:], keywordLower)
		if idx < 0 {
			break
		}

		// Check if it's a whole word
		startIsWordBoundary := idx == 0 || !isAlphaNumeric(rune(lower[idx-1]))
		endIdx := idx + len(keywordLower)
		endIsWordBoundary := endIdx >= len(lower) || !isAlphaNumeric(rune(lower[endIdx]))

		if startIsWordBoundary && endIsWordBoundary {
			for i := idx; i < idx+len(keyword); i++ {
				styles[i] = style
			}
		}

		idx += len(keywordLower)
		if idx >= len(lower) {
			break
		}
	}
}

// isAlphaNumeric checks if a character is alphanumeric or underscore
func isAlphaNumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
}

// getHighlightRules returns the highlighting rules for a language
func getHighlightRules(language string) map[string]tcell.Style {
    rules := map[string]tcell.Style{
        "keyword":  defaultStyle.Foreground(tcell.ColorYellow),
        "string":   defaultStyle.Foreground(tcell.ColorGreen),
        "comment":  defaultStyle.Foreground(tcell.ColorGray),
        "number":   defaultStyle.Foreground(tcell.ColorRed),
        "function": defaultStyle.Foreground(tcell.ColorBlue),
        "type":     defaultStyle.Foreground(tcell.ColorPurple),
    }

    return rules
}

// highlightMarkdown applies Markdown syntax highlighting
func (h *SyntaxHighlighter) highlightMarkdown(line string, styles []tcell.Style) {
    // Highlight headers
    if strings.HasPrefix(line, "#") {
        headerEnd := 0
        for i, c := range line {
            if c == '#' {
                styles[i] = h.rules["keyword"]
                headerEnd = i + 1
            } else {
                break
            }
        }
        
        // Highlight the rest of the header line
        if headerEnd > 0 && headerEnd < len(line) {
            for i := headerEnd; i < len(line); i++ {
                styles[i] = h.rules["function"]
            }
        }
    }
    
    // Highlight bold text
    h.highlightPairs(line, styles, "**", "**", h.rules["keyword"])
    h.highlightPairs(line, styles, "__", "__", h.rules["keyword"])
    
    // Highlight italic text
    h.highlightPairs(line, styles, "*", "*", h.rules["string"])
    h.highlightPairs(line, styles, "_", "_", h.rules["string"])
    
    // Highlight code blocks
    h.highlightPairs(line, styles, "`", "`", h.rules["number"])
    
    // Highlight links
    h.highlightPairs(line, styles, "[", "]", h.rules["function"])
    h.highlightPairs(line, styles, "(", ")", h.rules["comment"])
}

// highlightShell applies Shell script syntax highlighting
func (h *SyntaxHighlighter) highlightShell(line string, styles []tcell.Style) {
    // Shell keywords
    keywords := []string{
        "if", "then", "else", "elif", "fi", "case", "esac", "for", "while", 
        "do", "done", "in", "function", "return", "exit", "export", "source",
        "alias", "unalias", "echo", "read", "set", "unset", "shift",
    }
    
    for _, keyword := range keywords {
        h.highlightKeyword(line, styles, keyword, h.rules["keyword"])
    }
    
    // Highlight variables
    for i := 0; i < len(line); i++ {
        if line[i] == '$' && i+1 < len(line) {
            // Variable reference
            if line[i+1] == '{' {
                // ${VAR} format
                start := i
                for j := i + 2; j < len(line); j++ {
                    if line[j] == '}' {
                        for k := start; k <= j; k++ {
                            styles[k] = h.rules["type"]
                        }
                        i = j
                        break
                    }
                }
            } else {
                // $VAR format
                start := i
                for j := i + 1; j < len(line); j++ {
                    if !isAlphaNumeric(rune(line[j])) {
                        break
                    }
                    styles[j] = h.rules["type"]
                }
                styles[start] = h.rules["type"]
            }
        }
    }
    
    // Highlight comments
    commentIdx := strings.Index(line, "#")
    if commentIdx >= 0 {
        for i := commentIdx; i < len(line); i++ {
            styles[i] = h.rules["comment"]
        }
    }
    
    // Highlight strings
    h.highlightPairs(line, styles, "\"", "\"", h.rules["string"])
    h.highlightPairs(line, styles, "'", "'", h.rules["string"])
}

// highlightPairs highlights text between start and end markers
func (h *SyntaxHighlighter) highlightPairs(line string, styles []tcell.Style, start, end string, style tcell.Style) {
    startLen := len(start)
    endLen := len(end)
    
    pos := 0
    for {
        // Find start marker
        startPos := strings.Index(line[pos:], start)
        if startPos < 0 {
            break
        }
        startPos += pos
        
        // Find end marker
        endPos := strings.Index(line[startPos+startLen:], end)
        if endPos < 0 {
            break
        }
        endPos += startPos + startLen
        
        // Apply style to everything between markers (including markers)
        for i := startPos; i < endPos+endLen; i++ {
            styles[i] = style
        }
        
        pos = endPos + endLen
        if pos >= len(line) {
            break
        }
    }
}
