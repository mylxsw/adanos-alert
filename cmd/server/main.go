package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mylxsw/adanos-alert/internal/llm"

	"github.com/mylxsw/adanos-alert/pubsub"
	"github.com/mylxsw/adanos-alert/rpc"
	"github.com/mylxsw/adanos-alert/service"
	"github.com/mylxsw/asteria/formatter"
	"github.com/mylxsw/asteria/writer"
	"github.com/mylxsw/glacier/event"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/go-utils/str"

	"github.com/mylxsw/adanos-alert/api"
	"github.com/mylxsw/adanos-alert/configs"
	"github.com/mylxsw/adanos-alert/internal/action"
	"github.com/mylxsw/adanos-alert/internal/job"
	"github.com/mylxsw/adanos-alert/internal/queue"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/internal/repository/impl"
	"github.com/mylxsw/adanos-alert/migrate"
	"github.com/mylxsw/asteria/level"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/glacier/starter/app"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Version = "1.0"
var GitCommit = "5dbef13fb456f51a5d29464d"
var AsyncRunner = 3
var DEBUG = ""

func main() {
	if DEBUG == "true" {
		infra.DEBUG = true
	}

	ins := app.Create(fmt.Sprintf("%s (%s)", Version, GitCommit), AsyncRunner).
		WithLogger(log.Module("glacier")).
		WithYAMLFlag("conf")

	ins.AddStringFlag("listen", ":19999", "http server listen address")
	ins.AddStringFlag("api_token", "", "api token for http server")
	ins.AddStringFlag("grpc_listen", ":19998", "grpc server listen address")

	ins.AddStringFlag("mongo_uri", "mongodb://localhost:27017", "mongodb uri, 参考 https://docs.mongodb.com/manual/reference/connection-string/")
	ins.AddStringFlag("mongo_db", "adanos-alert", "mongodb database name")

	ins.AddStringFlag("preview_url", "http://localhost:19999/ui/groups/%s.html", "preview url")
	ins.AddStringFlag("report_url", "http://localhost:19999/ui/reports/%s.html", "report url")

	ins.AddBoolFlag("use_local_dashboard", "使用本地的 Dashboard，用于开发调试")
	ins.AddBoolFlag("no_job_mode", "启用该标识后，将会停止事件聚合和队列任务处理，用于开发调试")

	ins.AddBoolFlag("enable_migrate", "在应用启动时自动执行数据库迁移")
	ins.AddBoolFlag("re_migrate", "重新执行数据库迁移，重新迁移会移除已有的预定义模板")

	ins.AddStringFlag("aggregation_period", "5s", "事件聚合周期，每隔 aggregation_period 时间对新增事件执行聚合任务")
	ins.AddStringFlag("action_trigger_period", "5s", "动作触发周期，每隔 action_trigger_period 时间执行一次动作触发")
	ins.AddIntFlag("queue_job_max_retry_times", 3, "任务最大重试次数")
	ins.AddIntFlag("keep_period", 0, "保留多长时间的报警，如果全部保留，设置为0，单位为天，Adanos-Alert 会自动清理超过 keep_period 天的报警")
	ins.AddIntFlag("syslog_keep_period", 0, "保留多长时间的系统日志，如果全部保留，设置为0，单位为天，Adanos-Alert 会自动清理超过 syslog_keep_period 天的系统日志")
	ins.AddIntFlag("queue_worker_num", 3, "队列工作线程数")
	ins.AddStringFlag("query_timeout", "30s", "后端服务查询超时时间")
	ins.AddStringFlag("log_path", "", "日志文件输出目录（非文件名），默认为空，输出到标准输出")

	ins.AddStringFlag("aliyun_access_key", "", "阿里云语音通知接口 Access Key ID")
	ins.AddStringFlag("aliyun_access_secret", "", "阿里云语音通知接口 Access Secret")
	ins.AddStringFlag("aliyun_voice_called_show_number", "", "阿里云语音通知被叫显号")
	ins.AddStringFlag("aliyun_voice_tts_code", "", "阿里云语音通知模板，这里是模板ID，模板内容在阿里云申请，建议内容：\"您有一条名为 ${title} 的报警通知，请及时处理！\"")
	ins.AddStringFlag("aliyun_voice_tts_param", "title", "阿里云语音通知模板变量名")

	ins.AddStringFlag("jira_url", "", "Jira 服务器地址，如 http://127.0.0.1:8080")
	ins.AddStringFlag("jira_username", "", "Jira 用户名")
	ins.AddStringFlag("jira_password", "", "Jira 密码")

	ins.AddStringFlag("email_smtp_host", "", "邮件服务器地址")
	ins.AddIntFlag("email_smtp_port", 25, "邮件服务器端口")
	ins.AddStringFlag("email_smtp_username", "", "邮件服务器用户名")
	ins.AddStringFlag("email_smtp_password", "", "邮件服务器密码")

	ins.AddStringFlag("openai_endpoint", "https://api.openai.com/v1", "OpenAI API Endpoint")
	ins.AddStringFlag("openai_api_key", "", "OpenAI API Key")
	ins.AddStringFlag("openai_organization", "", "OpenAI Organization")

	ins.Init(func(f infra.FlagContext) error {
		stackWriter := writer.NewStackWriter()
		logPath := f.String("log_path")
		if logPath == "" {
			stackWriter.PushWithLevels(writer.NewStdoutWriter())
			return nil
		}

		log.All().LogFormatter(formatter.NewJSONWithTimeFormatter())
		stackWriter.PushWithLevels(writer.NewDefaultRotatingFileWriter(context.TODO(), func(le level.Level, module string) string {
			return filepath.Join(logPath, fmt.Sprintf("server-%s.%s.log", le.GetLevelName(), time.Now().Format("20060102")))
		}))

		stackWriter.PushWithLevels(
			NewErrorCollectorWriter(ins.Container()),
			level.Error,
			level.Emergency,
			level.Critical,
		)
		log.All().LogWriter(stackWriter)

		return nil
	})

	ins.Singleton(func(c infra.FlagContext) *configs.Config {
		aggregationPeriod, err := time.ParseDuration(c.String("aggregation_period"))
		if err != nil {
			log.Warningf("invalid argument [aggregation_period: %s], using default value", c.String("aggregation_period"))
			aggregationPeriod = 30 * time.Second
		}

		actionTriggerPeriod, err := time.ParseDuration(c.String("action_trigger_period"))
		if err != nil {
			log.Warningf("invalid argument [action_trigger_period: %s], using default value", c.String("action_trigger_period"))
			actionTriggerPeriod = 5 * time.Second
		}

		queryTimeout, err := time.ParseDuration(c.String("query_timeout"))
		if err != nil {
			log.Warningf("invalid argument [query_timeout: %s], using default value", c.String("query_timeout"))
			queryTimeout = 5 * time.Second
		}

		return &configs.Config{
			Listen:                c.String("listen"),
			GRPCListen:            c.String("grpc_listen"),
			GRPCToken:             c.String("grpc_token"),
			MongoURI:              c.String("mongo_uri"),
			MongoDB:               c.String("mongo_db"),
			UseLocalDashboard:     c.Bool("use_local_dashboard"),
			APIToken:              c.String("api_token"),
			AggregationPeriod:     aggregationPeriod,
			ActionTriggerPeriod:   actionTriggerPeriod,
			QueueJobMaxRetryTimes: c.Int("queue_job_max_retry_times"),
			QueueWorkerNum:        c.Int("queue_worker_num"),
			QueryTimeout:          queryTimeout,
			Migrate:               c.Bool("enable_migrate"),
			ReMigrate:             c.Bool("re_migrate"),
			PreviewURL:            c.String("preview_url"),
			ReportURL:             c.String("report_url"),
			KeepPeriod:            c.Int("keep_period"),
			SyslogKeepPeriod:      c.Int("syslog_keep_period"),
			NoJobMode:             c.Bool("no_job_mode"),
			AliyunVoiceCall: configs.AliyunVoiceCall{
				BaseURI:            "http://dyvmsapi.aliyuncs.com/",
				AccessKey:          c.String("aliyun_access_key"),
				AccessSecret:       c.String("aliyun_access_secret"),
				TTSCode:            c.String("aliyun_voice_tts_code"),
				TTSTemplateVarName: c.String("aliyun_voice_tts_param"),
				CalledShowNumber:   c.String("aliyun_voice_called_show_number"),
			},
			Jira: configs.Jira{
				BaseURL:  c.String("jira_url"),
				Username: c.String("jira_username"),
				Password: c.String("jira_password"),
			},
			EmailSMTP: configs.EmailSMTP{
				Host:     c.String("email_smtp_host"),
				Port:     c.Int("email_smtp_port"),
				Username: c.String("email_smtp_username"),
				Password: c.String("email_smtp_password"),
			},
			OpenAI: configs.OpenAI{
				Endpoint:     c.String("openai_endpoint"),
				APIKey:       c.String("openai_api_key"),
				Organization: c.String("openai_organization"),
			},
		}
	})

	ins.Singleton(func(ctx context.Context, conf *configs.Config) *mongo.Database {
		ctx, _ = context.WithTimeout(ctx, conf.QueryTimeout)
		conn, err := mongo.Connect(ctx, options.Client().
			ApplyURI(conf.MongoURI).
			SetServerSelectionTimeout(conf.QueryTimeout).
			SetConnectTimeout(conf.QueryTimeout).
			SetSocketTimeout(conf.QueryTimeout))
		if err != nil {
			log.Errorf("connect to mongodb failed: %s", err)
			return nil
		}

		return conn.Database(conf.MongoDB)
	})

	ins.Async(func(conf *configs.Config, em event.Manager) {
		if log.DebugEnabled() {
			log.WithFields(log.Fields{
				"config": conf,
			}).Debug("configuration")
		}

		_ = em.Publish(pubsub.SystemUpDownEvent{
			Up:        true,
			CreatedAt: time.Now(),
		})
	})

	ins.BeforeServerStop(func(cc infra.Resolver) error {
		return cc.Resolve(func(em event.Manager) {
			_ = em.Publish(pubsub.SystemUpDownEvent{
				Up:        false,
				CreatedAt: time.Now(),
			})
		})
	})

	ins.Provider(action.Provider{})
	ins.Provider(impl.Provider{})
	ins.Provider(api.Provider{})
	ins.Provider(job.Provider{})
	ins.Provider(queue.Provider{})
	ins.Provider(migrate.Provider{})
	ins.Provider(rpc.Provider{})
	ins.Provider(service.Provider{})
	ins.Provider(pubsub.Provider{})
	ins.Provider(llm.Provider{})

	if err := ins.Run(os.Args); err != nil {
		log.Errorf("exit with error: %s", err)
	}
}

