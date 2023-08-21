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
package main

import (
	"context"
	"log"
	"os"

	api "github.com/sanselme/helloworld/api/v1alpha2"
	"github.com/sanselme/helloworld/internal"

	"github.com/sanselme/helloworld/pkg/errors"
	"github.com/sanselme/helloworld/pkg/version"
	"github.com/spf13/cobra"
)

func main() {
	var endpoint string

	cmd := &cobra.Command{
		Use:     "helloctl",
		Short:   "Hello Client",
		Version: version.GetVersion(),
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			conn, err := internal.Dial(ctx, endpoint)
			if err != nil {
				errors.CheckErr(err)
			}
			go func() {
				<-ctx.Done()
				if err := conn.Close(); err != nil {
					errors.CheckErr(err)
				}
			}()

			client := api.NewGreeterServiceClient(conn)
			req := &api.SayHelloRequest{Name: args[0]}

			res, err := client.SayHello(ctx, req)
			if err != nil {
				errors.CheckErr(err)
			}

			log.Println(res.Message)
		},
	}

	cmd.Flags().StringVarP(&endpoint, "endpoint", "e", "localhost:8080", "Endpoint to connect to")

	err := cmd.Execute()
	if err != nil {
		errors.CheckErr(err)
	}

	os.Exit(0)
}
