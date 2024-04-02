package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
)

type IPubsubClient interface {
	NewClient(ctx context.Context, projectID string) (IPubsubClient, error)
	Close() error
	Subscription(ctx context.Context, subID string) (IPubsubSubscription, error)
	GetSubscription(subID string) (IPubsubSubscription, error)
}

type IPubsubSubscription interface {
	Receive(ctx context.Context, f func(context.Context, *pubsub.Message)) error
}

type IPubsubTopic interface {
	Publish(ctx context.Context, msg *pubsub.Message) IPubsubPublishResult
}

type IPubsubPublishResult interface {
	Get(ctx context.Context) (serverID string, err error)
}
