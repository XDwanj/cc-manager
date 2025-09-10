package config

import (
	"os"
	"path/filepath"
	"strings"
)

type Client struct {
	Name string
	Dir  string
}

type ConfigType struct {
	Name  string
	Files []FileSpec
}

type FileSpec struct {
	Client   string
	Pattern  string
	LinkName string
}

func (c *Client) ExpandDir() string {
	if strings.HasPrefix(c.Dir, "~/") {
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, c.Dir[2:])
	}
	return c.Dir
}

func (fs *FileSpec) BuildSourcePath(clientDir, configName string) string {
	pattern := fs.Pattern
	placeholder := strings.Replace(pattern, "*", configName, 1)
	return filepath.Join(clientDir, placeholder)
}

func (fs *FileSpec) ExtractConfigName(filename string) string {
	pattern := fs.Pattern
	if !strings.Contains(pattern, "*") {
		return ""
	}
	
	parts := strings.Split(pattern, "*")
	if len(parts) != 2 {
		return ""
	}
	
	prefix, suffix := parts[0], parts[1]
	if !strings.HasPrefix(filename, prefix) || !strings.HasSuffix(filename, suffix) {
		return ""
	}
	
	return strings.TrimSuffix(strings.TrimPrefix(filename, prefix), suffix)
}