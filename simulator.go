package main

import (
	"client/proto"
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
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
func (context *SimulatorContext) RegisterUsers() {
	for i := 0; i < 10; i++ {
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

// // create a subreddit
// func (context *SimulatorContext) CreateSubreddit() {
// 	context.RegisterUsers(0, 1)
// 	context.LoginUsers(0, 1)
// 	future := context.system.Root.RequestFuture(context.consumers[0], &messages.Subreddit{
// 		Name:        "subreddit_1",
// 		Description: "This is a test subreddit",
// 	}, 5*time.Second)
// 	res, err := future.Result()
// 	if err != nil {
// 		fmt.Println("Error: ", err)
// 	} else {
// 		context.redditID = res.(*messages.CreateSubredditResponse).SubredditId
// 	}
// }

// func (context *SimulatorContext) SimulateSubscriptions(userCount int) {
//     // Initialize random source with Zipf distribution
// 	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
//     zipf := rand.NewZipf(r, 2.0, 1.0, uint64(userCount ))

//     // Track subscription attempts
//     subscriptionCount := zipf.Uint64()

//     // Simulate subscriptions following Zipf distribution
//     for i := uint64(0); i < subscriptionCount; i++ {
// 		future := context.system.Root.RequestFuture(context.consumers[int(i)], &messages.Join{
//             SubredditId: context.redditID,
//         }, 5*time.Second)
        
//         if _, err := future.Result(); err != nil {
//             fmt.Printf("Subscription failed for user %d: %v\n", i, err)
//         }
// 		time.Sleep(100 * time.Millisecond)
//     }
// }
