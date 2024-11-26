# Adanos-alert

[![Build Status](https://www.travis-ci.org/mylxsw/adanos-alert.svg?branch=master)](https://www.travis-ci.org/mylxsw/adanos-alert)
[![Coverage Status](https://coveralls.io/repos/github/mylxsw/adanos-alert/badge.svg?branch=master)](https://coveralls.io/github/mylxsw/adanos-alert?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/mylxsw/adanos-alert)](https://goreportcard.com/report/github.com/mylxsw/adanos-alert)
[![codecov](https://codecov.io/gh/mylxsw/adanos-alert/branch/master/graph/badge.svg)](https://codecov.io/gh/mylxsw/adanos-alert)
[![Sourcegraph](https://sourcegraph.com/github.com/mylxsw/adanos-alert/-/badge.svg)](https://sourcegraph.com/github.com/mylxsw/adanos-alert?badge)
[![GitHub](https://img.shields.io/github/license/mylxsw/adanos-alert.svg)](https://github.com/mylxsw/adanos-alert)

Adanos-alert is an alert notification system developed using Golang. It focuses on notifications, offering rich event aggregation and notification mechanisms, addressing the shortcomings of common monitoring alert systems in terms of notifications.

## Features 

- [x] Customize the formatting of alarm events using Go template syntax.
  - [x] Provide multiple built-in templates for alarm notification events.
  - [x] Support custom alarm notification templates.
- [x] Grouping of events based on custom expressions.
- [x] Support for ignoring certain events based on custom expressions.
  - [x] Supports a maximum number of ignored events, after which an alarm is issued.
- [x] Support for establishing relationships between events using custom expressions, for easy tracing.
- [x] Aggregate events based on multiple time dimensions.
  - [x] Aggregation at fixed time intervals.
  - [x] Aggregation at multiple fixed points in time.
  - [x] Aggregation at multiple fixed time ranges, with different aggregation periods for each range.
- [x] Support for multiple alarm event delivery methods.
  - [x] Logstash error log output.
  - [x] HTTP API interfaces.
  - [x] Command-line pipelines (using adanos-proxy).
  - [x] Grafana alarm event delivery.
  - [x] Prometheus alarm event delivery.
  - [x] Prometheus Alertmanager alarm event delivery.
  - [x] Openfalcon alarm event delivery.
- [x] Multiple notification channel support.
  - [x] DingTalk group messages.
  - [x] Jira issue creation.
  - [x] Alibaba Cloud voice notifications.
  - [x] HTTP interface calls.
  - [ ] Email notifications.
  - [ ] WeChat notifications.
  - [ ] Alibaba Cloud SMS.
  - [ ] NetEase Cloud SMS.
- [x] Implement alarm event reception and dispatching separation through adanos-alert-agent for high availability of event alarm notifications.

## Environment

Adanos-alert uses MongoDB as its backend data storage. All data (alert rules, alert events, event groups, etc.) are stored in MongoDB. Therefore, it is necessary to [install the MongoDB database](https://www.mongodb.com/docs/v4.4/administration/install-community/), with a recommended version of 4.0 or higher (lower versions have not been tested and may have compatibility issues).

## Installation

The following configuration files are available for reference:

- The **./systemd/** directory contains startup commands for the server and agent managed under systemd
- **server.config.yaml** is an example configuration file for adanos-alert-server
- **agent.config.yaml** is an example configuration file for adanos-alert-agent
- **logstash.example.conf** is an example of Logstash outputting error logs to adanos
- **prometheus.example.yaml** is an example of Prometheus integrating alert information into adanos
- **prometheus.rules.example.conf** provides examples of Prometheus alert rules

adanos-alert-server command line options (config file)

- **conf**: configuration file path
- **shutdown_timeout**: set a shutdown timeout for each service (default: 5s) [$GLACIER_SHUTDOWN_TIMOUT]
- **listen**: http listen addr (default: ":19999")
- **grpc_listen**: GRPC Server listen address (default: ":19998")
- **grpc_token**: GRPC Server token (default: "000000")
- **preview_url**: Alert preview page url (default: "http://localhost:19999/ui/groups/%s.html") [$ADANOS_PREVIEW_URL]
- **report_url**: Alert report page url (default: "http://localhost:19999/ui/reports/%s.html") [$ADANOS_REPORT_URL]
- **mongo_uri**: Mongodb connection uri，refer to https://docs.mongodb.com/manual/reference/connection-string/ (default: "mongodb://localhost:27017") [$MONGODB_HOST]
- **mongo_db**: Mongodb database name (default: "adanos-alert") [$MONGODB_DB]
- **api_token**: API Token for api access control [$ADANOS_API_TOKEN]
- **use_local_dashboard**: whether using local dashboard, this is used when development
- **enable_migrate**: whether enable database migrate when app run
- **re_migrate**: whether re-execute migrate, re-migrate will remove existing predefined templates
- **aggregation_period**: aggregation job execute period (default: "5s") [$ADANOS_AGGREGATION_PERIOD]
- **action_trigger_period**: action trigger job execute period (default: "5s") [$ADANOS_ACTION_TRIGGER_PERIOD]
- **queue_job_max_retry_times**: set queue job max retry times (default: 3) [$ADANOS_QUEUE_JOB_MAX_RETRY_TIMES]
- **keep_period**: how long to keep the alarm. If all are kept, set it to 0. The unit is the day. Adanos-Alert automatically cleans up alarms that exceed the keep_period day (default: 0) [$ADANOS_KEEP_PERIOD]
- **syslog_keep_period**: how long to keep system logs. If all are kept, set it to 0. The unit is the day. Adanos-Alert automatically cleans up system logs that exceed the syslog_keep_period day (default: 0) [$ADANOS_SYSLOG_KEEP_PERIOD]
- **queue_worker_num**: set queue worker numbers (default: 3) [$ADANOS_QUEUE_WORKER_NUM]
- **query_timeout**: query timeout for backend service (default: "30s") [$ADANOS_QUERY_TIMEOUT]
- **aliyun_access_key**: Aliyun voice notification interface Access Key ID [$ADANOS_ALIYUN_ACCESS_KEY]
- **aliyun_access_secret**: Aliyun voice notification interface Access Secret [$ADANOS_ALIYUN_ACCESS_SECRET]
- **aliyun_voice_called_show_number**: Aliyun voice notification called display number
- **aliyun_voice_tts_code**: Aliyun voice notification template. This is the template ID. The template content is applied for in Aliyun. The recommended content is: "You have an alarm notification named ${title}. Please handle it in time!"
- **aliyun_voice_tts_param**: Aliyun voice notification template variable name (default: "title")
- **log_path**: The log file output directory (not the file name). The default is empty, which outputs to the standard output
- **jira_url**: Jira server address, such as http://127.0.0.1:8080 [$ADANOS_JIRA_URL]
- **jira_username**: Jira connection account [$ADANOS_JIRA_USERNAME]
- **jira_password**: Jira connection password [$ADANOS_JIRA_PASSWORD]
- **no_job_mode**: After this flag is enabled, the event aggregation and queue task processing will be stopped, which is used for development and debugging

adanos-alert-agent CLI options (config file)

- **conf**: configuration file path
- **shutdown_timeout**: set a shutdown timeout for each service (default: 5s) [$GLACIER_SHUTDOWN_TIMOUT]
- **server_addr**: server grpc listen address (default: "127.0.0.1:19998") [$ADANOS_SERVER_ADDR]
- **server_token**: API Token for grpc api access control (default: "000000") [$ADANOS_SERVER_TOKEN]
- **data_dir**: local database storage dir (default: "/tmp/adanos-agent")
- **listen**: listen address (default: "127.0.0.1:29999") [$ADANOS_AGENT_LISTEN_ADDR]
- **log_path**: log file output dir (not filename), empty for stdout by default

## Manual compilation

During compilation, the following tools are required:

- [esc](https://github.com/mjibson/esc) is a tool library used to bundle static files into Go programs.

Execute the following command to compile and run directly:

```bash
make run
```

Access the dashboard at http://localhost:19999

![Preview Image](https://ssl.aicode.cc/prometheus/20201025172345.png)


## Architecture

The relationship between the Adanos Alert platform and other systems

![The relationship between the Adanos Alert platform and other systems](https://ssl.aicode.cc/prometheus/20201025172918.png)

The relationship between the components of the Adanos Alert platform

![The relationship between the components of the Adanos Alert platform](https://ssl.aicode.cc/prometheus/20201025172846.png)

The internal structure of Adanos Alert Server

![The internal structure of Adanos Alert Server](https://ssl.aicode.cc/prometheus/20201025172817.png)

## Related Projects

- [adanos-mail-connector](https://github.com/mylxsw/adanos-mail-connector) can pretend to be an SMTP server, convert emails into Adanos events and send them to Adanos-alert Server
- [Glacier Framework](https://github.com/mylxsw/glacier) Go language application development framework
- [Go-IOC](https://github.com/mylxsw/go-ioc) Go language runtime dependency injection framework
- [Asteria](https://github.com/mylxsw/asteria) Structured logging library
