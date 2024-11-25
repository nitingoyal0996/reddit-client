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
	
	// time.Sleep(2 * time.Second)	
	// println("Scenario 1 - Many users")
	// SimulateUserCreation(system, 1000)
	// time.Sleep(2 * time.Second)
	// println("Scenario 2 - Many subreddits with zipf user distribution")
	// SimulateZipfDistribution(system, 2, 10)
	time.Sleep(2 * time.Second)
	println("Scenario 3 - Connection and disconnection")
	SimulateConnectionAndDisconnection(system, 5 * time.Second, 10)

	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt, os.Kill)
	<-finish
}

// Tests
// Scenario 1 - Many users
func SimulateUserCreation(system *actor.ActorSystem, userCount int) {
	Simulator := NewSimulatorContext(system)
	Simulator.RegisterUsers(0, userCount)
	time.Sleep(1 * time.Second)
	Simulator.LoginUsers(0, userCount)
}

// Scenario 2 - Many subreddits with zipf user distribution
func SimulateZipfDistribution(system *actor.ActorSystem, subCount int, userCount int) {
	Simulator := NewSimulatorContext(system)
	Simulator.RegisterUsers(0, userCount)
	time.Sleep(1 * time.Second)
	Simulator.LoginUsers(0, userCount)
	time.Sleep(1 * time.Second)
	for i := 0; i < subCount; i++ {
		Simulator.CreateOneSubreddit(i)
		time.Sleep(50 * time.Millisecond)
	}
	Simulator.GetZipfDistributionForMembers(subCount, userCount)
	time.Sleep(1 * time.Second)
	Simulator.AssignMembershipsToUsers()
	time.Sleep(1 * time.Second)
}

// Scenario 3 - Connection and disconnection
func SimulateConnectionAndDisconnection(system *actor.ActorSystem, simulationDuration time.Duration, userCount int) {
	Simulator := NewSimulatorContext(system)
	Simulator.RegisterUsers(0, userCount)
	Simulator.LiveConnectDisconnect(simulationDuration)
}