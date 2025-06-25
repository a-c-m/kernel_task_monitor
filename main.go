package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/getlantern/systray"
	"github.com/titanous/json5"
)

type Config struct {
	ESPURL string `json:"esp_url"`
	
	// Thermal thresholds (optional)
	Thresholds struct {
		Idle      float64 `json:"idle"`
		Light     float64 `json:"light"`
		Heavy     float64 `json:"heavy"`
		Throttle  float64 `json:"throttle"`
	} `json:"thresholds"`
	
	// Emoji indicators (optional)
	Emojis struct {
		Idle            string `json:"idle"`
		LightLoad       string `json:"light_load"`
		HeavyLoad       string `json:"heavy_load"`
		Throttling      string `json:"throttling"`
		HeavyThrottling string `json:"heavy_throttling"`
		Error           string `json:"error"`
	} `json:"emojis"`
}

var (
	config Config
	configPath string
	
	// Protected by mutex for thread safety
	mu sync.RWMutex
	lastError string
	lastCPU float64
	
	// Command-line flags
	debugMode = flag.Bool("debug", false, "Show CPU percentage always")
	updateInterval = flag.Float64("t", 5.0, "Update interval in seconds (min 0.5)")
	
	// HTTP client with timeout
	httpClient = &http.Client{
		Timeout: 5 * time.Second,
	}
)

func init() {
	flag.Parse()
	
	// Validate update interval
	if *updateInterval < 0.5 {
		*updateInterval = 0.5
	}
}

// Get thermal threshold with config override
func getThreshold(name string) float64 {
	switch name {
	case "idle":
		if config.Thresholds.Idle > 0 {
			return config.Thresholds.Idle
		}
		return ThermalThresholdIdle
	case "light":
		if config.Thresholds.Light > 0 {
			return config.Thresholds.Light
		}
		return ThermalThresholdLight
	case "heavy":
		if config.Thresholds.Heavy > 0 {
			return config.Thresholds.Heavy
		}
		return ThermalThresholdHeavy
	case "throttle":
		if config.Thresholds.Throttle > 0 {
			return config.Thresholds.Throttle
		}
		return ThermalThresholdThrottle
	default:
		return 0
	}
}

// Get thermal state with custom thresholds
func getThermalState(cpu float64) string {
	switch {
	case cpu <= getThreshold("idle"):
		return "Idle"
	case cpu <= getThreshold("light"):
		return "Light Load"
	case cpu <= getThreshold("heavy"):
		return "Heavy Load"
	case cpu <= getThreshold("throttle"):
		return "Throttling"
	default:
		return "Heavy Throttling"
	}
}

// Get configured emoji with fallback to default
func getEmoji(state string) string {
	switch state {
	case "Idle":
		if config.Emojis.Idle != "" {
			return config.Emojis.Idle
		}
		return "😴"
	case "Light Load":
		if config.Emojis.LightLoad != "" {
			return config.Emojis.LightLoad
		}
		return "😊"
	case "Heavy Load":
		if config.Emojis.HeavyLoad != "" {
			return config.Emojis.HeavyLoad
		}
		return "😅"
	case "Throttling":
		if config.Emojis.Throttling != "" {
			return config.Emojis.Throttling
		}
		return "🥵"
	case "Heavy Throttling":
		if config.Emojis.HeavyThrottling != "" {
			return config.Emojis.HeavyThrottling
		}
		return "🔥"
	case "Error":
		if config.Emojis.Error != "" {
			return config.Emojis.Error
		}
		return "❓"
	default:
		return "❓"
	}
}

func getKernelTaskUsage() float64 {
	cpu, err := GetKernelTaskCPU()
	
	mu.Lock()
	defer mu.Unlock()
	
	if err != nil {
		lastError = err.Error()
		if *debugMode {
			log.Printf("Error getting kernel_task: %v", err)
		}
		return -1.0
	}
	lastError = ""
	lastCPU = cpu
	
	if *debugMode {
		log.Printf("kernel_task: CPU %.1f%%", cpu)
	}
	
	return cpu
}

