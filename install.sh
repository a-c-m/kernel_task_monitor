#!/bin/bash

# Kernel Task Monitor Installation Script

set -e

echo "Kernel Task Monitor (KTM) Installer"
echo "==================================="
echo ""

# Check if running as root
if [ "$EUID" -eq 0 ]; then 
   echo "❌ Please don't run this as root/sudo"
   exit 1
fi

# Build the app first
if [ ! -d "Kernel Task Monitor.app" ]; then
    echo "Building KTM app..."
    ./build.sh
fi

# Check if app was built
if [ ! -d "Kernel Task Monitor.app" ]; then
    echo "❌ Failed to build app. Please run ./build.sh first"
    exit 1
fi

# Install to Applications
echo ""
echo "1. Installing app to /Applications..."
if [ -d "/Applications/Kernel Task Monitor.app" ]; then
    echo "   Removing old version..."
    rm -rf "/Applications/Kernel Task Monitor.app"
fi

cp -R "Kernel Task Monitor.app" /Applications/
echo "   ✅ App installed"

# Set up LaunchAgent for auto-start
echo ""
echo "2. Setting up auto-start..."
LAUNCH_AGENTS_DIR="$HOME/Library/LaunchAgents"
PLIST_NAME="com.ktm.kernel-task-monitor.plist"

# Create LaunchAgents directory if it doesn't exist
mkdir -p "$LAUNCH_AGENTS_DIR"

# Copy plist file
cp "$PLIST_NAME" "$LAUNCH_AGENTS_DIR/"
echo "   ✅ LaunchAgent installed"

# Load the LaunchAgent
echo ""
echo "3. Starting KTM..."
launchctl unload "$LAUNCH_AGENTS_DIR/$PLIST_NAME" 2>/dev/null || true
launchctl load -w "$LAUNCH_AGENTS_DIR/$PLIST_NAME"
echo "   ✅ KTM started"

echo ""
echo "Installation complete! 🎉"
echo ""
echo "KTM is now:"
echo "  • Installed in /Applications"
echo "  • Set to start automatically on login"
echo "  • Currently running"
echo ""
echo "To configure KTM:"
echo "  • Click the menu bar icon and select 'Configure...'"
echo "  • Or edit: ~/.kernel_task_monitor.json"
echo ""
echo "To uninstall:"
echo "  • Run: ./uninstall.sh"
echo ""