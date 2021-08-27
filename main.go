package main

import (
	"context"
	"github.com/BurntSushi/toml"
	"github.com/jimmale/godynamicdns/config"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main(){

	log.SetLevel(log.TraceLevel)
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)

	contents, err := ioutil.ReadFile("config.toml")
	if err != nil{
		log.Fatal(err.Error())
	}

	conf := config.Configuration{}

	toml.Decode(string(contents), &conf)
	for _, v := range conf.Domains{
		log.Tracef("Launching goroutine to keep %s updated", v.Hostname)
		go maintain(context.Background(), &v)
	}


	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	log.Info("awaiting exit signal")
	<-sigs
	log.Infof("Exiting")
}

func maintain(ctx context.Context, domain *config.Domain){
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
				update(domain.Hostname, domain.Username, domain.Password)
				log.Tracef("Updater for %s sleeping for %s", domain.Hostname, domain.Frequency.Duration.String())
				time.Sleep(domain.Frequency.Duration)
			}
		}
	}
}


func update(hostname string, username string, password string) error{
	log.Tracef("Initiating IP Update for %s", hostname)
	client := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       60*time.Second,
	}

	req, _ := http.NewRequest(http.MethodPost, "https://domains.google.com/nic/update", nil )

	q := req.URL.Query()
	q.Add("hostname", hostname)
	req.URL.RawQuery = q.Encode()
	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil{
		return err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	log.Infof("IP Update for %s complete: %s", hostname, string(respBody))
	return nil
}