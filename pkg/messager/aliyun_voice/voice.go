package aliyun_voice

import (
	"errors"
	"fmt"
	"time"

	"github.com/mylxsw/adanos-alert/configs"
	"github.com/mylxsw/asteria/log"
)

type VoiceCall struct {
	conf *configs.Config
}

func NewVoiceCall(conf *configs.Config) *VoiceCall {
	return &VoiceCall{conf: conf}
}

func (vc *VoiceCall) Send(title string, receivers []string) error {
	if len(receivers) == 0 {
		return nil
	}

	for _, receiver := range receivers {
		go vc.call(receiver, title)
	}

	return nil
}

func (vc *VoiceCall) call(receiver string, title string) {
	ap := CreateAliyunPOP(vc.conf.AliyunVoiceCall.AccessKey, vc.conf.AliyunVoiceCall.AccessSecret)
	ap.SetParam("Action", "SingleCallByTts")
	ap.SetParam("Version", "2017-05-25")
	ap.SetParam("RegionId", "cn-hangzhou")
	// 被叫显号
	ap.SetParam("CalledShowNumber", vc.conf.AliyunVoiceCall.CalledShowNumber)
	// 被叫号码
	ap.SetParam("CalledNumber", receiver)
	// TTS文本模板Code
	ap.SetParam("TtsCode", vc.conf.AliyunVoiceCall.TTSCode)
	// 替换TTS模板中变量的JSON串
	ap.SetParam("TtsParam", fmt.Sprintf(`{"%s":"%s"}`, vc.conf.AliyunVoiceCall.TTSTemplateVarName, title))
	// 音量
	// ap.SetParam("Volume", "100")
	// 播放次数（最多3次）
	// ap.SetParam("PlayTimes", "3")
	// 预留给调用方使用的ID, 最终会通过在回执消息中将此ID带回给调用方
	ap.SetParam("OutId", "1")

	// 发送通知
	var success = false
	var lastError error
	for i := 0; i < 3; i++ {
		resp, err := ap.Request(vc.conf.AliyunVoiceCall.BaseURI)
		if err != nil {
			log.Warningf("阿里云语音接口调用失败，5s 后自动重试：%s", err)
			lastError = err
			time.Sleep(5 * time.Second)

			continue
		}

		if resp.Code != "OK" {
			msg := "阿里云语音接口调用失败"
			log.WithFields(log.Fields{
				"receiver": receiver,
				"title":    title,
				"resp":     resp,
			}).Warningf(msg)

			lastError = errors.New("阿里云语音接口调用失败")

			break
		}

		success = true
		break
	}

	if !success {
		msg := fmt.Sprintf("阿里云语音通知失败，最后错误: %s", lastError)
		log.WithFields(log.Fields{
			"receiver": receiver,
			"title":    title,
		}).Error(msg)
	}
}
