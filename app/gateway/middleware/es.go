package middleware

import (
	"douyin/config"
	"douyin/pkg/eslogrus"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"

	"github.com/sirupsen/logrus"

	elasticsearch "github.com/elastic/go-elasticsearch"
)

var (
	EsClient *elasticsearch.Client
)

func EsHookLog() *eslogrus.ElasticHook {
	hook, err := eslogrus.NewElasticHook(EsClient, config.Conf.ElasticSearch.Host, logrus.DebugLevel, "ESSever")
	if err != nil {
		panic(err)
	}

	return hook
}

// InitEs 初始化es
func EsInit() {
	esConn := fmt.Sprintf("http://%s", config.Conf.ElasticSearch.Addr)
	cfg := elasticsearch.Config{
		Addresses: []string{esConn},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Panic(err)
	}
	EsClient = client
}

func ElasticMiddleware(hook *eslogrus.ElasticHook) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 在处理请求前记录日志
		logrus.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		}).Info("Incoming request")

		// 处理请求
		c.Next()

		// 在处理请求后记录日志
		logrus.WithFields(logrus.Fields{
			"status": c.Writer.Status(),
		}).Info("Request handled")

		// 将日志发送到 Elasticsearch
		logEntry := &logrus.Entry{Data: map[string]interface{}{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"status": c.Writer.Status(),
		}}
		hook.Fire(logEntry)
	}
}
