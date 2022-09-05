package healthz

import (
	"net/http"

	_ "github.com/PCManiac/logrus_init"
	"github.com/klyve/go-healthz"
	"github.com/sirupsen/logrus"
)

type HealthzServer interface {
	StartHealthz()
}

type server struct {
	metricsPort string
	provider    healthz.Checkable
}

func (s *server) StartHealthz() {
	healthzInstance := healthz.Instance{
		Logger: logrus.StandardLogger(),
	}

	if s.provider != nil {
		healthzInstance.Providers = []healthz.Provider{
			{
				Handle: s.provider,
				Name:   "server",
			},
		}
	}

	http.Handle("/healthz", healthzInstance.Healthz())
	http.Handle("/liveness", healthzInstance.Liveness())

	go http.ListenAndServe(s.metricsPort, nil)
}

func New(addr string, provider healthz.Checkable) HealthzServer {
	s := server{
		metricsPort: addr,
		provider:    provider,
	}

	s.StartHealthz()

	return &s
}
