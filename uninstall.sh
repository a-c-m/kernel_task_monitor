#!/bin/bash

# Kernel Task Monitor Uninstall Script

echo "Kernel Task Monitor (KTM) Uninstaller"
echo "====================================="
echo ""

# Stop and remove LaunchAgent
LAUNCH_AGENTS_DIR="$HOME/Library/LaunchAgents"
PLIST_NAME="com.ktm.kernel-task-monitor.plist"
PLIST_PATH="$LAUNCH_AGENTS_DIR/$PLIST_NAME"

if [ -f "$PLIST_PATH" ]; then
    echo "1. Stopping KTM..."
    launchctl unload "$PLIST_PATH" 2>/dev/null || true
    echo "   ✅ KTM stopped"
    
    echo ""
    echo "2. Removing auto-start..."
    rm -f "$PLIST_PATH"
    echo "   ✅ Auto-start removed"
else
    echo "1. Auto-start not found, skipping..."
fi

# Remove app from Applications
if [ -d "/Applications/Kernel Task Monitor.app" ]; then
    echo ""
    echo "3. Removing app from /Applications..."
    rm -rf "/Applications/Kernel Task Monitor.app"
    echo "   ✅ App removed"
else
    echo "3. App not found in /Applications, skipping..."
fi

# Ask about config file
echo ""
echo "4. Configuration file"
CONFIG_FILE="$HOME/.kernel_task_monitor.json"
if [ -f "$CONFIG_FILE" ]; then
    read -p "   Remove configuration file? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        rm -f "$CONFIG_FILE"
        rm -f "$CONFIG_FILE.README" 2>/dev/null || true
        echo "   ✅ Configuration removed"
    else
        echo "   ℹ️  Configuration kept at: $CONFIG_FILE"
    fi
else
    echo "   No configuration file found"
fi

echo ""
echo "Uninstall complete!"
echo ""
echo "KTM has been removed from your system."
echo "To reinstall, run: ./install.sh"
echo ""