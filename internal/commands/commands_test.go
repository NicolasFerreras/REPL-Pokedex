package commands

import (
	"errors"
	"testing"
)

func TestUserInput(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantErr   bool
		errMsg    string
	}{
		{
			name:    "empty input returns error",
			input:   "",
			wantErr: true,
			errMsg:  "no command entered",
		},
		{
			name:    "invalid command returns error",
			input:   "poke",
			wantErr: true,
			errMsg:  "unknown command: poke",
		},
		{
			name:    "another invalid command",
			input:   "catch",
			wantErr: true,
			errMsg:  "unknown command: catch",
		},
		{
			name:    "help command succeeds",
			input:   "help",
			wantErr: false,
		},
		{
			name:    "exit command returns ErrExit",
			input:   "exit",
			wantErr: true,
			errMsg:  "", // We check for ErrExit specifically
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UserInput(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("UserInput(%q) expected error, got nil", tt.input)
					return
				}
				// Special case: exit command returns ErrExit
				if tt.input == "exit" {
					if !errors.Is(err, ErrExit) {
						t.Errorf("UserInput(%q) expected ErrExit, got %v", tt.input, err)
					}
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("UserInput(%q) error = %q, want containing %q", tt.input, err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("UserInput(%q) unexpected error: %v", tt.input, err)
				}
			}
		})
	}
}

func TestCommandsMap(t *testing.T) {
	// Verify that required commands exist
	requiredCommands := []string{"help", "exit"}
	for _, cmd := range requiredCommands {
		if _, exists := Commands[cmd]; !exists {
			t.Errorf("required command %q not found in Commands map", cmd)
		}
	}

	// Verify that commands have non-nil callbacks
	for name, cmd := range Commands {
		if cmd.Callback == nil {
			t.Errorf("command %q has nil Callback", name)
		}
		if cmd.Name == "" {
			t.Errorf("command %q has empty Name", name)
		}
		if cmd.Description == "" {
			t.Errorf("command %q has empty Description", name)
		}
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
