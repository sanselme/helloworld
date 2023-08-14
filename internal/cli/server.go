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
package cli

import (
	"fmt"
	"log"
	"net"

	api "github.com/sanselme/helloworld/api/v1alpha1"

	"github.com/anselmes/util/pkg/host"
	"github.com/anselmes/util/pkg/util"
	"github.com/sanselme/helloworld/internal/service"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func NewServerCommand() *cobra.Command {
	ep := host.Endpoint{}

	cmd := &cobra.Command{
		Use:     "service",
		Aliases: []string{"svc"},
		Short:   "helloworld service",
		Run: func(cmd *cobra.Command, args []string) {
			uri = fmt.Sprintf("%s:%d", ep.Address, ep.Port)

			// Create a TCP listener
			lis, err := net.Listen("tcp", uri)
			if err != nil {
				util.CheckErr(err)
			}

			// Create gRPC server object
			server := grpc.NewServer()

			// Create service objects
			service := service.NewServer()

			// Attach services to the server
			api.RegisterGreeterServiceServer(server, service)

			// Start serving
			log.Println("server listening on", uri)
			log.Fatal(server.Serve(lis))
		},
	}

	cmd.Flags().StringVar(&ep.Address, "address", "localhost", "Address to listen on")
	cmd.Flags().IntVarP(&ep.Port, "port", "p", 8080, "port to listen on")

	return cmd
}
