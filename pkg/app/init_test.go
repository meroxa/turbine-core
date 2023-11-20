package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/meroxa/turbine-core/pkg/ir"

	"github.com/stretchr/testify/require"
)

func TestAppInit_duplicateFile(t *testing.T) {
	type fields struct {
		AppName  string
		Language ir.Lang
		Path     string
	}
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "duplicates the ruby file from the embedded fs",
			fields: fields{
				AppName:  "duplicatefile",
				Language: ir.Ruby,
				Path:     t.TempDir(),
			},
			args: args{
				fileName: "Gemfile",
			},
			wantErr: false,
		},
		{
			name: "duplicates the go file from the embedded fs",
			fields: fields{
				AppName:  "duplicatefile",
				Language: ir.GoLang,
				Path:     t.TempDir(),
			},
			args: args{
				fileName: "app_test.go",
			},
			wantErr: false,
		},
		{
			name: "duplicates the js file from the embedded fs",
			fields: fields{
				AppName:  "duplicatefile",
				Language: ir.JavaScript,
				Path:     t.TempDir(),
			},
			args: args{
				fileName: "index.js",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AppInit{
				appName: tt.fields.AppName,
			}
			srcPath := filepath.Join("templates", string(tt.fields.Language))
			dstPath := filepath.Join(tt.fields.Path, a.appName)

			os.MkdirAll(dstPath, 0o755)
			if err := a.duplicateFileInPath(srcPath, dstPath, tt.args.fileName); (err != nil) != tt.wantErr {
				t.Errorf("AppInit.duplicateFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			require.FileExists(t, filepath.Join(tt.fields.Path, tt.fields.AppName, tt.args.fileName))
		})
	}
}

func TestAppInit_listTemplateContentFromPath(t *testing.T) {
	const appName = "testapp"

	type fields struct {
		AppName  string
		Language ir.Lang
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		want1   []string
		wantErr bool
	}{
		{
			name: "lists files and dir content for ruby app embedded template",
			fields: fields{
				AppName:  appName,
				Language: ir.Ruby,
			},
			want:  []string{"Gemfile", "app.json", "app.rb"},
			want1: []string{"fixtures"},
		},
		{
			name: "lists files and dir content for go app embedded template",
			fields: fields{
				AppName:  appName,
				Language: ir.GoLang,
			},
			want:  []string{".gitignore", "README.md", "app.go", "app.json", "app_test.go"},
			want1: []string{"fixtures"},
		},
		{
			name: "lists files and dir content for js app embedded template",
			fields: fields{
				AppName:  appName,
				Language: ir.JavaScript,
			},
			want:  []string{".gitignore", "README.md", "app.json", "index.js", "index.test.js", "package.json"},
			want1: []string{"fixtures"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AppInit{
				appName: tt.fields.AppName,
			}
			got, got1, err := a.listTemplateContentFromPath(filepath.Join("templates", string(tt.fields.Language)))
			if (err != nil) != tt.wantErr {
				t.Errorf("AppInit.listTemplateContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, tt.want, got)
			require.Equal(t, tt.want1, got1)
		})
	}
}

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
			if err := Init(tt.fields.Path, tt.fields.AppName, tt.fields.Language); (err != nil) != tt.wantErr {
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