func loadConfig() {
	if configPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Printf("Error getting home directory: %v", err)
			home = "."
		}
		configPath = filepath.Join(home, ".kernel_task_monitor.json")
	}

	data, err := os.ReadFile(configPath)
	if err == nil {
		// Try JSON5 first (supports comments)
		if err := json5.Unmarshal(data, &config); err != nil {
			// Fall back to standard JSON
			if err := json.Unmarshal(data, &config); err != nil {
				log.Printf("Error parsing config file: %v", err)
			}
		}
	} else if !os.IsNotExist(err) {
		log.Printf("Error reading config file: %v", err)
	}
	
	// ESP32 integration is optional - leave empty to disable
}

func saveConfig() {
	data, err := json.Marshal(config)
	if err != nil {
		log.Printf("Error marshaling config: %v", err)
		return
	}
	
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		log.Printf("Error saving config file: %v", err)
	}
}

func createConfigTemplate() error {
	// Check if file already exists
	if _, err := os.Stat(configPath); err == nil {
		return nil // File exists, don't overwrite
	}
	
	// Create the file with a helpful JSON5 template (supports comments!)
	template := `{
  // Kernel Task Monitor (KTM) Configuration
  // ========================================
  // This file uses JSON5 format, which supports comments and trailing commas
  // 
  // IMPORTANT: Restart KTM after making changes to this file!
  
  // ESP32 Integration (optional)
  // ---------------------------
  // Set the URL of your ESP32 device to receive thermal data
  // Leave empty ("") to disable ESP32 integration
  // 
  // When enabled, KTM sends HTTP GET requests with these parameters:
  //   - kernel_task: CPU percentage (0.0-400.0+)
  //   - state: Current thermal state (Idle|Light_Load|Heavy_Load|Throttling|Heavy_Throttling)
  //
  // Example request:
  //   GET http://192.168.1.100/fanspeed?kernel_task=45.5&state=Heavy_Load
  //
  // URL examples:
  // - "http://192.168.1.100/fanspeed"  (typical home network)
  // - "http://esp32.local/fanspeed"    (mDNS hostname)
  // - ""                               (disabled)
  
  "esp_url": "",
  
  // Thermal Thresholds (optional)
  // -----------------------------
  // Customize when thermal states change (in CPU %)
  // Leave commented to use defaults
  
  // "thresholds": {
  //   "idle": 5,        // Below this = Idle (default: 5)
  //   "light": 20,      // Below this = Light Load (default: 20)
  //   "heavy": 50,      // Below this = Heavy Load (default: 50)
  //   "throttle": 100,  // Below this = Throttling (default: 100)
  //                     // Above throttle = Heavy Throttling
  // },
  
  // Custom Emojis (optional)
  // ------------------------
  // Customize the menu bar indicators
  // Leave commented to use defaults
  
  // "emojis": {
  //   "idle": "😴",              // Default: 😴 (sleeping)
  //   "light_load": "😊",        // Default: 😊 (happy)
  //   "heavy_load": "😅",        // Default: 😅 (sweating)
  //   "throttling": "🥵",        // Default: 🥵 (hot)
  //   "heavy_throttling": "🔥",  // Default: 🔥 (fire)
  //   "error": "❓",             // Default: ❓ (question)
  // },
}
`
	
	// Create the main config file
	return os.WriteFile(configPath, []byte(template), 0600)
}

func openConfigFile() {
	// Create template if file doesn't exist
	if err := createConfigTemplate(); err != nil {
		log.Printf("Error creating config template: %v", err)
	}
	
	// Open the file with the default editor
	cmd := exec.Command("open", configPath)
	if err := cmd.Start(); err != nil {
		log.Printf("Error opening config file: %v", err)
		log.Printf("Config file location: %s", configPath)
	}
}

