## Distribution IRC notifier plugin

Plugin for [IRC-Bot](https://github.com/greboid/irc-bot)

Receives notifications from a [Distribution](https://github.com/distribution/distribution) instance and outputs them to
a channel.

- go build github.com/csmith/irc-distribution/cmd/distribution

#### Configuration

At a bare minimum you also need to give it a channel, a bearer token to verify received notifications
on and an RPC token. You'll likely also want to specify the bot host.

Once configured the URL to configure in distribution would be <Bot URL>/distribution, for example:

```yaml
notifications:
  endpoints:
    - name: irc
      url: https://mybothost.example/distribution
      headers:
        Authorization: [Bearer TheConfiguredBearerToken]
      timeout: 500ms
      threshold: 5
      backoff: 1s
```

#### Example running

```
distribution -rpc-host bot -rpc-token <as configured on the bot> -channel #spam -bearer-token cUCrb7HJ
```
