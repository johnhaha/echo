# echo

echo is a golang message tool
currently support pub/sub and time job

## Install

```bash
go get github.com/johnhaha/echo@v0.2.5
```

## Usage

### Pub

```go
echo.Pub(CHANNEL,CONTENT)
```

### Pub Json

```go
echo.PubJson(CHANNEL,JSON)
```

### Sub

```go
echo.Sub(CONTEXT,CHANNEL,HANDLER)
```

## use suber to make multi sub

```go
suber := echo.NewSuber()
suber.Add(CHANNEL,CONSUMER)
suber.Sub(ctx)
```
