package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"time"

	"github.com/Bekhzood/telegramBotAPI/telegramBotpb"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
	low  = iota
	medium
	high
)

var database_url string
var telegramUrl string
var chatId string

//var messagesList []*telegramBotpb.Message

type MessageServer struct {
	conn *pgx.Conn
	telegramBotpb.MessageServiceServer
}

func main() {
	LoadEnv()
	messageServer := NewMessageServer()
	messageServer.conn = CreateConn()
	if err := messageServer.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func NewMessageServer() *MessageServer {
	return &MessageServer{}
}

func CreateConn() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), database_url)

	if err != nil {
		log.Fatalf("Unable connect to server: %v", err)
	}
	defer conn.Close(context.Background())
	return conn
}

func (server *MessageServer) Run() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	telegramBotpb.RegisterMessageServiceServer(s, server)
	log.Printf("server listening at %v", lis.Addr())

	return s.Serve(lis)
}

func LoadEnv() {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	database_url = "host=localhost dbname=restapi_dev sslmode=disable password=" + os.Getenv("DB_PASSWORD")
	telegramUrl = "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_BOT_TOKEN") + "/sendMessage"
	chatId = os.Getenv("TELEGRAM_CHAT_ID")
}

func (server *MessageServer) Send(ctx context.Context, req *telegramBotpb.SendMessageRequest) (*telegramBotpb.SendMessageResponse, error) {

	message := req.GetMessage()
	//messagesList = append(messagesList, message...)
	SortSlice(message)

	for i := 0; i < len(message); i++ {
		time.Sleep(5 * time.Second)
		tgResponse, err := http.PostForm(telegramUrl,
			url.Values{
				"chat_id": {chatId},
				"text":    {message[i].Text},
			})

		if err != nil {
			log.Printf("error when posting text to the chat: %s", err.Error())
			return nil, err
		}
		defer tgResponse.Body.Close()
		//message = RemoveIndex(message, i)
	}

	response := telegramBotpb.SendMessageResponse{
		Message: message,
	}

	return &response, nil
}

func RemoveIndex(slice []*telegramBotpb.Message, index int) []*telegramBotpb.Message {
	return append(slice[:index], slice[index+1:]...)
}

func SortSlice(slice []*telegramBotpb.Message) {
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].Priority > slice[j].Priority
	})
}
