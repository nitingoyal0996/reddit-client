package main

import (
	"client/proto"
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"golang.org/x/exp/rand"
)

type SimulatorContext struct {
	system    *actor.ActorSystem
	consumers map[int]*actor.PID
	reddit	  *actor.PID
	redditID  uint64
}

func NewSimulatorContext(system *actor.ActorSystem) *SimulatorContext {
	return &SimulatorContext{
		system:    system,
		consumers: make(map[int]*actor.PID),
	}
}
func (context *SimulatorContext) RegisterUsers(startIdx int, endIdx int) {
	for i := startIdx; i < endIdx; i++ {
		username := fmt.Sprintf("user_%d", i)
		email := fmt.Sprintf("test_%d@example.com", i)
		password := fmt.Sprintf("password_%d", i)

		consumerProps := actor.PropsFromProducer(func() actor.Actor {
			return &ConsumerActor{}
		})
		consumerActor := context.system.Root.Spawn(consumerProps)
		context.system.Root.Send(consumerActor, &proto.RegisterRequest{
			Username: username,
			Email:    email,
			Password: password,
		})
		context.consumers[i] = consumerActor
		time.Sleep(100 * time.Millisecond)
	}
}

func (context *SimulatorContext) LoginUsers(startIdx int, endIdx int) {
	for i := startIdx; i < endIdx; i++ {
		context.system.Root.Send(context.consumers[i], &proto.LoginRequest{
			Username: fmt.Sprintf("user_%d", i),
			Password: fmt.Sprintf("password_%d", i),
		})
	}
}

func (context *SimulatorContext) GetActiveConsumerCount() int {
    return len(context.consumers)
}

func (context *SimulatorContext) CreateOneSubreddit() {
	future:=context.system.Root.RequestFuture(context.consumers[1], &proto.CreateSubredditRequest{
		Token:       "",
		Name:        "subreddit_1",
		Description: "This is a test subreddit",
		CreatorId:   1,
	}, 5*time.Second)
	res, err := future.Result()
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		context.redditID = res.(*proto.CreateSubredditResponse).SubredditId
	}
	println("Subreddit created with ID: ", context.redditID)
}

func (context *SimulatorContext) GetZipfDistributionForMembers(subredditCount int, userCount int) []int {
    // Seed the random number generator
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	// Parameters for Zipf distribution
	s := 2.0  // Skew factor (Zipf exponent)
	v := 1.0  // Minimum value
	max := uint64(subredditCount)

	zipf := rand.NewZipf(r, s, v, max)

	// Generate raw Zipf distribution
	rawCounts := make([]uint64, subredditCount)
	var rawSum uint64 = 0
	for i := 0; i < subredditCount; i++ {
		rawCounts[i] = zipf.Uint64() + 1 // Ensure no subreddit gets 0 members
		rawSum += rawCounts[i]
	}

	// Normalize raw counts to match userCount
	scaledCounts := make([]int, subredditCount)
	totalScaled := 0
	for i := 0; i < subredditCount; i++ {
		scaledCounts[i] = int(float64(rawCounts[i]) / float64(rawSum) * float64(userCount))
		totalScaled += scaledCounts[i]
	}

	// Adjust for rounding errors
	diff := userCount - totalScaled
	if diff != 0 {
		scaledCounts[0] += diff // Add/subtract difference to the largest subreddit
	}
	// Print results
	for i, members := range scaledCounts {
		fmt.Printf("Subreddit %d: %d members\n", i+1, members)
	}

	return scaledCounts
}
