package event

import di "flamingo/core/flamingo/dependencyinjection"

type (
	// DefaultRouter is a default event routing implementation
	DefaultRouter struct {
		subscriber map[interface{}][]interface{}
	}
)

// Subscribe adds a subscriber to the list of subscribed subscribers
func (d *DefaultRouter) Subscribe(key, callback interface{}) {
	if d.subscriber == nil {
		d.subscriber = make(map[interface{}][]interface{})
	}
	d.subscriber[key] = append(d.subscriber[key], callback)
}

// Dispatch calls the event's Dispatch method on each subscriber
func (d *DefaultRouter) Dispatch(key interface{}, dispatcher Event) {
	for _, s := range d.subscriber[key] {
		dispatcher.Dispatch(s)
	}
}

// CompilerPass retrieves a list of event subscribers
func (d *DefaultRouter) CompilerPass(c *di.Container) {
	for _, s := range c.GetTagged("event.subscriber") {
		for _, event := range s.Value.(Subscriber).Events() {
			d.Subscribe(event, s.Value)
		}
	}
}
