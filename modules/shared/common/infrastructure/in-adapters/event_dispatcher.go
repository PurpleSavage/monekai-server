package commoninadapters

import (
	commonports "github.com/PurpleSavage/monekai-server/modules/shared/common/application/ports"
	commonentities "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/entities"
)


type ObserverBucket struct {
	events map[string][]commonports.Observer
}

func NewObserverBucket() commonports.ObserverBucketPort {
	return &ObserverBucket{
		events: map[string][]commonports.Observer{
			"sample_event":  {},
			"payment_event": {},
		},
	}
}

func (eb *ObserverBucket) AddObserver(observer commonports.Observer, eventname string) {
	eb.events[eventname] = append(eb.events[eventname], observer)
}

func (eb *ObserverBucket) RemoveObserver(observer commonports.Observer, eventname string) {
	for i, obs := range eb.events[eventname] {
		if obs == observer {
			eb.events[eventname] = append(eb.events[eventname][:i], eb.events[eventname][i+1:]...)
			return
		}
	}
}

func (eb *ObserverBucket) NotifyObservers(event commonentities.Event, eventname string) {
	for _, observer := range eb.events[eventname] {
		observer.ReceiveEvent(event)
	}
}