package handlers

import coreevents "github.com/geevheit/intelligence360/backend/api/internal/events"

type Registry struct {
	handlers []coreevents.EventHandler
}

func NewRegistry() *Registry {
	return &Registry{handlers: []coreevents.EventHandler{}}
}

func (r *Registry) Register(handler coreevents.EventHandler) {
	r.handlers = append(r.handlers, handler)
}

func (r *Registry) All() []coreevents.EventHandler {
	out := make([]coreevents.EventHandler, len(r.handlers))
	copy(out, r.handlers)
	return out
}

func (r *Registry) SubscribeAll(bus coreevents.EventBus) error {
	for _, handler := range r.handlers {
		if err := bus.Subscribe(handler); err != nil {
			return err
		}
	}
	return nil
}
