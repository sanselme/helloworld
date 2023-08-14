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
package service

import (
	"github.com/sanselme/helloworld/internal"
	"github.com/spf13/cobra"
)

func NewGatewayCommand() *cobra.Command {
	gw := internal.Gateway{}

	cmd := &cobra.Command{
		Use:     "gateway",
		Aliases: []string{"gw"},
		Short:   "helloworld gateway",
		Run: func(cmd *cobra.Command, args []string) {
			err := internal.RunGateway(cmd.Context(), gw)
			if err != nil {
				cmd.PrintErrln(err)
			}
		},
	}

	cmd.Flags().StringVar(&gw.Address, "address", "localhost", "Address to listen on")
	cmd.Flags().IntVarP(&gw.Port, "port", "p", 8081, "port to listen on")
	cmd.Flags().StringVar(&gw.Service.Address, "svc-addr", "localhost", "gRPC service address")
	cmd.Flags().IntVar(&gw.Service.Port, "svc-port", 8080, "gRPC service port")
	cmd.Flags().StringVar(&gw.OpenAPIDir, "oa-dir", "docs/openapi", "OpenAPI directory")

	return cmd
}
