package environment

// WebServer environment variable names
const (
	ENV_WEBSERVER_PORT = "WEBSERVER_PORT" // web server port (fallback: WEBAPIPORT)
	ENV_WEBSERVER_HOST = "WEBSERVER_HOST" // web server host (fallback: WEBAPIHOST)
	ENV_WEBSERVER_LOGS = "WEBSERVER_LOGS" // web server HTTP logs (fallback: HTTPLOGS)
)

// WebServerSettings holds all WebServer configuration loaded from environment
type WebServerSettings struct {
	Port uint32 `json:"port"`
	Host string `json:"host"`
	Logs bool   `json:"logs"`
}

// NewWebServerSettings creates a new WebServer settings by loading all values from environment
// with fallback compatibility to old variable names
func NewWebServerSettings() WebServerSettings {
	return WebServerSettings{
		Port: getWebServerPort(),
		Host: getWebServerHost(),
		Logs: getWebServerLogs(),
	}
}

// getWebServerPort gets the web server port with fallback compatibility
func getWebServerPort() uint32 {
	// Try new variable first
	if isEnvVarSet(ENV_WEBSERVER_PORT) {
		return getEnvOrDefaultUint32(ENV_WEBSERVER_PORT, 31000)
	}
	// Fallback to old variable
	return getEnvOrDefaultUint32("WEBAPIPORT", 31000)
}

// getWebServerHost gets the web server host with fallback compatibility
func getWebServerHost() string {
	// Try new variable first
	if host := getEnvOrDefaultString(ENV_WEBSERVER_HOST, ""); host != "" {
		return host
	}
	// Fallback to old variable
	return getEnvOrDefaultString("WEBAPIHOST", "")
}

// getWebServerHTTPLogs gets the web server HTTP logs setting with fallback compatibility
func getWebServerLogs() bool {
	// Try new variable first
	if isEnvVarSet(ENV_WEBSERVER_LOGS) {
		return getEnvOrDefaultBool(ENV_WEBSERVER_LOGS, false)
	}
	// Fallback to old variable
	return getEnvOrDefaultBool("HTTPLOGS", false)
}
