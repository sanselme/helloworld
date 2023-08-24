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
	"net"
	"net/http"

	"github.com/anselmes/util/pkg/host"
	"github.com/anselmes/util/pkg/util"
	"github.com/google/uuid"
	api "github.com/sanselme/helloworld/api"
	v1alpha2 "github.com/sanselme/helloworld/api/v1alpha2"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type service struct {
	Endpoint host.Endpoint
}

func (s *service) SayHello(
	ctx context.Context,
	in *v1alpha2.SayHelloRequest,
) (*v1alpha2.SayHelloResponse, error) {
	event := &api.CloudEvent{
		Id:          uuid.New().String(),
		Source:      "https://github.com/sanselme/helloworld/api/v1alpha2",
		SpecVersion: "1.0",
		Type:        "es.anselm.helloworld",
		Data:        &api.CloudEvent_TextData{TextData: fmt.Sprintf("%s world!", in.GetName())},
	}
	log.Println(event)

	return &v1alpha2.SayHelloResponse{Message: event}, nil
}

func (s *service) RunService(cmd *cobra.Command, args []string) error {
	uri := fmt.Sprintf("%s:%d", s.Endpoint.Address, s.Endpoint.Port)

	// Create a TCP listener
	l, err := net.Listen("tcp", uri)
	if err != nil {
		return err
	}
	defer func() {
		if err := l.Close(); err != nil {
			util.CheckErr(err)
		}
	}()

	// Create gRPC server object
	server := grpc.NewServer()
	v1alpha2.RegisterGreeterServiceServer(server, NewService())
	go func() {
		defer server.GracefulStop()
		<-cmd.Context().Done()
	}()

	// Start serving
	log.Println("listening on", uri)
	if err := server.Serve(l); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func NewService() *service {
	return &service{}
}
