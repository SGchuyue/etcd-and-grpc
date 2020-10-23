package example

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"log"
	"testing"
	"time"
	s "watch_etcd/server"
)

func TestService_SendMail(t *testing.T) {
	r := s.NewResolver([]string{
		"127.0.0.1:2134",
		"127.0.0.1:8082",
		"127.0.0.1:8080",
	}, "g.srv.mail")
	resolver.Register(r)
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	addr := fmt.Sprintf("%s:///%s", r.Scheme(), "etcd grpc")
	_, err := grpc.DialContext(
		ctx, addr, grpc.WithInsecure(),
		// grpc.WithBalancerName(roundrobin.Name),
		//指定初始化round_robin => balancer (后续可以自行定制balancer和 register、resolver 同样的方式)
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		grpc.WithBlock(),
	)
	// 这种方式也行
	//conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBalancerName("round_robin"))
	//conn, err := grpc.Dial(":8972", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	//c := pb.TestResponse(conn)
	/*resp, err := Send(context.TODO(), &pb.TestRequest{
		Send: "test,test",
		Text: "test,test",
	})
	log.Print(resp)*/
}
