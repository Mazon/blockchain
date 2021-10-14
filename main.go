package main

//hashGenesisBlock := "0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"
var (
	cfg *config
)

func main() {
	// Load configuration and parse command line.  This function also
	// initializes logging and configures it accordingly.
	tcfg, _, err := loadConfig()
	if err != nil {
		return
	}
	cfg = tcfg

	// Get a channel that will be closed when a shutdown signal has been
	// triggered either from an OS signal such as SIGINT (Ctrl+C) or from
	// another subsystem such as the RPC server.
	interrupt := interruptListener()

	// Create server and start it.
	server, err := newServer(activeNetParams.Params, interrupt)
	if err != nil {
		// TODO: this logging could do with some beautifying.
		//btcdLog.Errorf("Unable to start server on %v: %v",
		//		cfg.Listeners, err)
		return
	}
	defer func() {
		//btcdLog.Infof("Gracefully shutting down the server...")
		//server.Stop()
		//server.WaitForShutdown()
		//srvrLog.Infof("Server shutdown complete")
	}()
	server.Start()
	//if serverChan != nil {
	//	serverChan <- server
	//}

	// Wait until the interrupt signal is received from an OS signal or
	// shutdown is requested through one of the subsystems such as the RPC
	// server.
	<-interrupt
	return
	//mutex := &sync.Mutex{}

}
