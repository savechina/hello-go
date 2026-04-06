package trace

var events []string

// Record appends an init or demo event.
func Record(event string) {
	events = append(events, event)
}

// Events returns a copy of recorded events.
func Events() []string {
	result := make([]string, len(events))
	copy(result, events)
	return result
}
