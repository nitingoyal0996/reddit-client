package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/cluster/clusterproviders/automanaged"
	"github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"
	"github.com/asynkron/protoactor-go/remote"
)

func main() {
	// Set up actor system
	system := actor.NewActorSystem()

	// Prepare a remote env that listens to 8081
	remoteConfig := remote.Configure("127.0.0.1", 8081)
	cp := automanaged.NewWithConfig(1*time.Second, 6330, "localhost:6331")
	lookup := disthash.New()
	clusterConfig := cluster.Configure("reddit-cluster", cp, lookup, remoteConfig)
	c := cluster.New(system, clusterConfig)
	fmt.Println("Cluster client started")
	// Manage the cluster client's lifecycle
	c.StartClient() // Configure as a client
	defer c.Shutdown(false)
	
	// simulators
	time.Sleep(5 * time.Second)
	Simulator := NewSimulatorContext(system)
	
	// Tests
	Simulator.RegisterUsers()
	Simulator.LoginUsers()

	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt, os.Kill)
	<-finish
}