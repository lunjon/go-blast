package main

import (
    "flag"
    "fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
    "regexp"
    "strings"
    "sync"
    "time"

    "github.com/lunjon/go-blast/pkg/blaster"
)

func init() {
	// Blast file
	flag.StringVar(&file, "file", "", "Filepath to a blast configuration file.")

	// Request flags
	flag.StringVar(&url, "url", "", "The target URL.")
	flag.StringVar(&method, "method", http.MethodGet, "The HTTP method to use.")
	flag.StringVar(&body, "body", "", "JSON request body (only used in POST).")
    flag.Var(&headers, "header", "Headers to use in the request.")

	// Number of blasters, rate and duration
	flag.IntVar(&numBlasters, "num", defaultBlasters, "The number of blasters to run.")
	flag.IntVar(&rate, "rate", blaster.DefaultRate, "The rate of the requests.")
	flag.IntVar(&duration, "duration", blaster.DefaultDuration, "Time in seconds to run.")

	// Misc
	flag.BoolVar(&verbose, "verbose", false, "Output detailed logs.")
	flag.BoolVar(&verbose, "v", false, "Output detailed logs. (shortname)")
	flag.Parse()

	if !verbose {
		log.SetOutput(ioutil.Discard)
	}

	if numBlasters < minBlasters || numBlasters > maxBlasters {
		fmt.Printf("invalid number of blasters, must be an integer in the range %d-%d\n", minBlasters, maxBlasters)
		os.Exit(1)
	}

}

const (
	defaultBlasters = 1
	maxBlasters = 100
	minBlasters = 1
)

var (
    headerReg = regexp.MustCompile(`([\w-]+)\s?[:=]\s?(.+)`)
    numBlasters int
	file string
	url string
	method string
	body string
    headers = HeaderFlag{header: http.Header{}}
	rate int
	duration int
	verbose bool

)

func main() {
    // Parse flags to get configuration
	config := parseFlags()

	// Print configuration
	fmt.Printf("Number of blasters:\t%d\n", numBlasters)
	fmt.Printf("Request rate (req/s):\t%d\n", config.Rate)
	fmt.Printf("Duration:\t\t%v\n", config.Duration)
	fmt.Printf("Endpoint URL:\t\t%v\n", config.URL)

	if len(config.Header) > 0 {
		fmt.Println("Headers:")
		for k, v := range config.Header {
			fmt.Printf("\t%s: %s\n", k, strings.Join(v, "; "))
		}
	}

	start := time.Now()
    fmt.Printf("Starting:\t\t%s\n", start.Format(time.Stamp))

	// Initialize blasters and start
	var wg sync.WaitGroup
	wg.Add(numBlasters)

	blasters := make([]*blaster.Blaster, numBlasters)
	for n := 0; n < numBlasters; n++ {
		b, _ := blaster.NewBlaster(fmt.Sprintf("#%d", n), config, &wg)
		b.Start()
		blasters[n] = b
	}

	// Wait for the blasters to finish
	wg.Wait()
    end := time.Now()
	elapsed := time.Since(start)

	// Display the results
	totalRequests := 0
	successfulRequests := 0
	for _, b := range blasters {
		totalRequests += b.TotalRequests()
		successfulRequests += b.SuccessfulRequests()
	}

	blastersFormat := "blaster"
	if numBlasters != 1 {
		blastersFormat += "s"
	}

	fmt.Printf(
		"%d %s done %s (after %v) with %d/%d successful requests\n",
		numBlasters,
		blastersFormat,
        end.Format(time.Stamp),
		elapsed,
		successfulRequests,
		totalRequests)
}

type HeaderFlag struct {
	header http.Header
}

func (h *HeaderFlag) String() string {
    return ""
}

func (h *HeaderFlag) Set(v string) error {
    match := headerReg.FindAllStringSubmatch(v, -1)
    if match == nil {
        return fmt.Errorf("invalid header format: %v", v)
    }

    key := match[0][1]
    value := match[0][2]
    h.header.Add(key, value)
    return nil
}

func parseFlags() (config *blaster.Configuration) {
	var err error
	if file != "" {
        log.Printf("Loading from file: %s", file)

		// The --file flag was provided
		config, err = blaster.LoadFile(file)
		checkError(err, "failed to load blast file")

		// Allow rate to be overridden by the command line flag.
		if rate != blaster.DefaultRate {
			err = config.SetRate(rate)
			checkError(err, "failed to set rate")
		}

		// Also allow duration to be overridden by the command line flag.
		if duration != blaster.DefaultDuration {
			err = config.SetDuration(duration)
			checkError(err, "failed to set duration")
		}

		// Headers added to the command line should be added as well,
		// and they may also override any value from the file.
		config.UpdateHeader(headers.header)
	} else {
        log.Print("Using command line options as configurations")
		// Assume all configuration comes from the command line
		if url == "" {
			fmt.Println("A valid URL must be provided")
			os.Exit(1)
		}

		config, err = blaster.NewConfiguration(
			url,
			method,
			rate,
			duration,
			headers.header)
		checkError(err, "failed to create configuration")

		var b []byte
		if body != "" && config.HTTPMethod == http.MethodPost {
			b = []byte(body)
		}
		config.SetRequestBody(b)
	}

	return
}

func checkError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %v\n", msg, err)
		os.Exit(1)
	}
}
