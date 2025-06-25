# Changelog

## Code Review and Cleanup

### Major Changes

1. **Removed Rank Functionality**
   - Simplified code by removing process ranking
   - Renamed `kernel_task_rank.go` to `kernel_task.go`
   - Removed all rank-related variables and logic

2. **Thread Safety Improvements**
   - Added mutex protection for global variables
   - Prevents race conditions between goroutines
   - Protected `lastError` and `lastCPU` variables

3. **Error Handling Improvements**
   - Replaced deprecated `ioutil` with `os` package functions
   - Added proper error logging for all operations
   - Fixed HTTP response body leak on error paths

4. **Security Enhancements**
   - Changed config file permissions from 0644 to 0600 (user-only)
   - Added HTTP client with 5-second timeout
   - Proper resource cleanup with deferred closes

5. **Code Quality Improvements**
   - Added constants for thermal threshold values
   - Improved documentation for exported functions
   - Better code organization and clarity

### Technical Details

- **Constants Added:**
  - `ThermalThresholdIdle = 5.0`
  - `ThermalThresholdLight = 20.0`
  - `ThermalThresholdHeavy = 50.0`
  - `ThermalThresholdThrottle = 100.0`

- **Files Removed:**
  - `kernel_task_rank.go` (replaced with `kernel_task.go`)
  - `debug_temp`
  - `test-temp.go.bak`
  - `setup_sudo.sh`

### Performance Note
The app uses `top -l 2 -s 1` to get accurate kernel_task CPU usage, which takes approximately 1 second. The update interval accounts for this delay.