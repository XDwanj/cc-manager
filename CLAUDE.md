# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

cc-manager 是一个用于管理 Claude 配置的 Go CLI 工具，支持两种配置类型：
- **settings 配置**: Claude 的 settings.json 配置文件
- **CLAUDE 配置**: Claude 的 CLAUDE.md 指令文件

通过符号链接方式在不同配置间切换。

## 开发命令

- `go run main.go <args>` - 运行程序
- `go build` - 构建二进制文件  
- `go test ./...` - 运行测试

## 技术要求

- **Go 版本**: 1.22+
- **主要依赖**: Cobra CLI 框架 (v1.8.1)

## CLI 命令

### 列表命令
- `cc-manager ls` - 列出 settings 配置
- `cc-manager ls --claude` - 列出 CLAUDE 配置
- `cc-manager ls -d` - 显示详细信息（文件路径）

### 切换命令
- `cc-manager switch <名称>` - 切换 settings 配置
- `cc-manager switch <名称> --claude` - 切换 CLAUDE 配置
- 支持智能 tab 补全，根据 `--claude` 标志动态补全可用配置名

### 补全命令
- `cc-manager completion bash` - 生成 bash 补全脚本
- `cc-manager completion zsh` - 生成 zsh 补全脚本
- `cc-manager completion fish` - 生成 fish 补全脚本
- `cc-manager completion powershell` - 生成 PowerShell 补全脚本

### 全局选项
- `-v` - 启用详细日志输出

## 代码架构

**核心数据结构**: `ConfigInfo` (name, fullPath, isCurrent)

**关键组件**:
- `main.go`: 设置 slog 日志并启动 CLI
- `cmd/root.go`: Cobra 根命令和全局 `-v` 标志
- `cmd/ls.go`: 列表命令，支持 `--claude` 和 `-d` 标志
- `cmd/switch.go`: 切换命令，支持 `--claude` 标志和智能 tab 补全
- `cmd/completion.go`: 生成各种 shell 的自动补全脚本
- `internal/config/manager.go`: 核心管理器，处理两种配置类型

**工作原理**:
1. 扫描 `~/.claude/` 下的 `settings.*.json` 和 `CLAUDE.*.md` 文件
2. 通过符号链接 `settings.json` → `settings.<name>.json`
3. 通过符号链接 `CLAUDE.md` → `CLAUDE.<name>.md`
4. 切换时移除旧链接，创建新链接

## 开发约定

- 所有用户界面文本使用中文
- 使用 `slog` 进行结构化日志记录
- 错误消息包含上下文信息并使用 `fmt.Errorf` 包装