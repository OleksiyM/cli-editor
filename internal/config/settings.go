package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Settings represents editor settings
type Settings struct {
    Theme           string            `json:"theme"`
    TabSize         int               `json:"tabSize"`
    UseSpaces       bool              `json:"useSpaces"`
    LineNumbers     bool              `json:"lineNumbers"`
    SyntaxHighlight bool              `json:"syntaxHighlight"`
    SaveHistory     bool              `json:"saveHistory"`
    MaxFileSizeMB   int               `json:"maxFileSizeMB"`  // Maximum file size in MB
    KeyBindings     map[string]string `json:"keyBindings"`
}

// DefaultSettings returns default editor settings
func DefaultSettings() *Settings {
    return &Settings{
        Theme:           "default",
        TabSize:         4,
        UseSpaces:       true,
        LineNumbers:     true,
        SyntaxHighlight: true,
        SaveHistory:     true,
        MaxFileSizeMB:   50,  // Default maximum file size: 50 MB
        KeyBindings: map[string]string{
            "Ctrl+Q": "quit",
            "F1":     "help",
            "F2":     "save",
            "F3":     "open",
            "F4":     "history",
            "F7":     "search",
            "F9":     "settings",
            "F10":    "ai",
        },
    }
}

// LoadSettings loads settings from a file
func LoadSettings(path string) (*Settings, error) {
	// If file doesn't exist, return default settings
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultSettings(), nil
	}
	
	// Read file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	
	// Parse JSON
	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, err
	}
	
	// Ensure MaxFileSizeMB has a valid value
	if settings.MaxFileSizeMB <= 0 {
		settings.MaxFileSizeMB = DefaultSettings().MaxFileSizeMB
	}
	
	return &settings, nil
}

// SaveSettings saves settings to a file
func SaveSettings(settings *Settings, path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	// Marshal to JSON
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	
	// Write to file
	return ioutil.WriteFile(path, data, 0644)
}

// GetSettingsPath returns the path to the settings file
func GetSettingsPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".cli-editor/settings.json"
	}
	
	return filepath.Join(homeDir, ".cli-editor", "settings.json")
}