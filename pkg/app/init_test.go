package app

import (
	"github.com/meroxa/turbine-core/pkg/ir"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppInit_createAppDirectory(t *testing.T) {
	type fields struct {
		AppName  string
		Language string
		Path     string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "creates the app directory",
			fields: fields{
				AppName:  "createappdir",
				Language: string(ir.Ruby),
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
		Language string
		Path     string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "creates the fixtures directory",
			fields: fields{
				AppName:  "createfixtures",
				Language: string(ir.Ruby),
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

			a.createAppDirectory()
			if err := a.createFixtures(); (err != nil) != tt.wantErr {
				t.Errorf("AppInit.createFixtures() error = %v, wantErr %v", err, tt.wantErr)
			}
			require.DirExists(t, filepath.Join(tt.fields.Path, tt.fields.AppName, "fixtures"))
			require.FileExists(t, filepath.Join(tt.fields.Path, tt.fields.AppName, "fixtures", "demo.json"))
		})
	}
}

func TestAppInit_duplicateFile(t *testing.T) {
	type fields struct {
		AppName  string
		Language string
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
			name: "duplicates the file from the embedded fs",
			fields: fields{
				AppName:  "duplicatefile",
				Language: string(ir.Ruby),
				Path:     t.TempDir(),
			},
			args: args{
				fileName: "Gemfile",
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
			require.FileExists(t, filepath.Join(tt.fields.Path, tt.fields.AppName, "Gemfile"))
		})
	}
}

func TestAppInit_listTemplateContent(t *testing.T) {
	type fields struct {
		AppName  string
		Language string
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
				Language: string(ir.Ruby),
				Path:     t.TempDir(),
			},
			want:  []string{"Gemfile", "app.json", "app.rb"},
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
		Language string
		Path     string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "copies the ruby app template to the path",
			fields: fields{
				AppName:  "testapp",
				Language: string(ir.Ruby),
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
			if err := a.Init(); (err != nil) != tt.wantErr {
				t.Errorf("AppInit.Init() error = %v, wantErr %v", err, tt.wantErr)
			}

			require.DirExists(t, filepath.Join(tt.fields.Path, tt.fields.AppName))
			require.FileExists(t, filepath.Join(tt.fields.Path, tt.fields.AppName, "app.json"))
			require.FileExists(t, filepath.Join(tt.fields.Path, tt.fields.AppName, "app.rb"))
			require.FileExists(t, filepath.Join(tt.fields.Path, tt.fields.AppName, "Gemfile"))
			require.FileExists(t, filepath.Join(tt.fields.Path, tt.fields.AppName, "fixtures", "demo.json"))
		})
	}
}
