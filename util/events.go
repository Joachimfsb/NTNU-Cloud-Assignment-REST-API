package util

import (
	"fmt"
	"strings"
)

// Events are available notification events that a user can register to
const (
	EVENT_ALL      = ""
	EVENT_REGISTER = "REGISTER"
	EVENT_CHANGE   = "CHANGE"
	EVENT_DELETE   = "DELETE"
	EVENT_INVOKE   = "INVOKE"
)

// ValidateEvents ensures that the event type is valid.
// To learn more about available events, please look at 'EVENT_*' in assignment2.util.events
func ValidateEvents(event string) error {
	e := strings.ToUpper(event)
	if e != EVENT_REGISTER &&
		e != EVENT_CHANGE &&
		e != EVENT_DELETE &&
		e != EVENT_INVOKE &&
		e != EVENT_ALL {
		return fmt.Errorf("the event must be empty, 'REGISTER', 'CHANGE', 'DELETE or 'INVOKE'")
	} else {
		return nil
	}
}
