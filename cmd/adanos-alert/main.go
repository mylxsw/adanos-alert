package main

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/mylxsw/adanos-alert/rpc"
	"github.com/mylxsw/asteria/writer"

	"github.com/gorilla/mux"
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
	"github.com/mylxsw/glacier"
	"github.com/mylxsw/glacier/starter/application"
	"github.com/mylxsw/glacier/web"
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
		Name:   "mongo_uri",
		Usage:  "Mongodb connection uri",
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
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "aggregation_period",
		Usage:  "aggregation job execute period",
		EnvVar: "ADANOS_AGGREGATION_PERIOD",
		Value:  "30s",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "action_trigger_period",
		Usage:  "action trigger job execute period",
		EnvVar: "ADANOS_ACTION_TRIGGER_PERIOD",
		Value:  "10s",
	}))
	app.AddFlags(altsrc.NewIntFlag(cli.IntFlag{
		Name:   "queue_job_max_retry_times",
		Usage:  "set queue job max retry times",
		EnvVar: "ADANOS_QUEUE_JOB_MAX_RETRY_TIMES",
		Value:  3,
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
		Value:  "5s",
	}))

	app.WithHttpServer().TCPListenerAddr(":19999")

	app.BeforeServerStart(func(cc container.Container) error {
		stackWriter := writer.NewStackWriter()
		stackWriter.PushWithLevels(writer.NewStdoutWriter())
		stackWriter.PushWithLevels(
			NewErrorCollectorWriter(app.Container()),
			level.Error,
			level.Emergency,
			level.Critical,
		)
		log.All().LogWriter(stackWriter)

		return nil
	})

	app.Singleton(func(c glacier.FlagContext) *configs.Config {
		aggregationPeriod, err := time.ParseDuration(c.String("aggregation_period"))
		if err != nil {
			log.Warningf("invalid argument [aggregation_period: %s], using default value", c.String("aggregation_period"))
			aggregationPeriod = 30 * time.Second
		}

		actionTriggerPeriod, err := time.ParseDuration(c.String("action_trigger_period"))
		if err != nil {
			log.Warningf("invalid argument [action_trigger_period: %s], using default value", c.String("action_trigger_period"))
			actionTriggerPeriod = 10 * time.Second
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
			PreviewURL:            c.String("preview_url"),
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

	app.WebAppExceptionHandler(func(ctx web.Context, err interface{}) web.Response {
		log.Errorf("Stack: %s", debug.Stack())
		return nil
	})

	app.Main(func(conf *configs.Config, router *mux.Router) {
		log.WithFields(log.Fields{
			"config": conf,
		}).Debug("configuration")

		for _, r := range web.GetAllRoutes(router) {
			log.Debugf("route: %s -> %s | %s | %s", r.Name, r.Methods, r.PathTemplate, r.PathRegexp)
		}
	})

	app.Provider(action.ServiceProvider{})
	app.Provider(impl.ServiceProvider{})
	app.Provider(api.ServiceProvider{})
	app.Provider(job.ServiceProvider{})
	app.Provider(queue.ServiceProvider{})
	app.Provider(migrate.ServiceProvider{})
	app.Provider(rpc.Provider{})

	if err := app.Run(os.Args); err != nil {
		log.Errorf("exit with error: %s", err)
	}
}

type ErrorCollectorWriter struct {
	cc container.Container
}

func NewErrorCollectorWriter(cc container.Container) *ErrorCollectorWriter {
	return &ErrorCollectorWriter{cc: cc}
}

func (e *ErrorCollectorWriter) Write(le level.Level, module string, message string) error {
	return e.cc.ResolveWithError(func(msgRepo repository.MessageRepo) error {
		_, err := msgRepo.Add(repository.Message{
			Content: message,
			Meta:    repository.MessageMeta{"level": le.GetLevelName(), "module": module},
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
