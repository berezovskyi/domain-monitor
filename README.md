# domain-monitor

[![Test Build Go Project](https://github.com/berezovskyi/domain-monitor/actions/workflows/go.yml/badge.svg)](https://github.com/berezovskyi/domain-monitor/actions/workflows/go.yml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/berezovskyi/domain-monitor?style=flat)
![GitHub package.json version](https://img.shields.io/github/package-json/v/berezovskyi/domain-monitor?style=flat)

Self-hosted server to monitor WHOIS records for specified domains, to help alert you of unwanted/unexpected changes
to your domains and to help remind you when they need renewed.

If you use domain-monitor, make sure to abide by the terms of service of
the TLD WHOIS server. Most forbid high-volume queries, marketing usage,
and automated queries that are more than reasonably needed to register
domain names.

To help you abide by WHOIS server TOS, domain-monitor barely queries the
whois databases. domain-monitor acts in this way:

1. Creates a local cache of the whois info for the domain, including time queried
2. Updates a whois info only if one of these conditions is met:

- No cached whois info for the domain
- Reference info becomes 9 months old
- It's 3 months until the domain expiry date
- It's 2 months until the domain expiry date
- It's 1 month until the domain expiry date
- It's 2 weeks until the domain expiry date
- It's 1 week until the domain expiry date
- Manual update is requested
- A DNS query of the FQDN returns a different set of name servers than the cached whois info

## Installation

### Docker

The docker image has a single volume for the files generated by the server.

| Image Mount | Contains                                 |
| ----------- | ---------------------------------------- |
| /app/data   | config.yaml, domain.yaml and WHOIS cache |

| Exposed Ports | Used for        |
| ------------- | --------------- |
| 3124/tcp      | WEB GUI and API |

#### Using [Github Packages](https://github.com/berezovskyi/domain-monitor/packages/)

`docker run -p 127.0.0.1:3124:3124 -v ./data:/app/data ghcr.io/berezovskyi/domain-monitor:1`

Example docker-compose:

```yaml
services:
  dm:
    image: ghcr.io/berezovskyi/domain-monitor:1
    ports:
      - 127.0.0.1:3124:3124/tcp
    volumes:
      - ./data:/app/data:rw
version: "3.9"
```

#### Using [Docker Hub](https://hub.docker.com/repository/docker/berezovskyi/domain-monitor)

Image is just `berezovskyi/domain-monitor`, latest tag will be the most recent version, or pull by tagged version. The most recent old (non-go) version is tagged with `0.3.4`.

`docker run -p 127.0.0.1:3124:3124 -v ./data:/app/data docker.io/berezovskyi/domain-monitor`

### Running Locally

Should be OS-agnostic. Requires Go and NodeJS. Clone this repo and follow the [build instructions](#development). Then you can run the binary alongside the `assets` folder.

Configuration can be done via the configuration page of the web gui
(default http://localhost:3124)

## Config

There are two config files which you can edit yourself if you so choose.

### config.yaml

A sample is provided as `sample.config.yaml` and on first run if you don't have an existing `config.yaml`, domain-monitor will create one with all default values. Any changes you made in the webgui persist in `config.yaml`.

#### App Settings

_Port_

Set the port used by the http server

_Automated Whois Lookups_

Enable or disable automated whois lookups. If disabled, whois lookups will only be done when manually requested.

##### Sample App Config

```yaml
app:
  port: 3124
  automateWHOISRefresh: yes
```

#### Alerts

_admin_

Set what email should receive alerts from domain-monitor

_sendalerts_

Boolean, if false prevents all alerts from being sent.

_Send Alert at 2 Months to Expiry_

Boolean, if true, an alert will be sent when a domain is 2 months from expiry.

_Send Alert at 1 Month to Expiry_

Boolean, if true, an alert will be sent when a domain is 1 month from expiry.

_Send Alert at 2 Weeks to Expiry_

Boolean, if true, an alert will be sent when a domain is 2 weeks from expiry.

_Send Alert at 1 Week to Expiry_

Boolean, if true, an alert will be sent when a domain is 1 week from expiry.

_Send Alert at 3 Days to Expiry_

Boolean, if true, an alert will be sent when a domain is 3 days from expiry.

_Send Daily Expiry Alerts_

Boolean, if true, an alert will be sent every day for domains that expire within a week.

##### Sample Alerts Config

```yaml
alerts:
  admin: support@example.com
  sendalerts: true
  send2MonthAlert: false
  send1MonthAlert: true
  send2WeekAlert: false
  send1WeekAlert: false
  send3DayAlert: true
  sendDailyExpiryAlert: false
```

#### SMTP

Set smtp settings for domain-monitor to use to send email alerts.

_Host_

The hostname of the SMTP server

_Port_

The port of the SMTP server

_Secure_

Boolean, if true, will force the SMTP connection to use TLS. If false, will use opportunistic TLS (so fallback to plaintext if the server doesn't support TLS)

_Authuser_

The username to use to authenticate with the SMTP server

_Authpass_

The password to use to authenticate with the SMTP server

_Enabled_

Boolean, if false, domain-monitor will not send any email alerts

_FromName_

The name to use in the "from" field of the email alerts

_FromEmail_

The email address to use in the "from" field of the email alerts

##### Sample SMTP Config

```yaml
smtp:
  host: localhost
  port: 25
  secure: false
  authuser: domain-alert@example.com
  authpass: SECRET-PASS
  enabled: false
  fromName: Domain Monitor
  fromEmail: domain-monitor@example.com
```

#### Scheduler

Set some schedule options for the WHOIS lookups.

_WHOIS Cache Stale Interval_

The number of days after which the WHOIS cache is considered stale and a new lookup will be done.

_Use Standard WHOIS Refresh Schedule_

Boolean, if true, domain-monitor will use a standard schedule for WHOIS lookups. If false, it will still perform the automated WHOIS lookup for stale, new domains and DNS changes, but will not perform regular lookups.

##### Sample Scheduler Config

```yaml
scheduler:
  whoisCacheStaleInterval: 190
  useStandardWHOISRefreshSchedule: true
```

### domain.yaml

Contains a single object (domains) which is a list of domains to
monitor. Each domain has the following properties:

| Property | Type   | Description                                                                            |
| -------- | ------ | -------------------------------------------------------------------------------------- |
| name     | string | Descriptive name for the domain entry                                                  |
| fqdn     | string | FQDN for the domain in question. This is just `host.tld`                               |
| alerts   | bool   | If true, email alerts will be sent for this domain                                     |
| enabled  | bool   | If true, whois lookups will be done (on the schedule described above) for this domain. |

## Development

Requirements:

- [golang](https://golang.org/)
- [nodejs](https://nodejs.org/)
- [pnpm](https://pnpm.io/)
- [air](https://github.com/air-verse/air) `go install github.com/air-verse/air@latest`
- [templ](https://github.com/a-h/templ/) `go install github.com/a-h/templ/cmd/templ@latest`

### Development server

Build steps:

1. Install npm dependencies and "build" (copy the JS libraries to the `assets` folder).

```sh
pnpm i && pnpm build
```

2. Compile the templates

```sh
templ generate
```

3. Build the go binary

```sh
go build cmd/main.go
```

The binary requires the `assets` folder to be in the same directory as the binary.

To run the development server, use the following command:

```sh
air
```

Air will take care of all the build steps whenever a change is detected.