func main() {
	loadConfig()
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("--%")
	tooltip := fmt.Sprintf("kernel_task CPU Monitor (update: %.1fs)", *updateInterval)
	if *debugMode {
		tooltip += " [DEBUG]"
	}
	systray.SetTooltip(tooltip)

	mConfigure := systray.AddMenuItem("Configure...", "Open configuration file")
	currentURLText := "ESP32: Disabled"
	if config.ESPURL != "" {
		currentURLText = "ESP32 URL: " + config.ESPURL
	}
	mCurrentURL := systray.AddMenuItem(currentURLText, "Current ESP32 configuration")
	mCurrentURL.Disable()
	systray.AddSeparator()
	
	// Error status menu items (hidden initially)
	mError := systray.AddMenuItem("", "")
	mError.Hide()
	mError.Disable()
	mHelp1 := systray.AddMenuItem("", "")
	mHelp1.Hide()
	mHelp2 := systray.AddMenuItem("", "")
	mHelp2.Hide()
	mHelp3 := systray.AddMenuItem("", "")
	mHelp3.Hide()
	
	// Add thermal status menu item
	systray.AddSeparator()
	mThermalStatus := systray.AddMenuItem("State: Unknown", "Thermal state based on kernel_task")
	mThermalStatus.Disable()
	
	// Show current settings
	mSettings := systray.AddMenuItem(fmt.Sprintf("Update: %.1fs%s", *updateInterval, func() string {
		if *debugMode {
			return ", Debug ON"
		}
		return ""
	}()), "Current settings")
	mSettings.Disable()
	
	systray.AddSeparator()
	mNote := systray.AddMenuItem("Changes require restart", "Configuration changes require restarting KTM")
	mNote.Disable()
	
	mQuit := systray.AddMenuItem("Quit", "Quit the app")

	// Update kernel_task usage in title
	go func() {
		for {
			usage := getKernelTaskUsage()
			titleText := ""
			
			if usage >= 0 {
				// Hide error messages when working
				if !mError.Disabled() {
					mError.Hide()
					mHelp1.Hide()
					mHelp2.Hide()
					mHelp3.Hide()
				}
				
				mu.RLock()
				currentCPU := lastCPU
				mu.RUnlock()
				
				thermalState := getThermalState(currentCPU)
				
				// Add thermal indicator emoji
				titleText = getEmoji(thermalState)
				
				// Add CPU percentage in debug mode or when significant
				if *debugMode || currentCPU > getThreshold("light") {
					titleText += fmt.Sprintf(" %.0f%%", currentCPU)
				}
				
				// Update thermal status menu
				statusText := fmt.Sprintf("State: %s (CPU: %.0f%%)", thermalState, currentCPU)
				mThermalStatus.SetTitle(statusText)
			} else {
				titleText = getEmoji("Error")
			}
			
			systray.SetTitle(titleText)
			
			// Show error and help if reading failed
			mu.RLock()
			currentError := lastError
			mu.RUnlock()
			
			if usage < 0 && currentError != "" {
				mError.SetTitle("⚠️ Error: " + currentError)
				mError.Show()
				
				if strings.Contains(currentError, "privileged") || strings.Contains(currentError, "sudo") {
					mHelp1.SetTitle("To fix: Run with sudo")
					mHelp1.Show()
					mHelp2.SetTitle("Terminal: sudo " + os.Args[0])
					mHelp2.Show()
					mHelp3.SetTitle("Or: Configure sudo NOPASSWD for powermetrics")
					mHelp3.Show()
				} else {
					mHelp1.SetTitle("Cannot read kernel_task CPU usage")
					mHelp1.Show()
					mHelp2.SetTitle("This should not happen - kernel_task always exists")
					mHelp2.Show()
					mHelp3.SetTitle("Try restarting the app")
					mHelp3.Show()
				}
			}

			mu.RLock()
			sendCPU := lastCPU
			mu.RUnlock()
			
			if config.ESPURL != "" && sendCPU >= 0 {
				// Send kernel_task CPU percentage
				url := fmt.Sprintf("%s?kernel_task=%.1f", config.ESPURL, sendCPU)
				
				// Always include state
				thermalState := getThermalState(sendCPU)
				url += fmt.Sprintf("&state=%s", strings.ReplaceAll(thermalState, " ", "_"))
				
				resp, err := httpClient.Get(url)
				if err != nil {
					if *debugMode {
						log.Printf("Error sending data to ESP32: %v", err)
					}
				} else {
					defer resp.Body.Close()
					if *debugMode {
						log.Printf("Sent to ESP32: kernel_task=%.1f, state=%s", sendCPU, thermalState)
					}
				}
			}
			
			// Sleep for the interval minus the time taken by top command (approximately 1 second)
			sleepDuration := *updateInterval - 1.0
			if sleepDuration < 0.5 {
				sleepDuration = 0.5
			}
			time.Sleep(time.Duration(sleepDuration * float64(time.Second)))
		}
	}()

	// Handle menu clicks
	go func() {
		for {
			select {
			case <-mConfigure.ClickedCh:
				openConfigFile()
				log.Println("Opening configuration file...")
				log.Println("After editing, restart KTM for changes to take effect.")
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func onExit() {
	// Cleanup
}
