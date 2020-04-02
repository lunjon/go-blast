package blaster

import (
    "encoding/json"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "log"
    "net/http"
)

// blastFile is the structure of a BlastFile, i.e. a YAML file
// that contains the configuration for a blast.
type blastFile struct {
    Rate     int `yaml:"rate"`
    Duration int `yaml:"duration"`
    Request  struct {
        URL     string `yaml:"url"`
        Method  string `yaml:"method"`
        Headers []struct {
            Name  string `yaml:"name"`
            Value string `yaml:"value"`
        } `yaml:"headers"`
        Body map[string]interface{} `yaml:"body"`
    }
}

// LoadFile returns the resulting configuration in the file.
// filepath must be a path to a YAML file conforming to the
// structure of a blast configuration file.
func LoadFile(filename string) (*Configuration, error) {
    b, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    c := &blastFile{}
    err = yaml.Unmarshal(b, c)
    if err != nil {
        return nil, err
    }

    header := http.Header{}
    for _, h := range c.Request.Headers {
        log.Printf("%s: %s", h.Name, h.Value)
        header.Add(h.Name, h.Value)
    }

    config, err :=  NewConfiguration(
        c.Request.URL,
        c.Request.Method,
        c.Rate,
        c.Duration,
        header)
    if err != nil {
        return nil, err
    }

    var body []byte
    if c.Request.Body != nil {
        body, err = json.Marshal(c.Request.Body)
    }
    config.SetRequestBody(body)
    return config, err
}
