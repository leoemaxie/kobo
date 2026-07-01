package account

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidTransitions(t *testing.T) {
	tests := []struct {
		name          string
		currentState  State
		event         Event
		expectedState State
		expectError   bool
	}{
		// Valid transitions
		{"Pending -> Active", StatePending, EventProvisionSuccess, StateActive, false},
		{"Pending -> Failed", StatePending, EventProvisionFail, StateFailed, false},
		{"Active -> Limited", StateActive, EventLimitExceeded, StateLimited, false},
		{"Active -> Closing", StateActive, EventCloseInitiated, StateClosing, false},
		{"Limited -> Active", StateLimited, EventLimitResolved, StateActive, false},
		{"Limited -> Closing", StateLimited, EventCloseInitiated, StateClosing, false},
		{"Closing -> Closed", StateClosing, EventZeroBalance, StateClosed, false},
		{"Closed -> Active (Reopen)", StateClosed, EventReopenInitiated, StateActive, false},

		// Invalid transitions
		{"Pending -> Closed", StatePending, EventZeroBalance, StatePending, true},
		{"Active -> Active", StateActive, EventProvisionSuccess, StateActive, true},
		{"Closing -> Active", StateClosing, EventLimitResolved, StateClosing, true},
		{"Closed -> Pending", StateClosed, EventProvisionSuccess, StateClosed, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newState, err := ValidTransition(tt.currentState, tt.event)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedState, newState)
			}
		})
	}
}
