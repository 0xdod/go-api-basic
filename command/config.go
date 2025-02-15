package command

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/gilcrest/diy-go-api/datastore"
	"github.com/gilcrest/diy-go-api/domain/errs"
)

const (
	// local JSON Config File path - relative to project root
	localJSONConfigFile = "./config/local.json"
	// staging JSON Config File path - relative to project root
	stagingJSONConfigFile = "./config/staging.json"
	// production JSON Config File path - relative to project root
	productionJSONConfigFile = "./config/production.json"
	// genesisRequestFile is the local JSON Genesis Request File path
	// (relative to project root)
	genesisRequestFile = "./config/genesis/request.json"
)

// ConfigFile defines the configuration file. It is the superset of
// fields for the various environments/builds. For example, when setting
// the local environment based on the ConfigFile, you do not need
// to fill any of the GCP fields.
type ConfigFile struct {
	Config struct {
		HTTPServer struct {
			ListenPort int `json:"listenPort"`
		} `json:"httpServer"`
		Logger struct {
			MinLogLevel   string `json:"minLogLevel"`
			LogLevel      string `json:"logLevel"`
			LogErrorStack bool   `json:"logErrorStack"`
		} `json:"logger"`
		Database struct {
			Host       string `json:"host"`
			Port       int    `json:"port"`
			Name       string `json:"name"`
			User       string `json:"user"`
			Password   string `json:"password"`
			SearchPath string `json:"searchPath"`
		} `json:"database"`
		EncryptionKey string `json:"encryptionKey"`
		GCP           struct {
			ProjectID        string `json:"projectID"`
			ArtifactRegistry struct {
				RepoLocation string `json:"repoLocation"`
				RepoName     string `json:"repoName"`
				ImageID      string `json:"imageID"`
				Tag          string `json:"tag"`
			} `json:"artifactRegistry"`
			CloudSQL struct {
				InstanceName           string `json:"instanceName"`
				InstanceConnectionName string `json:"instanceConnectionName"`
			} `json:"cloudSQL"`
			CloudRun struct {
				ServiceName string `json:"serviceName"`
			} `json:"cloudRun"`
		} `json:"gcp"`
	} `json:"config"`
}

// LoadEnv conditionally sets the environment from a config file
// relative to whichever environment is being set. If Existing is
// passed as EnvConfig, the current environment is used and not overridden.
func LoadEnv(env Env) (err error) {
	var f ConfigFile
	f, err = NewConfigFile(env)
	if err != nil {
		return err
	}

	err = overrideEnv(f)
	if err != nil {
		return err
	}
	return nil
}

// overrideEnv sets the environment
func overrideEnv(f ConfigFile) error {
	var err error

	// minimum accepted log level
	err = os.Setenv(logLevelMinEnv, f.Config.Logger.MinLogLevel)
	if err != nil {
		return err
	}

	// log level
	err = os.Setenv(loglevelEnv, f.Config.Logger.LogLevel)
	if err != nil {
		return err
	}

	// log error stack
	err = os.Setenv(logErrorStackEnv, fmt.Sprintf("%t", f.Config.Logger.LogErrorStack))
	if err != nil {
		return err
	}

	// server port
	err = os.Setenv(portEnv, strconv.Itoa(f.Config.HTTPServer.ListenPort))
	if err != nil {
		return err
	}

	// database host
	err = os.Setenv(datastore.DBHostEnv, f.Config.Database.Host)
	if err != nil {
		return err
	}

	// database port
	err = os.Setenv(datastore.DBPortEnv, strconv.Itoa(f.Config.Database.Port))
	if err != nil {
		return err
	}

	// database name
	err = os.Setenv(datastore.DBNameEnv, f.Config.Database.Name)
	if err != nil {
		return err
	}

	// database user
	err = os.Setenv(datastore.DBUserEnv, f.Config.Database.User)
	if err != nil {
		return err
	}

	// database user password
	err = os.Setenv(datastore.DBPasswordEnv, f.Config.Database.Password)
	if err != nil {
		return err
	}

	// database search path
	err = os.Setenv(datastore.DBSearchPathEnv, f.Config.Database.SearchPath)
	if err != nil {
		return err
	}

	// encryption key
	err = os.Setenv(encryptKeyEnv, f.Config.EncryptionKey)
	if err != nil {
		return err
	}

	return nil
}

