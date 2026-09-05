# CLI Editor - Project Context

## Project Overview

This is a terminal-based text editor written in Go with the following key features:

- Cross-platform terminal interface using tcell library
- UTF-8 support for international characters
- Modal dialogs for user interactions
- Syntax highlighting for code files
- File history tracking
- Search functionality
- Configurable settings (tab size, line numbers, etc.)
- AI assistant integration (planned)

## Project Structure

```
cli-editor/
├── cmd/
│   └── editor/
│       └── main.go          # Entry point
├── internal/
│   ├── ai/
│   │   └── assistant.go     # AI assistant functionality
│   ├── config/
│   │   └── settings.go      # Configuration management
│   ├── highlight/
│   │   └── syntax.go        # Syntax highlighting
│   └── ui/
│       ├── editor.go        # Main editor logic
│       ├── dialog.go        # Dialog UI components
│       ├── keybindings.go   # Keyboard shortcuts
│       └── statusbar.go     # Status bar management
├── pkg/
│   └── utils/
│       └── file.go          # Utility functions
├── go.mod                   # Go module definition
└── go.sum                   # Go module checksums
```

## Key Technologies

- **Go** (1.19+) - Main programming language
- **tcell/v2** - Terminal UI library for cross-platform support
- **go-runewidth** - Unicode character width handling

## Building and Running

### Prerequisites
- Go 1.19 or later

### Build
```bash
go build -o editor cmd/editor/main.go
```

### Run
```bash
# Start with empty editor
./editor

# Open a specific file
./editor filename.txt
```

## Core Features

### Editor Functionality
- Basic text editing (insert, delete, navigation)
- UTF-8 character support
- Tab handling (configurable spaces or tab characters)
- Line numbers (toggleable)
- Syntax highlighting (toggleable)
- File history tracking
- Large file protection (default 50MB limit)
- Page up/down navigation

### UI Components
- Modal dialogs for user interactions with consistent structure:
  - Text content
  - Empty line
  - Input field (if applicable)
  - Empty line
  - Buttons
  - Empty line
- Status bar with cursor position
- Function key shortcuts (F1-F10)
- Ctrl key alternatives for function keys
- Scrollbar for large files

### Configuration
Settings are stored in `~/.cli-editor/settings.json`:
- Tab size (default: 4)
- Use spaces instead of tabs (default: true)
- Line numbers (default: true)
- Syntax highlighting (default: true)
- Save file history (default: true)
- Max file size limit (default: 50MB)

## Development Conventions

### Code Structure
- Follow standard Go project layout
- Use internal packages for project-specific code
- Separate UI logic from business logic
- Use interfaces for decoupling components

### UI Patterns
- Modal dialogs for user input and error messages
- Function key shortcuts for main actions
- Status bar for feedback messages
- Consistent error handling with user-friendly messages

### Error Handling
- Use dialogs for user-facing errors instead of console output
- Graceful degradation when features aren't available
- Clear error messages with actionable information

## Key Components

### Editor (internal/ui/editor.go)
Main editor component handling:
- File loading/saving
- Text buffer management
- Cursor navigation
- UI rendering
- Event handling

### Dialog (internal/ui/dialog.go)
Modal dialog system supporting:
- Informational messages
- User input prompts
- Error displays
- Confirmation dialogs
- Consistent layout with proper spacing

### Configuration (internal/config/settings.go)
Settings management:
- Default values
- JSON serialization
- File persistence
- Validation

## Function Key Shortcuts

- F1: Help
- F2: Save file
- F3: Open file
- F4: File history
- F7: Search
- F9: Settings
- F10: AI Assistant

### Ctrl Key Alternatives
- Ctrl+H: Help
- Ctrl+S: Save
- Ctrl+O: Open
- Ctrl+P: Settings
- Ctrl+A: AI Assistant
- Ctrl+Q: Quit

## Recent Changes

### Dialog Layout Improvements
- Updated dialog structure to follow consistent pattern:
  - Text content
  - Empty line
  - Input field (if applicable)
  - Empty line
  - Buttons
  - Empty line
- Fixed content rendering to ensure all text is properly displayed
- Improved height calculation for dialogs with multiple content lines

### Navigation Enhancements
- Added Page Up/Page Down support for faster document navigation
- Improved cursor positioning when navigating between lines of different lengths

## Testing

Currently, the project doesn't have a comprehensive test suite. Tests would need to:
- Mock tcell.Screen for UI testing
- Test file operations with temporary files
- Verify configuration loading/saving
- Validate syntax highlighting rules

To run existing tests:
```bash
go test ./...
```

## Future Improvements

### 1. Enhanced Search Functionality
**Description**: Implement more advanced search features like:
- Find and replace
- Regular expression support
- Search history
- Case-sensitive/insensitive options

**Complexity**: Medium
**Resource Requirements**: 2-3 days development time
**Benefits**: Significantly improves text editing capabilities

### 2. Multiple Buffer Support
**Description**: Allow opening and switching between multiple files:
- Tabbed interface or buffer list
- Buffer switching shortcuts
- Session management (save/restore open files)

**Complexity**: High
**Resource Requirements**: 1-2 weeks development time
**Benefits**: Enables working on multiple files simultaneously

### 3. Plugin/Extension System
**Description**: Allow extending functionality through plugins:
- Plugin loading mechanism
- Standardized API for plugins
- Plugin repository or marketplace

**Complexity**: High
**Resource Requirements**: 2-3 weeks development time
**Benefits**: Extensibility without bloating core application

### 4. Improved AI Assistant
**Description**: Fully implement the AI assistant functionality:
- Integration with actual AI services (OpenAI, etc.)
- Context-aware suggestions
- Code completion and refactoring assistance

**Complexity**: Medium-High
**Resource Requirements**: 1-2 weeks development time
**Benefits**: Modern AI-powered development assistance

### 5. Macro Recording
**Description**: Allow users to record and replay sequences of actions:
- Start/stop recording functionality
- Macro storage and management
- Playback with customization options

**Complexity**: Medium
**Resource Requirements**: 3-5 days development time
**Benefits**: Increases productivity for repetitive tasks

### 6. Enhanced Syntax Highlighting
**Description**: Improve syntax highlighting capabilities:
- More language support
- Customizable color themes
- Better parsing for complex languages

**Complexity**: Medium
**Resource Requirements**: 1-2 weeks development time
**Benefits**: Better code readability and user experience

### 7. Undo/Redo System
**Description**: Implement comprehensive undo/redo functionality:
- Multi-level undo/redo stack
- Grouped operations (e.g., multi-line deletes as single operation)
- Visual indication of modification status

**Complexity**: Medium
**Resource Requirements**: 3-5 days development time
**Benefits**: Professional text editing capabilities

### 8. Terminal Integration
**Description**: Allow executing terminal commands from within the editor:
- Built-in terminal panel
- Command history
- Integration with editor buffer (pipe text to/from commands)

**Complexity**: High
**Resource Requirements**: 1-2 weeks development time
**Benefits**: IDE-like integrated development environment

### 9. Performance Optimizations
**Description**: Optimize editor performance for large files:
- Virtual scrolling for extremely large files
- Incremental parsing for syntax highlighting
- Memory usage improvements

**Complexity**: Medium
**Resource Requirements**: 1-2 weeks development time
**Benefits**: Better handling of large files and improved responsiveness

### 10. Collaboration Features
**Description**: Enable real-time collaborative editing:
- Multi-user editing with conflict resolution
- Chat functionality
- Presence indicators

**Complexity**: High
**Resource Requirements**: 3-4 weeks development time
**Benefits**: Enables collaborative development directly in the editor