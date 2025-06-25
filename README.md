# Kernel Task Monitor (KTM)

A macOS menu bar application that monitors system thermal state by tracking `kernel_task` CPU usage.

![Status: Working](https://img.shields.io/badge/status-working-green)
![Platform: macOS](https://img.shields.io/badge/platform-macOS-blue)
![Language: Go](https://img.shields.io/badge/language-Go-00ADD8)
![License: MIT](https://img.shields.io/badge/license-MIT-green)

## Overview

On Apple Silicon Macs, direct temperature sensor access is restricted. KTM uses `kernel_task` CPU usage as a proxy for thermal state, as this system process increases its CPU usage to manage thermals when the system gets hot.

Built to scratch a personal itch during a hot UK summer. Publishing in case others get use out of it. Use at your own risk.

## Features

- 🌡️ Real-time thermal state monitoring via menu bar emoji
- 📊 Shows CPU percentage when thermal activity is significant
- 🔄 Configurable update interval
- 🔌 Optional ESP32 integration for external cooling
- 🐛 Debug mode for troubleshooting
- 🔒 Secure configuration file handling

## Installation

### Download Pre-built App

1. Go to [Releases](https://github.com/a-c-m/kernel_task_monitor/releases)
2. Download the latest `.dmg` file
3. Open the DMG and drag "Kernel Task Monitor" to Applications
4. Launch from Applications folder

### Build from Source

```bash
# Clone the repository
git clone https://github.com/a-c-m/kernel_task_monitor.git
cd kernel_task_monitor

# Install KTM (builds app, installs to /Applications, sets up auto-start)
./install.sh
```

That's it! KTM is now running and will start automatically on login.

### Manual Build

```bash
# Build the app
./build.sh

# The app will be created as "Kernel Task Monitor.app"
# Drag it to /Applications to install

# Or run directly:
open "Kernel Task Monitor.app"
```

### Development Mode

```bash
# Build just the binary
go build -o kernel_task_monitor

# Run from terminal
./kernel_task_monitor --debug
```

## Menu Bar Display

KTM shows your system's thermal state using emojis:

- 😴 **Idle** - System is cool (< 5% CPU)
- 😊 **Light Load** - Normal operation (5-20% CPU)
- 😅 **Heavy Load** - Working hard (20-50% CPU)
- 🥵 **Throttling** - Thermal throttling active (50-100% CPU)
- 🔥 **Heavy Throttling** - Severe throttling (> 100% CPU)

## Command Line Options

```bash
./kernel_task_monitor [options]

Options:
  --debug    Always show CPU percentage (not just when > 20%)
  -t float   Update interval in seconds (default 5.0, minimum 0.5)
```

## Configuration

KTM stores its configuration in `~/.kernel_task_monitor.json`. 

To configure KTM:
1. Click the menu bar icon
2. Select "Configure..." 
3. The config file will open in your default text editor
4. If the file doesn't exist, it will be created with a commented template

The config file uses **JSON5** format, which supports:
- Comments (// and /* */)
- Trailing commas
- Unquoted keys
- Single quotes

**Important:** Restart KTM after changing the configuration!

### Configuration Options

1. **ESP32 Integration** - Send thermal data to external device
2. **Custom Thresholds** - Define your own thermal state boundaries
3. **Custom Emojis** - Personalize the menu bar indicators

Example configuration:
```json5
{
  // ESP32 URL - leave empty to disable
  "esp_url": "http://192.168.1.100/fanspeed",
  
  // Custom thresholds (optional)
  "thresholds": {
    "idle": 10,      // More relaxed idle threshold
    "light": 30,     // Different boundaries
    "heavy": 60,
    "throttle": 120
  },
  
  // Custom emojis (optional)
  "emojis": {
    "idle": "🧊",    // Ice cube for idle
    "light_load": "✅", // Check mark for normal
    "heavy_load": "⚡", // Lightning for busy
    "throttling": "🌡️", // Thermometer
    "heavy_throttling": "☢️", // Nuclear warning
    "error": "⚠️"
  }
}
```

### ESP32 Integration (Optional)

If you have an ESP32-based cooling system, KTM can send thermal data:

```
GET http://your-esp32/endpoint?kernel_task=45.5&state=Heavy_Load
```

To disable ESP32 integration, leave `esp_url` empty:
```json
{
  "esp_url": ""
}
```

## How It Works

1. Uses `top -l 2 -s 1` to sample kernel_task CPU usage
2. Determines thermal state based on CPU percentage
3. Updates menu bar emoji and optional percentage display
4. Optionally sends data to ESP32 for external cooling control

## Requirements

- macOS (tested on Apple Silicon)
- Go 1.18 or later (for building)
- No special permissions required

## Building from Source

```bash
# Install dependencies
go mod download

# Build
go build -o kernel_task_monitor

# Run
./kernel_task_monitor
```

## Testing

To verify KTM is working correctly, you can stress your CPU using any method (e.g., running intensive tasks). You should see the emoji change from 😊 → 😅 → 🥵 → 🔥 as the system heats up.

## Uninstalling

To completely remove KTM:

```bash
./uninstall.sh
```

This will:
- Stop KTM
- Remove it from auto-start
- Delete the app from /Applications
- Optionally remove your configuration

## Auto-Start Management

KTM automatically starts on login. To manage this:

```bash
# Disable auto-start
launchctl unload ~/Library/LaunchAgents/com.ktm.kernel-task-monitor.plist

# Re-enable auto-start
launchctl load ~/Library/LaunchAgents/com.ktm.kernel-task-monitor.plist
```

## Troubleshooting

- **Shows ❓**: The `top` command failed. Check terminal access.
- **Always shows 😴**: Your system is running cool - this is normal!
- **ESP32 not connecting**: Verify URL and network connectivity.
- **Not starting automatically**: Run `./install.sh` to set up auto-start.

## Files and Locations

- **App**: `/Applications/Kernel Task Monitor.app`
- **Config**: `~/.kernel_task_monitor.json`
- **Auto-start**: `~/Library/LaunchAgents/com.ktm.kernel-task-monitor.plist`
- **Logs** (when debugging): `/tmp/ktm.out` and `/tmp/ktm.err`

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

PRs are welcome for bug fixes and improvements. Feature requests are less likely to be accepted - this tool was built to scratch a personal itch and is published in case others find it useful.

Please:
- Test your changes thoroughly
- Follow the existing code style
- Update documentation as needed