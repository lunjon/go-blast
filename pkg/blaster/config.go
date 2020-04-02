package blaster

import (
    "bytes"
    "fmt"
    "io"
    "log"
    "net/http"
    "net/url"
    "strings"
    "time"
)

const (
    // DefaultRate is the default number of requests per second to send
    DefaultRate = 10
    // MinRate is the minimum rate of a blaster measured in hertz
    MinRate = 1
    // MaxRate is the maximum rate of a blaster measured in hertz
    MaxRate = 100
    // DefaultDuration is the default number of seconds to run the blasters
    DefaultDuration = 60
    // MinDuration is the minimum duration of a blast measured in seconds
    MinDuration = 5
    // MaxDuration is the maximum duration of a blast measured in seconds
    MaxDuration = 900 // 900 seconds = 15 minutes
)

var (
    supportedHTTPMethods = []string{
        http.MethodGet,
        http.MethodDelete,
        http.MethodPost,
    }
)

// Configuration holds the configuration for a blast
type Configuration struct {
    Rate        int
    Duration    time.Duration
    URL         *url.URL
    HTTPMethod  string
    Header      http.Header
    requestBody []byte
    valid bool
}

// NewConfiguration creates a new blaster configuration
// using the given parameters. The values can be changed
// later using the corresponding Set* method.
func NewConfiguration(
    url string,
    method string,
    rate int,
    durationSeconds int,
    header http.Header) (config *Configuration, err error) {
    config = &Configuration{}
    if err = config.SetURL(url); err != nil {
        return
    }
    if err = config.SetMethod(method); err != nil {
        return
    }
    if err = config.SetRate(rate); err != nil {
        return
    }
    if err = config.SetDuration(durationSeconds); err != nil {
        return
    }
    if header == nil {
        header = http.Header{}
    }
    config.Header = header
    config.valid = true
    return
}

// SetRate sets the request rate of each blaster.
func (c *Configuration) SetRate(rate int) error {
    if rate == 0 {
        rate = DefaultRate
    }

    if rate < MinRate || rate > MaxRate {
        return fmt.Errorf("rate must be either zero or a positive integer in the range %d-%d", MinRate, MaxRate)
    }

    c.Rate = rate
    return nil
}

// SetDuration sets the duration in seconds on how to long blast.
func (c *Configuration) SetDuration(duration int) error {
    if duration == 0 {
        duration = DefaultDuration
    }

    if duration < MinDuration || c.Duration > MaxDuration*time.Second {
        return fmt.Errorf("duration must be a positive integer in the range %d-%d", MinDuration, MaxDuration)
    }

    c.Duration = time.Duration(duration) * time.Second
    return nil
}

// SetURL is useful when testing
func (c *Configuration) SetURL(u string) (err error) {
    c.URL, err = url.ParseRequestURI(u)
    return
}

func (c *Configuration) SetMethod(method string) error {
    if method == "" {
        log.Printf("Using default HTTP methid: %s", http.MethodGet)
        c.HTTPMethod = http.MethodGet
        return nil
    }

    method = strings.ToUpper(method)
    for _, m := range supportedHTTPMethods {
        if method == m {
            c.HTTPMethod = m
            return nil
        }
    }

    return fmt.Errorf("unsupported HTTP method: %s", method)
}

// Valid returns true if this configuration has been created
// using NewConfiguration function and validating.
func (c *Configuration) Valid() bool {
    return c.valid
}

// SetRequestBody sets the request body of this configuration.
func (c *Configuration) SetRequestBody(body []byte) {
    c.requestBody = body
}

// BuildRequest returns the corresponding request object that
// that this configuration describes.
func (c *Configuration) BuildRequest() ( req*http.Request, err error) {
    if !c.valid {
        err = fmt.Errorf("invalid configuration, use NewConfiguration to create")
        return
    }

    var body io.Reader
    if c.requestBody != nil {
        body = bytes.NewReader(c.requestBody)
    }

    req, err = http.NewRequest(c.HTTPMethod, c.URL.String(), body)
    if req != nil {
        req.Header = c.Header
    }

    return
}

// UpdateHeader add all entries that do not exist in the configuration
// and override any existing values.
func (c *Configuration) UpdateHeader(header http.Header) {
    for key := range header {
        value := header.Get(key)
        if _, exists := c.Header[key]; exists {
            log.Printf("Overwriting header with new value: %s = %s", key, value)
        } else {
            log.Printf("Adding header: %s = %s", key, value)
        }
        c.Header.Set(key, value)
    }
}
