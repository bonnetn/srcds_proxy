package utils

type DoneEvent struct{}

func IsDone(done <-chan DoneEvent) bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}
