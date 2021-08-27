package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jimmale/godynamicdns/config"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"
)

var goodRegex = regexp.MustCompile("good .*")
var nochgRegex = regexp.MustCompile("nochg .*")

func main() {
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)

	contents, err := ioutil.ReadFile("config.toml")
	if err != nil {
		log.Fatal(err.Error())
	}

	conf := config.Configuration{}

	_, err = toml.Decode(string(contents), &conf)

	if err != nil {
		log.Fatalf("Could not parse configuration file: %s", err.Error())
	}

	if conf.Debug {
		log.SetLevel(log.TraceLevel)
	}

	for _, v := range conf.Domains {
		log.Tracef("Launching goroutine to keep %s updated", v.Hostname)
		go maintain(context.Background(), &v)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	log.Info("awaiting exit signal")
	<-sigs
	log.Infof("Exiting")
}

func maintain(ctx context.Context, domain *config.Domain) {
	for {

		select {
		case <-ctx.Done():
			{
				log.Trace("maintain -> ctx.Done case")
				return
			}
		default:
			{
				log.Trace("maintain -> default case")
				err := update(domain.Hostname, domain.Username, domain.Password)
				if err != nil{
					log.Error(err.Error())
					return
				}
				log.Tracef("Updater for %s sleeping for %s", domain.Hostname, domain.Frequency.Duration.String())
				time.Sleep(domain.Frequency.Duration)
			}
		}
	}
}

func update(hostname string, username string, password string) error {
	log.Tracef("Initiating IP Update for %s", hostname)
	client := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       60 * time.Second,
	}

	req, _ := http.NewRequest(http.MethodPost, "https://domains.google.com/nic/update", nil)

	q := req.URL.Query()
	q.Add("hostname", hostname)
	req.URL.RawQuery = q.Encode()
	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	switch {
	case nochgRegex.Match(respBody):
		{
			log.Tracef("No change to IP for %s: %s", hostname, respBody)
		}

	case goodRegex.Match(respBody):
		{
			log.Infof("Successfully updated %s: %s", hostname, respBody)
			break
		}
		// TODO handle other responses https://support.google.com/domains/answer/6147083#zippy=%2Cuse-the-api-to-update-your-dynamic-dns-record
	default:
		{
		return errors.New(fmt.Sprintf("Received error while updating IP: %s", respBody))
		}
	}

	return nil
}
