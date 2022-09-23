# Adanos-alert

[![Build Status](https://www.travis-ci.org/mylxsw/adanos-alert.svg?branch=master)](https://www.travis-ci.org/mylxsw/adanos-alert)
[![Coverage Status](https://coveralls.io/repos/github/mylxsw/adanos-alert/badge.svg?branch=master)](https://coveralls.io/github/mylxsw/adanos-alert?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/mylxsw/adanos-alert)](https://goreportcard.com/report/github.com/mylxsw/adanos-alert)
[![codecov](https://codecov.io/gh/mylxsw/adanos-alert/branch/master/graph/badge.svg)](https://codecov.io/gh/mylxsw/adanos-alert)
[![Sourcegraph](https://sourcegraph.com/github.com/mylxsw/adanos-alert/-/badge.svg)](https://sourcegraph.com/github.com/mylxsw/adanos-alert?badge)
[![GitHub](https://img.shields.io/github/license/mylxsw/adanos-alert.svg)](https://github.com/mylxsw/adanos-alert)

Adanos-alert 是一款使用 Golang 开发的告警通知系统，Adanos-alert 将关注点聚焦在了通知上，提供了丰富的事件聚合和通知机制，弥补了常见监控告警系统在通知方面的不足。


## Features

- [x] 使用 Go 模板语法定制告警事件格式，展示效果
  - [x] 预置多种告警通知事件模板
  - [x] 支持自定义告警通知模板
- [x] 基于自定义表达式对事件进行分组
- [x] 支持基于自定义表达式忽略特定的事件
  - [x] 支持最大忽略事件数，超过阈值后告警
- [x] 支持基于自定义表达式为事件建立关联关系，方便追溯
- [x] 基于多种时间维度聚合事件
  - [x] 按照固定的时间周期（时间间隔）聚合
  - [x] 按照多个固定的时间点聚合
  - [x] 按照多个固定的时间范围，不同时间范围指定不同的聚合周期
- [x] 支持多种告警事件接入方式
  - [x] Logstash 错误日志输出
  - [x] HTTP API 接口
  - [x] 命令行管道（使用 adanos-proxy）
  - [x] Grafana 告警事件接入
  - [x] Prometheus 告警事件接入
  - [x] Prometheus Alertmanager 告警事件接入
  - [x] Openfalcon 告警事件接入
- [x] 多种通知通道支持
  - [x] 钉钉群消息
  - [x] Jira 问题创建
  - [x] 阿里云语音通知
  - [x] HTTP 接口调用
  - [ ] 邮件通知
  - [ ] 微信通知
  - [ ] 阿里云短信
  - [ ] 网易云信短信
- [x] 通过 adanos-alert-agent 实现告警事件接收与调度分离，为事件告警通知提供高可用支持

## 环境依赖

Adanos-alert 使用 MongoDB 作为数据存储后端，所有的数据（告警规则、告警事件、事件组等）均存储在 MongoDB 中，因此需要先 [安装 MongoDB 数据库](https://www.mongodb.com/docs/v4.4/administration/install-community/)，版本建议在 4.0 以上（低版本未测试过，可能存在兼容问题）。

## 手动编译

编译时依赖以下工具

- [esc](https://github.com/mjibson/esc) 是用于将静态文件打包到 Go 程序中的一个工具库

使用下面的命令直接完成编译并运行

```bash
make run
```

Dashboard 访问地址 http://localhost:19999

![预览图](https://ssl.aicode.cc/prometheus/20201025172345.png)

## Architecture

Adanos Alert 平台与其他系统之间的关系

![Adanos Alert 平台与其他系统之间的关系](https://ssl.aicode.cc/prometheus/20201025172918.png)

Adanos Alert 平台各组件之间的关系

![Adanos Alert 平台各组件之间的关系](https://ssl.aicode.cc/prometheus/20201025172846.png)

Adanos Alert Server 内部结构

![Adanos Alert Server 内部结构](https://ssl.aicode.cc/prometheus/20201025172817.png)

## Related Projects

- [adanos-mail-connector](https://github.com/mylxsw/adanos-mail-connector) 可以伪装成为 SMTP 服务器，将邮件转换为 Adanos 事件发送给 Adanos-alert Server
- [Glacier Framework](https://github.com/mylxsw/glacier) Go 语言应用开发框架
- [Container](https://github.com/mylxsw/container) Go 语言运行时依赖注入框架
- [Asteria](https://github.com/mylxsw/asteria) 结构化日志库
