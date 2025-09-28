package Providers

import (
	"QiNiuCloud/QiNiuCloud/pkg/AsyncModelGenerationTaskManager/event"
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"encoding/json"
	"fmt"
	ai3d "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ai3d/v20250513"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"net/http"
)

type HunyuanTo3D struct {
	l          logger.ZapLogger
	client     *ai3d.Client
	credential common.Credential
	producer   event.ModelProviderResultProducer
}

const Endpoint = "ai3d.tencentcloudapi.com "

func (h *HunyuanTo3D) SubmitTask(client *http.Client, bizData interface{}) *Result {
	request := ai3d.NewSubmitHunyuanTo3DJobRequest()
	request.Prompt = common.StringPtr("生成一只鸡")
	request.ResultFormat = common.StringPtr("OBJ")
	response, err := h.client.SubmitHunyuanTo3DJob(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return &Result{
			Provider: "HunyuanTo3D",
			Err:      err,
			Msg:      "HunyuanTo3D SDK error",
		}
	}
	if err != nil {
		return &Result{
			Provider: "HunyuanTo3D",
			Err:      err,
			Msg:      "HunyuanTo3D Submit task failed",
		}
	}
	return &Result{
		Provider:  "HunyuanTo3D",
		Err:       nil,
		Msg:       "HunyuanTo3D Submit task successfully",
		JobId:     *response.Response.JobId,
		RequestId: *response.Response.RequestId,
	}
}

func (h *HunyuanTo3D) QueryTask(client *http.Client, txId string, jobId string, token string) *Result {
	request := ai3d.NewQueryHunyuanTo3DJobRequest()
	request.JobId = common.StringPtr(jobId)
	response, err := h.client.QueryHunyuanTo3DJob(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return &Result{
			Provider: "HunyuanTo3D",
			Err:      err,
			Msg:      "HunyuanTo3D SDK error",
		}
	}
	if err != nil {
		return &Result{
			Provider: "HunyuanTo3D",
			Err:      err,
			Msg:      "HunyuanTo3D Submit task failed",
		}
	}
	var r Resp
	err = json.Unmarshal([]byte(response.ToJsonString()), &r)
	if err != nil {
		return &Result{
			Provider: "HunyuanTo3D",
			Err:      err,
			Msg:      "HunyuanTo3D Submit task response unmarshal failed",
		}
	}
	if *response.Response.Status == "DONE" {
		err = h.ReportResult(txId, token, &ModelGenerationTaskResult{
			jobiD: "HUNYUAN" + r.Response.RequestId,
			token: token,
			url:   r.Response.ResultFile3Ds[0].Url,
			thumb: r.Response.ResultFile3Ds[0].PreviewImageUrl,
		})
		if err != nil {
			return &Result{
				Provider: "HunyuanTo3D",
				Err:      err,
				Msg:      fmt.Sprintf("HunyuanTo3D Submit task result failed, request id%s", r.Response.RequestId),
			}
		}
	}
	return &Result{
		Provider: "HunyuanTo3D",
		Err:      nil,
		Msg:      r.Response.Status,
	}
}
func (h *HunyuanTo3D) GetProviderName() string {
	return "HunyuanTo3D"
}
func (h *HunyuanTo3D) ReportResult(txId string, token string, bizData interface{}) error {
	evt := event.AddEvent{
		DATA: bizData,
	}
	err := h.producer.AddEvent(evt)
	if err != nil {
		h.l.Error(err.Error())
		return err
	}
	return nil
}

type Resp struct {
	Response struct {
		ErrorCode     string `json:"ErrorCode"`
		ErrorMessage  string `json:"ErrorMessage"`
		RequestId     string `json:"RequestId"`
		ResultFile3Ds []struct {
			PreviewImageUrl string `json:"PreviewImageUrl"`
			Type            string `json:"Type"`
			Url             string `json:"Url"`
		} `json:"ResultFile3Ds"`
		Status string `json:"Status"`
	} `json:"Response"`
}

func NewHunyuanTo3D(l logger.ZapLogger, client *ai3d.Client) ProviderSpecificGenerator {
	//credential := common.NewCredential(
	//	os.Getenv("TENCENTCLOUD_SECRET_ID"),
	//	os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	//)
	return &HunyuanTo3D{
		l:      l,
		client: client,
	}
}
