package blastertest

import (
    "github.com/lunjon/go-blast/pkg/blaster"
    "net/http"
    "testing"
)

func TestNewConfiguration(t *testing.T) {
    tests := []struct {
        name       string
        url             string
        method          string
        rate            int
        durationSeconds int
    }{
        {"default values", "http://localhost", "", 0, 0},
        {"localhost, post", "http://localhost", "post", 0, 0},
        {"rate = 10", "http://localhost", "post", 10, 0},
        {"duration = 10", "http://localhost", "delete", 0, 10},
        {"https url", "https://google.com", "get", 0, 0},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            config, err := blaster.NewConfiguration(tt.url, tt.method, tt.rate, tt.durationSeconds, http.Header{})
            if err != nil {
                t.Errorf("NewConfiguration() error = %v", err)
                return
            }

            if config == nil {
                t.Errorf("Expected config to be not nil")
                return
            }

            if !config.Valid() {
                t.Errorf("Expected config.Valid() to return true")
                return
            }

            _, err = config.BuildRequest()
            if err != nil {
                t.Errorf("BuildRequest() error = %v", err)
                return
            }
        })
    }
}

func TestNewConfigurationInvalid(t *testing.T) {
    tests := []struct {
        name       string
        url             string
        method          string
        rate            int
        durationSeconds int
    }{
        {"whitespace URL", "   ", "get", 0, 0},
        {"missing protocol", "localhost", "get", 0, 0},
        {"invalid method", "http://localhost", "lol", 0, 0},
        {"negative rate", "http://localhost", "get", -1, 0},
        {"negative duration", "http://localhost", "get", 0, -1},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := blaster.NewConfiguration(tt.url, tt.method, tt.rate, tt.durationSeconds, http.Header{})
            if err == nil {
                t.Errorf("NewConfiguration() error = %v", err)
                return
            }
        })
    }
}
