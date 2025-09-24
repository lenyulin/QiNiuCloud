package textshrink

import (
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type Shrink interface {
	Shrink(ctx context.Context, text string) (string, error)
}
type Shrinker struct {
	clientPool *http.Client
	l          logger.ZapLogger
}

const (
	OPENAI_BASE_URL = "https://openai.qiniu.com/v1/chat/completions"
	OPENAI_API_KEY  = "sk-ac3498e9cd6bceb54a6f9f1e62518d85772a93a98e77c20652faa22ee9643764"
	PROMPT          = "请你精简这句话，提取它的主要表达内容："
	MODEL           = "doubao-seed-1.6-flash"
	ROLE            = "user"
)

func NewShrink(clientPool *http.Client, l logger.ZapLogger) Shrink {
	return &Shrinker{
		clientPool: clientPool,
		l:          l,
	}
}

func (s *Shrinker) Shrink(ctx context.Context, text string) (string, error) {
	reqData := request{
		Messages: []reqMsg{
			{
				Role:    ROLE,
				Content: PROMPT + text,
			},
		},
		Model:  MODEL,
		Stream: false,
	}
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		s.l.Error("Marshal Json Failed",
			logger.Field{
				Key: "error",
				Val: err.Error()},
			logger.Field{
				Key: "data",
				Val: reqData},
		)
		return "", ErrShrinkerMarshalData
	}
	req, err := http.NewRequest("POST", OPENAI_BASE_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		s.l.Error("Creat Http Request Failed", logger.Field{
			Key: "error",
			Val: err.Error(),
		})
		return "", ErrShrinkerSendHttpRequestFailed
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+OPENAI_API_KEY)
	resp, err := s.clientPool.Do(req)
	if err != nil {
		s.l.Error("Send Http Request Failed", logger.Field{
			Key: "error",
			Val: err.Error(),
		})
		return "", ErrShrinkerNotGetDataFromRemoteAPI
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			s.l.Warn("Close Body Failed", logger.Field{
				Key: "warn",
				Val: err.Error(),
			})
		}
	}(resp.Body)
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		s.l.Error("Bad Http Request",
			logger.Field{
				Key: "status",
				Val: resp.StatusCode},
			logger.Field{
				Key: "response",
				Val: resp.Body},
		)
		return "", ErrShrinkerNotGetDataFromRemoteAPI
	}
	var result response
	err = json.Unmarshal(body, &result)
	if err != nil {
		s.l.Error("Unmarshal Response Failed",
			logger.Field{
				Key: "error",
				Val: body})
		return "", ErrShrinkerMarshalData
	}
	for _, c := range result.Choices {
		return c.Message.Content, nil
	}
	return "", ErrShrinkerNotAvailable
}
