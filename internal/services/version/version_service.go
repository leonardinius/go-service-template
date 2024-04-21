package version

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	"github.com/leonardinius/go-service-template/internal/services/serviceotel"

	versionv1 "github.com/leonardinius/go-service-template/internal/apigen/version/v1"
	versionv1connect "github.com/leonardinius/go-service-template/internal/apigen/version/v1/versionv1connect"
)

type versionServiceServer struct {
	versionv1connect.UnimplementedVersionServiceHandler
}

func (versionServiceServer) GetVersion(
	context.Context,
	*connect.Request[versionv1.GetVersionRequest],
) (*connect.Response[versionv1.GetVersionResponse], error) {
	return connect.NewResponse(&versionv1.GetVersionResponse{
		Version: &versionv1.Version{
			Vcs:         versionv1.VcsType_VCS_TYPE_GIT,
			ServiceName: ServiceName,
			RefName:     RefName,
			Commit:      Commit,
			BuildTime:   BuildTime,
			FullVersion: FullVersion,
		},
	}), nil
}

func NewVersionServiceHandler() (string, http.Handler) {
	return versionv1connect.NewVersionServiceHandler(
		&versionServiceServer{},
		connect.WithInterceptors(serviceotel.DefaultServicesInterceptors()...))
}
