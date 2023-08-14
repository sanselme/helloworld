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
*/package cli

import (
	"fmt"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	api "github.com/sanselme/helloworld/api/v1alpha1"
	"github.com/sanselme/helloworld/docs"

	"github.com/anselmes/util/pkg/host"
	"github.com/anselmes/util/pkg/util"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGatewayCommand() *cobra.Command {
	ep := host.Endpoint{}
	svc := host.Endpoint{}

	cmd := &cobra.Command{
		Use:     "gateway",
		Aliases: []string{"gw"},
		Short:   "helloworld gateway",
		Run: func(cmd *cobra.Command, args []string) {
			uri = fmt.Sprintf("%s:%d", ep.Address, ep.Port)
			service := fmt.Sprintf("%s:%d", svc.Address, svc.Port)

			// Create a client connection to the gRPC server we just started
			// This is where the gRPC-Gateway proxies the requests
			conn, err := grpc.DialContext(
				cmd.Context(),
				service,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			if err != nil {
				util.CheckErr(err)
			}

			// Register generated routes to mux
			mux := runtime.NewServeMux()
			err = api.RegisterGreeterServiceHandler(cmd.Context(), mux, conn)
			if err != nil {
				util.CheckErr(err)
			}

			gw := &http.Server{
				Addr: uri,
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if strings.HasPrefix(r.URL.Path, "/api") {
						mux.ServeHTTP(w, r)
						return
					}
					openapi().ServeHTTP(w, r)
				}),
			}

			log.Println("gateway listening on", uri)
			log.Fatal(gw.ListenAndServe())
		},
	}

	cmd.Flags().StringVar(&ep.Address, "address", "localhost", "Address to listen on")
	cmd.Flags().IntVarP(&ep.Port, "port", "p", 8081, "port to listen on")
	cmd.Flags().StringVar(&svc.Address, "svc-addr", "localhost", "gRPC service address")
	cmd.Flags().IntVar(&svc.Port, "svc-port", 8080, "gRPC service port")

	return cmd
}

func openapi() http.Handler {
	err := mime.AddExtensionType(".svg", "image/svg+xml")
	if err != nil {
		util.CheckErr(err)
	}

	sfs, err := fs.Sub(docs.OpenAPI, "openapi")
	if err != nil {
		util.CheckErr(err)
	}

	return http.FileServer(http.FS(sfs))
}
