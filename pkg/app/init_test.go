package app

import (
	"path/filepath"
	"testing"

	"github.com/meroxa/turbine-core/v2/pkg/ir"
	"github.com/stretchr/testify/require"
)

type directory struct {
	name    string
	subDirs []directory
	files   []string
}

func TestAppInit_Init(t *testing.T) {
	const appName = "testapp"

	type fields struct {
		AppName  string
		Language ir.Lang
		Path     string
	}
	tests := []struct {
		name      string
		fields    fields
		wantFiles directory
		wantErr   bool
	}{
		{
			name: "copies the ruby app template to the path",
			fields: fields{
				AppName:  appName,
				Language: ir.Ruby,
				Path:     t.TempDir(),
			},
			wantFiles: directory{
				name: appName,
				subDirs: []directory{
					{
						name:  "fixtures",
						files: []string{"demo.json"},
					},
				},
				files: []string{"app.json", "app.rb", "Gemfile"},
			},
			wantErr: false,
		},
		{
			name: "copies the go app template to the path",
			fields: fields{
				AppName:  appName,
				Language: ir.GoLang,
				Path:     t.TempDir(),
			},
			wantFiles: directory{
				name: appName,
				subDirs: []directory{
					{
						name:  "fixtures",
						files: []string{"demo-no-cdc.json"},
					},
				},
				files: []string{"app.json", "app_test.go", "app.go", "README.md"},
			},
			wantErr: false,
		},
		{
			name: "copies the js app template to the path",
			fields: fields{
				AppName:  appName,
				Language: ir.JavaScript,
				Path:     t.TempDir(),
			},
			wantFiles: directory{
				name: appName,
				subDirs: []directory{
					{
						name:  "fixtures",
						files: []string{"demo-no-cdc.json"},
					},
				},
				files: []string{"app.json", "package.json", "index.js", "index.test.js", "README.md"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAppInit(
				tt.fields.AppName,
				tt.fields.Language,
				tt.fields.Path,
			)
			if err := a.Init(); (err != nil) != tt.wantErr {
				t.Errorf("AppInit.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
			assertDirectory(t, tt.fields.Path, tt.wantFiles)
		})
	}
}

// assertDirectory will continue checking for files and subdirectories until there's none left.
func assertDirectory(t *testing.T, basePath string, dir directory) {
	require.DirExists(t, filepath.Join(basePath, dir.name))

	for _, file := range dir.files {
		require.FileExists(t, filepath.Join(basePath, dir.name, file))
	}

	for _, subDir := range dir.subDirs {
		assertDirectory(t, filepath.Join(basePath, dir.name), subDir)
	}
}
