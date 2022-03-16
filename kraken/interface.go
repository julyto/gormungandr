package kraken

import "github.com/CanalTP/gonavitia/pbnavitia"

type KrakenTimeout struct {
	message string
}

func NewKrakenTimeout(message string) *KrakenTimeout {
	return &KrakenTimeout{
		message: message,
	}
}

func (e *KrakenTimeout) Error() string {
	return e.message
}

// Low level interface to use kraken.
type Kraken interface {
	// Send a generic request to kraken and return the response
	Call(request *pbnavitia.Request) (*pbnavitia.Response, error)
	UpdateKrakenZMQ(addr string)
}
