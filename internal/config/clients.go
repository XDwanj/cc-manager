package config

import (
	"fmt"
)

var Clients = []Client{
	{Name: "claude", Dir: "~/.claude"},
	{Name: "codex", Dir: "~/.codex"},
	{Name: "gemini", Dir: "~/.gemini"},
}

var ConfigTypes = []ConfigType{
	{
		Name: "config",
		Files: []FileSpec{
			{Client: "claude", Pattern: "settings.*.json", LinkName: "settings.json"},
			{Client: "codex", Pattern: "config.*.toml", LinkName: "config.toml"},
			{Client: "gemini", Pattern: "settings.*.json", LinkName: "settings.json"},
		},
	},
	{
		Name: "agents",
		Files: []FileSpec{
			{Client: "claude", Pattern: "CLAUDE.*.md", LinkName: "CLAUDE.md"},
			{Client: "codex", Pattern: "AGENTS.*.md", LinkName: "AGENTS.md"},
			{Client: "gemini", Pattern: "GEMINI.*.md", LinkName: "GEMINI.md"},
		},
	},
}

func FindClient(clientName string) (*Client, error) {
	for i := range Clients {
		if Clients[i].Name == clientName {
			return &Clients[i], nil
		}
	}
	return nil, fmt.Errorf("客户端不存在: %s", clientName)
}

func FindConfigType(typeName string) (*ConfigType, error) {
	for i := range ConfigTypes {
		if ConfigTypes[i].Name == typeName {
			return &ConfigTypes[i], nil
		}
	}
	return nil, fmt.Errorf("配置类型不存在: %s", typeName)
}

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

func GetClientNames() []string {
	names := make([]string, len(Clients))
	for i, client := range Clients {
		names[i] = client.Name
	}
	return names
}

func GetConfigTypeNames() []string {
	names := make([]string, len(ConfigTypes))
	for i, configType := range ConfigTypes {
		names[i] = configType.Name
	}
	return names
}

func GetDefaultClient() *Client {
	return &Clients[0] // claude
}

func GetDefaultConfigType() *ConfigType {
	return &ConfigTypes[0] // config
}