package main

import (
	"context"
	"fmt"
	pb "github.com/pingjing0628/mail-server/proto"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

func main() {
	// 連線到遠端 gRPC 伺服器。
	conn, err := grpc.Dial("server:9999", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("連線失敗：%v", err)
	}
	defer conn.Close()

	// 建立新的 Mail 客戶端，所以等一下就能夠使用 Mail 的所有方法。
	c := pb.NewMailClient(conn)

	// 傳送新請求到遠端 gRPC 伺服器 Mail 中，並呼叫 Send 函式
	mr := pb.MailRequest{
		From:    "fish1063345@gmail.com",
		To:      []string{"m10509216@gapps.ntust.edu.tw"},
		Cc:      []string{},
		Subject: "How to use gRPC",
		Body:    "Just done",
		Type:    "text/html",
	}
	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		ret, err := c.Send(context.Background(), &mr)
		if err != nil {
			log.Fatalf("無法執行 Send 函式：%v", err)
		} else {
			fmt.Fprintf(w, "Send %s", ret.Code)
		}
	})

	http.ListenAndServe(":9000", nil)
}