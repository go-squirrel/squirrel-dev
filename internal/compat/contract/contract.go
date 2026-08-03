package contract

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

//go:embed legacy_contract.json
var legacyContractJSON []byte

// Contract is the machine-readable compatibility boundary captured from
// squirrel-dev-old before the structural migration starts.
type Contract struct {
	SchemaVersion int               `json:"schemaVersion"`
	Source        Source            `json:"source"`
	Services      []Service         `json:"services"`
	CLI           []CLI             `json:"cli"`
	Configs       []Config          `json:"configs"`
	Responses     ResponseContracts `json:"responses"`
	Databases     []Database        `json:"databases"`
	Async         []AsyncBehavior   `json:"asyncBehaviors"`
	HTTPSamples   []HTTPSample      `json:"httpSamples"`
	KnownGaps     []string          `json:"knownGaps"`
}

type Source struct {
	Repository string `json:"repository"`
	Commit     string `json:"commit"`
}

type Service struct {
	Name   string  `json:"name"`
	Routes []Route `json:"routes"`
}

type Route struct {
	Method    string `json:"method"`
	Path      string `json:"path"`
	Transport string `json:"transport,omitempty"`
	Auth      string `json:"auth"`
}

type CLI struct {
	Program             string       `json:"program"`
	DefaultConfig       string       `json:"defaultConfig"`
	VersionOutputPrefix string       `json:"versionOutputPrefix"`
	Commands            []CLICommand `json:"commands"`
}

type CLICommand struct {
	Name  string    `json:"name"`
	Flags []CLIFlag `json:"flags,omitempty"`
}

type CLIFlag struct {
	Name         string `json:"name"`
	Short        string `json:"short,omitempty"`
	DefaultValue any    `json:"default"`
}

type Config struct {
	Program          string   `json:"program"`
	DefaultFile      string   `json:"defaultFile"`
	RequiredSections []string `json:"requiredSections"`
	SQLitePaths      []string `json:"sqlitePaths,omitempty"`
}

type ResponseContracts struct {
	Envelope  []string       `json:"envelope"`
	Common    []ResponseCode `json:"common"`
	Agent     []ResponseCode `json:"agent"`
	APIServer []ResponseCode `json:"apiserver"`
}

type ResponseCode struct {
	Code    int    `json:"code"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

type Database struct {
	Program     string   `json:"program"`
	LogicalName string   `json:"logicalName"`
	DefaultPath string   `json:"defaultPath,omitempty"`
	Tables      []string `json:"tables"`
	Migrations  []string `json:"migrations"`
}

type AsyncBehavior struct {
	Program string `json:"program"`
	Name    string `json:"name"`
	Detail  string `json:"detail"`
}

type HTTPSample struct {
	Service     string `json:"service"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	Status      int    `json:"status"`
	ContentType string `json:"contentType"`
	Body        string `json:"body"`
}

// LoadLegacy returns a fresh copy of the captured legacy contract.
func LoadLegacy() (*Contract, error) {
	var value Contract
	if err := json.Unmarshal(legacyContractJSON, &value); err != nil {
		return nil, fmt.Errorf("decode legacy compatibility contract: %w", err)
	}
	return &value, nil
}

// Validate checks invariants that must remain true while the implementation is
// migrated module by module.
func Validate(value *Contract) error {
	if value == nil {
		return errors.New("contract is nil")
	}
	if value.SchemaVersion != 1 {
		return fmt.Errorf("unsupported schema version: %d", value.SchemaVersion)
	}
	if value.Source.Repository == "" || value.Source.Commit == "" {
		return errors.New("source repository and commit are required")
	}

	serviceNames := make(map[string]struct{}, len(value.Services))
	for _, service := range value.Services {
		if service.Name == "" {
			return errors.New("service name is required")
		}
		if _, exists := serviceNames[service.Name]; exists {
			return fmt.Errorf("duplicate service: %s", service.Name)
		}
		serviceNames[service.Name] = struct{}{}

		routes := make(map[string]struct{}, len(service.Routes))
		for _, route := range service.Routes {
			if route.Method == "" || !strings.HasPrefix(route.Path, "/") {
				return fmt.Errorf("invalid route in %s: %s %s", service.Name, route.Method, route.Path)
			}
			key := route.Method + " " + route.Path
			if _, exists := routes[key]; exists {
				return fmt.Errorf("duplicate route in %s: %s", service.Name, key)
			}
			routes[key] = struct{}{}
		}
	}

	programs := make(map[string]struct{}, len(value.CLI))
	for _, cli := range value.CLI {
		if cli.Program == "" || cli.DefaultConfig == "" {
			return errors.New("CLI program and default config are required")
		}
		if _, exists := programs[cli.Program]; exists {
			return fmt.Errorf("duplicate CLI program: %s", cli.Program)
		}
		programs[cli.Program] = struct{}{}
	}

	if err := validateResponseCodes("common", value.Responses.Common, nil); err != nil {
		return err
	}
	common := make(map[int]struct{}, len(value.Responses.Common))
	for _, item := range value.Responses.Common {
		common[item.Code] = struct{}{}
	}
	if err := validateResponseCodes("agent", value.Responses.Agent, common); err != nil {
		return err
	}
	if err := validateResponseCodes("apiserver", value.Responses.APIServer, common); err != nil {
		return err
	}

	databases := make(map[string]struct{}, len(value.Databases))
	for _, database := range value.Databases {
		key := database.Program + "/" + database.LogicalName
		if database.Program == "" || database.LogicalName == "" || len(database.Tables) == 0 {
			return fmt.Errorf("invalid database contract: %s", key)
		}
		if _, exists := databases[key]; exists {
			return fmt.Errorf("duplicate database contract: %s", key)
		}
		databases[key] = struct{}{}
	}

	return nil
}

func validateResponseCodes(scope string, values []ResponseCode, reserved map[int]struct{}) error {
	codes := make(map[int]struct{}, len(values))
	for _, item := range values {
		if item.Name == "" || item.Message == "" {
			return fmt.Errorf("invalid response code in %s: %d", scope, item.Code)
		}
		if _, exists := codes[item.Code]; exists {
			return fmt.Errorf("duplicate response code in %s: %d", scope, item.Code)
		}
		if _, exists := reserved[item.Code]; exists {
			return fmt.Errorf("response code in %s conflicts with common code: %d", scope, item.Code)
		}
		codes[item.Code] = struct{}{}
	}
	return nil
}

func (value *Contract) Service(name string) (Service, bool) {
	for _, service := range value.Services {
		if service.Name == name {
			return service, true
		}
	}
	return Service{}, false
}

func (value *Contract) Config(program string) (Config, bool) {
	for _, config := range value.Configs {
		if config.Program == program {
			return config, true
		}
	}
	return Config{}, false
}

func (value *Contract) ResponseMessage(scope string, code int) (string, bool) {
	if message, ok := findResponseMessage(value.Responses.Common, code); ok {
		return message, true
	}
	switch scope {
	case "agent":
		return findResponseMessage(value.Responses.Agent, code)
	case "apiserver":
		return findResponseMessage(value.Responses.APIServer, code)
	}
	return "", false
}

func findResponseMessage(values []ResponseCode, code int) (string, bool) {
	for _, item := range values {
		if item.Code == code {
			return item.Message, true
		}
	}
	return "", false
}
