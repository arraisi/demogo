package pubsub

import (
	"context"
	"errors"
	"github.com/arraisi/demogo/config"
	"time"

	"cloud.google.com/go/pubsub"
	vkit "cloud.google.com/go/pubsub/apiv1"
	gax "github.com/googleapis/gax-go/v2"
	"google.golang.org/grpc/codes"
)

type Subscription struct {
	ID      string
	SubInst *pubsub.Subscription
}

type Topic struct {
	ID        string
	TopicInst *pubsub.Topic
}

type Client struct {
	Topics        map[string]*Topic
	Subscriptions map[string]*Subscription
	Client        *pubsub.Client
}

func NewPubsubClient() *Client {
	return &Client{
		Topics:        make(map[string]*Topic),
		Subscriptions: make(map[string]*Subscription),
	}
}

func (c *Client) NewClient(ctx context.Context, projectID string) (IPubsubClient, error) {
	var err error
	c.Client, err = pubsub.NewClient(ctx, projectID)
	return c, err

}

func (c *Client) Subscription(ctx context.Context, subID string) (IPubsubSubscription, error) {
	s := c.Client.Subscription(subID)
	s.ReceiveSettings.NumGoroutines = 1
	ps := &Subscription{
		ID:      subID,
		SubInst: s,
	}
	c.Subscriptions[subID] = ps
	return c.Subscriptions[subID], nil
}

func (c *Client) GetSubscription(subID string) (IPubsubSubscription, error) {
	if c.Subscriptions[subID] == nil {
		return nil, errors.New("subs not found")
	}
	return c.Subscriptions[subID], nil
}

func (c *Client) Topic(ctx context.Context, topicID string) IPubsubTopic {
	t := c.Client.Topic(topicID)
	pt := &Topic{
		ID:        topicID,
		TopicInst: t,
	}
	c.Topics[topicID] = pt
	return c.Topics[topicID]
}

func (c *Client) Close() error {
	return c.Client.Close()
}

func (s *Subscription) Receive(ctx context.Context, f func(context.Context, *pubsub.Message)) error {
	return s.SubInst.Receive(ctx, f)
}

func (t *Topic) Publish(ctx context.Context, msg *pubsub.Message) IPubsubPublishResult {
	return &PublishResult{
		PRInst: t.TopicInst.Publish(ctx, msg),
	}
}

func (c *Client) NewClientWithRetry(ctx context.Context, projectID string, conf *config.Config) (IPubsubClient, error) {
	config := &pubsub.ClientConfig{
		PublisherCallOptions: &vkit.PublisherCallOptions{
			Publish: []gax.CallOption{
				gax.WithRetry(func() gax.Retryer {
					return gax.OnCodes([]codes.Code{
						codes.Aborted,
						codes.Canceled,
						codes.Internal,
						codes.ResourceExhausted,
						codes.Unknown,
						codes.Unavailable,
						codes.DeadlineExceeded,
					}, gax.Backoff{
						Initial:    conf.PubSub.BackoffInitial * time.Millisecond,
						Max:        conf.PubSub.BackoffMax * time.Second,
						Multiplier: conf.PubSub.BackoffMultiplier,
					})
				}),
			},
		},
	}

	cc, err := pubsub.NewClientWithConfig(ctx, projectID, config)
	if err != nil {
		return c, err
	}

	c.Client = cc
	return c, nil
}

type PublishResult struct {
	PRInst *pubsub.PublishResult
}

func (pr *PublishResult) Get(ctx context.Context) (serverID string, err error) {
	serverID, err = pr.PRInst.Get(ctx)
	return serverID, err
}
