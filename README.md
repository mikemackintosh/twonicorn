# Twonicorn

Simply put, Twonicorn is an API translation service. There are many services that exist on the web that we would like to communicate natively with another system. Due to no standard API schema or Webhook data structures, this is almost impossible for the service we find ourselves configuring.

Think of systems that don't have native integrations within their app, or those that grant excessively permissive roles to their applications. Twonicorn head-butt's this problem by giving you the ability to configure what the destination API request should look like. Using [Go HTML template](https://golang.org/pkg/html/template/) language, you can substitute your received webhook payload into your destination schema. 

With Twonicorn, you can magically make all your non-native application integrations disappear like fairy dust. 

## Key Features
- Simple and Reliable
- Configurations as Code
- No Dependencies Required

## Configuration
The following example can be found within [`/example-config.yml`](/example-config.yml).

```
curl_to_slack:             # Name of your configuration entry, must be unique
  target: http             # Currently the default and only target
  url: https://hooks.sla.. # The URL you would like to send your translated request to
  method: post             # The method you would like to send your request to
  headers:                 # Any additional headers, like auth headers, to add to your
    content-type: "application/json"                              # translated request
  owner: mike@domain       # All code should have an owner, let's start a good habbit
  translate:               # Starts the translation block
    text: "{{.blah}}"
    channel: "#teehee"
    username: "Twonicorn"
```

When making the following request:

```
curl -d'{"blah":"So much awesomenessssssss"}' \
     -H'Content-type: application/json' \
     localhost:8080/hooks/d12aaaff28608037b/6dd5b
```

The JSON key of `.blah` will be transmogrified into `{"text": "{{.blah}}", ...}`. This would result in the following payload being sent to the targeted upstream service:

```
{"text": "So much awesomeness", "channel": "#teehee", "username": "Twonicorn"}
```

## Usage

### -c, Config File
***Default: config.yml***
The configuration file to use. There is an included `make configtest` for validating your configuration. Using the `-keys` flag when running `cmd/configtest/configtest.go` will provide you with your `/hooks/<key>` to send your requests to.

### -p, Listener
***Default: localhost:8080***
Use this option to change the listening host and port in `<ip>:<port>` format.

### -d, Debugging
***Default: false***
When enabled, debugging will print out debugging information such as upstream requests and responses, and configuration-related breakpoints. 

### -o, Output Logging Variations
***Default: JSON***

As configured, all output is done to `os.Stderr`. This allows logging to selectively by captured or redirected. There are plans to make this more controllable by the operator. 

The following logging types are supported:
  - Apache
  - JSON
  - Yaml
  - Key=Value
  
Example JSON output:

```
{
  "RequestID": "bn37qg7budgk2hdfpp50",
  "Host": "localhost:8080",
  "RemoteAddr": "127.0.0.1:64732",
  "Method": "POST",
  "RequestURI": "/hooks/d12aaaff28608037b/6dd5b",
  "Proto": "HTTP/1.1",
  "Status": 404,
  "ContentLen": 48,
  "UserAgent": "Mozilla/5.0 (Windows NT 6.1; rv:31.0) Gecko/20100101 Firefox/31.0",
  "Time": "09/Nov/2019:00:48:00 -0800",
  "Referer": "",
  "Duration": 375983061
}
```
