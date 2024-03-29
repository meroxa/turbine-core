package app_test

import (
	"os"
	"path"
	"testing"

	"github.com/meroxa/turbine-core/v2/pkg/app"
	"github.com/stretchr/testify/assert"
)

func Test_ReadConfig(t *testing.T) {
	testCases := []struct {
		desc    string
		appName string
		appPath string
		errmsg  string
	}{
		{
			desc:    "reads a valid app config without app name",
			appName: "testapp",
			appPath: setupAppJson(t),
		},
		{
			desc:    "reads a valid app config with app name provided",
			appName: "new-name-app",
			appPath: setupAppJson(t),
		},
		{
			desc:    "reads a valid app config with deprecated fields",
			appName: "new-name-app",
			appPath: setupAppJsonWithDeprecatedFields(t),
		},
		{
			desc:    "fails to read an config with missing app name",
			appPath: setupAppJsonMissingField(t),
			errmsg:  "application name is required",
		},
		{
			desc:    "fails to read bad app json",
			appPath: setupBadAppJson(t),
			errmsg:  "invalid character",
		},
		{
			desc:    "fails when app.json is missing",
			appPath: t.TempDir(),
			errmsg:  "no such file or directory",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ac, err := app.ReadConfig(tc.appName, tc.appPath)
			if tc.errmsg != "" && assert.Error(t, err) {
				assert.ErrorContains(t, err, tc.errmsg)
			}
			if tc.errmsg == "" && assert.NoError(t, err) {
				assert.Equal(t, ac.Name, tc.appName)
			}
		})
	}
}

func setupAppJson(t *testing.T) string {
	tmpdir := t.TempDir()
	if err := os.WriteFile(
		path.Join(tmpdir, "app.json"),
		[]byte(`{
				  "name": "testapp",
				  "language": "golang",
				  "resources": {
				    "source_name": "fixtures/demo-cdc.json"
				  }
				}`),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	return tmpdir
}

func setupAppJsonWithDeprecatedFields(t *testing.T) string {
	tmpdir := t.TempDir()
	if err := os.WriteFile(
		path.Join(tmpdir, "app.json"),
		[]byte(`{
				  "name": "testapp",
				  "language": "golang",
				  "resources": {
				    "source_name": "fixtures/demo-cdc.json"
				  }
				}`),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	return tmpdir
}

func setupAppJsonMissingField(t *testing.T) string {
	tmpdir := t.TempDir()
	if err := os.WriteFile(
		path.Join(tmpdir, "app.json"),
		[]byte(`{
				  "language": "golang",
				  "resources": {
				    "source_name": "fixtures/demo-cdc.json"
				  }
				}`),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	return tmpdir
}

func setupBadAppJson(t *testing.T) string {
	tmpdir := t.TempDir()
	if err := os.WriteFile(
		path.Join(tmpdir, "app.json"),
		[]byte(`invalid-json`),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	return tmpdir
}
