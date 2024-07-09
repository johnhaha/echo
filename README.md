# echo

echo is a golang event message tool
support pub/sub, time job and trigger

## Install

```bash
go get github.com/johnhaha/echo@v0.2.7
```

## Usage

### set up handler

```go
echo.SetEventHandler(CHANNEL,HANDLER)
```

### Pub

```go
echo.PubEvent(CHANNEL,CONTENT)
```

### Pub Json

```go
echo.PubEventJson(CHANNEL,JSON)
```

### Listen to event

```go
echo.StartEventListener(CONTEXT)
```

### use pubsub to make multi sub

```go
suber := echo.NewPubSub()
suber.Pub(DATA)
suber.Sub(ctx,GROUP,CONSUMER)
```

### set up timer event storage

```go
echo.SetStorage(storage)
```

the timer event will be load if you load timer event after you set storage

```go
echo.LoadTimerEvent()
```
