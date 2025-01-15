package server

import (
	"centris-api/internal/repository"
	"context"
)

func GetCompleteBroker(s *Server, ctx context.Context, broker repository.Broker) (CompleteBroker, error) {
	var completeBroker CompleteBroker

	completeBroker.Broker = broker

	broker_phones, err := s.queries.GetAllBrokerPhonesByBrokerId(ctx, broker.ID)
	if err != nil {
		return completeBroker, err
	}
	completeBroker.Broker_Phones = broker_phones

	broker_links, err := s.queries.GetAllBrokerLinksByBrokerId(ctx, broker.ID)
	if err != nil {
		return completeBroker, err
	}
	completeBroker.Broker_Links = broker_links

	return completeBroker, nil
}
