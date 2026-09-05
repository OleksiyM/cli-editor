#!/bin/bash

# Script to clean up VS Code bloat on macOS and Linux (Debian-based and Fedora-based)
# Removes cached extensions, logs, and unnecessary files
# Preserves user settings and active extensions
# Run periodically (e.g., via cron) or manually

# Exit on error
set -e

# Detect OS
OS=$(uname -s)
case "$OS" in
    Linux*)     PLATFORM=Linux ;;
    Darwin*)    PLATFORM=macOS ;;
    *)          echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Define VS Code directories
if [ "$PLATFORM" = "macOS" ]; then
    VSCODE_APP="/Applications/Visual Studio Code.app"
    USER_DIR="$HOME/Library/Application Support/Code"
    CACHE_DIR="$HOME/Library/Caches/com.microsoft.VSCode"
    LOG_DIR="$USER_DIR/logs"
    CACHE_VSIX="$USER_DIR/CachedExtensionVSIXs"
    WEB_STORAGE="$USER_DIR/WebStorage"
elif [ "$PLATFORM" = "Linux" ]; then
    VSCODE_APP="/usr/share/code"
    USER_DIR="$HOME/.config/Code"
    CACHE_DIR="$HOME/.cache/code"
    LOG_DIR="$USER_DIR/logs"
    CACHE_VSIX="$USER_DIR/CachedExtensionVSIXs"
    WEB_STORAGE="$USER_DIR/WebStorage"
fi

# Check if VS Code is installed
if [ ! -d "$USER_DIR" ] && [ ! -d "$VSCODE_APP" ]; then
    echo "VS Code not found. Exiting."
    exit 1
fi

# Function to clean cache
clean_cache() {
    echo "Cleaning VS Code cache..."
    if [ -d "$CACHE_DIR" ]; then
        rm -rf "$CACHE_DIR"/*
        echo "Cache cleared: $CACHE_DIR"
    else
        echo "No cache directory found."
    fi
    if [ -d "$CACHE_VSIX" ]; then
        rm -rf "$CACHE_VSIX"/*
        # Remove all files and specifically named files in the cache directory
        find "$CACHE_VSIX" -type f -name ".*" -delete
        echo "Cache VSIx cleared: $CACHE_VSIX"
    else
        echo "No VSIx cache directory found."
    fi
}

# Function to clean logs
clean_logs() {
    echo "Cleaning VS Code logs..."
    if [ -d "$LOG_DIR" ]; then
        find "$LOG_DIR" -type f -name "*.log" -delete
        echo "Logs cleared: $LOG_DIR"
    else
        echo "No log directory found."
    fi
}

# Function to clean uninstalled extensions
clean_extensions() {
    echo "Cleaning uninstalled extensions..."
    EXT_DIR="$USER_DIR/User/extensions"
    if [ -d "$EXT_DIR" ]; then
        # Get list of installed extensions
        INSTALLED_EXTS=$(code --list-extensions)
        # Remove extension folders not matching installed extensions
        for ext_folder in "$EXT_DIR"/*; do
            if [ -d "$ext_folder" ]; then
                ext_id=$(basename "$ext_folder" | cut -d'-' -f1-2) # Extract publisher.extension
                if ! echo "$INSTALLED_EXTS" | grep -qi "$ext_id"; then
                    rm -rf "$ext_folder"
                    echo "Removed uninstalled extension: $ext_folder"
                fi
            fi
        done
    else
        echo "No extensions directory found."
    fi
}

# Function to clean old VS Code update files
clean_update_files() {
    echo "Checking for old VS Code update files..."
    TEMP_DIR="/tmp/vscode*"
    if ls $TEMP_DIR >/dev/null 2>&1; then
        rm -rf $TEMP_DIR
        echo "Cleared temporary update files: $TEMP_DIR"
    else
        echo "No temporary update files found."
    fi
}

# Function to clean orphaned packages (Linux-specific)
clean_orphaned_packages() {
    if [ "$PLATFORM" = "Linux" ]; then
        echo "Checking for orphaned VS Code packages..."
        if command -v apt >/dev/null 2>&1; then
            echo "Detected Debian-based system (apt)."
            sudo apt autoremove -y
            sudo apt autoclean
            echo "Orphaned packages and apt cache cleaned."
        elif command -v dnf >/dev/null 2>&1; then
            echo "Detected Fedora-based system (dnf)."
            sudo dnf autoremove -y
            sudo dnf clean all
            echo "Orphaned packages and dnf cache cleaned."
        else
            echo "No supported package manager (apt or dnf) found. Skipping orphaned package cleanup."
        fi
    fi
}

# Function to clean WebStorage
clean_webstorage() {
    echo "Cleaning WebStorage..."
    if [ -d "$WEB_STORAGE" ]; then
        rm -rf "$WEB_STORAGE"/*
        echo "WebStorage cleared: $WEB_STORAGE"
    else
        echo "No WebStorage directory found."
    fi
}

# Main cleanup routine
echo "Before cleanup"
du -sh "$USER_DIR"

echo "Starting VS Code cleanup on $PLATFORM..."
clean_cache
clean_logs
clean_extensions
clean_update_files
clean_orphaned_packages
clean_webstorage
echo "After cleanup"
du -sh "$USER_DIR"
echo "Cleanup completed!"


echo "---------------"
echo "Big files +50M:"
find "$USER_DIR" -type f -size +50M
echo "---------------"
echo "WorkspaceStorage:"
du -sh "$USER_DIR/User/workspaceStorage"
# echo "---------------"
# echo "Installed plugins:"
# code --list-extensions
