package app

import (
	"path/filepath"
	"testing"

	"github.com/meroxa/turbine-core/pkg/ir"

	"github.com/stretchr/testify/require"
)

func TestAppInit_createAppDirectory(t *testing.T) {
	type fields struct {
		AppName  string
		Language ir.Lang
		Path     string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "creates the ruby app directory",
			fields: fields{
				AppName:  "createappdir",
				Language: ir.Ruby,
				Path:     t.TempDir(),
			},
			wantErr: false,
		},
		{
			name: "creates the go app directory",
			fields: fields{
				AppName:  "createappdir",
				Language: ir.GoLang,
				Path:     t.TempDir(),
			},
			wantErr: false,
		},
		{
			name: "creates the js app directory",
			fields: fields{
				AppName:  "createappdir",
				Language: ir.JavaScript,
				Path:     t.TempDir(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AppInit{
				AppName:  tt.fields.AppName,
				Language: tt.fields.Language,
				Path:     tt.fields.Path,
			}
			if err := a.createAppDirectory(); (err != nil) != tt.wantErr {
				t.Errorf("AppInit.createAppDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}

			require.DirExists(t, filepath.Join(tt.fields.Path, tt.fields.AppName))
		})
	}
}

func TestAppInit_createFixtures(t *testing.T) {
	type fields struct {
		AppName  string
		Language ir.Lang
		Path     string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "creates the ruby fixtures directory",
			fields: fields{
				AppName:  "createfixtures",
				Language: ir.Ruby,
				Path:     t.TempDir(),
			},
			want:    "demo.json",
			wantErr: false,
		},
		{
			name: "creates the go fixtures directory",
			fields: fields{
				AppName:  "createfixtures",
				Language: ir.GoLang,
				Path:     t.TempDir(),
			},
			want:    "demo-cdc.json",
			wantErr: false,
		},
		{
			name: "creates the js fixtures directory",
			fields: fields{
				AppName:  "createfixtures",
				Language: ir.JavaScript,
				Path:     t.TempDir(),
			},
			want:    "demo-cdc.json",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AppInit{
				AppName:  tt.fields.AppName,
				Language: tt.fields.Language,
				Path:     tt.fields.Path,
			}

			a.createAppDirectory()
			if err := a.createFixtures(); (err != nil) != tt.wantErr {
				t.Errorf("AppInit.createFixtures() error = %v, wantErr %v", err, tt.wantErr)
			}
			require.DirExists(t, filepath.Join(tt.fields.Path, tt.fields.AppName, "fixtures"))
			require.FileExists(t, filepath.Join(tt.fields.Path, tt.fields.AppName, "fixtures", tt.want))
		})
	}
}

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
				AppName:  tt.fields.AppName,
				Language: tt.fields.Language,
				Path:     tt.fields.Path,
			}
			a.createAppDirectory()
			if err := a.duplicateFile(tt.args.fileName); (err != nil) != tt.wantErr {
				t.Errorf("AppInit.duplicateFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			require.FileExists(t, filepath.Join(tt.fields.Path, tt.fields.AppName, tt.args.fileName))
		})
	}
}

func TestAppInit_listTemplateContent(t *testing.T) {
	type fields struct {
		AppName  string
		Language ir.Lang
		Path     string
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
				AppName:  "testapp",
				Language: ir.Ruby,
				Path:     t.TempDir(),
			},
			want:  []string{"Gemfile", "app.json", "app.rb"},
			want1: []string{"fixtures"},
		},
		{
			name: "lists files and dir content for go app embedded template",
			fields: fields{
				AppName:  "testapp",
				Language: ir.GoLang,
				Path:     t.TempDir(),
			},
			want:  []string{"README.md", "app.go", "app.json", "app_test.go"},
			want1: []string{"fixtures"},
		},
		{
			name: "lists files and dir content for js app embedded template",
			fields: fields{
				AppName:  "testapp",
				Language: ir.JavaScript,
				Path:     t.TempDir(),
			},
			want:  []string{"README.md", "app.json", "index.js", "index.test.js", "package.json"},
			want1: []string{"fixtures"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AppInit{
				AppName:  tt.fields.AppName,
				Language: tt.fields.Language,
				Path:     tt.fields.Path,
			}
			got, got1, err := a.listTemplateContent()
			if (err != nil) != tt.wantErr {
				t.Errorf("AppInit.listTemplateContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, tt.want, got)
			require.Equal(t, tt.want1, got1)
		})
	}
}

func TestAppInit_Init(t *testing.T) {
	type fields struct {
		AppName  string
		Language ir.Lang
		Path     string
	}
	tests := []struct {
		name            string
		fields          fields
		wantFiles       []string
		wantFixtureFile string
		wantErr         bool
	}{
		{
			name: "copies the ruby app template to the path",
			fields: fields{
				AppName:  "testapp",
				Language: ir.Ruby,
				Path:     t.TempDir(),
			},
			wantFiles:       []string{"app.json", "app.rb", "Gemfile"},
			wantFixtureFile: "demo.json",
			wantErr:         false,
		},
		{
			name: "copies the go app template to the path",
			fields: fields{
				AppName:  "testapp",
				Language: ir.GoLang,
				Path:     t.TempDir(),
			},
			wantFiles:       []string{"app.json", "app_test.go", "app.go", "README.md"},
			wantFixtureFile: "demo-no-cdc.json",
			wantErr:         false,
		},
		{
			name: "copies the go app template to the path",
			fields: fields{
				AppName:  "testapp",
				Language: ir.JavaScript,
				Path:     t.TempDir(),
			},
			wantFiles:       []string{"app.json", "package.json", "index.js", "index.test.js", "README.md"},
			wantFixtureFile: "demo-no-cdc.json",
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AppInit{
				AppName:  tt.fields.AppName,
				Language: tt.fields.Language,
				Path:     tt.fields.Path,
			}
			if err := a.Init(); (err != nil) != tt.wantErr {
				t.Errorf("AppInit.Init() error = %v, wantErr %v", err, tt.wantErr)
			}

			require.DirExists(t, filepath.Join(tt.fields.Path, tt.fields.AppName))
			for _, f := range tt.wantFiles {
				require.FileExists(t, filepath.Join(tt.fields.Path, tt.fields.AppName, f))
			}
			require.FileExists(t, filepath.Join(tt.fields.Path, tt.fields.AppName, "fixtures", tt.wantFixtureFile))
		})
	}
}
