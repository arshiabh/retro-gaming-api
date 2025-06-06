package events

type SignedUpEvent struct {
	EventType string `json:"event_type"`
	UserID    string `json:"user_id"`
	Timestamp string `json:"timestamp"`
}

type LoginEvent struct {
	EventType string `json:"event_type"`
	UserID    string `json:"user_id"`
	Timestamp string `json:"timestamp"`
}
