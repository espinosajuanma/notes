package notes

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed configs
var configs embed.FS

type Config struct {
	Template    string
	UserConfig  string
	Initialized bool
	Name        string `json:"name"`
	Author      struct {
		Name   string `json:"name"`
		Github string `json:"github"`
	} `json:"author"`
	Version    string `json:"version"`
	Categories []struct {
		Title       string            `json:"title"`
		Name        string            `json:"name"`
		Hidden      bool              `json:"hidden"`
		Transitions map[string]string `json:"transitions"`
		Weight      int               `json:"weight"`
		Starter     bool              `json:"starter"`
	} `json:"categories"`
}

func (c *Config) Init() error {
	if c.Template != "" {
		return c.ImportTemplate()
	}
	if c.UserConfig != "" {
		return c.ImportUserConfig()
	}
	return fmt.Errorf("need a template or user config file")
}

func (c *Config) ImportUserConfig() error {
	_, err := os.Stat(c.UserConfig)
	if err != nil {
		return nil
	}
	data, err := os.ReadFile(c.UserConfig)
	if err != nil {
		return err
	}
	config, err := ParseConfig(data)
	if err != nil {
		return err
	}
	c.Categories = append(c.Categories, config.Categories...)
	return nil
}

func (c *Config) ImportTemplate() error {
	filePath := "configs/" + c.Template + ".json"
	data, err := configs.ReadFile(filePath)
	if err != nil {
		return err
	}
	config, err := ParseConfig(data)
	if err != nil {
		return err
	}
	c.Categories = append(c.Categories, config.Categories...)
	return nil
}

func ParseConfig(data []byte) (*Config, error) {
	var t Config
	err := json.Unmarshal(data, &t)
	if err != nil {
		return &Config{}, err
	}
	t.Initialized = true
	return &t, nil
}

func (c *Config) GetTemplates() ([]string, error) {
	names := []string{}
	files, err := configs.ReadDir("configs")
	if err != nil {
		return names, err
	}
	for _, file := range files {
		name := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		names = append(names, name)
	}
	return names, nil
}

func (c *Config) Export() ([]byte, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}
