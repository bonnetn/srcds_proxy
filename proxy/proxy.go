package proxy

import (
	"log"
	"srcds_proxy/proxy/srcds"
	"srcds_proxy/proxy/handler"
	"srcds_proxy/proxy/config"
	"sync"
)

func Launch() error {
	var (
		done    = make(chan struct{})
		h       = handler.NewRequestProcessorHandler(done)
		running = true
	)
	log.Println("INFO: Starting proxy...")
	for err := range doServe(done, h) {
		if err != nil {
			log.Print("ERROR: ", err)
		} else {
			log.Print("INFO: Worker exited gracefully.")
		}

		// Kill all workers and close all connections if one worker crashes.
		if running {
			log.Print("WARN: A worker crashed, killing all other workers.")
			close(done)
			running = false
		}
	}
	log.Println("INFO: Proxy stopped.")

	return nil
}

func doServe(done <-chan struct{}, handler srcds.Handler) <-chan error {
	var (
		resultChan = make(chan error)
		wg         = sync.WaitGroup{}
	)

	conn, err := srcds.Listen(config.ListenFullAddr)
	if err != nil {
		log.Println("ERROR: Could not listen: ", err)
		retChan := make(chan error)
		go func() { retChan <- err }()
		return retChan
	}

	// Close the result chan when all the workers have stopped.
	wg.Add(config.WorkerCount)
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Ensure workers terminate when done event received.
	go func() {
		<-done
		conn.Close()
	}()

	for i := 0; i < config.WorkerCount; i++ {
		go func() {
			resultChan <- srcds.Serve(done, *conn, handler, 0)
			wg.Done()
		}()
	}

	return resultChan
}
