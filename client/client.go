package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Bekhzood/telegramBotAPI/telegramBotpb"
	"google.golang.org/grpc"
)

const (
	low = iota
	medium
	high
)

func main() {
	fmt.Println("starting client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("could not connect: %v", err)
	}

	defer conn.Close()

	client := telegramBotpb.NewMessageServiceClient(conn)
	SendMessage(client)
}

func SendMessage(client telegramBotpb.MessageServiceClient) {
	request := &telegramBotpb.SendMessageRequest{
		Message: []*telegramBotpb.Message{
			{
				Text:     "First low World",
				Priority: low,
			},
			{
				Text:     "Second medium World",
				Priority: medium,
			},
			{
				Text:     "Third high World",
				Priority: high,
			},
			{
				Text:     "Forth low World",
				Priority: low,
			},
		},
	}
	response, err := client.Send(context.Background(), request)
	if err != nil {
		log.Fatalf("error while calling Send RPC: %v", err)
	}
	log.Printf("Response from Message: %v", response.Message)
}
