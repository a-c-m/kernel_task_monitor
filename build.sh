#!/bin/bash

# Kernel Task Monitor Build Script
# Creates a macOS app bundle

set -e

echo "Building Kernel Task Monitor..."

# Build the binary
echo "1. Compiling Go binary..."
go build -o kernel_task_monitor

# Create app bundle structure
echo "2. Creating app bundle..."
APP_NAME="Kernel Task Monitor.app"
rm -rf "$APP_NAME"
mkdir -p "$APP_NAME/Contents/MacOS"
mkdir -p "$APP_NAME/Contents/Resources"

# Copy binary
echo "3. Copying binary..."
cp kernel_task_monitor "$APP_NAME/Contents/MacOS/"

# Create Info.plist
echo "4. Creating Info.plist..."
cat > "$APP_NAME/Contents/Info.plist" << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleDevelopmentRegion</key>
    <string>en</string>
    <key>CFBundleExecutable</key>
    <string>kernel_task_monitor</string>
    <key>CFBundleIdentifier</key>
    <string>com.ktm.kernel-task-monitor</string>
    <key>CFBundleInfoDictionaryVersion</key>
    <string>6.0</string>
    <key>CFBundleName</key>
    <string>Kernel Task Monitor</string>
    <key>CFBundleDisplayName</key>
    <string>KTM</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0.0</string>
    <key>CFBundleVersion</key>
    <string>1</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.15</string>
    <key>LSUIElement</key>
    <true/>
    <key>NSHighResolutionCapable</key>
    <true/>
    <key>NSPrincipalClass</key>
    <string>NSApplication</string>
    <key>NSHumanReadableCopyright</key>
    <string>Copyright © 2024. All rights reserved.</string>
</dict>
</plist>
EOF

# Create a simple icon if it doesn't exist
if [ -f "icon.icns" ]; then
    echo "5. Copying icon..."
    cp icon.icns "$APP_NAME/Contents/Resources/"
else
    echo "5. No icon.icns found, skipping..."
fi

# Sign the app (optional, requires Developer ID)
if command -v codesign &> /dev/null && [ -n "$CODESIGN_IDENTITY" ]; then
    echo "6. Signing app..."
    codesign --force --deep --sign "$CODESIGN_IDENTITY" "$APP_NAME"
else
    echo "6. Skipping code signing (set CODESIGN_IDENTITY to enable)"
fi

echo ""
echo "✅ Build complete!"
echo ""
echo "The app is ready at: $APP_NAME"
echo ""
echo "To install:"
echo "  1. Drag '$APP_NAME' to /Applications"
echo "  2. Run: ./install.sh (to set up auto-start)"
echo ""
echo "To run now:"
echo "  open '$APP_NAME'"