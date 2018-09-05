package worker

const MAX_DATAGRAM_SIZE = 4096

type Worker interface {
	Run() error
	Join()
}
