package oss

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/sms/bytes"
	"github.com/qiniu/go-sdk/v7/storage"
	"strconv"
	"time"
)

var (
	cfg     = storage.Config{}
	mac     = auth.Credentials{}
	upToken string
)

func Init_oss() {
	bucket := "jincheng211"
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	// 从七牛云官网获取
	accessKey := "0ewqxYH5BZXV_mTxfxYhKfarPqkkyX20KaJB7Me2"
	secretKey := "fXqt51eAuda1D0f8fv1FW5h8jiNQL8bnnESIeLXG"
	mac = *qbox.NewMac(accessKey, secretKey)
	upToken = putPolicy.UploadToken(&mac)

	// 空间对应的机房
	cfg.Region = &storage.ZoneHuadong
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
}

func PutVideo(video []byte, cover []byte, timeNow int64, id int64) (err error) {
	videoKey := "douyinVideoList/" + strconv.FormatInt(id, 10) + "/" + strconv.FormatInt(timeNow, 10) + ".mp4"
	coverKey := "douyinVideoList/" + strconv.FormatInt(id, 10) + "/" + strconv.FormatInt(timeNow, 10) + ".jpg"

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}

	// 视频上传
	videoLen := int64(len(video))
	err = formUploader.Put(context.Background(), &ret, upToken, videoKey, bytes.NewReader(video), videoLen, &putExtra)
	if err != nil {
		fmt.Println("视频上传失败")
		return nil
	}

	// 封面上传
	coverLen := int64(len(cover))
	err = formUploader.Put(context.Background(), &ret, upToken, coverKey, bytes.NewReader(cover), coverLen, &putExtra)
	if err != nil {
		fmt.Println("图片上传失败")
		return nil
	}

	return nil
}

func GetVideo(url string) string {

	domain := "http://rzddmop0p.hd-bkt.clouddn.com"
	key := url

	deadline := time.Now().Add(time.Second * 3600 * 24 * 30).Unix() //1小时有效期
	privateAccessURL := storage.MakePrivateURL(&mac, domain, key, deadline)

	return privateAccessURL
}
