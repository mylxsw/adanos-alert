package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ledisdb/ledisdb/ledis"
	"github.com/mylxsw/adanos-alert/agent/api"
	"github.com/mylxsw/adanos-alert/agent/config"
	"github.com/mylxsw/adanos-alert/agent/job"
	"github.com/mylxsw/adanos-alert/agent/store"
	"github.com/mylxsw/adanos-alert/internal/extension"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/pkg/misc"
	"github.com/mylxsw/adanos-alert/rpc/protocol"
	"github.com/mylxsw/asteria/formatter"
	"github.com/mylxsw/asteria/level"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/asteria/writer"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/starter/application"
	"github.com/urfave/cli"
	"github.com/urfave/cli/altsrc"
	"google.golang.org/grpc"

	lediscfg "github.com/ledisdb/ledisdb/config"
)

// Version GitCommit 编译参数，编译时通过编译选项 ldflags 指定
var Version = "1.0"
var GitCommit = "5dbef13fb456f51a5d29464d"
var AsyncRunner = 3
var DEBUG = ""

func main() {
	if DEBUG == "true" {
		infra.DEBUG = true
	}

	app := application.Create(fmt.Sprintf("%s (%s)", Version, GitCommit), AsyncRunner).WithLogger(log.Module("glacier"))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "server_addr",
		Usage:  "server grpc listen address",
		EnvVar: "ADANOS_SERVER_ADDR",
		Value:  "127.0.0.1:19998",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "server_token",
		Usage:  "API Token for grpc api access control",
		EnvVar: "ADANOS_SERVER_TOKEN",
		Value:  "000000",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:  "data_dir",
		Usage: "本地数据库存储目录",
		Value: "/tmp/adanos-agent",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "listen",
		Usage:  "listen address",
		EnvVar: "ADANOS_AGENT_LISTEN_ADDR",
		Value:  "127.0.0.1:29999",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:  "log_path",
		Usage: "日志文件输出目录（非文件名），默认为空，输出到标准输出",
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
				return filepath.Join(logPath, fmt.Sprintf("agent-%s.%s.log", le.GetLevelName(), time.Now().Format("20060102")))
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

	// Config
	app.Singleton(func(c infra.FlagContext) *config.Config {
		return &config.Config{
			DataDir:     c.String("data_dir"),
			ServerAddr:  c.String("server_addr"),
			ServerToken: c.String("server_token"),
			Listen:      c.String("listen"),
			LogPath:     c.String("log_path"),
		}
	})

	// Ledis DB
	app.Singleton(func(conf *config.Config) (*ledis.Ledis, error) {
		cfg := lediscfg.NewConfigDefault()
		cfg.DataDir = conf.DataDir
		cfg.Databases = 1

		return ledis.Open(cfg)
	})
	app.Singleton(func(ld *ledis.Ledis) (*ledis.DB, error) {
		return ld.Select(0)
	})

	// GRPC
	app.Singleton(func(conf *config.Config) (grpc.ClientConnInterface, error) {
		return grpc.Dial(conf.ServerAddr, grpc.WithInsecure(), grpc.WithPerRPCCredentials(NewAuthAPI(conf.ServerToken)))
	})

	app.Async(func(conf *config.Config, db *ledis.DB) {
		agentID, err := db.Get([]byte("agent-id"))
		if err != nil || agentID == nil {
			_ = db.Set([]byte("agent-id"), []byte(misc.UUID()))

			agentID, _ = db.Get([]byte("agent-id"))
		}

		if log.DebugEnabled() {
			log.WithFields(log.Fields{
				"config":   conf,
				"agent_id": string(agentID),
			}).Debug("configuration")
		}
	})

	app.Provider(api.Provider{})
	app.Provider(store.Provider{})
	app.Provider(job.Provider{})

	if err := app.Run(os.Args); err != nil {
		log.Errorf("exit with error: %s", err)
	}
}

// ErrorCollectorWriter Agent 错误日志采集器
type ErrorCollectorWriter struct {
	cc infra.Resolver
}

// NewErrorCollectorWriter 创建一个错误日志采集器
func NewErrorCollectorWriter(cc infra.Resolver) *ErrorCollectorWriter {
	return &ErrorCollectorWriter{cc: cc}
}

func (e *ErrorCollectorWriter) Write(le level.Level, module string, message string) error {
	return e.cc.ResolveWithError(func(msgStore store.EventStore) error {
		ips, _ := misc.GetLanIPs()
		data, _ := json.Marshal(extension.CommonEvent{
			Content: message,
			Meta: repository.EventMeta{
				"module": module,
				"level":  le.GetLevelName(),
				"ip":     ips,
			},
			Tags:   []string{"agent"},
			Origin: "agent",
		})

		req := protocol.MessageRequest{Data: string(data)}
		return msgStore.Enqueue(&req)
	})
}

func (e *ErrorCollectorWriter) ReOpen() error {
	return nil
}

func (e *ErrorCollectorWriter) Close() error {
	return nil
}

// AuthAPI 权限插件
type AuthAPI struct {
	token string
}

// NewAuthAPI 创建一个Auth API
func NewAuthAPI(token string) *AuthAPI {
	return &AuthAPI{token: token}
}

func (a *AuthAPI) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"token": a.token,
	}, nil
}

func (a *AuthAPI) RequireTransportSecurity() bool {
	return false
}
