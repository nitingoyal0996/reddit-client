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
	redditIDs []uint64
	redditMembers map[uint64]int
}

func NewSimulatorContext(system *actor.ActorSystem) *SimulatorContext {
	return &SimulatorContext{
		system:    system,
		consumers: make(map[int]*actor.PID),
		redditIDs: []uint64{},
		redditMembers: make(map[uint64]int),
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
		// stops chocking
		time.Sleep(100 * time.Millisecond)
	}
}

func (context *SimulatorContext) LoginUsers(startIdx int, endIdx int) {
	for i := startIdx; i < endIdx; i++ {
		context.system.Root.Send(context.consumers[i], &proto.LoginRequest{
			Username: fmt.Sprintf("user_%d", i),
			Password: fmt.Sprintf("password_%d", i),
		})
		// stop chocking
		time.Sleep(100 * time.Millisecond)
	}
}

func (context *SimulatorContext) GetActiveConsumerCount() int {
	return len(context.consumers)
}

func (context *SimulatorContext) CreateOneSubreddit(name int) {
	future := context.system.Root.RequestFuture(context.consumers[1], &proto.CreateSubredditRequest{
		Name:        fmt.Sprintf("subreddit_%d", name),
		Description: fmt.Sprintf("subreddit_%d description", name),
	}, 5*time.Second)
	res, err := future.Result()
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		subredditID := res.(*proto.CreateSubredditResponse).SubredditId
		context.redditIDs = append(context.redditIDs, subredditID)
	}
	println("Subreddit created with ID: ", context.redditIDs[len(context.redditIDs)-1])
}

func (context *SimulatorContext) GetZipfDistributionForMembers(subredditCount int, userCount int) {
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
		context.redditMembers[context.redditIDs[i]] = members
	}
}

func (context *SimulatorContext) AssignMembershipsToUsers() {
	// loop over the redditMembers map and assign memberships to users
	for subredditID, memberCount := range context.redditMembers {
		print("SubredditID: ", subredditID)
		for i := 0; i < memberCount; i++ {
			userIdx := i % len(context.consumers)
			print("UserIdx: ", userIdx)
			context.system.Root.Send(context.consumers[userIdx], &proto.SubscriptionRequest{
				SubredditId: subredditID,
			})
			// create 2 posts - every user who joins a subreddit creates 2 posts
			context.CreateOnePost(int(subredditID), userIdx)
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func (context *SimulatorContext) CreateOnePost(subredditId int, consumerIdx int) {
	consumer := context.consumers[consumerIdx]
	future := context.system.Root.RequestFuture(consumer, &proto.CreatePostRequest{
		Title:       fmt.Sprintf("post_%d", consumerIdx),
		Content:     fmt.Sprintf("post_%d content", consumerIdx),
		SubredditId: uint64(subredditId),
	}, 5*time.Second)
	res, err := future.Result()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	postResponse, ok := res.(*proto.CreatePostResponse)
	if !ok {
		fmt.Println("Error: failed to cast response to CreatePostResponse")
	} else {
		fmt.Println("Error: ", postResponse.Error)
	}
	println("Post created with ID: ", context.redditIDs[len(context.redditIDs)-1])
}

// GenerateUniqueRandomNumbers generates a shuffled list of numbers in the range [min, max].
func GenerateUniqueRandomNumbers(min, max int) ([]int, error) {
	if min > max {
		return nil, fmt.Errorf("invalid range: min (%d) > max (%d)", min, max)
	}

	// Create a slice with numbers from min to max
	numbers := make([]int, max-min+1)
	for i := range numbers {
		numbers[i] = min + i
	}

	// Shuffle the numbers
	rand.Seed(uint64(time.Now().UnixNano()))
	rand.Shuffle(len(numbers), func(i, j int) { numbers[i], numbers[j] = numbers[j], numbers[i] })

	return numbers, nil
}


func (context *SimulatorContext) LiveConnectDisconnect(simulationDuration time.Duration) {

	indices, err := GenerateUniqueRandomNumbers(0, len(context.consumers)-1)
	println("Indices: ", indices)
	// iterate over indices
	if err != nil {
		fmt.Println("Error: ", err)
	}
	for _, idx := range indices {
		startTime := time.Now()
		for time.Since(startTime) < simulationDuration {
			// Simulate login
			context.system.Root.Send(context.consumers[idx], &proto.LoginRequest{
				Username: fmt.Sprintf("user_%d", idx),
				Password: fmt.Sprintf("password_%d", idx),
			})
			onlineTime := time.Duration(rand.Intn(4)+1) * time.Second 
			time.Sleep(onlineTime)
			// Simulate logout
			context.system.Root.Send(context.consumers[idx], &proto.LogoutRequest{})
			offlineTime := time.Duration(rand.Intn(4)+1) * time.Second
			time.Sleep(offlineTime)
		}
	}
}