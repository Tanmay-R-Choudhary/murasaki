package main

type ServerState struct {
	isShuttingDown bool
}

func (s *ServerState) getShutdown() bool {
	return s.isShuttingDown
}

func (s *ServerState) setShutdown(shutdown bool) {
	s.isShuttingDown = shutdown
}

func NewServerState() *ServerState {
	return &ServerState{
		isShuttingDown: false,
	}
}
