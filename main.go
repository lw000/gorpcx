package main

import (
	"context"
	"github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"log"
	"time"
)

type Args struct {
	A int
	B int
}

type Reply struct {
	C int
}

type Arith struct{}

func (t *Arith) Mul(ctx context.Context, args Args, reply *Reply) error {
	reply.C = args.A * args.B
	log.Println("C=", reply.C)
	return nil
}

func runClient() {
	d := client.NewPeer2PeerDiscovery("tcp@"+"localhost:8972", "")
	c := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer func() {
		_ = c.Close()
	}()

	args := &Args{
		A: 10,
		B: 20,
	}

	for {
		reply := &Reply{}
		err := c.Call(context.Background(), "Mul", args, reply)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		}

		log.Printf("%d * %d = %d", args.A, args.B, reply.C)
		time.Sleep(time.Millisecond)
	}
}

func addRegistryPlugin(s *server.Server) {
	r := &serverplugin.EtcdRegisterPlugin{
		ServiceAddress: "tcp@" + "localhost:8972",
		EtcdServers:    []string{"127.0.0.1:2379"},
		BasePath:       "/go-rpcx-server/",
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		log.Fatal(err)
	}
	s.Plugins.Add(r)
}

func main() {
	time.AfterFunc(time.Second*2, func() {
		runClient()
	})
	s := server.NewServer()
	//addRegistryPlugin(s)
	// s.Register(new(Arith), "")
	s.RegisterName("Arith", new(Arith), "")
	err := s.Serve("tcp", ":8972")
	if err != nil {
		log.Panic(err)
	}
}
