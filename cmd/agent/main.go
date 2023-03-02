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
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/starter/app"
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

	ins := app.Create(fmt.Sprintf("%s (%s)", Version, GitCommit), AsyncRunner).WithLogger(log.Module("glacier"))

	ins.AddFlags(app.StringEnvFlag("server_addr", "127.0.0.1:19998", "server grpc listen address", "ADANOS_SERVER_ADDR"))
	ins.AddFlags(app.StringEnvFlag("server_token", "000000", "API Token for grpc api access control", "ADANOS_SERVER_TOKEN"))
	ins.AddStringFlag("data_dir", "/tmp/adanos-agent", "本地数据库存储目录")
	ins.AddStringFlag("log_path", "", "日志文件输出目录（非文件名），默认为空，输出到标准输出")
	ins.AddFlags(app.StringEnvFlag("listen", "127.0.0.1:29999", "agent listen address", "ADANOS_AGENT_LISTEN_ADDR"))

	ins.Init(func(c infra.FlagContext) error {
		stackWriter := writer.NewStackWriter()
		logPath := c.String("log_path")
		if logPath == "" {
			log.All().LogFormatter(formatter.NewJSONFormatter())
			stackWriter.PushWithLevels(writer.NewStdoutWriter())

			return nil
		}

		log.All().LogFormatter(formatter.NewJSONWithTimeFormatter())
		stackWriter.PushWithLevels(writer.NewDefaultRotatingFileWriter(context.TODO(), func(le level.Level, module string) string {
			return filepath.Join(logPath, fmt.Sprintf("agent-%s.%s.log", le.GetLevelName(), time.Now().Format("20060102")))
		}))

		stackWriter.PushWithLevels(
			NewErrorCollectorWriter(ins.Resolver()),
			level.Error,
			level.Emergency,
			level.Critical,
		)
		log.All().LogWriter(stackWriter)

		return nil
	})

	// Config
	ins.Singleton(func(c infra.FlagContext) *config.Config {
		return &config.Config{
			DataDir:     c.String("data_dir"),
			ServerAddr:  c.String("server_addr"),
			ServerToken: c.String("server_token"),
			Listen:      c.String("listen"),
			LogPath:     c.String("log_path"),
		}
	})

	// Ledis DB
	ins.Singleton(func(conf *config.Config) (*ledis.Ledis, error) {
		cfg := lediscfg.NewConfigDefault()
		cfg.DataDir = conf.DataDir
		cfg.Databases = 1

		return ledis.Open(cfg)
	})
	ins.Singleton(func(ld *ledis.Ledis) (*ledis.DB, error) {
		return ld.Select(0)
	})

	// GRPC
	ins.Singleton(func(conf *config.Config) (grpc.ClientConnInterface, error) {
		return grpc.Dial(conf.ServerAddr, grpc.WithInsecure(), grpc.WithPerRPCCredentials(NewAuthAPI(conf.ServerToken)))
	})

	ins.Async(func(conf *config.Config, db *ledis.DB) {
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

	ins.Provider(api.Provider{})
	ins.Provider(store.Provider{})
	ins.Provider(job.Provider{})

	if err := ins.Run(os.Args); err != nil {
		log.Errorf("exit with error: %s", err)
	}
}

// ErrorCollectorWriter Agent 错误日志采集器
type ErrorCollectorWriter struct {
	resolver infra.Resolver
}

// NewErrorCollectorWriter 创建一个错误日志采集器
func NewErrorCollectorWriter(resolver infra.Resolver) *ErrorCollectorWriter {
	return &ErrorCollectorWriter{resolver: resolver}
}

func (e *ErrorCollectorWriter) Write(le level.Level, module string, message string) error {
	return e.resolver.Resolve(func(msgStore store.EventStore) error {
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
