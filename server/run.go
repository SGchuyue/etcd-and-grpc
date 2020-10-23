package server

import (
	"context"
	"encoding/json"
	"github.com/SGchuyue/logger/logger"
	"github.com/pkg/errors"
	"go.etcd.io/etcd/clientv3"
	"time"
)

type Info struct {
	Name string `json:"name"` // 服务名
	IP   string `json:"ip"`   // ip地址
	Type string `json:"type"` // 服务类型
}

type Service struct {
	ServiceInfo Info             // 服务信息
	end         chan error       // 服务结束
	leaseId     clientv3.LeaseID // 租约id
	cli         *clientv3.Client // 服务连接
}

// NewService 创建一个连接服务
func NewService(info Info, endpoints []string) (service *Service) {
	// 引入日志包
	logger.InitLogger("../debug.log", 2, 2, 3, true)
	// 创建连接服务
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 30 * time.Second,
	})
	if err != nil {
		logger.Error("连接发生错误：", err)
		return nil
	}
	service = &Service{
		ServiceInfo: info,
		cli:         cli,
	}
	return service
}

// Run 注册服务启动
func (service *Service) Run() (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	info := &service.ServiceInfo
	key := info.Name + "/" + info.IP + "/" + service.ServiceInfo.Type
	val, _ := json.Marshal(info)
	// 创建一个租约
	resp, err := service.cli.Grant(context.TODO(), 20)
	if err != nil {
		logger.Error("设置租约时出现错误:", err)
		return nil, err
	}
	_, err = service.cli.Put(context.TODO(), key, string(val), clientv3.WithLease(resp.ID))
	if err != nil {
		logger.Error("进行服务注册时发生错误:", err)
		return nil, err
	}
	service.leaseId = resp.ID
	return service.cli.KeepAlive(context.TODO(), resp.ID)
}

// Watch 监听服务
func (service *Service) Watch() (err error) {
	ch, err := service.Run()
	if err != nil {
		logger.Error("服务注册时发生错误", err)
		return
	}
	for {
		select {
		case err := <-service.end:
			return err
		case <-service.cli.Ctx().Done():
			return errors.New("服务已断开")
		case resp, ok := <-ch:
			// 监听租约
			if !ok {
				logger.Debug("注册服务运行已结束")
				return service.revoke()
			}
			logger.Debugf("%s 服务运行了ttl:%d", service.getKey(), resp.TTL)
		}
	}
}

// revoke 撤销取消给定的租约
func (service *Service) revoke() error {
	_, err := service.cli.Revoke(context.TODO(), service.leaseId)
	if err != nil {
		logger.Error("撤销取消给定的租约发生错误:", err)
	}
	logger.Debugf("服务:%s 即将停止\n", service.getKey())
	return err
}

// getKey 获取key的值
func (service *Service) getKey() string {
	return service.ServiceInfo.Name + "/" + service.ServiceInfo.IP + "/" + service.ServiceInfo.Type
}

// End 服务结束
func (service *Service) End() {
	service.end <- nil
}
