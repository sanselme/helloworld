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

	"github.com/anselmes/util/pkg/cli"
	"github.com/anselmes/util/pkg/host"
	"github.com/anselmes/util/pkg/util"
	"github.com/anselmes/util/pkg/version"
	"github.com/sanselme/helloworld/internal/service"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func Command() *cobra.Command {
	host := host.Endpoint{}

	cmd := &cobra.Command{
		Use:     "helloworld",
		Short:   "Hello world",
		Version: version.GetVersion(),
		Run: func(cmd *cobra.Command, args []string) {
			lis, err := net.Listen("tcp", fmt.Sprintf(":%d", host.Port))
			if err != nil {
				util.CheckErr(err)
			}

			server := grpc.NewServer()
			service := service.NewServer()
			api.RegisterGreeterServiceServer(server, service)

			log.Println("Starting server on port", host.Port)
			log.Fatal(server.Serve(lis))
		},
	}

	cmd.Flags().IntVarP(&host.Port, "port", "p", 8080, "Port to listen on")

	return cmd
}

func Run() int {
	cmd := Command()
	cli.Run(cmd)
	return 0
}
