package grpc

import (
	locationv1 "github.com/zchelalo/neuraclinic-location/gen/go/location/v1"
	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/application"
)

type LocationService struct {
	locationv1.UnimplementedLocationServiceServer
	app *application.Service
}

func NewLocationService(app *application.Service) *LocationService {
	return &LocationService{app: app}
}
