package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

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
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/starter/application"
	"github.com/urfave/cli"
	"github.com/urfave/cli/altsrc"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Version = "1.0"
var GitCommit = "5dbef13fb456f51a5d29464d"

func main() {
	app := application.Create(fmt.Sprintf("%s (%s)", Version, GitCommit[:8]))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:  "listen",
		Usage: "http listen addr",
		Value: ":19999",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:  "grpc_listen",
		Usage: "GRPC Server listen address",
		Value: ":19998",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:  "grpc_token",
		Usage: "GRPC Server token",
		Value: "000000",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "preview_url",
		Usage:  "Alert preview page url",
		EnvVar: "ADANOS_PREVIEW_URL",
		Value:  "http://localhost:19999/ui/groups/%s.html",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "report_url",
		Usage:  "Alert report page url",
		EnvVar: "ADANOS_REPORT_URL",
		Value:  "http://localhost:19999/ui/reports/%s.html",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "mongo_uri",
		Usage:  "Mongodb connection uri，参考 https://docs.mongodb.com/manual/reference/connection-string/",
		EnvVar: "MONGODB_HOST",
		Value:  "mongodb://localhost:27017",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "mongo_db",
		Usage:  "Mongodb database name",
		EnvVar: "MONGODB_DB",
		Value:  "adanos-alert",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "api_token",
		Usage:  "API Token for api access control",
		EnvVar: "ADANOS_API_TOKEN",
		Value:  "",
	}))
	app.AddFlags(altsrc.NewBoolFlag(cli.BoolFlag{
		Name:  "use_local_dashboard",
		Usage: "whether using local dashboard, this is used when development",
	}))
	app.AddFlags(altsrc.NewBoolFlag(cli.BoolFlag{
		Name:  "enable_migrate",
		Usage: "whether enable database migrate when app run",
	}))
	app.AddFlags(altsrc.NewBoolFlag(cli.BoolFlag{
		Name:  "re_migrate",
		Usage: "是否重新执行迁移，重新迁移会移除已有的预定义模板",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "aggregation_period",
		Usage:  "aggregation job execute period",
		EnvVar: "ADANOS_AGGREGATION_PERIOD",
		Value:  "5s",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "action_trigger_period",
		Usage:  "action trigger job execute period",
		EnvVar: "ADANOS_ACTION_TRIGGER_PERIOD",
		Value:  "5s",
	}))
	app.AddFlags(altsrc.NewIntFlag(cli.IntFlag{
		Name:   "queue_job_max_retry_times",
		Usage:  "set queue job max retry times",
		EnvVar: "ADANOS_QUEUE_JOB_MAX_RETRY_TIMES",
		Value:  3,
	}))

	app.AddFlags(altsrc.NewIntFlag(cli.IntFlag{
		Name:   "keep_period",
		Usage:  "保留多长时间的报警，如果全部保留，设置为0，单位为天，Adanos-Alert 会自动清理超过 keep_period 天的报警",
		EnvVar: "ADANOS_KEEP_PERIOD",
		Value:  0,
	}))
	app.AddFlags(altsrc.NewIntFlag(cli.IntFlag{
		Name:   "audit_keep_period",
		Usage:  "保留多长时间的审计日志，如果全部保留，设置为0，单位为天，Adanos-Alert 会自动清理超过 audit_keep_period 天的审计日志",
		EnvVar: "ADANOS_AUDIT_KEEP_PERIOD",
		Value:  0,
	}))

	app.AddFlags(altsrc.NewIntFlag(cli.IntFlag{
		Name:   "queue_worker_num",
		Usage:  "set queue worker numbers",
		EnvVar: "ADANOS_QUEUE_WORKER_NUM",
		Value:  3,
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "query_timeout",
		Usage:  "query timeout for backend service",
		EnvVar: "ADANOS_QUERY_TIMEOUT",
		Value:  "30s",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "aliyun_access_key",
		EnvVar: "ADANOS_ALIYUN_ACCESS_KEY",
		Value:  "",
		Usage:  "阿里云语音通知接口 Access Key ID",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "aliyun_access_secret",
		EnvVar: "ADANOS_ALIYUN_ACCESS_SECRET",
		Value:  "",
		Usage:  "阿里云语音通知接口 Access Secret",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:  "aliyun_voice_called_show_number",
		Value: "",
		Usage: "阿里云语音通知被叫显号",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:  "aliyun_voice_tts_code",
		Value: "",
		Usage: "阿里云语音通知模板，这里是模板ID，模板内容在阿里云申请，建议内容：\"您有一条名为 ${title} 的报警通知，请及时处理！\"",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:  "aliyun_voice_tts_param",
		Value: "title",
		Usage: "阿里云语音通知模板变量名",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:  "log_path",
		Usage: "日志文件输出目录（非文件名），默认为空，输出到标准输出",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "jira_url",
		EnvVar: "ADANOS_JIRA_URL",
		Usage:  "Jira 服务器地址，如 http://127.0.0.1:8080",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "jira_username",
		EnvVar: "ADANOS_JIRA_USERNAME",
		Usage:  "Jira 连接账号",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "jira_password",
		EnvVar: "ADANOS_JIRA_PASSWORD",
		Usage:  "Jira 连接密码",
	}))
	app.AddFlags(altsrc.NewBoolFlag(cli.BoolFlag{
		Name:  "no_job_mode",
		Usage: "启用该标识后，将会停止事件聚合和队列任务处理，用于开发调试",
	}))

	app.BeforeServerStart(func(cc container.Container) error {
		stackWriter := writer.NewStackWriter()
		cc.MustResolve(func(ctx context.Context, c infra.FlagContext) {
			logPath := c.String("log_path")
			if logPath == "" {
				log.All().LogFormatter(formatter.NewJSONFormatter())
				stackWriter.PushWithLevels(writer.NewStdoutWriter())
				return
			}

			log.All().LogFormatter(formatter.NewJSONWithTimeFormatter())
			stackWriter.PushWithLevels(writer.NewDefaultRotatingFileWriter(ctx, func(le level.Level, module string) string {
				return filepath.Join(logPath, fmt.Sprintf("server-%s.%s.log", le.GetLevelName(), time.Now().Format("20060102")))
			}))
		})

		stackWriter.PushWithLevels(
			NewErrorCollectorWriter(cc),
			level.Error,
			level.Emergency,
			level.Critical,
		)
		log.All().LogWriter(stackWriter)

		return nil
	})

	app.Singleton(func(c infra.FlagContext) *configs.Config {
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
			AuditKeepPeriod:       c.Int("audit_keep_period"),
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
		}
	})

	app.Singleton(func(ctx context.Context, conf *configs.Config) *mongo.Database {
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

	app.Main(func(conf *configs.Config, em event.Manager) {
		rand.Seed(time.Now().Unix())

		if log.DebugEnabled() {
			log.WithFields(log.Fields{
				"config": conf,
			}).Debug("configuration")
		}

		em.Publish(pubsub.SystemUpDownEvent{
			Up:        true,
			CreatedAt: time.Now(),
		})
	})

	app.BeforeServerStop(func(cc infra.Resolver) error {
		return cc.Resolve(func(em event.Manager) {
			em.Publish(pubsub.SystemUpDownEvent{
				Up:        false,
				CreatedAt: time.Now(),
			})
		})
	})

	app.Provider(action.Provider{})
	app.Provider(impl.Provider{})
	app.Provider(api.Provider{})
	app.Provider(job.Provider{})
	app.Provider(queue.Provider{})
	app.Provider(migrate.Provider{})
	app.Provider(rpc.Provider{})
	app.Provider(service.Provider{})
	app.Provider(pubsub.Provider{})

	if err := app.Run(os.Args); err != nil {
		log.Errorf("exit with error: %s", err)
	}
}

type ErrorCollectorWriter struct {
	cc infra.Resolver
}

func NewErrorCollectorWriter(cc infra.Resolver) *ErrorCollectorWriter {
	return &ErrorCollectorWriter{cc: cc}
}

func (e *ErrorCollectorWriter) Write(le level.Level, module string, message string) error {
	return e.cc.ResolveWithError(func(evtRepo repository.EventRepo, auditRepo repository.AuditLogRepo) error {

		auditLogID, err := auditRepo.Add(repository.AuditLog{
			Type: repository.AuditLogTypeError,
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
			Meta:    repository.EventMeta{"level": le.GetLevelName(), "module": module, "audit_id": auditLogID.Hex()},
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
