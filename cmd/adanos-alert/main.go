package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mylxsw/adanos-alert/api"
	"github.com/mylxsw/adanos-alert/configs"
	"github.com/mylxsw/adanos-alert/internal/action"
	"github.com/mylxsw/adanos-alert/internal/job"
	"github.com/mylxsw/adanos-alert/internal/repository/impl"
	"github.com/mylxsw/asteria/formatter"
	"github.com/mylxsw/asteria/level"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/asteria/writer"
	"github.com/mylxsw/glacier"
	"github.com/urfave/cli"
	"github.com/urfave/cli/altsrc"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Version string
var GitCommit string

// ConnectionTimeout is a timeout setting for mongodb connection
const ConnectionTimeout = 5 * time.Second

func main() {
	app := glacier.Create(fmt.Sprintf("%s (%s)", Version, GitCommit))
	app.WithHttpServer(":19999")

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
		Name:  "api_token",
		Usage: "API Token for api access control",
		Value: "",
	}))
	app.AddFlags(altsrc.NewBoolTFlag(cli.BoolTFlag{
		Name:  "console_color",
		Usage: "log colorful for console",
	}))
	app.AddFlags(altsrc.NewBoolFlag(cli.BoolFlag{
		Name:  "use_local_dashboard",
		Usage: "whether using local dashboard, this is used when development",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:  "aggregation_period",
		Usage: "aggregation job execute period",
		Value: "30s",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:  "action_trigger_period",
		Usage: "action trigger job execute period",
		Value: "15s",
	}))

	app.BeforeInitialize(func(c *cli.Context) error {
		log.DefaultLogFormatter(formatter.NewDefaultFormatter(c.Bool("console_color")))
		return nil
	})

	app.Singleton(func(c *cli.Context) *configs.Config {
		aggregationPeriod, err := time.ParseDuration(c.String("aggregation_period"))
		if err != nil {
			log.Warningf("invalid argument [aggregation_period: %s], using default value", c.String("aggregation_period"))
			aggregationPeriod = 30 * time.Second
		}

		actionTriggerPeriod, err := time.ParseDuration(c.String("action_trigger_period"))
		if err != nil {
			log.Warningf("invalid argument [action_trigger_period: %s], using default value", c.String("action_trigger_period"))
			actionTriggerPeriod = 15 * time.Second
		}

		return &configs.Config{
			MongoURI:            c.String("mongo_uri"),
			MongoDB:             c.String("mongo_db"),
			UseLocalDashboard:   c.Bool("use_local_dashboard"),
			ConsoleColor:        c.Bool("console_color"),
			APIToken:            c.String("api_token"),
			AggregationPeriod:   aggregationPeriod,
			ActionTriggerPeriod: actionTriggerPeriod,
		}
	})

	app.Singleton(func(ctx context.Context, conf *configs.Config) *mongo.Database {
		ctx, _ = context.WithTimeout(ctx, ConnectionTimeout)
		conn, err := mongo.Connect(ctx, options.Client().
			ApplyURI(conf.MongoURI).
			SetServerSelectionTimeout(ConnectionTimeout).
			SetConnectTimeout(ConnectionTimeout).
			SetSocketTimeout(ConnectionTimeout))
		if err != nil {
			log.Errorf("connect to mongodb failed: %s", err)
			return nil
		}

		return conn.Database(conf.MongoDB)
	})

	app.Main(func(conf *configs.Config) {
		stackWriter := writer.NewStackWriter()
		stackWriter.PushWithLevels(writer.NewStdoutWriter())
		stackWriter.PushWithLevels(
			NewErrorCollectorWriter(),
			level.Error,
			level.Emergency,
			level.Critical,
		)

		log.All().LogWriter(stackWriter)

		log.WithFields(log.Fields{
			"config": conf,
		}).Debug("configuration")
	})

	app.Provider(action.ServiceProvider{})
	app.Provider(impl.ServiceProvider{})
	app.Provider(api.ServiceProvider{})
	app.Provider(job.ServiceProvider{})

	if err := app.Run(os.Args); err != nil {
		log.Errorf("exit with error: %s", err)
	}
}

type ErrorCollectorWriter struct{}

func NewErrorCollectorWriter() *ErrorCollectorWriter {
	return &ErrorCollectorWriter{}
}

func (e *ErrorCollectorWriter) Write(le level.Level, module string, message string) error {
	// return e.errorStore.Record(strings.ReplaceAll(message, "\n", color.TextWrap(color.Green, "↙")))
	return nil
}

func (e *ErrorCollectorWriter) ReOpen() error {
	return nil
}

func (e *ErrorCollectorWriter) Close() error {
	return nil
}
