package main

import (
	console "github.com/asynkron/goconsole"
	"github.com/asynkron/protoactor-go/actor"
	"log/slog"
)

type (
	hello      struct{ Who string }
	helloActor struct{}
)

// helloActorがmessageを受け取るときの挙動を定義
// Sendされたら、"Receive"が呼び出されるようになっている
func (state *helloActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *hello:
		// actor contextにlogを出せる
		// slog.String("key", value)でkey-valueの形でlogを出せる
		context.Logger().Info("Hello ", slog.String("who", msg.Who))
	}
}

func main() {
	// お決まりのactor systemの作成
	system := actor.NewActorSystem()
	// propertiesの作成
	// propsはactorの生成方法を定義する
	// Producerとは、Actorを生成する関数型
	// helloActor(Actor型)を返す無名関数をProducerに渡している
	// spawnするときに、この関数に従ってActorが生成される
	// 今回だと、spawnされたらhelloActorがreturnされる
	props := actor.PropsFromProducer(func() actor.Actor { return &helloActor{} })

	// spawn(生成)するお作法
	pid := system.Root.Spawn(props)

	// Send(相手, message)
	system.Root.Send(pid, &hello{Who: "Roger"})

	// systemが終了しないようにする定石
	_, _ = console.ReadLine()
}
