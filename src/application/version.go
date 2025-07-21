// Package application provides the application context and lifecycle management
package application

import (
	"fmt"
	"os"

	"github.com/deadelus/go-clean-app/src/cerr"
)

type VersionOption Option

// SetOptionVersion is an Option that sets the application version in the Engine.
// It allows the application version to be configured at runtime.
// This option can be used to set the version of the application when creating a new Engine instance.
// It is useful for applications that need to report their version or for logging purposes.
func SetOptionVersion(version string) VersionOption {
	return func(e *Engine) error {
		e.appVersion = version
		return nil
	}
}

// SetVersionFromSpecifiedEnv is an Option that sets the application version from an environment variable.
// It retrieves the version from the specified environment variable and sets it in the Engine.
// This option is useful for applications that want to configure their version dynamically based on environment variables.
func SetVersionFromSpecifiedEnv(envName string) VersionOption {
	return func(e *Engine) error {
		if version, exists := os.LookupEnv(envName); exists {
			e.appVersion = version
		} else {
			return fmt.Errorf("environment variable %s not set", cerr.ErrMissingConfig)
		}
		return nil
	}
}

// SetVersionFromEnv is an Option that sets the application version from the default environment variable.
// It retrieves the version from the environment variable named AppVersionEnvName and sets it in the Engine.
// This option is useful for applications that want to configure their version dynamically based on environment variables
func SetVersionFromEnv() VersionOption {
	return SetVersionFromSpecifiedEnv(AppVersionEnvName)
}

// SetDefaultVersion is an Option that sets a default version for the Engine.
/*
func setDefaultVersion() VersionOption {
	return func(e *Engine) error {
		e.appVersion = "1.0.0" // Default version
		return nil
	}
}
*/

// Version returns the version of the application.
// It retrieves the version that was set during the Engine creation or through the SetOptionVersion option
func (e *Engine) Version() string {
	return e.appVersion
}
