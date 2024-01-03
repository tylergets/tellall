# Tell All

Swiss Army knife for local and network notifications.

**Note -- This project is a work in progress and is subject to change. Especially the configuration file format.**

# TODO
- [ ] Complete Nix Home Manager Module

# Description

At its core, Tell All is a notification service designed for use amongst personal machines and home lab style servers. 

When running, it can:
 * Ingest notifications via HTTP
   * Compatible with the following clients:
     * Gotify (/gotify)
     * Ntfy (/ntfy)
 * Ingest notifications via MQTT
 * Service those notifications via:
   * Anything supported by [Shoutrrr](https://containrrr.dev/shoutrrr/v0.8/services/overview/)
   * Local desktop notifications via [Beeep](https://github.com/gen2brain/beeep)

## Use Cases
    * Send a notification to every device you use from any server/desktop you use
    * Send a notification to your phone when a long running task completes
    * Use as a bridge from Gotify/Ntfy to another 3rd Party Service (Slack, Discord, etc)

### Config Example
```yaml
# System name, defaults to the system's hostname if not set.
# Should be unique if you have multiple instances of Tell All running.
name: "HostName" 

# Prefix used for the MQTT connection, default is 'tellall'.
prefix: "tellall" 

# Debug mode, set to true for verbose logging.
debug: false 

# MQTT connection string, required to connect to the MQTT broker.
mqtt_connection: "tcp://localhost:1883"

# Path to the file containing the MQTT secure connection configuration.
# If this is set, the MQTT connection string will be ignored.
mqtt_secure_connection: "/path/to/secure/connection/file"

# HTTP server configuration.
http_server:
  # Enable or disable the HTTP server.
  enabled: true 

  # Port number on which the HTTP server will listen.
  port: "8080" 

  # Hostname for the HTTP server, defaults to 'localhost'.
  host: "localhost" 

# List of listeners, an array of strings representing different services to listen to.
listeners:
  - "notify-send://" # Desktop notifications
  # See https://containrrr.dev/shoutrrr/v0.8/services/overview/ for more details
  - "bark://devicekey@host"
  - "discord://token@id"
  - "smtp://username:password@host:port/?from=fromAddress&to=recipient1,recipient2"
  - "gotify://gotify-host/token"
  - "googlechat://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz"
  - "ifttt://key/?events=event1,event2&value1=value1&value2=value2&value3=value3"
  - "join://shoutrrr:api-key@join/?devices=device1,device2&icon=icon&title=title"
  - "mattermost://username@mattermost-host/token/channel"
  - "matrix://username:password@host:port/?rooms=!roomID1,roomAlias2"
  - "ntfy://username:password@ntfy.sh/topic"
  - "opsgenie://host/token?responders=responder1,responder2"
  - "pushbullet://api-token/device/#channel/email"
  - "pushover://shoutrrr:apiToken@userKey/?devices=device1,device2"
  - "rocketchat://username@rocketchat-host/token/channel|@recipient"
  - "slack://botname@token-a/token-b/token-c"
  - "teams://group@tenant/altId/groupOwner?host=organization.webhook.office.com"
  - "telegram://token@telegram?chats=@channel-1,chat-id-1"
  - "zulip://bot-mail:bot-key@zulip-domain/?stream=name-or-id&topic=name"

```

### Compatibility with Gotify/Ntfy

Tell All is compatible with both Gotify and Ntfy style requests, but it does not support all the features of either.

#### Gotify
```shell
curl "http://localhost:8080/gotify/message" -F "title=My Title" -F "message=Hello World"
```

#### Ntfy

Ntfy supports topics at this time, so changing all to the device name will send the notification to only that device.

```shell
curl \
  -H "Title: My Title" \
  -d "Hello World" \
  http://localhost:8080/ntfy/all
```

## FAQ

 * Is this a replacement for Gotify/Ntfy?
   * For some use cases it might be, but its not intended to ever be as feature rich. The original goal was to provide a lightweight networked notification based on MQTT for my computers at home, but I decided to add mobile support via Pushover and then expanded the project.
 * Is there any authentication?
   * The MQTT side can be secured by using the `mqtt_secure_connection` option, but the HTTP side you will need to bring your own authentication (or don't expose it to the internet). I'm open to adding it if there is a need for it. 
<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.
