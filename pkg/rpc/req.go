/**
 * @Author: zjj
 * @Date: 2024/12/12
 * @Desc:
**/

package rpc

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"gmicro/pkg/gerr"
	"gmicro/pkg/json"
	"gmicro/pkg/jsonpb"
	"gmicro/pkg/log"
	"gmicro/pkg/uctx"
	"google.golang.org/protobuf/proto"
	"net/url"
)

var defaultRestyClient *RestyClient

type RestyClient struct {
	*resty.Client
}

func init() {
	defaultRestyClient = NewRestyClient()
}

type Resp struct {
	Data    string `json:"data"`
	ErrCode int32  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Hint    string `json:"hint"`
}

func NewRestyClient() *RestyClient {
	return &RestyClient{
		resty.New(),
	}
}

func NewRequest() *resty.Request {
	return defaultRestyClient.R()
}

func DoRequest(_ uctx.IUCtx, _, path, method string, req, out proto.Message) error {
	var body []byte
	val, err := jsonpb.MarshalToString(req)
	body = []byte(val)
	if err != nil {
		return gerr.Wrap(err)
	}

	var headers = make(map[string]string)

	// todo 服务发现
	var target = fmt.Sprintf("%s://%s:%s/gateway", "http", "127.0.0.1", "20000")
	result, err := url.JoinPath(target, path)
	if err != nil {
		return gerr.Wrap(err)
	}

	resp, err := NewRequest().SetHeaders(headers).SetBody(body).Execute(method, result)
	if err != nil {
		return gerr.Wrap(err)
	}

	log.Infof("rpc resp: %s", string(resp.Body()))
	var respBody Resp
	err = json.Unmarshal(resp.Body(), &respBody)
	if err != nil {
		return gerr.Wrap(err)
	}

	if respBody.ErrCode > 0 {
		return gerr.NewErr(respBody.ErrCode, respBody.ErrMsg)
	}

	err = jsonpb.Unmarshal([]byte(respBody.Data), out)
	if err != nil {
		return gerr.Wrap(err)
	}

	return nil
}
