package cos

import (
	oss "QiNiuCloud/QiNiuCloud/pkg/ossx"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"strings"
)

var (
	ErrUploadToOSSFailed = errors.New("upload to cos failed")
)

type OssHandler struct {
	oss *cos.Client
}

func NewCOSHandler(oss *cos.Client) oss.OSSHandler {
	return &OssHandler{oss: oss}
}

const salt = "xiFgge1O4DqWs5og"

func (hdl *OssHandler) filenameToUniqueWithSalt(filename string) string {
	// 结合文件名和盐值生成哈希，降低碰撞概率
	data := []byte(filename + "|" + salt)
	fmt.Println(data)
	hash := sha256.Sum256(data)
	fmt.Println(hash)
	return hex.EncodeToString(hash[:])
}

// Upload filename example C:\Users\lyl69\GolandProjects\anqi\recoder\downloads\抖音直播\演员_陈安琪_2025-08-21_05-54-50_000.mp4
func (hdl *OssHandler) Upload(ctx context.Context, fileDir string, preview bool) (string, string, error) {
	splitDir := strings.Split(fileDir, "\\")
	rawName := splitDir[len(splitDir)-1]
	filename := hdl.filenameToUniqueWithSalt(splitDir[len(splitDir)-1])
	_, _, err := hdl.oss.Object.Upload(context.Background(), filename+".mp4", fileDir, nil)
	if err != nil {
		return "", "", err
	}
	return filename, rawName, nil
}

func (hdl *OssHandler) Find(ctx context.Context, uid int64) error {
	return nil
}
func (hdl *OssHandler) Delete(ctx context.Context, uid int64) error {
	return nil
}
