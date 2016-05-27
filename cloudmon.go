package cloudmon

import (
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/monitoring/v3"
)

const prefix = "custom.googleapis.com/"

// A Client is a cloud monitor client.
type Client interface {
	NewGauge(name string) (Gauge, error)
}

// Gauge is a cloud monitor metric that represent a single int value a specific time
type Gauge interface {
	Set(value int64) error
}

type client struct {
	oauthEmail      string
	oauthPrivateKey string
	projectID       string
}

type gouge struct {
	name   string
	client *client
}

type optionFunc func(*client)

// NewClient creates a new cloud monitor client
func NewClient(opts ...optionFunc) *client {
	c := &client{}
	for _, fn := range opts {
		fn(c)
	}

	return c
}

// OAuthSettings is a option function that can be sent as an argument to NewClient to setup oauth
func OAuthSettings(email, privateKey string) optionFunc {
	return func(c *client) {
		c.oauthEmail = email
		c.oauthPrivateKey = privateKey
	}
}

// ProjectID is a option function that can be sent as an argument to NewClient to set google project id
func ProjectID(projectID string) optionFunc {
	return func(c *client) {
		c.projectID = projectID
	}
}

// NewGauge creates a new int gauge in google cloud monitoring
func (c *client) NewGauge(name string) (Gauge, error) {
	cloud, err := cloudmonitorClient(c.oauthEmail, c.oauthPrivateKey)
	if err != nil {
		return nil, err
	}

	descriptor := &monitoring.MetricDescriptor{
		Type:       prefix + name,
		MetricKind: "GAUGE",
		ValueType:  "INT64",
	}

	_, err = cloud.Projects.MetricDescriptors.Create("projects/"+c.projectID, descriptor).Do()
	if err != nil {
		return nil, err
	}

	g := &gouge{
		name:   name,
		client: c,
	}

	return g, nil
}

// Set updates the value of the Gauge in google cloud monitoring
func (g *gouge) Set(value int64) error {
	cloud, err := cloudmonitorClient(g.client.oauthEmail, g.client.oauthPrivateKey)
	if err != nil {
		return err
	}

	_, err = cloud.Projects.TimeSeries.Create("projects/"+g.client.projectID, &monitoring.CreateTimeSeriesRequest{
		TimeSeries: []*monitoring.TimeSeries{
			&monitoring.TimeSeries{
				Points: []*monitoring.Point{
					&monitoring.Point{
						Interval: &monitoring.TimeInterval{
							StartTime: time.Now().Format(time.RFC3339),
							EndTime:   time.Now().Format(time.RFC3339),
						},
						Value: &monitoring.TypedValue{
							Int64Value:      value,
							ForceSendFields: []string{"Int64Value"},
						},
					},
				},
				Metric: &monitoring.Metric{
					Type: prefix + g.name,
				},
			},
		},
	}).Do()

	return err
}

func cloudmonitorClient(email, privateKey string) (*monitoring.Service, error) {
	conf := &jwt.Config{
		Email:      email,
		PrivateKey: []byte(privateKey),
		Scopes: []string{
			monitoring.MonitoringScope,
		},
		TokenURL: google.JWTTokenURL,
	}

	return monitoring.New(conf.Client(oauth2.NoContext))
}
