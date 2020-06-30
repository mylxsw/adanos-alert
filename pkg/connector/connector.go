package connector

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/mylxsw/adanos-alert/misc"
	"github.com/mylxsw/asteria/log"
	"github.com/pkg/errors"
)

// Send send a message to adanos servers
func Send(servers []string, token string, meta map[string]interface{}, tags []string, origin string, message string) error {
	commonMessage := misc.CommonMessage{
		Content: message,
		Meta:    meta,
		Tags:    tags,
		Origin:  origin,
	}
	data, _ := json.Marshal(commonMessage)

	var err error
	for _, s := range servers {
		if err = sendMessageToServer(commonMessage, data, s, token); err == nil {
			break
		}

		log.Warningf("send to server %s failed: %v", s, err)
	}

	return err
}

func sendMessageToServer(commonMessage misc.CommonMessage, data []byte, adanosServer, adanosToken string) error {
	reqURL := fmt.Sprintf("%s/api/messages/", strings.TrimRight(adanosServer, "/"))

	log.WithFields(log.Fields{
		"message": commonMessage,
	}).Debugf("request: %v", reqURL)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewReader(data))
	if err != nil {
		return errors.Wrap(err, "create request failed")
	}

	if adanosToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adanosToken))
	}

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "read response body failed")
	}

	log.Debugf("response: %v", string(respBody))
	return nil
}
