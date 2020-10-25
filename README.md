# Adanos-alert

[![Build Status](https://www.travis-ci.org/mylxsw/adanos-alert.svg?branch=master)](https://www.travis-ci.org/mylxsw/adanos-alert)
[![Coverage Status](https://coveralls.io/repos/github/mylxsw/adanos-alert/badge.svg?branch=master)](https://coveralls.io/github/mylxsw/adanos-alert?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/mylxsw/adanos-alert)](https://goreportcard.com/report/github.com/mylxsw/adanos-alert)
[![codecov](https://codecov.io/gh/mylxsw/adanos-alert/branch/master/graph/badge.svg)](https://codecov.io/gh/mylxsw/adanos-alert)
[![Sourcegraph](https://sourcegraph.com/github.com/mylxsw/adanos-alert/-/badge.svg)](https://sourcegraph.com/github.com/mylxsw/adanos-alert?badge)
[![GitHub](https://img.shields.io/github/license/mylxsw/adanos-alert.svg)](https://github.com/mylxsw/adanos-alert)


Adanos-alert is a alert manager with multi alert channel support

## Build 

使用下面的命令直接完成编译并运行

```bash
make run
```

Dashboard 访问地址 http://localhost:19999

![预览图](https://ssl.aicode.cc/prometheus/20201025172345.png)

## Dependency

- esc: https://github.com/mjibson/esc

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
