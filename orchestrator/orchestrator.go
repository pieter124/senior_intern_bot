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
	ticker        time.Ticker
	isDoneChannel chan bool
}

func New(poller poller.Poller, filter filter.Filter, sender sender.Sender, storer storer.Storer, ticker time.Ticker) (Orchestrator, error) {
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
	// We need to implement:
	// - polling (every 20 mins -> fetch from job boards),
	// - filtering (filter & classify into ACCEPT, REVIEW, and REJECT),
	// - sending (send message(s) to the internship channel(s)),
	// - storing (our db layer, stores results, we will read from this during the filtering stage),

	// So, poll -> filter -> send -> store.
	// We can have an orchestrator that handles all of this (inject poller, filter, sender, and storer services).
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
			senderData, err := orchestrator.sender.Send(filteredData)

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
