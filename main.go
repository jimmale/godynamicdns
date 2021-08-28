package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jimmale/godynamicdns/config"
	"github.com/jimmale/godynamicdns/licenseterms"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"
)

// These are set in the CICD pipeline.
var version string
var commit string
var date string
var builtBy string

var goodRegex = regexp.MustCompile("good .*")
var nochgRegex = regexp.MustCompile("nochg .*")

func main() {

	if version == "" {
		version = "snapshot"
	}
	if commit == "" {
		commit = "unknown commit"
	}
	if date == "" {
		date = "???"
	}
	if builtBy == "" {
		builtBy = "anonymous"
	}

	versionString := fmt.Sprintf("%s (%s) built on %s by %s", version, commit, date, builtBy)

	app := &cli.App{
		Name:    "godynamicdns",
		Usage:   "A Dynamic DNS Updater in Go",
		Version: versionString,
		Action:  mainAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "enable debug logging",
				Value: false,
			},
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "specify a configuration file",
				EnvVars:     nil,
				FilePath:    "",
				Required:    false,
				Hidden:      false,
				TakesFile:   false,
				Value:       "/etc/godynamicdns/config.toml",
				DefaultText: "",
				Destination: nil,
				HasBeenSet:  false,
			},
			&cli.BoolFlag{
				Name: "license",
				Usage: "print the license terms of this software and exit",
				Value: false,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
func mainAction(c *cli.Context) error {

	if c.Bool("license"){
		licenseterms.PrintLicenseTerms()
		os.Exit(0)
	}

	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)

	configFile := c.String("config")

	contents, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	conf := config.Configuration{}

	_, err = toml.Decode(string(contents), &conf)

	if err != nil {
		log.Fatalf("Could not parse configuration file: %s", err.Error())
	}

	//// connect to router
	//d, err := upnp.Discover()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// discover external IP
	//ip, err := d.ExternalIP()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println("Your external IP is:", ip)

	if conf.Debug || c.Bool("debug") {
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
	return nil
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
				err := update(domain, false)
				if err != nil {
					log.Error(err.Error())
					return
				}
				log.Tracef("Updater for %s sleeping for %s", domain.Hostname, domain.Frequency.Duration.String())
				time.Sleep(domain.Frequency.Duration)
			}
		}
	}
}

func update(domain *config.Domain, isRetry bool) error {
	log.Tracef("Initiating IP Update for %s", domain.Hostname)
	client := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       60 * time.Second,
	}

	req, _ := http.NewRequest(http.MethodPost, "https://domains.google.com/nic/update", nil)

	req.Header.Set("User-Agent", "godynamicdns")

	q := req.URL.Query()
	q.Add("hostname", domain.Hostname)
	req.URL.RawQuery = q.Encode()
	req.SetBasicAuth(domain.Username, domain.Password)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		if !isRetry {
			time.Sleep(5 * time.Minute)
			return update(domain, true)
		} else {
			return fmt.Errorf("could not update DNS. result code was %d", resp.StatusCode)
		}

	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	switch {
	case nochgRegex.Match(respBody):
		{
			log.Tracef("No change to IP for %s: %s", domain.Hostname, respBody)
		}

	case goodRegex.Match(respBody):
		{
			log.Infof("Successfully updated %s: %s", domain.Hostname, respBody)
			break
		}
	case string(respBody) == "nohost":
		{
			return fmt.Errorf("the hostname %s doesn't exist, or doesn't have Dynamic DNS enabled", domain.Hostname)
		}
	case string(respBody) == "badauth":
		{
			return fmt.Errorf("the username/password combination isn't valid for the specified host (%s)", domain.Hostname)
		}
	case string(respBody) == "notfqdn":
		{
			return fmt.Errorf("the supplied hostname (%s) isn't a valid fully-qualified domain name", domain.Hostname)
		}
	case string(respBody) == "badagent":
		{
			return errors.New("your Dynamic DNS client makes bad requests. Ensure the user agent is set in the request")
		}
	case string(respBody) == "abuse":
		{
			return errors.New("dynamic DNS access for the hostname has been blocked due to failure to interpret previous responses correctly")
		}
	case string(respBody) == "911":
		{
			if isRetry {
				return errors.New("an error occurred")
			} else {
				time.Sleep(5*time.Minute + 30*time.Second) // I know the docs say 5 minutes, but let's give it a bit more time.
				return update(domain, true)
			}
		}
	case string(respBody) == "conflict A":
		{
			return fmt.Errorf("a custom A resource record conflicts with the update to %s . Delete the indicated resource record within the DNS settings page and try the update again", domain.Hostname)
		}
	case string(respBody) == "conflict AAAA":
		{
			return fmt.Errorf("a custom AAAA resource record conflicts with the update to %s . Delete the indicated resource record within the DNS settings page and try the update again", domain.Hostname)
		}

	default:
		{
			return errors.New(fmt.Sprintf("Received error while updating IP: %s", respBody))
		}
	}
	return nil
}
