package commonentities

type Event[T any] struct {
	Name string
	Data []T
}
type Observer interface {
	receiveEvent(event Event[any])
}
type  ObserverBucket struct{
	events map[string][]Observer
} 
func NewEventBucket() *ObserverBucket{
	return &ObserverBucket{
		events: map[string][]Observer{
			"replicate_events": []Observer{},
			"payment:events": []Observer{},
		},
	}
}

func (eb *ObserverBucket) AddObserver(observer Observer, eventname string){
	eb.events[eventname] = append(eb.events[eventname], observer)
}

func (eb *ObserverBucket) RemoveObserver(observer Observer, eventname string){
	for i, obs := range eb.events[eventname] {
		if obs == observer {
			eb.events[eventname] = append(eb.events[eventname][:i], eb.events[eventname][i+1:]...)
			return
		}
	}
}

func (eb *ObserverBucket) NotifyObservers(event Event[any], eventname string){
	for _, observer := range eb.events[eventname] {
		observer.receiveEvent(event)
	}
}
