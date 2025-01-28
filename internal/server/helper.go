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

func calculateURLRange(podIndex, totalPods, totalURLs int) (int, int) {
	urlsPerPod := totalURLs / totalPods
	remainder := totalURLs % totalPods

	start := podIndex * urlsPerPod
	if podIndex < remainder {
		start += podIndex
	} else {
		start += remainder
	}

	end := start + urlsPerPod
	if podIndex < remainder {
		end += 1
	}

	return start, end
}
