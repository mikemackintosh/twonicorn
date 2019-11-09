# Twonicorn

Simply put, Twonicorn is an API translation service. There are many services that exist on the web that we would like to communicate natively with another system. Due to no standard API schema or Webhook data structures, this is almost impossible for the service we find ourselves configuring.

Think of systems that don't have native integrations within their app, or those that grant excessively permissive roles to their applications. Twonicorn head-butt's this problem by giving you the ability to configure what the destination API request should look like. Using [Go HTML template](https://golang.org/pkg/html/template/) language, you can substitute your received webhook payload into your destination schema. 

With Twonicorn, you can magically make all your non-native application integrations disappear like fairy dust. 

## Key Features
- Simple and Reliable
- Configurations as Code
- No Dependencies Required

## Usage

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
