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
	"os"

	"github.com/anselmes/util/pkg/util"
	"github.com/anselmes/util/pkg/version"
	"github.com/sanselme/helloworld/cmd/daemon/service"
	"github.com/spf13/cobra"
)

func command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "helloworld",
		Short:   "Hello world",
		Version: version.GetVersion(),
	}

	cmd.AddCommand(service.NewServiceCommand())
	cmd.AddCommand(service.NewGatewayCommand())

	return cmd
}

func execute() int {
	err := command().Execute()
	if err != nil {
		util.CheckErr(err)
	}
	return 0
}

func main() {
	os.Exit(execute())
}
