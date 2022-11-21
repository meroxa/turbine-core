package turbinecore

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/meroxa/turbine-core/pkg/ir"
)

const LanguageNotSupportedError = "Currently, we support \"javascript\", \"golang\", \"python\", and \"ruby (beta)\" "

type AppConfig struct {
	Name        string            `json:"name"`
	Environment string            `json:"environment"`
	Pipeline    string            `json:"pipeline"` // TODO: Eventually remove support for providing a pipeline if we need to
	Resources   map[string]string `json:"resources"`
	Language    string            `json:"language"`
}

// validateAppConfig will check if app.json contains information required
func (c *AppConfig) validateAppConfig() error {
	if c.Name == "" {
		return errors.New("application name is required to be specified in your app.json")
	}
	if err := c.validateLanguage(c.Language); err != nil {
		return err
	}
	return nil
}

//validate app.json language, make sure it is supported
func (c *AppConfig) validateLanguage(lang string) error {
	switch lang {
	case "go", string(ir.GoLang):
		return nil
	case "js", strings.ToLower(string(ir.JavaScript)):
		return nil
	case "py", strings.ToLower(string(ir.Python)), strings.ToLower(string(ir.Python3)):
		return nil
	case "rb", strings.ToLower(string(ir.Ruby)):
		return nil
	}
	return fmt.Errorf("language %q not supported. %s", lang, LanguageNotSupportedError)
}

// setPipelineName will check if Pipeline was specified via app.json
// otherwise, pipeline name will be set with the format of `turbine-pipeline-{Name}`
func (c *AppConfig) setPipelineName() {
	if c.Pipeline == "" {
		c.Pipeline = fmt.Sprintf("turbine-pipeline-%s", c.Name)
	}
}

var ReadAppConfig = func(appName, appPath string) (AppConfig, error) {
	if appPath == "" {
		exePath, err := os.Executable()
		if err != nil {
			log.Fatalf("unable to locate executable path; error: %s", err)
		}
		appPath = path.Dir(exePath)
	}

	b, err := os.ReadFile(appPath + "/" + "app.json")
	if err != nil {
		return AppConfig{}, err
	}

	var ac AppConfig
	err = json.Unmarshal(b, &ac)
	if err != nil {
		return AppConfig{}, err
	}

	if appName != "" {
		ac.Name = appName
	}
	err = ac.validateAppConfig()
	if err != nil {
		return AppConfig{}, err
	}

	ac.setPipelineName()
	return ac, nil
}
