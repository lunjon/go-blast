package blastertest

import (
    "fmt"
    "github.com/lunjon/go-blast/pkg/blaster"
    "github.com/stretchr/testify/suite"
    "math/rand"
    "net/http"
    "net/http/httptest"
    "sync"
)

type testHandler struct {
}

func (h *testHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
    // Chance of producing an error status code
    if rand.Intn(100) > 95 {
        writer.WriteHeader(http.StatusForbidden)
        return
    }

    msg := "awesome"
    _, err := writer.Write([]byte(msg))
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
    }
}

type blasterTestSuite struct {
    suite.Suite
    server *httptest.Server

    configuration *blaster.Configuration
    wg            *sync.WaitGroup
    blasters      []*blaster.Blaster
}

// Test lifecycle hooks ------------------------------------------------------------------------------------------------

func (suite *blasterTestSuite) SetupSuite() {
    handler := testHandler{}
    suite.server = httptest.NewServer(&handler)

    var err error
    suite.configuration, err = blaster.NewConfiguration(
        suite.server.URL,
        http.MethodGet,
        0, // Default rate
        5, // 5 seconds
        http.Header{})
    if err != nil {
        suite.FailNowf("Failed to create configuration", "%v", err)
    }
}

func (suite *blasterTestSuite) TearDownTest() {
    // Try to stop them all
    for _, b := range suite.blasters {
        b.Stop()
    }

    // Reset the slice of blasters
    suite.blasters = nil
}

func (suite *blasterTestSuite) TearDownSuite() {
    suite.server.Close()
}

// Convenience functions -----------------------------------------------------------------------------------------------

func (suite *blasterTestSuite) setup(numBlasters, rate, duration int) {
    if err := suite.configuration.SetRate(rate); err != nil {
        suite.FailNowf("Failed to set rate", "%v", err)
    }

    if err := suite.configuration.SetDuration(duration); err != nil {
        suite.FailNowf("Failed to set duration", "%v", err)
    }

    suite.wg = &sync.WaitGroup{}
    suite.wg.Add(numBlasters)

    suite.blasters = make([]*blaster.Blaster, numBlasters)
    for i := 0; i < numBlasters; i++ {
        id := fmt.Sprintf("test#%d", i)
        b, err := blaster.NewBlaster(id, suite.configuration, suite.wg)
        if err != nil {
            suite.FailNowf("Failed to create new blaster", "%v", err)
        }
        suite.blasters[i] = b
    }
}

func (suite *blasterTestSuite) start() {
    for _, b := range suite.blasters {
        b.Start()
    }
}

func (suite *blasterTestSuite) wait() {
    suite.wg.Wait()
}

func (suite *blasterTestSuite) totalRequests() (sum int) {
    for _, b := range suite.blasters {
        sum += b.TotalRequests()
    }
    return
}

func (suite *blasterTestSuite) successfulRequests() (sum int) {
    for _, b := range suite.blasters {
        sum += b.SuccessfulRequests()
    }
    return
}

// Assert that the number of total requests are within p % of
// the expected total. For instance, if expecting at least 80 %
// of the total requests after running the blasters call this
// method with p = 0.80.
func (suite *blasterTestSuite) assertTotalRequestsWithin(p float64) {
    numBlasters := len(suite.blasters)
    approximateTotal := numBlasters * suite.configuration.Rate * int(suite.configuration.Duration.Seconds())

    leastExpectedTotal := int(float64(approximateTotal) * p)
    totalRequestsSent := suite.totalRequests()
    suite.Greater(totalRequestsSent, leastExpectedTotal)
}
