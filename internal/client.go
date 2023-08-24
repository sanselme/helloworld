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
	"log"

	v1alpha2 "github.com/sanselme/helloworld/api/v1alpha2"
	"github.com/sanselme/helloworld/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Dial(ctx context.Context, ep string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	return grpc.DialContext(
		ctx,
		ep,
		opts...,
	)
}

func RunClient(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	ep, err := cmd.Flags().GetString("endpoint")
	if err != nil {
		return err
	}

	conn, err := Dial(ctx, ep)
	if err != nil {
		errors.CheckErr(err)
	}
	go func() {
		<-ctx.Done()
		if err := conn.Close(); err != nil {
			errors.CheckErr(err)
		}
	}()

	client := v1alpha2.NewGreeterServiceClient(conn)
	req := &v1alpha2.SayHelloRequest{Name: args[0]}

	res, err := client.SayHello(ctx, req)
	if err != nil {
		errors.CheckErr(err)
	}

	log.Println(res.Message)
	return nil
}
