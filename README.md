# Mulan 木兰

[![Go Report Card](https://goreportcard.com/badge/github.com/virzz/mulan)](https://goreportcard.com/report/github.com/virzz/mulan)
[![GoDoc](https://godoc.org/github.com/virzz/mulan?status.svg)](https://godoc.org/github.com/virzz/mulan)
[![License](https://img.shields.io/github/license/virzz/mulan.svg)](https://github.com/virzz/mulan/blob/main/LICENSE)

Mulan 是一个基于 Go 语言的轻量级开发框架，旨在简化应用程序开发流程，提供一套全面且高效的工具集。

## 功能特性

- **模块化设计**：核心组件可独立使用，也可无缝集成
- **Web服务**：基于 Gin 的 HTTP 服务，支持路由、中间件和各种常用功能
- **数据库操作**：集成 GORM 提供 ORM 功能
- **Redis 支持**：提供 Redis 客户端封装
- **配置管理**：基于 Viper 的灵活配置系统
- **命令行工具**：基于 Cobra 的 CLI 工具
- ~~**日志系统**：基于 Zap 的结构化日志~~ [mulan-ext/log](https://github.com/mulan-ext/log)
- ~~**接口认证**：提供 API Key 和会话认证~~ [mulan-ext/auth](https://github.com/mulan-ext/auth)
- ~~**验证码系统**：内置图形验证码生成功能~~ [mulan-ext/captcha](https://github.com/mulan-ext/captcha)

## 目录结构

```
mulan/
├── app/            # 应用核心组件和配置管理
├── code/           # 状态码和错误管理
├── db/             # 数据库操作和模型接口
├── rdb/            # Redis 客户端封装
├── req/            # 请求处理相关
├── rsp/            # 响应处理相关
│   └── apperr/     # 状态码和错误管理
├── tests/          # 测试代码
└── web/            # Web 服务和路由管理
```

## 安装

### 前置条件

- Go 1.24+
- MySQL 或 PostgreSQL（可选）
- Redis（可选）

### 安装步骤

```bash
go get https://github.com/virzz/mulan.git
```

## 使用方法

[webx-template](about:blank)

## 贡献

欢迎提交问题和合并请求，共同改进这个项目。

## 许可证

本项目基于 [MIT](LICENSE) 开源。
