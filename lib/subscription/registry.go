package subscription

type Registry struct {
	Subscriptions []ISubscription
}


func (r *Registry) Add(subscription ISubscription) error {

	r.Subscriptions = append(r.Subscriptions, subscription)
	return nil
}

func (r *Registry) OnEvent(msg Message) {

	for _, subscription := range r.Subscriptions {
		subscription.OnEvent(msg)
	}
}