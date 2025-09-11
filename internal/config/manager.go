package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

type Manager struct {
}

type ConfigInfo struct {
	Name      string
	FullPath  string
	IsCurrent bool
}

func NewManager() *Manager {
	return &Manager{}
}

// Switch 切换指定客户端和类型的配置文件
// 移除现有符号链接并创建新的符号链接指向目标配置文件
func (m *Manager) Switch(clientName, typeName, configName string) error {
	client, err := FindClient(clientName)
	if err != nil {
		return err
	}

	fileSpec, err := FindFileSpec(clientName, typeName)
	if err != nil {
		return err
	}

	clientDir := client.ExpandDir()
	sourceFile := fileSpec.BuildSourcePath(clientDir, configName)
	targetFile := filepath.Join(clientDir, fileSpec.LinkName)

	if !m.configExists(sourceFile) {
		return fmt.Errorf("配置文件不存在: %s", sourceFile)
	}

	realPath, err := filepath.Abs(sourceFile)
	if err != nil {
		return fmt.Errorf("获取配置文件绝对路径失败: %w", err)
	}

	if err := m.removeExistingSymlink(targetFile); err != nil {
		return fmt.Errorf("移除现有符号链接失败: %w", err)
	}

	if err := os.Symlink(realPath, targetFile); err != nil {
		return fmt.Errorf("创建符号链接失败: %w", err)
	}

	slog.Debug("Created symlink", "client", clientName, "type", typeName, "config", configName, "from", targetFile, "to", realPath)
	return nil
}

// List 列出指定客户端和类型的所有可用配置
// 扫描客户端目录下匹配模式的文件，返回配置信息列表
func (m *Manager) List(clientName, typeName string) ([]ConfigInfo, error) {
	client, err := FindClient(clientName)
	if err != nil {
		return nil, err
	}

	fileSpec, err := FindFileSpec(clientName, typeName)
	if err != nil {
		return nil, err
	}

	clientDir := client.ExpandDir()
	pattern := filepath.Join(clientDir, fileSpec.Pattern)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("扫描配置文件失败: %w", err)
	}

	currentConfig, _ := m.GetCurrent(clientName, typeName)
	configs := make([]ConfigInfo, 0, len(matches))

	for _, match := range matches {
		basename := filepath.Base(match)
		configName := fileSpec.ExtractConfigName(basename)
		if configName != "" {
			configs = append(configs, ConfigInfo{
				Name:      configName,
				FullPath:  match,
				IsCurrent: configName == currentConfig,
			})
		}
	}

	slog.Debug("Found configurations", "client", clientName, "type", typeName, "count", len(configs))
	return configs, nil
}

// GetCurrent 获取当前生效的配置名称
// 通过读取符号链接的目标文件路径解析出配置名
func (m *Manager) GetCurrent(clientName, typeName string) (string, error) {
	client, err := FindClient(clientName)
	if err != nil {
		return "", err
	}

	fileSpec, err := FindFileSpec(clientName, typeName)
	if err != nil {
		return "", err
	}

	clientDir := client.ExpandDir()
	linkPath := filepath.Join(clientDir, fileSpec.LinkName)
	
	target, err := os.Readlink(linkPath)
	if err != nil {
		return "", fmt.Errorf("读取符号链接失败: %w", err)
	}

	basename := filepath.Base(target)
	configName := fileSpec.ExtractConfigName(basename)
	if configName == "" {
		return "", fmt.Errorf("无法解析配置名称")
	}

	return configName, nil
}

// configExists 检查配置文件是否存在
// 使用 os.Stat 检查文件系统中的文件状态
func (m *Manager) configExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// removeExistingSymlink 移除已存在的符号链接
// 如果路径存在且是符号链接则删除，否则静默成功
func (m *Manager) removeExistingSymlink(path string) error {
	if _, err := os.Lstat(path); err == nil {
		slog.Debug("Removing existing symlink", "path", path)
		return os.Remove(path)
	}
	return nil
}