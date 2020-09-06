package main

import (
	"context"
	//"crypto/tls"
	pb "github.com/pingjing0628/mail-server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	gomail "gopkg.in/gomail.v2"
	"log"
	"net"
	"time"
)

type server struct {}

var ch = make(chan *gomail.Message)

func (s *server) Send(ctx context.Context, mail *pb.MailRequest) (*pb.MailStatus, error) {
	m := gomail.NewMessage()
	m.SetHeader("From", mail.From)
	m.SetHeader("To", mail.To...)
	m.SetHeader("Subject", mail.Subject)
	m.SetBody(mail.Type, mail.Body)
	ch <- m

	return &pb.MailStatus{Status: int32(0), Code: ""}, nil
}

func main()  {
	listen, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatalf("cannot listen to this port: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterMailServer(s, &server{})
	reflection.Register(s)

	go func() {
		d := gomail.NewDialer("mail.gandi.net", 587, "service@lumos.tw", "Zxasjk15236")

		//d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

		var s gomail.SendCloser
		var err error
		open := false

		for {
			select {
			case m, ok := <- ch:
				if !ok {
					return
				}

				if !open {
					if s, err = d.Dial(); err != nil {
						panic(err)
					}

					open = true
				}

				if err := gomail.Send(s, m); err != nil {
					log.Print(err)
				}
			// Close the connection to the SMTP server if no email was sent in the last 30 seconds.
			case <- time.After(30 * time.Second):
				if open {
					if err := s.Close(); err != nil {
						panic(err)
					}

					open = false
				}
			}
		}
	}()

	if err := s.Serve(listen); err != nil {
		log.Fatalf("cannot provide service: %v", err)
		close(ch)
	}
}
