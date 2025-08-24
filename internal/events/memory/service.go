package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/gentra/decorator-arch-go/internal/eventhandler"
	"github.com/gentra/decorator-arch-go/internal/events"
)

// service implements events.Service interface using in-memory storage
type service struct {
	events        []events.Event
	subscriptions map[string]*events.EventSubscription
	handlers      map[string][]eventhandler.Service
	mu            sync.RWMutex
	config        events.EventConfig
}

// NewService creates a new in-memory event service
func NewService(config events.EventConfig) events.Service {
	if !config.IsValid() {
		config = events.DefaultEventConfig()
	}

	return &service{
		events:        make([]events.Event, 0),
		subscriptions: make(map[string]*events.EventSubscription),
		handlers:      make(map[string][]eventhandler.Service),
		config:        config,
	}
}

// Publish publishes an event
func (s *service) Publish(ctx context.Context, event events.Event) error {
	if !event.IsValid() {
		return events.ErrInvalidEvent
	}

	// Set timestamp if not provided
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Generate ID if not provided
	if event.ID == "" {
		event.ID = uuid.New().String()
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Store the event
	s.events = append(s.events, event)

	// Handle the event asynchronously
	go s.handleEvent(ctx, event)

	return nil
}

// PublishBatch publishes multiple events
func (s *service) PublishBatch(ctx context.Context, eventList []events.Event) error {
	for _, event := range eventList {
		if err := s.Publish(ctx, event); err != nil {
			return fmt.Errorf("failed to publish event %s: %w", event.ID, err)
		}
	}
	return nil
}

// Subscribe subscribes to events by topics
func (s *service) Subscribe(ctx context.Context, topics []string, handler eventhandler.Service) error {
	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}

	subscriptionID := uuid.New().String()

	subscription := &events.EventSubscription{
		ID:        subscriptionID,
		Topics:    topics,
		Handler:   handler,
		CreatedAt: time.Now(),
		Active:    true,
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.subscriptions[subscriptionID] = subscription

	// Register handler for each event type it handles
	for _, eventType := range handler.GetHandledEventTypes() {
		s.handlers[eventType] = append(s.handlers[eventType], handler)
	}

	return nil
}

// Unsubscribe removes a subscription
func (s *service) Unsubscribe(ctx context.Context, subscriptionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	subscription, exists := s.subscriptions[subscriptionID]
	if !exists {
		return fmt.Errorf("subscription %s not found", subscriptionID)
	}

	// Mark subscription as inactive
	subscription.Active = false

	// Remove handlers
	for _, eventType := range subscription.Handler.GetHandledEventTypes() {
		handlers := s.handlers[eventType]
		for i, handler := range handlers {
			if handler == subscription.Handler {
				s.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}

	delete(s.subscriptions, subscriptionID)
	return nil
}

// GetEvents retrieves events based on filters
func (s *service) GetEvents(ctx context.Context, filters events.EventFilters) ([]events.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []events.Event

	for _, event := range s.events {
		if s.matchesFilters(event, filters) {
			result = append(result, event)
		}
	}

	// Apply pagination
	if filters.Offset > 0 && filters.Offset < len(result) {
		result = result[filters.Offset:]
	}

	if filters.Limit > 0 && filters.Limit < len(result) {
		result = result[:filters.Limit]
	}

	return result, nil
}

// GetEventsByAggregate retrieves events for a specific aggregate
func (s *service) GetEventsByAggregate(ctx context.Context, aggregateID string, limit int) ([]events.Event, error) {
	filters := events.EventFilters{
		AggregateID: aggregateID,
		Limit:       limit,
	}

	return s.GetEvents(ctx, filters)
}

// ReplayEvents replays events for an aggregate
func (s *service) ReplayEvents(ctx context.Context, aggregateID string, fromVersion int, handler eventhandler.Service) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, event := range s.events {
		if event.AggregateID == aggregateID && event.Version >= fromVersion {
			if err := handler.Handle(ctx, event); err != nil {
				return fmt.Errorf("failed to replay event %s: %w", event.ID, err)
			}
		}
	}

	return nil
}

// handleEvent processes an event by calling registered handlers
func (s *service) handleEvent(ctx context.Context, event events.Event) {
	s.mu.RLock()
	handlers, exists := s.handlers[event.Type]
	s.mu.RUnlock()

	if !exists {
		return // No handlers for this event type
	}

	for _, handler := range handlers {
		go func(h eventhandler.Service) {
			if err := h.Handle(ctx, event); err != nil {
				// In a real implementation, you might want to log this error
				// or implement retry logic
				fmt.Printf("Error handling event %s: %v\n", event.ID, err)
			}
		}(handler)
	}
}

// matchesFilters checks if an event matches the given filters
func (s *service) matchesFilters(event events.Event, filters events.EventFilters) bool {
	// Check event types
	if len(filters.EventTypes) > 0 {
		found := false
		for _, eventType := range filters.EventTypes {
			if event.Type == eventType {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check aggregate ID
	if filters.AggregateID != "" && event.AggregateID != filters.AggregateID {
		return false
	}

	// Check aggregate types
	if len(filters.AggregateTypes) > 0 {
		found := false
		for _, aggregateType := range filters.AggregateTypes {
			if event.AggregateType == aggregateType {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check time range
	if filters.StartTime != nil && event.Timestamp.Before(*filters.StartTime) {
		return false
	}

	if filters.EndTime != nil && event.Timestamp.After(*filters.EndTime) {
		return false
	}

	// Check user ID
	if filters.UserID != "" && event.Metadata.UserID != filters.UserID {
		return false
	}

	// Check correlation ID
	if filters.CorrelationID != "" && event.Metadata.CorrelationID != filters.CorrelationID {
		return false
	}

	return true
}
