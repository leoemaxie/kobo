package account

import (
	"fmt"
)

type Event string

const (
	EventProvisionSuccess Event = "provision_success"
	EventProvisionFail    Event = "provision_fail"
	EventLimitExceeded    Event = "limit_exceeded"
	EventLimitResolved    Event = "limit_resolved"
	EventCloseInitiated   Event = "close_initiated"
	EventZeroBalance      Event = "zero_balance"
	EventReopenInitiated  Event = "reopen_initiated"
)

// ValidTransition evaluates whether a transition from the current state triggered
// by the given event is valid. It returns the new state and an error if invalid.
func ValidTransition(current State, event Event) (State, error) {
	switch current {
	case StatePending:
		switch event {
		case EventProvisionSuccess:
			return StateActive, nil
		case EventProvisionFail:
			return StateFailed, nil
		}
	case StateActive:
		switch event {
		case EventLimitExceeded:
			return StateLimited, nil
		case EventCloseInitiated:
			return StateClosing, nil
		}
	case StateLimited:
		switch event {
		case EventLimitResolved:
			return StateActive, nil
		case EventCloseInitiated:
			return StateClosing, nil
		}
	case StateClosing:
		switch event {
		case EventZeroBalance:
			return StateClosed, nil
		}
	case StateClosed:
		switch event {
		case EventReopenInitiated:
			// Reopening requires reprovisioning, but goes to active immediately or pending?
			// The instructions say "Transitions CLOSED -> ACTIVE".
			// We will trigger reprovisioning and once successful, set it to ACTIVE.
			// Alternatively we could transition to PENDING and let it flow. But let's stick to the spec.
			// Actually, let's treat EventReopenInitiated as transitioning to PENDING so it can be reprovisioned.
			// But openapi says: Transitions CLOSED -> ACTIVE against the same identity ID.
			// Let's return StateActive and the service will synchronously reprovision it before persisting the state.
			return StateActive, nil
		}
	}
	return current, fmt.Errorf("invalid transition from state %q via event %q", current, event)
}
