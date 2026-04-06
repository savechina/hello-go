package pointers

import "testing"

func TestRenameWithPointer(t *testing.T) {
	tests := []struct {
		name    string
		initial string
		next    string
		want    string
		wantOK  bool
		useNil  bool
	}{
		{name: "rename value", initial: "old", next: "new", want: "new", wantOK: true},
		{name: "rename nil", next: "new", want: "", wantOK: false, useNil: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.useNil {
				if got := renameWithPointer(nil, tt.next); got != tt.wantOK {
					t.Fatalf("renameWithPointer(nil) = %t, want %t", got, tt.wantOK)
				}
				return
			}

			value := tt.initial
			if got := renameWithPointer(&value, tt.next); got != tt.wantOK || value != tt.want {
				t.Fatalf("renameWithPointer() ok=%t value=%q, want ok=%t value=%q", got, value, tt.wantOK, tt.want)
			}
		})
	}
}

func TestSwapValues(t *testing.T) {
	tests := []struct {
		name      string
		left      int
		right     int
		wantLeft  int
		wantRight int
		wantOK    bool
		useNil    bool
	}{
		{name: "swap works", left: 1, right: 2, wantLeft: 2, wantRight: 1, wantOK: true},
		{name: "swap nil", wantOK: false, useNil: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.useNil {
				right := 2
				if got := swapValues(nil, &right); got != tt.wantOK {
					t.Fatalf("swapValues(nil, &right) = %t, want %t", got, tt.wantOK)
				}
				return
			}

			left := tt.left
			right := tt.right
			if got := swapValues(&left, &right); got != tt.wantOK || left != tt.wantLeft || right != tt.wantRight {
				t.Fatalf("swapValues() ok=%t left=%d right=%d, want ok=%t left=%d right=%d", got, left, right, tt.wantOK, tt.wantLeft, tt.wantRight)
			}
		})
	}
}

func TestWalletBalanceAndNilSafety(t *testing.T) {
	tests := []struct {
		name   string
		start  int
		add    int
		want   int
		useNil bool
	}{
		{name: "deposit updates balance", start: 20, add: 5, want: 25},
		{name: "nil wallet safe", want: 0, useNil: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.useNil {
				var item *wallet
				item.Deposit(10)
				if got := item.Balance(); got != tt.want {
					t.Fatalf("nil wallet Balance() = %d, want %d", got, tt.want)
				}
				return
			}

			item := &wallet{balance: tt.start}
			item.Deposit(tt.add)
			if got := item.Balance(); got != tt.want {
				t.Fatalf("Balance() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestSafeLearnerName(t *testing.T) {
	tests := []struct {
		name string
		item *learner
		want string
	}{
		{name: "real learner", item: &learner{name: "gopher"}, want: "gopher"},
		{name: "nil learner", item: nil, want: "nil learner"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := safeLearnerName(tt.item); got != tt.want {
				t.Fatalf("safeLearnerName() = %q, want %q", got, tt.want)
			}
		})
	}
}
