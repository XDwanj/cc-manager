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

func (m *Manager) SwitchClaude(claudeName string) error {
	sourceFile := filepath.Join(m.claudeDir, fmt.Sprintf("CLAUDE.%s.md", claudeName))
	targetFile := filepath.Join(m.claudeDir, "CLAUDE.md")

	if !m.configExists(sourceFile) {
		return fmt.Errorf("CLAUDE 配置文件不存在: %s", sourceFile)
	}

	realPath, err := filepath.Abs(sourceFile)
	if err != nil {
		return fmt.Errorf("获取 CLAUDE 配置文件绝对路径失败: %w", err)
	}

	if err := m.removeExistingSymlink(targetFile); err != nil {
		return fmt.Errorf("移除现有 CLAUDE 符号链接失败: %w", err)
	}

	if err := os.Symlink(realPath, targetFile); err != nil {
		return fmt.Errorf("创建 CLAUDE 符号链接失败: %w", err)
	}

	slog.Debug("Created CLAUDE symlink", "from", targetFile, "to", realPath)
	return nil
}

func (m *Manager) ListClaudes() ([]ConfigInfo, error) {
	pattern := filepath.Join(m.claudeDir, "CLAUDE.*.md")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("扫描 CLAUDE 配置文件失败: %w", err)
	}

	currentClaude, _ := m.GetCurrentClaude()
	configs := make([]ConfigInfo, 0, len(matches))

	for _, match := range matches {
		basename := filepath.Base(match)
		if strings.HasPrefix(basename, "CLAUDE.") && strings.HasSuffix(basename, ".md") {
			claudeName := strings.TrimPrefix(strings.TrimSuffix(basename, ".md"), "CLAUDE.")
			
			configs = append(configs, ConfigInfo{
				Name:      claudeName,
				FullPath:  match,
				IsCurrent: claudeName == currentClaude,
			})
		}
	}

	slog.Debug("Found CLAUDE configurations", "count", len(configs))
	return configs, nil
}

func (m *Manager) GetCurrentClaude() (string, error) {
	claudePath := filepath.Join(m.claudeDir, "CLAUDE.md")
	
	target, err := os.Readlink(claudePath)
	if err != nil {
		return "", fmt.Errorf("读取 CLAUDE 符号链接失败: %w", err)
	}

	basename := filepath.Base(target)
	if strings.HasPrefix(basename, "CLAUDE.") && strings.HasSuffix(basename, ".md") {
		claudeName := strings.TrimPrefix(strings.TrimSuffix(basename, ".md"), "CLAUDE.")
		return claudeName, nil
	}

	return "", fmt.Errorf("无法解析 CLAUDE 配置名称")
}