package environment

import logrus "github.com/sirupsen/logrus"

// General environment variable names
const (
	ENV_MIGRATIONS               = "MIGRATIONS"               // enable migrations (can also be a path)
	ENV_TITLE                    = "APP_TITLE"                // application title for whatsapp id
	ENV_REMOVEDIGIT9             = "REMOVEDIGIT9"             // remove digit 9 from phones
	ENV_SYNOPSISLENGTH           = "SYNOPSISLENGTH"           // synopsis length for messages
	ENV_CACHELENGTH              = "CACHELENGTH"              // cache max items
	ENV_CACHEDAYS                = "CACHEDAYS"                // cache max days
	ENV_CONVERT_WAVE_TO_OGG      = "CONVERT_WAVE_TO_OGG"      // convert wave to OGG
	ENV_COMPATIBLE_MIME_AS_AUDIO = "COMPATIBLE_MIME_AS_AUDIO" // treat compatible MIME as audio
	ENV_ACCOUNTSETUP             = "ACCOUNTSETUP"             // enable or disable account creation
	ENV_TESTING                  = "TESTING"                  // testing mode
	ENV_LOGLEVEL                 = "LOGLEVEL"                 // general log level
	ENV_CONVERT_PNG_TO_JPG       = "CONVERT_PNG_TO_JPG"       // convert PNG to JPG (not implemented yet)
)

// GeneralConfig holds all general application configuration loaded from environment
type GeneralSettings struct {
	Migrations            string `json:"migrations"`
	AppTitle              string `json:"app_title"`
	RemoveDigit9          bool   `json:"remove_digit_9"`
	SynopsisLength        uint32 `json:"synopsis_length"`
	CacheLength           uint64 `json:"cache_length"`
	CacheDays             uint32 `json:"cache_days"`
	ConvertWaveToOGG      bool   `json:"convert_wave_to_ogg"`
	CompatibleMIMEAsAudio bool   `json:"compatible_mime_as_audio"`
	AccountSetup          bool   `json:"account_setup"`
	Testing               bool   `json:"testing"`
	LogLevel              string `json:"log_level"`
	ConvertPNGToJPG       bool   `json:"convert_png_to_jpg"`
}

// NewGeneralSettings creates a new general settings by loading all values from environment
func NewGeneralSettings() GeneralSettings {
	return GeneralSettings{
		Migrations:            getEnvOrDefaultString(ENV_MIGRATIONS, "true"),
		AppTitle:              getEnvOrDefaultString(ENV_TITLE, ""),
		RemoveDigit9:          getEnvOrDefaultBool(ENV_REMOVEDIGIT9, false),
		SynopsisLength:        getEnvOrDefaultUint32(ENV_SYNOPSISLENGTH, 150),
		CacheLength:           getEnvOrDefaultUint64(ENV_CACHELENGTH, 0),
		CacheDays:             getEnvOrDefaultUint32(ENV_CACHEDAYS, 0),
		ConvertWaveToOGG:      getEnvOrDefaultBool(ENV_CONVERT_WAVE_TO_OGG, true),
		CompatibleMIMEAsAudio: getEnvOrDefaultBool(ENV_COMPATIBLE_MIME_AS_AUDIO, true),
		AccountSetup:          getEnvOrDefaultBool(ENV_ACCOUNTSETUP, true),
		Testing:               getEnvOrDefaultBool(ENV_TESTING, false),
		LogLevel:              getEnvOrDefaultString(ENV_LOGLEVEL, ""),
		ConvertPNGToJPG:       getEnvOrDefaultBool(ENV_CONVERT_PNG_TO_JPG, false),
	}
}

// UseCompatibleMIMEsAsAudio returns combined result of ConvertWaveToOGG and CompatibleMIMEAsAudio
func (config *GeneralSettings) UseCompatibleMIMEsAsAudio() bool {
	return config.ConvertWaveToOGG && config.CompatibleMIMEAsAudio
}

// Migrate checks if database migrations should be enabled based on the Migrations setting
func (config *GeneralSettings) Migrate() bool {
	// If it's "false", return false. Otherwise (including custom paths), return true
	return config.Migrations != "false"
}

// MigrationPath returns the custom path for database migrations if specified
func (config *GeneralSettings) MigrationPath() string {
	// If it's "true" or "false", return empty (use default behavior)
	if config.Migrations == "true" || config.Migrations == "false" {
		return ""
	}
	// Otherwise, return the custom path
	return config.Migrations
}

// LogLevelFromLogrus returns a parsed logrus.Level from the config.
// If the environment variable is empty, returns the provided default level.
func (config *GeneralSettings) LogLevelFromLogrus(defaultLevel logrus.Level) logrus.Level {
	if len(config.LogLevel) == 0 {
		return defaultLevel
	}

	envLevel, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		panic("Invalid log level: " + config.LogLevel +
			". Valid levels are: panic, fatal, error, warn, info, debug, trace")
	}
	return envLevel
}
