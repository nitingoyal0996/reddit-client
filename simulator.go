package main

import (
	"client/messages"
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
)

type SimulatorContext struct {
	system    *actor.ActorSystem
	consumers map[int]*actor.PID
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
		context.system.Root.Send(consumerActor, &messages.Register{
			Username: username,
			Email:    email,
			Password: password,
		})
		context.consumers[i] = consumerActor
	}
}

func (context *SimulatorContext) LoginUsers() {
	for i := 0; i < 10; i++ {
		context.system.Root.Send(context.consumers[i], &messages.Login{
			Username: fmt.Sprintf("user_%d", i),
			Password: fmt.Sprintf("password_%d", i),
		})
	}
}


// func (context *SimulatorContext) GenerateSubredditsZipf(count int) error {
//     // Create Zipf distribution
//     zipf := rand.NewZipf(rand.New(rand.NewSource(uint64(time.Now().UnixNano()))), 1.8, 1, uint64(count))
    
//     for i := 0; i < count; i++ {
//         subscribers := zipf.Uint64()
//         subreddit := &models.Subreddit{
//             Name:           fmt.Sprintf("subreddit_%d", i),
//             SubscriberCount: int(subscribers),
//             PostCount:      calculatePostCount(subscribers),
//             CreatedAt:      randomDate(),
//         }
//         // Insert into database
//     }
//     return nil
// }
