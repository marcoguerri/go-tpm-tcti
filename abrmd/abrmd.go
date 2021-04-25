package abrmd

import (
	"fmt"

	"os"

	"github.com/godbus/dbus/v5"
)

const dstService = "com.intel.tss2.Tabrmd"
const dstObject = "/com/intel/tss2/Tabrmd/Tcti"

type connState int

const (
	// WaitReceive represents a state where the connection expects a Read call
	WaitReceive connState = iota
	// WaitTransmit represents a state where the connection expects a Write call
	WaitTransmit
)

const (
	createConnection = "com.intel.tss2.TctiTabrmd.CreateConnection"
)

// Broker is a TCTI transport towards tpm2-abrmd daemon
type Broker struct {
	conn  *os.File
	state connState
}

// Read reads a buffer from broker daemon over b.conn
func (b *Broker) Read(buff []byte) (int, error) {
	if b.state != WaitReceive {
		return 0, fmt.Errorf("conn not in WaitReceive state: %v", b.state)
	}
	n, err := b.conn.Read(buff)
	if err == nil {
		b.state = WaitTransmit
	}
	return n, err
}

// Write writes a buffer to broker daemon over b.conn
func (b *Broker) Write(buff []byte) (int, error) {
	if b.state != WaitTransmit {
		return 0, fmt.Errorf("conn not in WaitTransmit state: %v", b.state)
	}
	n, err := b.conn.Write(buff)
	if err == nil {
		b.state = WaitReceive
	}
	return n, err
}

// Close closes connection to broker daemon
func (b *Broker) Close() error {
	return b.conn.Close()
}

// NewBroker creates a broker object and establishes a connection to broker daemon via dbus
func NewBroker() (*Broker, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, fmt.Errorf("could not open connection with session bus: %v", err)
	}
	fds := conn.SupportsUnixFDs()
	if !fds {
		return nil, fmt.Errorf("connection doesn't support file descriptors")
	}

	obj := conn.Object(dstService, dstObject)
	call := obj.Call(createConnection, 0)
	if call.Err != nil {
		return nil, fmt.Errorf("could not create connection towards %s: %v", dstService, call.Err)
	}
	if len(call.Body) < 2 {
		return nil, fmt.Errorf("expected at least (fds, id) from %v call, got %d return values", createConnection, len(call.Body))
	}
	fd, ok := call.Body[0].([]dbus.UnixFD)
	if !ok {
		return nil, fmt.Errorf("expected []dbus.UnixFD as first return value for %v call, got %T", createConnection, call.Body[0])
	}
	if len(fd) != 1 {
		return nil, fmt.Errorf("expected one file descriptor from %v, got %d", createConnection, len(fd))
	}
	fdConn := fd[0]
	file := os.NewFile(uintptr(fdConn), "tpm2-abrmd")
	if file == nil {
		return nil, fmt.Errorf("could not create File object from file descriptor %v", fdConn)
	}
	return &Broker{conn: file, state: WaitTransmit}, nil
}
