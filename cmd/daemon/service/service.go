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

func NewServiceCommand() *cobra.Command {
	svc := internal.NewService()

	cmd := &cobra.Command{
		Use:     "service",
		Aliases: []string{"svc"},
		Short:   "helloworld service",
		RunE:    svc.RunService,
	}

	cmd.Flags().StringVar(&svc.Endpoint.Address, "address", "localhost", "Address to listen on")
	cmd.Flags().IntVarP(&svc.Endpoint.Port, "port", "p", 8080, "port to listen on")

	return cmd
}
