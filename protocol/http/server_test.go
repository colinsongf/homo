package http

import (
	"github.com/aiicy/aiicy/logger"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIHttp(t *testing.T) {
	s, err := NewServer(ServerInfo{Address: "tcp://127.0.0.1:0"}, func(u, p string) bool {
		return u == "u" && p == "p"
	}, logger.S)
	assert.NoError(t, err)
	defer s.Close()

	s.Handle(func(params Params, reqBody []byte) ([]byte, error) {
		assert.Equal(t, params["arg"], "1")
		assert.Len(t, reqBody, 0)
		return []byte{'a', 'b', 'c'}, nil
	}, "GET", "/test/get", "arg", "{arg}")
	s.Handle(func(params Params, reqBody []byte) ([]byte, error) {
		assert.Equal(t, params["arg"], "2")
		assert.Equal(t, reqBody, []byte{'a', 'b', 'c'})
		return reqBody[:2], nil
	}, "PUT", "/test/put", "arg", "{arg}")
	s.Handle(func(params Params, reqBody []byte) ([]byte, error) {
		assert.Equal(t, params["arg"], "3")
		assert.Equal(t, reqBody, []byte{'a', 'b', 'c'})
		return reqBody[:1], nil
	}, "POST", "/test/post", "arg", "{arg}")

	err = s.Start()
	assert.NoError(t, err)

	addr := "tcp://" + s.addr
	c, err := NewClient(ClientInfo{Address: addr})
	assert.NoError(t, err)

	resBody, err := c.Get("/test/get?arg=%d", 1)
	assert.EqualError(t, err, "[401] account unauthorized")
	assert.Nil(t, resBody)
	resBody, err = c.Put(nil, "/test/put?arg=%d", 1)
	assert.EqualError(t, err, "[401] account unauthorized")
	assert.Nil(t, resBody)
	resBody, err = c.Post(nil, "/test/post?arg=%d", 1)
	assert.EqualError(t, err, "[401] account unauthorized")
	assert.Nil(t, resBody)

	c, err = NewClient(ClientInfo{Address: addr, Username: "u", Password: "p"})
	assert.NoError(t, err)

	resBody, err = c.Get("/test/get?arg=%d", 1)
	assert.NoError(t, err)
	assert.Equal(t, "abc", string(resBody))
	resBody, err = c.Put([]byte{'a', 'b', 'c'}, "/test/put?arg=%d", 2)
	assert.NoError(t, err)
	assert.Equal(t, "ab", string(resBody))
	resBody, err = c.Post([]byte{'a', 'b', 'c'}, "/test/post?arg=%d", 3)
	assert.NoError(t, err)
	assert.Equal(t, "a", string(resBody))
}
