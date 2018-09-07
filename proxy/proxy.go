package proxy

import (
	"log"
	"srcds_proxy/proxy/srcds"
	"srcds_proxy/proxy/handler"
	"srcds_proxy/proxy/conntrack"
	"srcds_proxy/proxy/config"
	"net"
	"runtime"
	"sync"
	"context"
)

func doServe(ctx context.Context, handler srcds.Handler, conn *net.UDPConn) <-chan error {
	var (
		resultChan = make(chan error)
		wg         = sync.WaitGroup{}
		numCPU     = runtime.NumCPU()
	)

	wg.Add(numCPU)
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for i := 0; i < numCPU; i++ {
		go func() {
			resultChan <- srcds.Serve(ctx, *conn, handler)
			wg.Done()
		}()
	}

	return resultChan
}

func Launch() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		connectionTable = conntrack.NewConnectionTable()
		h               = handler.NewRequestProcessorHandler(connectionTable)
	)

	log.Println("Listening on ", config.ListenFullAddr)
	conn, err := srcds.Listen(ctx, config.ListenFullAddr)
	if err != nil {
		log.Println("Could not listen: ", err)
		return err
	}

	log.Println("Starting proxy...")
	for err := range doServe(ctx, h, conn) {
		if err != nil {
			log.Print("ERROR: ", err)
		}
		cancel() // Kill all workers and close all connections if one worker crashes.
	}
	log.Println("Proxy stopped.")

	return nil
}
