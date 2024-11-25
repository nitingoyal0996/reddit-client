package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/cluster/clusterproviders/automanaged"
	"github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"
	"github.com/asynkron/protoactor-go/remote"
)

func main() {
	// read command line arguments
	if len(os.Args) < 2 || len(os.Args) > 4 {
		fmt.Println("Usage: ./client <Simulation #?> [userCount] [subredditCount]")
		os.Exit(1)
	}

	simulation, err := strconv.Atoi(os.Args[1])
	if err != nil || (simulation != 1 && simulation != 2 && simulation != 3) {
		fmt.Println("Usage: ./client <Simulation #?> [userCount] [subredditCount]")
		os.Exit(1)
	}

	userCount := 1000
	if len(os.Args) > 2 {
		userCount, err = strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid user count. Please provide a valid number.")
			os.Exit(1)
		}
	}

	subredditCount := 10
	if len(os.Args) > 3 {
		subredditCount, err = strconv.Atoi(os.Args[3])
		if err != nil || subredditCount >= userCount {
			fmt.Println("Invalid subreddit count. Please provide a valid number less than user count.")
			os.Exit(1)
		}
	}

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

	switch simulation {
	case 1:
		time.Sleep(2 * time.Second)
		println("Scenario 1 - Many users")
		SimulateUserCreation(system, userCount)
	case 2:
		time.Sleep(2 * time.Second)
		println("Scenario 2 - Many subreddits with zipf user distribution")
		SimulateZipfDistribution(system, subredditCount, userCount)
	case 3:
		time.Sleep(2 * time.Second)
		println("Scenario 3 - Connection and disconnection")
		duration := 5*time.Second
		SimulateConnectionAndDisconnection(system, duration, userCount)
	default:
		fmt.Println("Invalid simulation number. Please use 1, 2, or 3.")
		os.Exit(1)
	}
	println("Simulation completed. Exiting...")
	os.Exit(0)

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