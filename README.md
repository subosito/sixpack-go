# sixpack-go

Go client library for SeatGeek's Sixpack AB testing framework.

## Usage

Here's the basic example:

```go
// import "github.com/subosito/sixpack-go/sixpack"

session, err := sixpack.NewSession(Options{})
if err != nil {
	t.Error(err)
}

// Participate in a test (create the test if necesssary)
res, err := session.Participate("new-test", []string{"alternative-1", "alternative-2"}, "")
if err == nil {
	fmt.Printf("%+v\n", res)
}

// Convert
rec, err := session.Convert("new-test")
if err == nil {
	fmt.Printf("%+v\n", rec)
}
```

Each session has a `ClientID` associates with it that must be preseved across requests.

Session initializes requires `Options` which you can use to customize the session.

```go
// import "net/url"

opts := sixpack.Options{
	BaseUrl: url.Parse("http://sixpack.server.com"),
}

session := sixpack.NewSession(opts)
```

