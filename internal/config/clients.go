package config

import (
	"fmt"
)

var Clients = []Client{
	{Name: "claude", Dir: "~/.claude"},
	{Name: "codex", Dir: "~/.codex"},
	{Name: "gemini", Dir: "~/.gemini"},
	{Name: "qwen", Dir: "~/.qwen"},
}

var ConfigTypes = []ConfigType{
	{
		Name: "config",
		Files: []FileSpec{
			{Client: "claude", Pattern: "settings.*.json", LinkName: "settings.json"},
			{Client: "codex", Pattern: "config.*.toml", LinkName: "config.toml"},
			{Client: "gemini", Pattern: "settings.*.json", LinkName: "settings.json"},
			{Client: "qwen", Pattern: "settings.*.json", LinkName: "settings.json"},
		},
	},
	{
		Name: "agents",
		Files: []FileSpec{
			{Client: "claude", Pattern: "CLAUDE.*.md", LinkName: "CLAUDE.md"},
			{Client: "codex", Pattern: "AGENTS.*.md", LinkName: "AGENTS.md"},
			{Client: "gemini", Pattern: "GEMINI.*.md", LinkName: "GEMINI.md"},
			{Client: "qwen", Pattern: "QWEN.*.md", LinkName: "QWEN.md"},
		},
	},
}

// FindClient 根据客户端名称查找客户端配置
// 返回匹配的客户端指针，如果不存在则返回错误
func FindClient(clientName string) (*Client, error) {
	for i := range Clients {
		if Clients[i].Name == clientName {
			return &Clients[i], nil
		}
	}
	return nil, fmt.Errorf("客户端不存在: %s", clientName)
}

// FindConfigType 根据类型名称查找配置类型
// 返回匹配的配置类型指针，如果不存在则返回错误
func FindConfigType(typeName string) (*ConfigType, error) {
	for i := range ConfigTypes {
		if ConfigTypes[i].Name == typeName {
			return &ConfigTypes[i], nil
		}
	}
	return nil, fmt.Errorf("配置类型不存在: %s", typeName)
}

// FindFileSpec 查找指定客户端和类型的文件规格
// 首先查找配置类型，然后在其中查找匹配的客户端文件规格
func FindFileSpec(clientName, typeName string) (*FileSpec, error) {
	configType, err := FindConfigType(typeName)
	if err != nil {
		return nil, err
	}
	
	for i := range configType.Files {
		if configType.Files[i].Client == clientName {
			return &configType.Files[i], nil
		}
	}
	
	return nil, fmt.Errorf("客户端 %s 不支持类型 %s", clientName, typeName)
}

// GetClientNames 获取所有客户端名称列表
// 用于命令行补全和显示可用选项
func GetClientNames() []string {
	names := make([]string, len(Clients))
	for i, client := range Clients {
		names[i] = client.Name
	}
	return names
}

// GetConfigTypeNames 获取所有配置类型名称列表
// 用于命令行补全和显示可用选项
func GetConfigTypeNames() []string {
	names := make([]string, len(ConfigTypes))
	for i, configType := range ConfigTypes {
		names[i] = configType.Name
	}
	return names
}

// GetDefaultClient 获取默认客户端（claude）
// 当用户未指定客户端时使用
func GetDefaultClient() *Client {
	return &Clients[0] // claude
}

// GetDefaultConfigType 获取默认配置类型（config）
// 当用户未指定配置类型时使用
func GetDefaultConfigType() *ConfigType {
	return &ConfigTypes[0] // config
}