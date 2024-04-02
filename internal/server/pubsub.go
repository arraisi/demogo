package server

import (
	"demogo/config"
	ipubsub "demogo/pkg/pubsub"
)

type ServiceSubscription interface {
	Subscribe()
}

type subGateway struct {
	Pubsub ipubsub.IPubsubClient
	Config *config.Config
}

func NewSubGateway(
	pubsub ipubsub.IPubsubClient,
	config *config.Config,
) ServiceSubscription {
	s := &subGateway{
		Pubsub: pubsub,
		Config: config,
	}

	return s
}

func (s *subGateway) Subscribe() {
	// Add all the subscription here
	// ex. go s.consumeOrderCompleted()
}
