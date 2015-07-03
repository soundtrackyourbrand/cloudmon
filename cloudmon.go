package cloudmon

import (
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/cloudmonitoring/v2beta2"
)

const prefix = "custom.cloudmonitoring.googleapis.com/"

type Client interface {
	CreateMetric(name string) error
	WriteInt(name string, value int64) error
}

type client struct {
	oauthEmail      string
	oauthPrivateKey string
	projectID       string
}

func (c *client) CreateMetric(name string) error {
	cloud, err := cloudmonitorClient(c.oauthEmail, c.oauthPrivateKey)
	if err != nil {
		return err
	}

	req := &cloudmonitoring.MetricDescriptor{
		Name:    prefix + name,
		Project: c.projectID,
		TypeDescriptor: &cloudmonitoring.MetricDescriptorTypeDescriptor{
			MetricType: "gauge",
			ValueType:  "int64",
		},
	}

	_, err = cloud.MetricDescriptors.Create(c.projectID, req).Do()

	if err != nil {
		return err
	}

	return nil
}

func (c *client) WriteInt(name string, value int64) error {
	cloud, err := cloudmonitorClient(c.oauthEmail, c.oauthPrivateKey)
	if err != nil {
		return err
	}

	_, err = cloud.Timeseries.Write(c.projectID, &cloudmonitoring.WriteTimeseriesRequest{
		Timeseries: []*cloudmonitoring.TimeseriesPoint{
			&cloudmonitoring.TimeseriesPoint{
				Point: &cloudmonitoring.Point{
					Int64Value: value,
					Start:      time.Now().Format(time.RFC3339),
					End:        time.Now().Format(time.RFC3339),
				},
				TimeseriesDesc: &cloudmonitoring.TimeseriesDescriptor{
					Metric:  prefix + name,
					Project: c.projectID,
				},
			},
		},
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

type optionFunc func(*client)

func NewClient(opts ...optionFunc) *client {
	c := &client{}
	for _, fn := range opts {
		fn(c)
	}

	return c
}

func OAuthSettings(email, privateKey string) optionFunc {
	return func(c *client) {
		c.oauthEmail = email
		c.oauthPrivateKey = privateKey
	}
}

func ProjectID(projectID string) optionFunc {
	return func(c *client) {
		c.projectID = projectID
	}
}

func cloudmonitorClient(email, privateKey string) (*cloudmonitoring.Service, error) {
	conf := &jwt.Config{
		Email:      email,
		PrivateKey: []byte(privateKey),
		Scopes: []string{
			cloudmonitoring.MonitoringScope,
		},
		TokenURL: google.JWTTokenURL,
	}

	return cloudmonitoring.New(conf.Client(oauth2.NoContext))
}
