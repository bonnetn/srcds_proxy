package utils

import "testing"

func TestIsDone(t *testing.T) {
	var closedChan = make(chan DoneEvent)
	close(closedChan)

	type args struct {
		done <-chan DoneEvent
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "channel_closed",
			args: args{
				done: closedChan,
			},
			want: true,
		},
		{
			name: "channel_open",
			args: args{
				done: make(chan DoneEvent),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDone(tt.args.done); got != tt.want {
				t.Errorf("IsDone() = %v, want %v", got, tt.want)
			}
		})
	}
}
