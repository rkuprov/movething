package cfg

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

type Task struct {
	SearchPattern        string `yaml:"search_pattern"`
	FileExtensionPattern string `yaml:"file_extension_pattern"`
	RenamePattern        string `yaml:"rename_pattern"`
	SearchDirectory      string `yaml:"search_directory"`
	DestinationDirectory string `yaml:"destination_directory"`
}

type Tasks struct {
	Tasks []Task `yaml:"tasks"`
}

func GetConfig(ctx context.Context) ([]Task, error) {
	switch runtime.GOOS {
	case "windows":
		// todo: add windows support
		return nil, nil
	case "linux":
		return getLinuxConfigs(ctx)
	case "darwin":
		return getDarwinConfigs(ctx)
	default:
		return nil, errors.New("unsupported platform")
	}
}

func getLinuxConfigs(_ context.Context) ([]Task, error) {
	bytes, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), ".config", "movething", "config.yaml"))
	if err != nil {
		return nil, err
	}

	var ret Tasks
	err = yaml.Unmarshal(bytes, &ret)
	if err != nil {
		return nil, err
	}

	return ret.Tasks, nil
}

func getDarwinConfigs(_ context.Context) ([]Task, error) {
	bytes, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), ".config", "movething", "config.yaml"))
	if err != nil {
		return nil, err
	}

	var ret Tasks
	err = yaml.Unmarshal(bytes, &ret)
	if err != nil {
		return nil, err
	}

	return ret.Tasks, nil
}
