package client

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

// FCMClient sends push notifications through Firebase Cloud Messaging
type FCMClient interface {
	// Send pushes one notification to every token and returns the tokens FCM rejected as unregistered or invalid
	Send(ctx context.Context, tokens []string, title, body, imageURL string) (invalidTokens []string, err error)
}

type fcmClient struct {
	client *messaging.Client
}

func NewFCMClient() FCMClient {
	opt := option.WithCredentialsFile("cert/firebase/firebase-adminsdk.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic(err)
	}
	client, err := app.Messaging(context.Background())
	if err != nil {
		panic(err)
	}
	return &fcmClient{client: client}
}

func (f *fcmClient) Send(ctx context.Context, tokens []string, title, body, imageURL string) ([]string, error) {
	br, err := f.client.SendEachForMulticast(ctx, &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title:    title,
			Body:     body,
			ImageURL: imageURL,
		},
		Tokens: tokens,
	})
	if err != nil {
		return nil, err
	}
	var invalidTokens []string
	if br.FailureCount > 0 {
		for i, resp := range br.Responses {
			if messaging.IsUnregistered(resp.Error) || messaging.IsInvalidArgument(resp.Error) {
				invalidTokens = append(invalidTokens, tokens[i])
			}
		}
	}
	return invalidTokens, nil
}
