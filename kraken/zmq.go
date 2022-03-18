package kraken

import (
	"fmt"
	"time"

	"github.com/CanalTP/gonavitia/pbnavitia"
	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sony/gobreaker"
)

type KrakenZMQ struct {
	Name    string
	Addr    string
	Timeout time.Duration

	cb         *gobreaker.CircuitBreaker
	observer   krakenObserver
	socketPool *pool
}

func NewKrakenZMQ(name, addr string, timeout time.Duration) Kraken {
	kraken := &KrakenZMQ{
		Name:       name,
		Timeout:    timeout,
		Addr:       addr,
		socketPool: newPool(addr, 100),
	}
	var st gobreaker.Settings
	st.Name = name
	kraken.cb = gobreaker.NewCircuitBreaker(st)
	return kraken
}

func (k *KrakenZMQ) UpdateKrakenZMQ(addr string) {
	if addr == k.Addr {
		return
	}
	close(k.socketPool.pool)
	k.socketPool = newPool(addr, 100)
	k.Addr = addr
}

func (k *KrakenZMQ) call(request []byte) ([]byte, error) {
	var err error
	var requester *zmq.Socket
	requester, err = k.socketPool.Borrow()
	if err != nil {
		return nil, err
	}
	if _, err = requester.SendBytes(request, 0); err != nil {
		closeSocket(requester)
		return nil, errors.Wrap(err, "error while sending")
	}
	poller := zmq.NewPoller()
	poller.Add(requester, zmq.POLLIN)
	p, err := poller.Poll(k.Timeout)
	if err != nil {
		closeSocket(requester)
		return nil, errors.Wrap(err, "error during polling")
	}
	if len(p) < 1 {
		closeSocket(requester)
		return nil, NewKrakenTimeout(fmt.Sprintf("kraken %s timeout", k.Name))
	}
	rawResp, err := p[0].Socket.RecvBytes(0)
	if err != nil {
		closeSocket(requester)
		return nil, errors.Wrap(err, "error while receiving response")
	}
	k.socketPool.Return(requester)
	return rawResp, nil
}

func (k *KrakenZMQ) Call(request *pbnavitia.Request) (*pbnavitia.Response, error) {
	data, err := proto.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "error while marshalling")
	}
	rawResponse, err := k.cb.Execute(func() (interface{}, error) {
		o := k.observer.StartRequest(request.RequestedApi.String())
		defer o.Finish()
		return k.call(data)
	})
	if err != nil {
		k.observer.OnError(request.RequestedApi.String(), err)
		return nil, err
	}

	response := &pbnavitia.Response{}
	if err = proto.Unmarshal(rawResponse.([]byte), response); err != nil {
		return nil, errors.Wrap(err, "error while unmarshalling response")
	}
	return response, nil
}

func closeSocket(s *zmq.Socket) {
	if err := s.Close(); err != nil {
		logrus.Warnf("error while closing the socket %s", err)
	}
}

func NewSocket(addr string) (*zmq.Socket, error) {
	socket, err := zmq.NewSocket(zmq.REQ)
	if err != nil {
		return nil, errors.Wrap(err, "error while creating ZMQ socket")
	}
	if err = socket.Connect(addr); err != nil {
		return nil, errors.Wrap(err, "error while connecting")
	}
	if err = socket.SetLinger(0); err != nil {
		return nil, errors.Wrap(err, "error while set linger")
	}
	return socket, nil
}

//hold established zocket to kraken instance
type pool struct {
	pool chan *zmq.Socket
	Addr string
}

// NewPool creates a new pool of socket.
func newPool(addr string, max int) *pool {
	return &pool{
		Addr: addr,
		pool: make(chan *zmq.Socket, max),
	}
}

// Borrow a Socket from the pool.
func (p *pool) Borrow() (*zmq.Socket, error) {
	var s *zmq.Socket
	select {
	case s = <-p.pool:
		return s, nil
	default:
		return NewSocket(p.Addr)
	}
}

// Return returns a socket to the pool.
func (p *pool) Return(s *zmq.Socket) {
	select {
	case p.pool <- s:
	default:
		//pool is full, closing socket
		closeSocket(s)
	}
}
