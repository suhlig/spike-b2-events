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

and test it locally with

```command
$ curl \
    --data @b2/events.json \
    --header 'Content-Type: application/json' \
    --header 'X-Bz-Event-Notification-Signature: v1=8ac50cc33d141ea5450a1c3e7c85c2cc98f7ba8b5043ecee56ecb61385d04920' \
  http://127.0.0.1:62057/
```

Replace the URL with the output of `tailscale funnel 62057`, if you use that one.

An end-to-end test can be performed with

```command
$ echo 3 > /tmp/three
$ b2 file upload suhlig-spike-event-notifications /tmp/three three
```
