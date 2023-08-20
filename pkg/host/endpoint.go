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
package host

import "fmt"

type Endpoint struct {
	Address string
	Path    string
	Port    int
	Scheme  string
}

func (e *Endpoint) GetURI() string {
	return fmt.Sprintf("%s:%d", e.Address, e.Port)
}

func (e *Endpoint) GetURL() string {
	if e.Scheme == "" {
		e.Scheme = "http"
	}

	if e.Path == "" {
		e.Path = "/"
	}

	return fmt.Sprintf("%s://%s:%d%s", e.Scheme, e.Address, e.Port, e.Path)
}
