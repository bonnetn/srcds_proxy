package proxy

import (
	"log"
	"srcds_proxy/proxy/srcds"
	"srcds_proxy/proxy/handler"
	"srcds_proxy/proxy/conntrack"
	"srcds_proxy/proxy/config"
	"net"
	"sync"
)

func Launch() error {
	log.Println("INFO: Listening on ", config.ListenFullAddr)
	conn, err := srcds.Listen(config.ListenFullAddr)
	if err != nil {
		log.Println("ERROR: Could not listen: ", err)
		return err
	}

	var (
		connectionTable = conntrack.NewConnectionTable()
		h               = handler.NewRequestProcessorHandler(connectionTable)
		done            = make(chan struct{})
		running         = true
	)
	log.Println("INFO: Starting proxy...")
	for err := range doServe(done, h, conn) {
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

func doServe(done <-chan struct{}, handler srcds.Handler, conn *net.UDPConn) <-chan error {
	var (
		resultChan = make(chan error)
		wg         = sync.WaitGroup{}
	)

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
			resultChan <- srcds.Serve(done, *conn, handler)
			wg.Done()
		}()
	}

	return resultChan
}
