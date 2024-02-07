# Figaro

[![Build Status](https://github.com/alexdebril/figaro/actions/workflows/go.yml/badge.svg)](https://github.com/alexdebril/figaro/actions/workflows/go.yml/)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexdebril/figaro)](https://goreportcard.com/report/github.com/alexdebril/figaro)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=alexdebril_figaro&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=alexdebril_figaro)

Small library to make easy Go 

*Easy come, easy go, will you let me go?*

The library takes its name from a famous song wrote by a famous artist as the goal of those packages is to make "easy go".

## Packages

### http/response

This package provides helpers designed to help in writing HTTP response. The principle is easy: create a new Response fed with a set of options and inject the `http.responseWriter` into it to write the response.

```go
type Message struct {
    Date    time.Time `json:"date"`
    Message string    `json:"message"`
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    message := &Message{
        Date:    time.Now(),
        Message: "this is the new message",
    }
    resp := response.NewJsonResponse(message)
    resp.Write(w)
}
```

See ? it's easy.