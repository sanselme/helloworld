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
	"net/http"

	events "github.com/cloudevents/sdk-go/v2"
	"github.com/sanselme/helloworld/pkg/host"
	"github.com/spf13/cobra"
)

type Event struct {
	Endpoint host.Endpoint
}

func (ev *Event) RunEvent(cmd *cobra.Command, _ []string) error {
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	proto, err := events.NewHTTP()
	if err != nil {
		return err
	}

	// Create a new client
	event, err := events.NewClient(proto)
	if err != nil {
		return err
	}

	// Start the receiver
	log.Println("listening on", ev.Endpoint.GetURI())
	if err := event.StartReceiver(ctx, receive); err != nil {
		return err
	}

	return nil
}

func NewEvent() *Event {
	return &Event{}
}

func receive(ctx context.Context, event events.Event) events.Result {
	log.Println(event)
	if event.Type() != "es.anselm.helloworld" {
		return events.NewHTTPResult(http.StatusBadRequest, "invalid type of %s", event.Type())
	}
	return events.NewHTTPResult(http.StatusOK, "ok")
}