type ErrorCollectorWriter struct {
	resolver infra.Resolver
}

func NewErrorCollectorWriter(resolver infra.Resolver) *ErrorCollectorWriter {
	return &ErrorCollectorWriter{resolver: resolver}
}

func (e *ErrorCollectorWriter) Write(le level.Level, module string, message string) error {
	return e.resolver.Resolve(func(evtRepo repository.EventRepo, syslogRepo repository.SyslogRepo) error {

		syslogID, err := syslogRepo.Add(repository.Syslog{
			Type: repository.SyslogTypeError,
			Context: map[string]interface{}{
				"level":  le.GetLevelName(),
				"module": module,
			},
			Body: str.Cutoff(500, message),
		})
		if err != nil {
			return err
		}

		_, err = evtRepo.Add(repository.Event{
			Content: message,
			Meta:    repository.EventMeta{"level": le.GetLevelName(), "module": module, "syslog_id": syslogID.Hex()},
			Tags:    []string{"internal-error"},
			Origin:  "internal",
		})

		return err
	})
	// return e.errorStore.Record(strings.ReplaceAll(message, "\n", color.TextWrap(color.Green, "↙")))
}

func (e *ErrorCollectorWriter) ReOpen() error {
	return nil
}

func (e *ErrorCollectorWriter) Close() error {
	return nil
}
