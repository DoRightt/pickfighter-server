package service

import (
	"net"
	"testing"
	"time"

	grpchandler "pickfighter.com/fighters/internal/handler/grpc"
	"pickfighter.com/fighters/pkg/version"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	apiService := New()

	assert.Equal(t, version.Name, apiService.ServiceName)

	assert.NotNil(t, apiService.Server)
}

func TestInit(t *testing.T) {
	apiService := New()

	handler := &grpchandler.Handler{}

	apiService.Init(handler)

	assert.Equal(t, handler, apiService.Handler)
	assert.Equal(t, handler, apiService.Handler)

	serviceInfo := apiService.Server.GetServiceInfo()
	_, ok := serviceInfo["FightersService"]
	assert.True(t, ok, "FightersService should be registered")
}

func TestRun(t *testing.T) {
	viper.Set("http.port", "9090")
	viper.Set("http.addr", "localhost:9090")

	apiService := New()

	go func() {
		err := apiService.Run()
		if err != nil {
			return
		}
	}()

	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("tcp", "localhost:9090")
	assert.NoError(t, err, "Server should be run on port 9090")
	if conn != nil {
		conn.Close()
	}

	apiService.Server.Stop()
}
