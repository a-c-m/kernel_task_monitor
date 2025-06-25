# Kernel Task Monitor (KTM) - Technical Documentation

A menu bar app that monitors `kernel_task` CPU usage as a proxy for thermal state on Apple Silicon Macs, with optional ESP32 integration for external cooling control.

## How It Works

On macOS, `kernel_task` is the system process responsible for thermal management. When your Mac gets hot, `kernel_task` increases its CPU usage to force other processes to slow down, effectively throttling the system.

### CPU Usage Ranges

- **😴 Idle**: 1-5% (Normal operation, cool system)
- **😊 Light Load**: 5-20% (System working normally)
- **😅 Heavy Load**: 20-50% (System under load but managing)
- **🥵 Throttling**: 50-100% (Thermal throttling active)
- **🔥 Heavy Throttling**: 100%+ (Severe thermal throttling)

## Features

- Shows emoji in menu bar indicating thermal state
- Shows CPU percentage when kernel_task usage > 20%
- Monitors kernel_task CPU usage accurately using `top -l 2`
- Optional ESP32 integration for external cooling control
- Configurable update interval (default 5 seconds)
- Debug mode for detailed logging

## Menu Bar Display

- **😴** - System idle (kernel_task < 5%)
- **😊** - Normal operation (kernel_task 5-20%)
- **😅** - Working hard (kernel_task 20-50%)
- **🥵** - Thermal throttling (kernel_task 50-100%)
- **🔥** - Heavy throttling (kernel_task > 100%)
- **❓** - Error reading kernel_task

When kernel_task CPU > 20%, the percentage is shown: "😅 35%"

## ESP32 Integration (Optional)

When configured with an ESP32 URL, sends thermal data via HTTP GET:
```
GET http://your-esp32/endpoint?kernel_task=45.5&state=Heavy_Load
```

Parameters:
- `kernel_task`: CPU percentage (0-400+)
- `state`: Current thermal state (Idle, Light_Load, Heavy_Load, Throttling, Heavy_Throttling)

To disable ESP32 integration, leave the `esp_url` field empty in the config file.

## Installation & Usage

```bash
# Build
go build -o kernel_task_monitor

# Run with default settings (5 second updates)
./kernel_task_monitor

# Run in debug mode (always show CPU percentage)
./kernel_task_monitor --debug

# Run with custom update interval (e.g., 2 seconds)
./kernel_task_monitor -t 2

# Combine flags (debug mode with 0.5 second updates)
./kernel_task_monitor --debug -t 0.5

# Configure ESP32 URL (optional)
# Click on the menu bar icon → "Set ESP32 URL..."
# (Currently logs config file location to console)
```

### Command-line Options

- `--debug`: Always show kernel_task CPU percentage in menu bar (not just when > 20%)
- `-t <seconds>`: Set update interval in seconds (minimum 0.5, default 5.0)

## Configuration

Saves to `~/.kernel_task_monitor.json`:
```json
{
  "esp_url": "http://192.168.1.100/fanspeed"
}
```

To disable ESP32 integration:
```json
{
  "esp_url": ""
}
```

## Testing Thermal Throttling

To test the monitor:
```bash
# Run CPU stress test
yes > /dev/null &

# Or use the included stress test
./stress_test.sh
```

Watch as kernel_task CPU usage increases and the emoji changes from 😊 → 😅 → 🥵 → 🔥

## Why kernel_task Instead of Temperature?

On Apple Silicon Macs (M1/M2/M3), direct temperature sensor access is restricted. However, `kernel_task` CPU usage is a perfect proxy:
- It directly indicates when thermal throttling is occurring
- It's always accessible without special permissions
- It shows the actual system response to thermal conditions
- High kernel_task usage = system is actively cooling itself

Note: The app uses `top -l 2 -s 1` to get accurate kernel_task CPU usage, which takes about 1 second to sample. This is accounted for in the update interval timing.

## Troubleshooting

- **❓ appears**: Check if `top` command works in terminal
- **Always shows 😴**: System is running cool (this is good!)
- **ESP32 not receiving data**: Check URL configuration and network (or disable if not needed)

## How It Helps

Your ESP32 can use the kernel_task data to control cooling:
- < 20%: Normal fan speed
- 20-50%: Increase fan speed
- 50-100%: High fan speed
- > 100%: Maximum cooling needed