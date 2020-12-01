package configs

import (
	"encoding/json"
	"time"

	"github.com/mylxsw/container"
)

type Config struct {
	PreviewURL string `json:"preview_url"`
	ReportURL  string `json:"report_url"`
	Listen     string `json:"listen"`
	GRPCListen string `json:"grpc_listen"`
	GRPCToken  string `json:"-"`

	MongoURI          string `json:"mongo_uri"`
	MongoDB           string `json:"mongo_db"`
	APIToken          string `json:"-"`
	UseLocalDashboard bool   `json:"use_local_dashboard"`

	AggregationPeriod     time.Duration `json:"aggregation_period"`
	ActionTriggerPeriod   time.Duration `json:"action_trigger_period"`
	QueueJobMaxRetryTimes int           `json:"queue_job_max_retry_times"`
	QueueWorkerNum        int           `json:"queue_worker_num"`
	QueryTimeout          time.Duration `json:"query_timeout"`

	KeepPeriod      int `json:"keep_period"`
	AuditKeepPeriod int `json:"audit_keep_period"`

	Migrate         bool            `json:"migrate"`
	ReMigrate       bool            `json:"re_migrate"`

	AliyunVoiceCall AliyunVoiceCall `json:"aliyun_voice_call"`
	EmailSMTP       EmailSMTP       `json:"email_smtp"`
	Jira            Jira            `json:"jira"`
}

type EmailSMTP struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"-"`
}

type AliyunVoiceCall struct {
	AccessKey          string `json:"-"`
	AccessSecret       string `json:"-"`
	CalledShowNumber   string `json:"called_show_number"`
	TTSCode            string `json:"tts_code"`
	TTSTemplateVarName string `json:"tts_template_var_name"`
	BaseURI            string `json:"base_uri"`
}

type Jira struct {
	BaseURL  string `json:"base_url"`
	Username string `json:"username"`
	Password string `json:"-"`
}

func (conf *Config) Serialize() string {
	rs, _ := json.Marshal(conf)
	return string(rs)
}

// Get return config object from container
func Get(cc container.Container) *Config {
	return cc.MustGet(&Config{}).(*Config)
}
