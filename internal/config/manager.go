package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type Manager struct {
	claudeDir string
}

type ConfigInfo struct {
	Name      string
	FullPath  string
	IsCurrent bool
}

func NewManager() *Manager {
	homeDir, _ := os.UserHomeDir()
	return &Manager{
		claudeDir: filepath.Join(homeDir, ".claude"),
	}
}

func (m *Manager) SwitchConfig(configName string) error {
	sourceFile := filepath.Join(m.claudeDir, fmt.Sprintf("settings.%s.json", configName))
	targetFile := filepath.Join(m.claudeDir, "settings.json")

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

	slog.Debug("Created symlink", "from", targetFile, "to", realPath)
	return nil
}

func (m *Manager) configExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (m *Manager) removeExistingSymlink(path string) error {
	if _, err := os.Lstat(path); err == nil {
		slog.Debug("Removing existing symlink", "path", path)
		return os.Remove(path)
	}
	return nil
}

func (m *Manager) ListConfigs() ([]ConfigInfo, error) {
	pattern := filepath.Join(m.claudeDir, "settings.*.json")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("扫描配置文件失败: %w", err)
	}

	currentConfig, _ := m.GetCurrentConfig()
	configs := make([]ConfigInfo, 0, len(matches))

	for _, match := range matches {
		basename := filepath.Base(match)
		if strings.HasPrefix(basename, "settings.") && strings.HasSuffix(basename, ".json") {
			configName := strings.TrimPrefix(strings.TrimSuffix(basename, ".json"), "settings.")
			
			configs = append(configs, ConfigInfo{
				Name:      configName,
				FullPath:  match,
				IsCurrent: configName == currentConfig,
			})
		}
	}

	slog.Debug("Found configurations", "count", len(configs))
	return configs, nil
}

func (m *Manager) GetCurrentConfig() (string, error) {
	settingsPath := filepath.Join(m.claudeDir, "settings.json")
	
	target, err := os.Readlink(settingsPath)
	if err != nil {
		return "", fmt.Errorf("读取符号链接失败: %w", err)
	}

	basename := filepath.Base(target)
	if strings.HasPrefix(basename, "settings.") && strings.HasSuffix(basename, ".json") {
		configName := strings.TrimPrefix(strings.TrimSuffix(basename, ".json"), "settings.")
		return configName, nil
	}

	return "", fmt.Errorf("无法解析配置名称")
}