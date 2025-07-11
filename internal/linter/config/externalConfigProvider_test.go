package config_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/maramkhaledn/protolint/internal/setting_test"

	"github.com/maramkhaledn/protolint/internal/linter/config"
	"github.com/maramkhaledn/protolint/linter/rule"
)

func TestGetExternalConfig(t *testing.T) {
	for _, test := range []struct {
		name               string
		inputFilePath      string
		inputDirPath       string
		cwdPath            string
		wantExternalConfig *config.ExternalConfig
		wantExistErr       bool
	}{
		{
			name:         "invalid config file",
			inputDirPath: setting_test.TestDataPath("invalidconfig"),
			wantExistErr: true,
		},
		{
			name: "not found a config file",
		},
		{
			name:         "valid config file",
			inputDirPath: setting_test.TestDataPath("validconfig"),
			wantExternalConfig: &config.ExternalConfig{
				SourcePath: setting_test.TestDataPath("validconfig", "protolint.yaml"),
				Lint: config.Lint{
					Ignores: []config.Ignore{
						{
							ID: "ENUM_FIELD_NAMES_UPPER_SNAKE_CASE",
							Files: []string{
								"path/to/foo.proto",
								"path/to/bar.proto",
							},
						},
						{
							ID: "ENUM_NAMES_UPPER_CAMEL_CASE",
							Files: []string{
								"path/to/foo.proto",
							},
						},
					},
					Rules: struct {
						NoDefault  bool     `yaml:"no_default" json:"no_default" toml:"no_default"`
						AllDefault bool     `yaml:"all_default" json:"all_default" toml:"all_default"`
						Add        []string `yaml:"add" json:"add" toml:"add"`
						Remove     []string `yaml:"remove" json:"remove" toml:"remove"`
					}{
						NoDefault: true,
						Add: []string{
							"FIELD_NAMES_LOWER_SNAKE_CASE",
							"MESSAGE_NAMES_UPPER_CAMEL_CASE",
						},
						Remove: []string{
							"RPC_NAMES_UPPER_CAMEL_CASE",
						},
					},
					RulesOption: config.RulesOption{
						MaxLineLength: config.MaxLineLengthOption{
							CustomizableSeverityOption: config.CustomizableSeverityOption{
								Severity: rule.SeverityNote,
							},
							MaxChars: 80,
							TabChars: 2,
						},
						Indent: config.IndentOption{
							CustomizableSeverityOption: config.CustomizableSeverityOption{
								Severity: rule.SeverityWarning,
							},
							Style:   "\t",
							Newline: "\n",
						},
					},
				},
			},
		},
		{
			name:         "load .protolint.yaml",
			inputDirPath: setting_test.TestDataPath("validconfig", "hidden"),
			wantExternalConfig: &config.ExternalConfig{
				SourcePath: setting_test.TestDataPath("validconfig", "hidden", ".protolint.yaml"),
				Lint: config.Lint{
					RulesOption: config.RulesOption{
						Indent: config.IndentOption{
							Style:   "\t",
							Newline: "\n",
						},
					},
				},
			},
		},
		{
			name:          "load my_protolint.yaml",
			inputFilePath: setting_test.TestDataPath("validconfig", "particular_name", "my_protolint.yaml"),
			wantExternalConfig: &config.ExternalConfig{
				SourcePath: setting_test.TestDataPath("validconfig", "particular_name", "my_protolint.yaml"),
				Lint: config.Lint{
					RulesOption: config.RulesOption{
						Indent: config.IndentOption{
							Style:   "\t",
							Newline: "\n",
						},
					},
				},
			},
		},
		{
			name:         "load protolint.yml",
			inputDirPath: setting_test.TestDataPath("validconfig", "yml"),
			wantExternalConfig: &config.ExternalConfig{
				SourcePath: setting_test.TestDataPath("validconfig", "yml", "protolint.yml"),
				Lint: config.Lint{
					RulesOption: config.RulesOption{
						Indent: config.IndentOption{
							Style:   "\t",
							Newline: "\n",
						},
					},
				},
			},
		},
		{
			name:         "load .protolint.yml",
			inputDirPath: setting_test.TestDataPath("validconfig", "yml_hidden"),
			wantExternalConfig: &config.ExternalConfig{
				SourcePath: setting_test.TestDataPath("validconfig", "yml_hidden", ".protolint.yml"),
				Lint: config.Lint{
					RulesOption: config.RulesOption{
						Indent: config.IndentOption{
							Style:   "\t",
							Newline: "\n",
						},
					},
				},
			},
		},
		{
			name:    "load .protolint.yml at cwd automatically",
			cwdPath: setting_test.TestDataPath("validconfig", "default"),
			wantExternalConfig: &config.ExternalConfig{
				SourcePath: ".protolint.yml",
				Lint: config.Lint{
					RulesOption: config.RulesOption{
						Indent: config.IndentOption{
							Style:   "\t",
							Newline: "\n",
						},
					},
				},
			},
		},
		{
			name:    "prefer .protolint.yml at cwd to one at its parent dir",
			cwdPath: setting_test.TestDataPath("validconfig", "default", "child"),
			wantExternalConfig: &config.ExternalConfig{
				SourcePath: ".protolint.yaml",
				Lint: config.Lint{
					RulesOption: config.RulesOption{
						Indent: config.IndentOption{
							Style:   "\t",
							Newline: "\n",
						},
					},
				},
			},
		},
		{
			name:    "locate .protolint.yml at the parent when not found at cwd",
			cwdPath: setting_test.TestDataPath("validconfig", "default", "empty_child"),
			wantExternalConfig: &config.ExternalConfig{
				SourcePath: setting_test.TestDataPath("validconfig", "default", ".protolint.yml"),
				Lint: config.Lint{
					RulesOption: config.RulesOption{
						Indent: config.IndentOption{
							Style:   "\t",
							Newline: "\n",
						},
					},
				},
			},
		},
		{
			name:    "locate .protolint.yml at the grand parent when not found at cwd",
			cwdPath: setting_test.TestDataPath("validconfig", "default", "empty_child", "empty_grand_child"),
			wantExternalConfig: &config.ExternalConfig{
				SourcePath: setting_test.TestDataPath("validconfig", "default", ".protolint.yml"),
				Lint: config.Lint{
					RulesOption: config.RulesOption{
						Indent: config.IndentOption{
							Style:   "\t",
							Newline: "\n",
						},
					},
				},
			},
		},
		{
			name:         "not found a config file even so inputDirPath is set",
			inputDirPath: setting_test.TestDataPath("validconfig", "particular_name"),
			wantExistErr: true,
		},
		{
			name:          "not found a config file even so inputFilePath is set",
			inputFilePath: setting_test.TestDataPath("validconfig", "particular_name", "not_found.yaml"),
			wantExistErr:  true,
		},
		{
			name:    "found a package.json with 'protolint' in it",
			cwdPath: setting_test.TestDataPath("js_config", "package"),
			wantExternalConfig: &config.ExternalConfig{
				SourcePath: "package.json",
				Lint: config.Lint{
					RulesOption: config.RulesOption{
						Indent: config.IndentOption{
							Style:   "\t",
							Newline: "\n",
						},
					},
				},
			},
		},
		{
			name:               "found a package.json without 'protolint' in it",
			cwdPath:            setting_test.TestDataPath("js_config", "package_no_protolint"),
			wantExternalConfig: nil,
		},
		{
			name:               "found a pyproject.toml without 'tools.protolint' in it",
			cwdPath:            setting_test.TestDataPath("py_project", "pyproject_no_protolint"),
			wantExternalConfig: nil,
		},
		{
			name:               "found a pyproject.toml without 'tools' in it",
			cwdPath:            setting_test.TestDataPath("py_project", "pyproject_no_tools"),
			wantExternalConfig: nil,
		},
		{
			name:    "found a package.json with 'protolint' and other stuff in it",
			cwdPath: setting_test.TestDataPath("js_config", "non_pure_package"),
			wantExternalConfig: &config.ExternalConfig{
				SourcePath: "package.json",
				Lint: config.Lint{
					RulesOption: config.RulesOption{
						Indent: config.IndentOption{
							Style:   "\t",
							Newline: "\n",
						},
					},
				},
			},
		},
		{
			name:    "found a pyproject.toml with 'tools.protolint' and other stuff in it",
			cwdPath: setting_test.TestDataPath("py_project", "with_pyproject"),
			wantExternalConfig: &config.ExternalConfig{
				SourcePath: "pyproject.toml",
				Lint: config.Lint{
					RulesOption: config.RulesOption{
						Indent: config.IndentOption{
							Style:   "\t",
							Newline: "\n",
						},
					},
				},
			},
		},
		{
			name:    "found a package.json with 'protolint' and other stuff in it, but superseded by sibling protolint.yaml",
			cwdPath: setting_test.TestDataPath("js_config", "package_with_yaml"),
			wantExternalConfig: &config.ExternalConfig{
				SourcePath: "protolint.yaml",
				Lint: config.Lint{
					RulesOption: config.RulesOption{
						Indent: config.IndentOption{
							Style:   "\t",
							Newline: "\n",
						},
					},
				},
			},
		},
		{
			name:    "found a pyproject.toml with 'toolsprotolint' and other stuff in it, but superseded by sibling protolint.yaml",
			cwdPath: setting_test.TestDataPath("py_project", "project_with_yaml"),
			wantExternalConfig: &config.ExternalConfig{
				SourcePath: "protolint.yaml",
				Lint: config.Lint{
					RulesOption: config.RulesOption{
						Indent: config.IndentOption{
							Style:   "\t",
							Newline: "\n",
						},
					},
				},
			},
		},
		{
			name:    "found a package.json with 'protolint' and other stuff in it, but superseded by parent protolint.yaml",
			cwdPath: setting_test.TestDataPath("js_config", "package_with_yaml_parent", "child"),
			wantExternalConfig: &config.ExternalConfig{
				SourcePath: setting_test.TestDataPath("js_config", "package_with_yaml_parent", "protolint.yaml"),
				Lint: config.Lint{
					RulesOption: config.RulesOption{
						Indent: config.IndentOption{
							Style:   "\t",
							Newline: "\n",
						},
					},
				},
			},
		},
		{
			name:    "found a pyproject.toml with 'toolsprotolint' and other stuff in it, but superseded by parent protolint.yaml",
			cwdPath: setting_test.TestDataPath("py_project", "project_with_yaml_parent", "child"),
			wantExternalConfig: &config.ExternalConfig{
				SourcePath: setting_test.TestDataPath("py_project", "project_with_yaml_parent", "protolint.yaml"),
				Lint: config.Lint{
					RulesOption: config.RulesOption{
						Indent: config.IndentOption{
							Style:   "\t",
							Newline: "\n",
						},
					},
				},
			},
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if len(test.cwdPath) != 0 {
				err := os.Chdir(test.cwdPath)
				if err != nil {
					t.Errorf("got err %v", err)
					return
				}
			}

			got, err := config.GetExternalConfig(test.inputFilePath, test.inputDirPath)
			if test.wantExistErr {
				if err == nil {
					t.Errorf("got err nil, but want err")
				}
				return
			}
			if err != nil {
				t.Errorf("got err %v, but want nil", err)
				return
			}

			if !reflect.DeepEqual(got, test.wantExternalConfig) {
				t.Errorf("got %v, but want %v", got, test.wantExternalConfig)
			}
		})
	}
}
