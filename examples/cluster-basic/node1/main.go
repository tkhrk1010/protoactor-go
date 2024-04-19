package main

import (
	"cluster-basic/shared"
	"fmt"
	console "github.com/asynkron/goconsole"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/cluster/clusterproviders/consul"
	"github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"
	"github.com/asynkron/protoactor-go/remote"
)

func main() {
	// clusterを開始し、シードノードを作成
	c := startNode()

	// enterを押されたら次に進む
	fmt.Print("\nBoot other nodes and press Enter\n")
	console.ReadLine()

	// kind: helloのactorをid: abcでcluster内に作成してpidを取得
	// kindはAkkaで言うroleのようなもの。多分
	// 
	pid := c.Get("abc", "hello")
	fmt.Printf("Got pid %v", pid)

	// Requestは、Sendと非常に似ているが、受信アクターが応答を送信できるように、送信元PIDも含まれる
	// Request is very similar to Send, but it also includes the sender PID so that the receiving actor can send a reply. 
	// https://proto.actor/docs/pid/#request
	res, _ := c.Request("abc", "hello", &shared.HelloRequest{Name: "Roger"})
	fmt.Printf("Got response %v", res)

	fmt.Println()
	console.ReadLine()
	c.Shutdown(true)
}

func startNode() *cluster.Cluster {
	system := actor.NewActorSystem()

	// consulとは、分散システムのためのキー/値ストアとサービス登録機能を提供するソフトウェア
	// ここでは、clusterを構築するために必要な人、というイメージでいい
	// https://proto.actor/docs/bootcamp/unit-8/lesson-3/#using-consul-in-cluster
	provider, _ := consul.New()
	// disthashは、クラスタメンバーの検索を行うためのインターフェース
	// cluster memberとは、クラスタ内のactorのこと
	lookup := disthash.New()
	// remote.Configureは、リモートアクターの設定を行う
	// ここでは、localhost:0のアドレスを設定している
	config := remote.Configure("localhost", 0)
	// 引数は、cluster名、provider、lookup、config
	clusterConfig := cluster.Configure("my-cluster", provider, lookup, config)
	c := cluster.New(system, clusterConfig)
	// StartMemberとは、cluster内のactorを開始するメソッド
	c.StartMember()

	// StartNodeといいつつ、clusterを返す
	return c
}
