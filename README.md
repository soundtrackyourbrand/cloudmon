# cloudmon [![GoDoc](https://godoc.org/github.com/soundtrackyourbrand/cloudmon?status.svg)](https://godoc.org/github.com/soundtrackyourbrand/cloudmon)
A small wrapper for parts of Google Cloud Monitoring API

# Example

```go
c := cloudmon.NewClient(
	cloudmon.ProjectID("google-project-id"),
	cloudmon.OAuthSettings(os.Getenv("OAUTH_EMAIL"), os.Getenv("OAUTH_PRIVATE_KEY")),
)

g, err := c.NewGauge("name-of-gauge")
if err != nil {
	log.Fatal(err)
}

if err := g.Set(500); err != nil {
  log.Fatal(err)
}
```
