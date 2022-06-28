# godynamicdns
A Dynamic DNS Updater in Go.

It fits my very narrow use case, hopefully it's useful to someone else too.

## Supported Providers
- Google Domains (Not Google Cloud DNS)

## Supported Platforms

| OS      | 386 | amd64 | arm6 | arm64 |
| ---     | --- | ----  | ---  | ----  |
| Linux   |     | ✅     | ✅    | ✅     |
| Windows |     |       |      |       |
| MacOS   |     |       |      |       |


## Building
### Requirements:
- GoLang 1.16+
- [go-licenses](https://github.com/google/go-licenses)
- [Goreleaser](https://goreleaser.com/) (optional, to build linux packages)

### Build instructions
```
go-licenses save github.com/jimmale/godynamicdns --save_path="./terms/terms/"
goreleaser release --rm-dist --snapshot
```


## Milestones
### MTP
- Configuration file ✅
- Make the HTTP POST to update the ip, not providing any IP and letting the server sort it out ✅
- Do so on a 60 minute loop ✅

### MVP
- systemd unit / pid file ✅
- deb/rpm packages ✅
- Configuration file ✅
- Retry logic ✅
- Debug logging flag ✅

### MVP+
- License flag ✅
- Run under systemd unit as something other than root
- Warn of dangerous file permissions on configuration file
- Config file in `/etc/godynamicdns/config.toml` ✅
- Get Public IP (and lease expiration) via UPnP
- Schedule IP updates based off of lease expiration
- Schedule IP updates when UPnP indicates that the public IP has changed
- Dry run
- Documentation
- SIGHUP for configuration reload
- version numbers in build ✅
- Update once every 24h + jitter

## Example configuration

`/etc/godynamicdns/config.toml`
```toml
debug = true

[[domain]]
username = "Bruce.Wayne"
password = "iamb@man"
hostname = "batcave.wayneindustries.com"
frequency = "60m"
```

## Dependencies
### Buildtime
| Library                                                         | License | Purpose                           |
| -------                                                         | ------- | -------                           |
| [Sirupsen/Logrus](https://github.com/Sirupsen/logrus)           | MIT     | Pretty Logging                    |
| [BurntSushi/toml](https://github.com/BurntSushi/toml)           | MIT     | Config File Parsing               |
| [urfave/cli](https://github.com/urfave/cli)                     | MIT     | Command line parameter management |
| [NebulousLabs/go-upnp](https://gitlab.com/NebulousLabs/go-upnp) | MIT     | Discovering external IP           |

### Runtime
- Working CA Certificate store (see [here](https://stackoverflow.com/a/40051432)) to build a secure connection to Google
- Systemd is recommended, but you can absolutely run it without.