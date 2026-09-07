package orchestrator

import (
	"log"
	filter "senior_intern_bot/filter"
	poller "senior_intern_bot/poller"
	sender "senior_intern_bot/sender"
	storer "senior_intern_bot/storer"
	"time"
)

type Orchestrator struct {
	poller        poller.Poller
	filter        filter.Filter
	sender        sender.Sender
	storer        storer.Storer
	ticker        *time.Ticker
	isDoneChannel chan bool
}

func New(poller poller.Poller, filter filter.Filter, sender sender.Sender, storer storer.Storer, ticker *time.Ticker) (Orchestrator, error) {
	return Orchestrator{
		poller:        poller,
		filter:        filter,
		sender:        sender,
		storer:        storer,
		ticker:        ticker,
		isDoneChannel: make(chan bool),
	}, nil
}

func (orchestrator *Orchestrator) MainLoop() {
	for {
		select {
		case <-orchestrator.isDoneChannel:
			return
		case <-orchestrator.ticker.C:
			// Poll
			polledData, err := orchestrator.poller.Poll()
			if err != nil {
				log.Println("error polling: ", err)
				continue
			}

			// Filter
			filteredData, err := orchestrator.filter.Filter(polledData)
			if err != nil {
				log.Println("error filtering: ", err)
			}

			// Sender
			err = orchestrator.sender.Send(filteredData)
			if err != nil {
				log.Println("error sending: ", err)
			}
			// Storer (to implement)

		}
	}
}

func (orchestrator *Orchestrator) Start() error {
	go orchestrator.MainLoop()
	return nil
}

func (orchestrator *Orchestrator) Stop() error {
	orchestrator.ticker.Stop()
	orchestrator.isDoneChannel <- true
	return nil
}
