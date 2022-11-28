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

## 安装运行

以下配置文件可供参考

- **./systemd/** 目录包含 server 和 agent 在 systemd 管理下的启动命令
- **server.config.yaml** 是 adanos-alert-server 的配置文件示例
- **agent.config.yaml** 是 adanos-alert-agent 的配置文件示例
- **logstash.example.conf** 是 Logstash 将错误日志输出到 adanos 的配置示例
- **prometheus.example.yaml** 是 Prometheus 将告警信息接入到 adanos 的配置示例
- **prometheus.rules.example.conf** 是 Prometheus 告警规则示例

adanos-alert-server 命令行选项（配置文件）

- **conf**： configuration file path
- **shutdown_timeout**： set a shutdown timeout for each service (default: 5s) [$GLACIER_SHUTDOWN_TIMOUT]
- **listen**： http listen addr (default: ":19999")
- **grpc_listen**： GRPC Server listen address (default: ":19998")
- **grpc_token**： GRPC Server token (default: "000000")
- **preview_url**： Alert preview page url (default: "http://localhost:19999/ui/groups/%s.html") [$ADANOS_PREVIEW_URL]
- **report_url**： Alert report page url (default: "http://localhost:19999/ui/reports/%s.html") [$ADANOS_REPORT_URL]
- **mongo_uri**： Mongodb connection uri，参考 https://docs.mongodb.com/manual/reference/connection-string/ (default: "mongodb://localhost:27017") [$MONGODB_HOST]
- **mongo_db**： Mongodb database name (default: "adanos-alert") [$MONGODB_DB]
- **api_token**： API Token for api access control [$ADANOS_API_TOKEN]
- **use_local_dashboard**： whether using local dashboard, this is used when development
- **enable_migrate**： whether enable database migrate when app run
- **re_migrate**： 是否重新执行迁移，重新迁移会移除已有的预定义模板
- **aggregation_period**： aggregation job execute period (default: "5s") [$ADANOS_AGGREGATION_PERIOD]
- **action_trigger_period**： action trigger job execute period (default: "5s") [$ADANOS_ACTION_TRIGGER_PERIOD]
- **queue_job_max_retry_times**： set queue job max retry times (default: 3) [$ADANOS_QUEUE_JOB_MAX_RETRY_TIMES]
- **keep_period**： 保留多长时间的报警，如果全部保留，设置为0，单位为天，Adanos-Alert 会自动清理超过 keep_period 天的报警 (default: 0) [$ADANOS_KEEP_PERIOD]
- **syslog_keep_period**： 保留多长时间的系统日志，如果全部保留，设置为0，单位为天，Adanos-Alert 会自动清理超过 syslog_keep_period 天的系统日志 (default: 0) [$ADANOS_SYSLOG_KEEP_PERIOD]
- **queue_worker_num**： set queue worker numbers (default: 3) [$ADANOS_QUEUE_WORKER_NUM]
- **query_timeout**： query timeout for backend service (default: "30s") [$ADANOS_QUERY_TIMEOUT]
- **aliyun_access_key**： 阿里云语音通知接口 Access Key ID [$ADANOS_ALIYUN_ACCESS_KEY]
- **aliyun_access_secret**： 阿里云语音通知接口 Access Secret [$ADANOS_ALIYUN_ACCESS_SECRET]
- **aliyun_voice_called_show_number**： 阿里云语音通知被叫显号
- **aliyun_voice_tts_code**： 阿里云语音通知模板，这里是模板ID，模板内容在阿里云申请，建议内容："您有一条名为 ${title} 的报警通知，请及时处理！"
- **aliyun_voice_tts_param**： 阿里云语音通知模板变量名 (default: "title")
- **log_path**： 日志文件输出目录（非文件名），默认为空，输出到标准输出
- **jira_url**： Jira 服务器地址，如 http://127.0.0.1:8080 [$ADANOS_JIRA_URL]
- **jira_username**： Jira 连接账号 [$ADANOS_JIRA_USERNAME]
- **jira_password**： Jira 连接密码 [$ADANOS_JIRA_PASSWORD]
- **no_job_mode**： 启用该标识后，将会停止事件聚合和队列任务处理，用于开发调试

adanos-alert-agent 命令行选项（配置文件）

- **conf**： configuration file path
- **shutdown_timeout**： set a shutdown timeout for each service (default: 5s) [$GLACIER_SHUTDOWN_TIMOUT]
- **server_addr**： server grpc listen address (default: "127.0.0.1:19998") [$ADANOS_SERVER_ADDR]
- **server_token**： API Token for grpc api access control (default: "000000") [$ADANOS_SERVER_TOKEN]
- **data_dir**： 本地数据库存储目录 (default: "/tmp/adanos-agent")
- **listen**： listen address (default: "127.0.0.1:29999") [$ADANOS_AGENT_LISTEN_ADDR]
- **log_path**： 日志文件输出目录（非文件名），默认为空，输出到标准输出

## 手动编译

编译时依赖以下工具

- [esc](https://github.com/mjibson/esc) 是用于将静态文件打包到 Go 程序中的一个工具库

使用下面的命令直接完成编译并运行

```bash
make run
```

Dashboard 访问地址 http://localhost:19999

![预览图](https://ssl.aicode.cc/prometheus/20201025172345.png)

## 系统接入

TODO

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
- [Go-IOC](https://github.com/mylxsw/go-ioc) Go 语言运行时依赖注入框架
- [Asteria](https://github.com/mylxsw/asteria) 结构化日志库
