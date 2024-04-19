package main

import (
	"fmt"

	"cluster-basic/shared"

	console "github.com/asynkron/goconsole"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/cluster/clusterproviders/consul"
	"github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"
	"github.com/asynkron/protoactor-go/remote"
)

func main() {
	cluster := startNode()

	fmt.Print("\nBoot other nodes and press Enter\n")
	console.ReadLine()

	// enterを押したら、clusterが終了する？
	cluster.Shutdown(true)
}

func startNode() *cluster.Cluster {
	// node1とは異なるactor systemを作成
	system := actor.NewActorSystem()

	provider, _ := consul.New()
	lookup := disthash.New()
	// node1とは異なるマシンなので、同じlocalhost:0でも別のアドレスとして扱われる？
	config := remote.Configure("localhost", 0)

	props := actor.PropsFromFunc(func(ctx actor.Context) {
		switch msg := ctx.Message().(type) {
		case *actor.Started:
			fmt.Printf("Started %v", msg)
		case *shared.HelloRequest:
			fmt.Printf("Hello %v\n", msg.Name)
			ctx.Respond(&shared.HelloResponse{})
		}
	})

	// kindはnode1で定義したものと同じものを指定している
	helloKind := cluster.NewKind("hello", props)
	// cludter.WithKindsで、指定したkindをclusterに登録している
	// my-clusterという名前はnode1と共通のため、node1と同じclusterに所属している
	// consulが同じnetwork内にいる前提。
	clusterConfig := cluster.Configure("my-cluster", provider, lookup, config, cluster.WithKinds(helloKind))
	c := cluster.New(system, clusterConfig)

	c.StartMember()
	return c
}
