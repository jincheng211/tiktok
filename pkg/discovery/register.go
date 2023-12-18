package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

// 服务提供方
// https://blog.csdn.net/Mr_YanMingXin/article/details/126452076(教程)

type Register struct {
	EtcdAddrs   []string
	DialTimeout int

	closeCh     chan struct{}
	leasesID    clientv3.LeaseID
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse

	srvInfo Server
	srvTTL  int64
	cli     *clientv3.Client
	logger  *logrus.Logger
}

// NewRegister create a register based on etcd
func NewRegister(etcdAddrs []string, logger *logrus.Logger) *Register {
	return &Register{
		EtcdAddrs:   etcdAddrs,
		DialTimeout: 3,
		logger:      logger,
	}
}

// 将服务地址注册到etcd中
func (r *Register) Register(srvInfo Server, ttl int64) (chan<- struct{}, error) {
	if strings.Split(srvInfo.Addr, ":")[0] == "" {
		return nil, errors.New("invalid ip address")
	}
	var err error
	if r.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.EtcdAddrs,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	}); err != nil {
		return nil, err
	}

	r.srvInfo = srvInfo
	r.srvTTL = ttl

	if err = r.register(); err != nil {
		return nil, err
	}
	r.closeCh = make(chan struct{})

	go r.keepAlive()

	return r.closeCh, nil
}

func (r *Register) register() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.DialTimeout)*time.Second)
	defer cancel()

	leaseResp, err := r.cli.Grant(ctx, r.srvTTL)
	if err != nil {
		return err
	}

	r.leasesID = leaseResp.ID

	if r.keepAliveCh, err = r.cli.KeepAlive(context.Background(), r.leasesID); err != nil {
		return err
	}

	data, err := json.Marshal(r.srvInfo)
	if err != nil {
		return nil
	}

	// Put将一个键值对放入etcd
	_, err = r.cli.Put(context.Background(), BuildRegisterPath(r.srvInfo), string(data), clientv3.WithLease(r.leasesID))
	return err
}

// stop register
func (r *Register) Stop() {
	r.closeCh <- struct{}{}
}

// 删除节点
func (r *Register) unregister() error {
	_, err := r.cli.Delete(context.Background(), BuildRegisterPath(r.srvInfo))
	return err
}

// 保持服务器与etcd的长连接
func (r *Register) keepAlive() {
	// 定时器
	ticker := time.NewTicker(time.Duration(r.srvTTL) * time.Second)

	for {
		select {
		case <-r.closeCh: // 关闭通道
			if err := r.unregister(); err != nil {
				r.logger.Error("unregister failed, error: ", err)
			}
			// Revoke revokes the given lease.
			if _, err := r.cli.Revoke(context.Background(), r.leasesID); err != nil {
				r.logger.Error("revoke failed, error: ", err)
			}

		case res := <-r.keepAliveCh:
			if res == nil {
				if err := r.register(); err != nil {
					r.logger.Error("register failed, error: ", err)
				}
			}

		case <-ticker.C: // 触发器触发重新检测
			if r.keepAliveCh == nil {
				if err := r.register(); err != nil {
					r.logger.Error("register failed, error: ", err)
				}
			}
		}
	}
}
