package main

import (
	"client/proto"
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
)

type ConsumerActor struct {
	id       uint
	username string
	token    string
}

func (consumer *ConsumerActor) Receive(context actor.Context) {
	msg := context.Message()
	fmt.Printf("Received message of type: %T\n", msg)
	// print consumer pid
	fmt.Printf("Consumer PID: %v\n", context.Self())
	switch actorMsg := msg.(type) {
	case *proto.RegisterRequest:
		fmt.Printf("Register the user")
		consumer.Register(context, actorMsg)
	case *proto.LoginRequest:
		fmt.Printf("Login the user")
		// simulation connection
		consumer.Login(context, actorMsg)
	case *proto.Logout:
		fmt.Printf("Logout the user")
		// simulate disconnection
	case *proto.Subreddit:
		fmt.Printf("Create Subreddit")
	case *proto.Join:
		fmt.Printf("Join Subreddit")
	case *proto.Post:
		fmt.Printf("Create Post")
	default:
		fmt.Println("Unknown Message: ", consumer.token)
	}
}

func (consumer *ConsumerActor) Register(context actor.Context, actorMsg *proto.RegisterRequest) {
	authActor := cluster.GetCluster(context.ActorSystem()).Get("auth", "Auth")
	future := context.RequestFuture(authActor, &proto.RegisterRequest{
		Username: actorMsg.Username,
		Email:    actorMsg.Email,
		Password: actorMsg.Password,
	}, 5*time.Second)
	res, err := future.Result()
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		// save the username, Id 
		consumer.id = uint(res.(*proto.RegisterResponse).Id)
		consumer.username = actorMsg.Username
	}
}

func (consumer *ConsumerActor) Login(context actor.Context, actorMsg *proto.LoginRequest) {
	authActor := cluster.GetCluster(context.ActorSystem()).Get("auth", "Auth")
	future := context.RequestFuture(authActor, &proto.LoginRequest{
		Username: consumer.username,
		Password: actorMsg.Password,
	}, 5*time.Second)
	res, err := future.Result()
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		consumer.token = res.(*proto.LoginResponse).Token
	}
}

// func (consumer *ConsumerActor) CreateSubreddit(context actor.Context, actorMsg *proto.Subreddit) {
// 	subredditActor := cluster.GetCluster(context.ActorSystem()).Get("subreddit", "Subreddit")
// 	future := context.RequestFuture(subredditActor, &messages.CreateSubredditRequest{
// 		CreatorId:   uint64(consumer.id),
// 		Token:       consumer.token,
// 		Name:        actorMsg.Name,
// 		Description: actorMsg.Description,
// 	}, 5*time.Second)
// 	res, err := future.Result()
// 	if err != nil {
// 		fmt.Println("Error: ", err)
// 	} else {
// 		subredditResponse, ok := res.(*proto.CreateSubredditResponse)
// 		if !ok {
// 			fmt.Println("Invalid response type")
// 			return
// 		}
// 		if subredditResponse.Error != "" {
// 			fmt.Println("Error: ", subredditResponse.Error)
// 			return
// 		}
// 		fmt.Printf("Subreddit created with ID: %d\n", subredditResponse.SubredditId)
// 		context.Respond(&messages.CreateSubredditResponse{SubredditId: subredditResponse.SubredditId})
// 	}
// }

// func (consumer *ConsumerActor) JoinSubreddit(context actor.Context, actorMsg *proto.Join) {
// 	subredditActor := cluster.GetCluster(context.ActorSystem()).Get("subreddit", "Subreddit")
// 	future := context.RequestFuture(subredditActor, &messages.SubscriptionRequest{
// 		UserId:      uint64(consumer.id),
// 		Token:       consumer.token,
// 		SubredditId: actorMsg.SubredditId,
// 	}, 5*time.Second)
// 	res, err := future.Result()
// 	if err != nil {
// 		fmt.Println("Error: ", err)
// 	} else {
// 		joinResponse, ok := res.(*proto.SubscriptionResponse)
// 		if !ok {
// 			fmt.Println("Invalid response type")
// 			return
// 		}
// 		if !joinResponse.Success {
// 			fmt.Println("Error: ", joinResponse.Message)
// 			return
// 		}
// 		fmt.Printf("Subreddit joined with ID: %d\n", actorMsg.SubredditId)
// 		context.Respond(&messages.SubscriptionResponse{Success: true, Message: ""})
// 	}
// }
