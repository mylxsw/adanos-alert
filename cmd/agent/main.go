package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/gorilla/mux"
	"github.com/ledisdb/ledisdb/ledis"
	"github.com/mylxsw/adanos-alert/agent/api"
	"github.com/mylxsw/adanos-alert/agent/config"
	"github.com/mylxsw/adanos-alert/agent/job"
	"github.com/mylxsw/adanos-alert/agent/store"
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/misc"
	"github.com/mylxsw/adanos-alert/rpc/protocol"
	"github.com/mylxsw/asteria/level"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/asteria/writer"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier"
	"github.com/mylxsw/glacier/starter/application"
	"github.com/mylxsw/glacier/web"
	"github.com/mylxsw/go-toolkit/network"
	"github.com/urfave/cli"
	"github.com/urfave/cli/altsrc"
	"google.golang.org/grpc"

	lediscfg "github.com/ledisdb/ledisdb/config"
)

var Version = "1.0"
var GitCommit = "5dbef13fb456f51a5d29464d"

func main() {
	app := application.Create(fmt.Sprintf("%s (%s)", Version, GitCommit[:8]))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "server_addr",
		Usage:  "server listen address",
		EnvVar: "ADANOS_SERVER_ADDR",
		Value:  "127.0.0.1:19998",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:   "server_token",
		Usage:  "API Token for api access control",
		EnvVar: "ADANOS_SERVER_TOKEN",
		Value:  "000000",
	}))
	app.AddFlags(altsrc.NewStringFlag(cli.StringFlag{
		Name:  "data_dir",
		Usage: "本地数据库存储目录",
		Value: "/tmp/adanos-agent",
	}))

	app.WithHttpServer().TCPListenerAddr("127.0.0.1:29999")

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

	// Config
	app.Singleton(func(c glacier.FlagContext) *config.Config {
		return &config.Config{
			DataDir:     c.String("data_dir"),
			ServerAddr:  c.String("server_addr"),
			ServerToken: c.String("server_token"),
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

	app.WebAppExceptionHandler(func(ctx web.Context, err interface{}) web.Response {
		log.Errorf("Stack: %s", debug.Stack())
		return nil
	})

	app.Main(func(conf *config.Config, router *mux.Router) {
		log.WithFields(log.Fields{
			"config": conf,
		}).Debug("configuration")

		for _, r := range web.GetAllRoutes(router) {
			log.Debugf("route: %s -> %s | %s | %s", r.Name, r.Methods, r.PathTemplate, r.PathRegexp)
		}
	})

	app.Provider(api.ServiceProvider{})
	app.Provider(store.ServiceProvider{})
	app.Provider(job.ServiceProvider{})

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
	return e.cc.ResolveWithError(func(msgStore store.MessageStore) error {
		ips, _ := network.GetLanIPs()
		data, _ := json.Marshal(misc.CommonMessage{
			Content: message,
			Meta: repository.MessageMeta{
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
