package constant

// Required environment variables
const (
	// EnvAppDomain application domain
	// APP_DOMAIN - port to listen to. example: https://0.0.0.0:433
	EnvAppDomain = "APP_DOMAIN"

	// APP_NAME application name
	EnvAppName = "APP_NAME"

	// EnvAppEnvironment application environment
	// APP_ENV - server environment configuration. [DEVELOPMENT, TESTING, PRODUCTION]
	EnvAppEnvironment = "APP_ENV"

	EnvDbMigrationDir = "DB_MIGRATION_DIR"

	EnvDbDriver = "DB_DRIVER"

	EnvDbOpen = "DB_OPEN"

	EnvAppVersion = "APP_VERSION"

	// EnvAppLogFolder ...
	EnvAppLogFolder = "LOG_FOLDER"

	EnvSwaggerPath = "SWAGGER_PATH"
)

// Application information/identifier
const (
	// Name application name
	Name = "image-storage"

	// Realm application's realm
	Realm = "image-storage"

	// Version current application version
	Version = "0.1"

	// Domain application's domain
	Domain = "http://0.0.0.0:443"
)

// RequiredEnvironmentVars these are environment variables that is needed by the application
var RequiredEnvironmentVars = []string{
	EnvAppDomain,
	EnvAppEnvironment,
	EnvDbDriver,
	EnvDbOpen,
	EnvAppLogFolder,
}
