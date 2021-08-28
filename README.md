# godynamicdns
A Dynamic DNS Updater in Go

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
- [Goreleaser](https://goreleaser.com/) (optional, to build linux packages)

### Build instructions
```
goreleaser release --rm-dist --snapshot
```


## Milestones
### MTP
- Configuration file ✅
- Make the HTTP POST to update the ip, not providing any IP and letting the server sort it out ✅
- Do so on a 60 minute loop ✅

### MVP
- systemd unit / pid file
- deb/rpm packages ✅
- Configuration file ✅
- Retry logic
- Debug logging flag ✅

### MVP+
- Config file in `/etc` or local directory
- Get Public IP (and lease expiration) via UPnP
- Schedule IP updates based off of lease expiration
- Dry run
- Documentation
- SIGHUP for configuration reload
- version numbers in build

## Example configuration

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