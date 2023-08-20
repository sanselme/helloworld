/*
Copyright (c) 2023 Schubert Anselme <schubert@anselm.es>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	api "github.com/sanselme/helloworld/api/v1alpha1"
	"github.com/sanselme/helloworld/pkg/handler"
	"github.com/spf13/cobra"

	"github.com/anselmes/util/pkg/host"
	"github.com/anselmes/util/pkg/util"
	"google.golang.org/grpc"
)

type (
	GatewayOption func(*gateway)
	gateway       struct {
		Endpoint   host.Endpoint
		Service    host.Endpoint
		OpenAPIDir string
	}
)

func (gw *gateway) RunGateway(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	svc := fmt.Sprintf("%s:%d", gw.Service.Address, gw.Service.Port)

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := Dial(ctx, svc, grpc.WithBlock())
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		if err := conn.Close(); err != nil {
			util.CheckErr(err)
		}
	}()

	// Register generated routes
	rt := runtime.NewServeMux()
	err = api.RegisterGreeterServiceHandler(ctx, rt, conn)
	if err != nil {
		return err
	}

	// Register custom routes
	mux := http.NewServeMux()
	mux.Handle("/", handler.AllowCORS(rt))
	mux.HandleFunc("/openapiv2/", handler.OpenAPIServer(gw.OpenAPIDir))
	mux.HandleFunc("/healthz", handler.HealthServer(conn))

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	uri := fmt.Sprintf("%s:%d", gw.Endpoint.Address, gw.Endpoint.Port)
	s := &http.Server{Addr: uri, Handler: mux}
	go func() {
		<-ctx.Done()
		log.Println("shutting down...")
		if err := s.Shutdown(context.Background()); err != nil {
			util.CheckErr(err)
		}
	}()

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	log.Println("listening on", uri)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func NewGateway(opts ...GatewayOption) *gateway {
	gw := &gateway{}
	for _, op := range opts {
		op(gw)
	}
	return gw
}
