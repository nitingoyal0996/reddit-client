package main

import (
	"client/messages"
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
)

type ConsumerActor struct {
	id uint
	username string
	token string
	reddits []string
	posts []string
	messages []string
}

func (consumer *ConsumerActor) Receive(context actor.Context) {
	msg := context.Message()
	fmt.Printf("Received message of type: %T\n", msg)
	switch actorMsg := msg.(type) {
	case *messages.Register:
		fmt.Printf("Register the user")
		consumer.Register(context, actorMsg)
	case *messages.Login:
		fmt.Printf("Login the user")
		consumer.Login(context, actorMsg)
	case *messages.Logout:
		fmt.Printf("Logout the user")
	case *messages.Subreddit:
		fmt.Printf("Create Subreddit")
	case *messages.Join:
		fmt.Printf("Join Subreddit")
	case **messages.Leave:
		fmt.Printf("Leave Subreddit")
	case *messages.Post:
		fmt.Printf("Create Post")
	case *messages.Comment:
		fmt.Printf("Create Comment")
	default:
		// print username
		fmt.Println("Unknown Message: ", consumer.token)
	}
}

func (consumer *ConsumerActor) Register(context actor.Context, actorMsg *messages.Register) {
	authActor := cluster.GetCluster(context.ActorSystem()).Get("auth", "Auth")
	future := context.RequestFuture(authActor, &messages.RegisterRequest{
		Username: actorMsg.Username,
		Email:    actorMsg.Email,
		Password: actorMsg.Password,
	}, 5*time.Second)
	res, err := future.Result()
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		// save the username, Id 
		consumer.id = uint(res.(*messages.RegisterResponse).Id)
		consumer.username = actorMsg.Username
	}
}

func (consumer *ConsumerActor) Login(context actor.Context, actorMsg *messages.Login) {
	authActor := cluster.GetCluster(context.ActorSystem()).Get("auth", "Auth")
	future := context.RequestFuture(authActor, &messages.LoginRequest{
		Username: consumer.username,
		Password: actorMsg.Password,
	}, 5*time.Second)
	res, err := future.Result()
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		consumer.token = res.(*messages.LoginResponse).Token
	}
}
