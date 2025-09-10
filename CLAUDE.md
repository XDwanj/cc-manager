# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

cc-manager 是一个多客户端配置管理器，通过符号链接方式管理AI工具的配置文件。

**支持的客户端**：
- `claude` → `~/.claude/`
- `codex` → `~/.codex/`  
- `gemini` → `~/.gemini/`

**支持的配置类型**：
- `config`: 客户端设置文件 (settings.json, config.toml 等)
- `agents`: 指令文件 (CLAUDE.md, AGENTS.md, GEMINI.md 等)

## 开发命令

- `go run main.go <args>` - 运行程序
- `go build` - 构建二进制文件  
- `go test ./...` - 运行测试

## 技术要求

- **Go 版本**: 1.22+
- **主要依赖**: Cobra CLI 框架 (v1.8.1)

## CLI 命令架构

### 全局选项
- `--client <name>` - 指定客户端 (claude/codex/gemini，默认: claude)
- `-v` - 启用详细日志输出

### 命令结构
```
cc-manager [--client=<client>] <command> [--type=<type>] [args...]

默认值: --client=claude --type=config
```

### 列表命令
```bash
cc-manager ls                                    # 列出 claude config
cc-manager --client=claude ls --type=agents     # 列出 claude agents
cc-manager --client=codex ls --type=config      # 列出 codex config
cc-manager ls -d                                # 显示详细信息（文件路径）
```

### 切换命令
```bash
cc-manager switch <名称>                        # 切换 claude config
cc-manager --client=claude switch <名称> --type=agents  # 切换 claude agents
cc-manager --client=codex switch <名称>         # 切换 codex config
```

### 补全命令
- `cc-manager completion bash/zsh/fish/powershell`
- 支持智能 tab 补全：客户端名、配置类型、配置名

## 代码架构

**核心设计**：Client × ConfigType → FileSpec 矩阵

**关键数据结构**：
- `Client`: 客户端定义 (name, dir)
- `ConfigType`: 配置类型定义 (name, files)  
- `FileSpec`: 文件规格 (client, pattern, linkName)
- `ConfigInfo`: 配置信息 (name, fullPath, isCurrent)

**关键组件**：
- `main.go`: slog日志初始化
- `cmd/root.go`: 全局标志 (`--client`, `-v`) 和补全
- `cmd/ls.go`: 列表命令，多维度筛选
- `cmd/switch.go`: 切换命令，智能参数补全
- `internal/config/clients.go`: 客户端和类型定义，查找函数
- `internal/config/client.go`: Client/FileSpec 结构和路径处理
- `internal/config/manager.go`: 符号链接操作的核心逻辑

**工作原理**：
1. 根据客户端名在 `clients.go` 查找 Client 定义
2. 根据类型名查找对应的 FileSpec (pattern, linkName)
3. 扫描客户端目录下匹配 pattern 的文件
4. 切换时移除现有符号链接，创建新链接指向目标配置

**文件模式映射**：
```
claude + config  → settings.*.json → settings.json
claude + agents  → CLAUDE.*.md     → CLAUDE.md
codex + config   → config.*.toml   → config.toml  
codex + agents   → AGENTS.*.md     → AGENTS.md
gemini + config  → settings.*.json → settings.json
gemini + agents  → GEMINI.*.md     → GEMINI.md
```

## 开发约定

- 用户界面文本使用中文
- 使用 `slog` 结构化日志，支持 `-v` 调试级别
- 错误消息用 `fmt.Errorf` 包装上下文