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

// ExpandDir 展开客户端目录中的 ~ 为用户主目录路径
// 支持以 ~/ 开头的相对路径表示法，自动替换为绝对路径
func (c *Client) ExpandDir() string {
	if strings.HasPrefix(c.Dir, "~/") {
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, c.Dir[2:])
	}
	return c.Dir
}

// BuildSourcePath 根据配置名构建源文件完整路径
// 将文件模式中的通配符 * 替换为具体的配置名
func (fs *FileSpec) BuildSourcePath(clientDir, configName string) string {
	pattern := fs.Pattern
	placeholder := strings.Replace(pattern, "*", configName, 1)
	return filepath.Join(clientDir, placeholder)
}

// ExtractConfigName 从文件名中提取配置名
// 根据文件模式从通配符 * 位置提取配置名称部分
// 如果文件名不匹配模式，返回空字符串
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