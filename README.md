# godynamicdns
A Dynamic DNS Updater in Go

## Supported Providers
- Google Domains (Not Google Cloud DNS)

## Milestones
### MTP
- Configuration file
- Make the HTTP POST to update the ip, not providing any IP and letting the server sort it out
- Do so on a 60 minute loop

### MVP
- systemd unit / pid file
- deb/rpm packages
- Configuration file
- Retry logic
- Debug logging flag

### MVP+
- Get Public IP (and lease expiration) via UPnP
- Schedule IP updates based off of lease expiration
- Dry run
- Documentation
- SIGHUP for configuration reload


##Example configuration
```toml
debug = true

[[domain]]
username = "Bruce.Wayne"
password = "iamb@man"
hostname = "batcave.wayneindustries.com"
frequency = "60m"
```


## Libraries Used
| Library | License | Purpose | 
| ------- | ------- | ------- | 
| [Sirupsen/Logrus](https://github.com/Sirupsen/logrus) | MIT | Pretty Logging | 
| [BurntSushi/toml](https://github.com/BurntSushi/toml) | MIT | Config File Parsing | 
| [NebulousLabs/go-upnp](https://gitlab.com/NebulousLabs/go-upnp) | MIT | Discovering external IP |

