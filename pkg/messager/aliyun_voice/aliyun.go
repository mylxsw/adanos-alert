package aliyun_voice

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

// AliyunResponse is a response object for aliyun request
type AliyunResponse struct {
	Message   string `json:"Message"`
	RequestID string `json:"RequestId"`
	Code      string `json:"Code"`
}

// AliyunPOP is an implemention for Aliyun POP protocol
type AliyunPOP struct {
	keyID   string
	secret  string
	request map[string]string
}

// CreateAliyunPOP create AliyunPOP protocol request
func CreateAliyunPOP(keyID, secret string) *AliyunPOP {
	request := make(map[string]string)

	request["SignatureMethod"] = "HMAC-SHA1"
	request["SignatureNonce"] = uuid.NewV4().String()
	request["AccessKeyId"] = keyID
	request["SignatureVersion"] = "1.0"
	request["Format"] = "JSON"

	timezone, _ := time.LoadLocation("GMT0")
	request["Timestamp"] = time.Now().In(timezone).Format("2006-01-02T15:04:05Z")

	return &AliyunPOP{
		keyID:   keyID,
		secret:  secret,
		request: request,
	}
}

func (ap *AliyunPOP) urlEncode(param string) string {
	param = url.QueryEscape(param)
	param = strings.Replace(param, "+", "%20", -1)
	param = strings.Replace(param, "*", "%2A", -1)
	param = strings.Replace(param, "%7E", "~", -1)
	return param
}

func (ap *AliyunPOP) sign(strToSign string) string {
	mac := hmac.New(sha1.New, []byte(ap.secret+"&"))
	mac.Write([]byte(strToSign))
	signData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(signData)
}

func (ap *AliyunPOP) createSignature(queryString string) string {
	signStr := ap.sign("GET" + "&" + ap.urlEncode("/") + "&" + ap.urlEncode(queryString))
	return ap.urlEncode(signStr)
}

func (ap *AliyunPOP) createQueryString() string {
	delete(ap.request, "Signature")

	indexes := make([]string, 0)
	for k := range ap.request {
		indexes = append(indexes, k)
	}
	sort.Strings(indexes)

	sortedQueryString := ""
	for _, v := range indexes {
		sortedQueryString = sortedQueryString + "&" + ap.urlEncode(v) + "=" + ap.urlEncode(ap.request[v])
	}

	return sortedQueryString[1:]
}

// SetParam set request params
func (ap *AliyunPOP) SetParam(key, value string) *AliyunPOP {
	ap.request[key] = value
	return ap
}

// GenerateURL generate request url
func (ap *AliyunPOP) GenerateURL(baseURL string) string {
	queryString := ap.createQueryString()
	return fmt.Sprintf("%s?Signature=%s&%s", baseURL, ap.createSignature(queryString), queryString)
}

// Request send a request to aliyun server
func (ap *AliyunPOP) Request(baseUri string) (*AliyunResponse, error) {
	urlStr := ap.GenerateURL(baseUri)

	client := http.Client{
		Timeout: 3 * time.Second,
	}
	response, err := client.Get(urlStr)
	if err != nil {
		return nil, fmt.Errorf("request failed: %s", err)
	}

	body := response.Body
	defer func() { _ = body.Close() }()

	resp, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %s", err)
	}

	var aliyunResp AliyunResponse
	if err := json.Unmarshal(resp, &aliyunResp); err != nil {
		return nil, fmt.Errorf("can not parse aliyun response to json")
	}

	return &aliyunResp, nil
}

