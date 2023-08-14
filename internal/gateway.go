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
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	api "github.com/sanselme/helloworld/api/v1alpha1"
	"github.com/sanselme/helloworld/pkg/handler"

	"github.com/anselmes/util/pkg/host"
	"github.com/anselmes/util/pkg/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Gateway struct {
	Address    string
	Port       int
	Service    host.Endpoint
	OpenAPIDir string
}

func RunGateway(ctx context.Context, gw Gateway) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	svc := fmt.Sprintf("%s:%d", gw.Service.Address, gw.Service.Port)

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := dial(ctx, svc)
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		if err := conn.Close(); err != nil {
			util.CheckErr(err)
		}
	}()

	mux := runtime.NewServeMux()

	// Register generated routes to mux
	err = api.RegisterGreeterServiceHandler(ctx, mux, conn)
	if err != nil {
		util.CheckErr(err)
	}

	uri := fmt.Sprintf("%s:%d", gw.Address, gw.Port)
	s := &http.Server{
		Addr: uri,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/openapiv2/") {
				handler.OpenAPIServer(gw.OpenAPIDir).ServeHTTP(w, r)
				return
			}

			if strings.HasPrefix(r.URL.Path, "/healthz") {
				handler.HealthServer(conn).ServeHTTP(w, r)
				return
			}

			handler.AlloCORS(mux)
		}),
	}
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

func dial(ctx context.Context, ep string) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		ctx,
		ep,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}
