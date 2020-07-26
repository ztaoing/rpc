/**
* @Author:zhoutao
* @Date:2020/7/25 下午3:46
 */

package bapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	APP_KEY    = "eddycjy"
	APP_SECRET = "go-programming-tour-book"
)

type AccessToken struct {
	Token string `json:"token"`
}

type API struct {
	URL    string
	Client *http.Client
}

func NewAPI(url string) *API {
	return &API{
		URL:    url,
		Client: http.DefaultClient,
	}
}

func (a *API) httpGet(ctx context.Context, path string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", a.URL, path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	//http的链路追踪和grpc的链路追踪的设置逻辑是一样的
	//链路追踪
	span, _ := opentracing.StartSpanFromContext(ctx, "HTTP GET"+a.URL, opentracing.Tag{Key: string(ext.Component), Value: "HTTP"})
	//设置tag
	span.SetTag("url", url)

	_ = opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	req = req.WithContext(context.Background())

	client := http.Client{Timeout: time.Second * 60}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	defer span.Finish()

	//读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, err

}
func (a *API) getAccessToken(ctx context.Context) (string, error) {
	body, err := a.httpGet(ctx, fmt.Sprintf("%s?app_key=%s&app_secret=%s", "auth", APP_KEY, APP_SECRET))
	if err != nil {
		return "", err
	}

	var accessToken AccessToken
	err = json.Unmarshal(body, &accessToken)
	if err != nil {
		return "", err
	}
	return accessToken.Token, nil
}

//从博客系统中获取标签列表
func (a *API) GetTagList(ctx context.Context, name string) ([]byte, error) {
	token, err := a.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	body, err := a.httpGet(ctx, fmt.Sprintf("%s?token=%s&name=%s", "api/v1/tags", token, name))
	if err != nil {
		return nil, err
	}
	return body, err
}