// NewConfigFile initializes a ConfigFile struct from a JSON file at a
// predetermined file path for each environment (paths are relative to project root)
//
// Production: ./config/production.json
//
// Staging:    ./config/staging.json
//
// Local:      ./config/local.json
func NewConfigFile(env Env) (ConfigFile, error) {
	var (
		b   []byte
		err error
	)
	switch env {
	case Existing:
		return ConfigFile{}, nil
	case Local:
		b, err = os.ReadFile(localJSONConfigFile)
		if err != nil {
			return ConfigFile{}, err
		}
	case Staging:
		b, err = os.ReadFile(stagingJSONConfigFile)
		if err != nil {
			return ConfigFile{}, err
		}
	case Production:
		b, err = os.ReadFile(productionJSONConfigFile)
		if err != nil {
			return ConfigFile{}, err
		}
	default:
		return ConfigFile{}, errs.E("Invalid environment")
	}

	f := ConfigFile{}
	err = json.Unmarshal(b, &f)
	if err != nil {
		return ConfigFile{}, err
	}

	return f, nil
}

// Env defines the environment
type Env uint8

// Kinds of errors.
//
// The values of the error kinds are common between both
// clients and servers. Do not reorder this list or remove
// any items since that will change their values.
// New items must be added only to the end.
const (
	Existing   Env = iota // Existing environment - current environment is not overridden
	Local                 // Local environment (Local machine)
	Staging               // Staging environment (GCP)
	Production            // Production environment (GCP)

	Invalid Env = 99 // Invalid defines an invalid environment option
)

func (e Env) String() string {
	switch e {
	case Existing:
		return "existing"
	case Local:
		return "local"
	case Staging:
		return "staging"
	case Production:
		return "production"
	case Invalid:
		return "invalid"
	}
	return "unknown_env_config"
}

// ParseEnv converts an env string into an Env value.
// returns Invalid if the input string does not match known values.
func ParseEnv(envStr string) Env {
	switch envStr {
	case "existing":
		return Existing
	case "local":
		return Local
	case "staging":
		return Staging
	case "prod":
		return Production
	default:
		return Invalid
	}
}

// ConfigCueFilePaths defines the paths for config files processed through CUE.
type ConfigCueFilePaths struct {
	// Input defines the list of paths for files to be taken as input for CUE
	Input []string
	// Output defines the path for the JSON output of CUE
	Output string
}

// CUEPaths returns the ConfigCueFilePaths given the environment.
// Paths are relative to the project root.
func CUEPaths(env Env) (ConfigCueFilePaths, error) {
	const (
		schemaInput  = "./config/cue/schema.cue"
		localInput   = "./config/cue/local.cue"
		stagingInput = "./config/cue/staging.cue"
		prodInput    = "./config/cue/production.cue"
	)

	switch env {
	case Local:
		return ConfigCueFilePaths{
			Input:  []string{schemaInput, localInput},
			Output: localJSONConfigFile,
		}, nil
	case Staging:
		return ConfigCueFilePaths{
			Input:  []string{schemaInput, stagingInput},
			Output: stagingJSONConfigFile,
		}, nil
	case Production:
		return ConfigCueFilePaths{
			Input:  []string{schemaInput, prodInput},
			Output: productionJSONConfigFile,
		}, nil
	default:
		return ConfigCueFilePaths{}, errs.E(fmt.Sprintf("There is no path configuration for the %s environment", env))
	}
}

// CUEGenesisPaths returns the ConfigCueFilePaths for the Genesis config.
// Paths are relative to the project root.
func CUEGenesisPaths() ConfigCueFilePaths {
	const (
		schemaInput  = "./config/genesis/cue/schema.cue"
		genesisInput = "./config/genesis/cue/genesis.cue"
	)

	return ConfigCueFilePaths{
		Input:  []string{schemaInput, genesisInput},
		Output: genesisRequestFile,
	}
}
