# Spike on B2 Event Notifications

B2 can send [event notifications](https://www.backblaze.com/bb/docs/cloud-storage-event-notifications) (currently in private beta). This is a webook server capable of receiving them.

# Iterate

Use [entr](https://eradman.com/entrproject/):

```command
$ while sleep 0.1; do
  fd | entr -r -d go run .
done
```

# Integration Test

Assuming the shared HMAC secret is `eqdIy6lTcbDXsGFJnyD5eK83sCQUIrPH`, we can start the server:

```command
$ B2_EVENT_NOTIFICATIONS_SHARED_SECRET=eqdIy6lTcbDXsGFJnyD5eK83sCQUIrPH go run .
```

and test it with

```command
$ curl \
    --data @b2/events.json \
    --header 'Content-Type: application/json' \
    --header 'X-Bz-Event-Notification-Signature: v1=ed43d3f65391cdeef81783321a5c3408977050e4747085023fd6ca02cdf71543' \
  http://127.0.0.1:62057/
```

Replace the URL with the output of `tailscale funnel 62057`, if you use that one.
