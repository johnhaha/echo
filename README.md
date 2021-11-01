# echo

echo is a golang runtime message tool
It's currently support message pub/sub

## Install

```bash
go get github.com/johnhaha/echo
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
