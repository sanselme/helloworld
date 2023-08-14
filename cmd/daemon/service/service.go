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
	"context"

	"github.com/anselmes/util/pkg/host"
	"github.com/sanselme/helloworld/internal"
	"github.com/spf13/cobra"
)

func NewServiceCommand() *cobra.Command {
	ctx := context.Background()
	ep := host.Endpoint{}

	cmd := &cobra.Command{
		Use:     "service",
		Aliases: []string{"svc"},
		Short:   "helloworld service",
		Run: func(cmd *cobra.Command, args []string) {
			err := internal.RunService(ctx, ep)
			if err != nil {
				cmd.PrintErrln(err)
			}
		},
	}

	cmd.Flags().StringVar(&ep.Address, "address", "localhost", "Address to listen on")
	cmd.Flags().IntVarP(&ep.Port, "port", "p", 8080, "port to listen on")

	return cmd
}
