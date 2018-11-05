package signal

import (
	"time"

	"github.com/rtctunnel/rtctunnel/crypt"
	"github.com/mr-tron/base58"
)

type Channel interface {
	Send(key, data string) error
	Recv(key string) (data string, err error)
}

type config struct {
	timeout time.Duration
	period  time.Duration
	channel Channel
}

func defaultConfig() *config {
	return &config{
		period:  time.Second * 30,
		channel: NewOperatorChannel("localhost:9451"),
	}
}

type Option func(cfg *config)

func WithChannel(channel Channel) Option {
	return func(cfg *config) {
		cfg.channel = channel
	}
}

func Send(keypair crypt.KeyPair, peerPublicKey crypt.Key, data []byte, options ...Option) error {
	cfg := defaultConfig()
	for _, o := range options {
		o(cfg)
	}
	encrypted := keypair.Encrypt(peerPublicKey, data)
	address := peerPublicKey.String() + "/" + keypair.Public.String()
	encoded := base58.Encode(encrypted)
	return cfg.channel.Send(address, encoded)
}

func Recv(keypair crypt.KeyPair, peerPublicKey crypt.Key, options ...Option) (data []byte, err error) {
	cfg := defaultConfig()
	for _, o := range options {
		o(cfg)
	}
	address := keypair.Public.String() + "/" + peerPublicKey.String()
	encoded, err := cfg.channel.Recv(address)
	if err != nil {
		return nil, err
	}
	decoded, err := base58.Decode(encoded)
	if err != nil {
		return nil, err
	}
	decrypted, err := keypair.Decrypt(peerPublicKey, decoded)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}
