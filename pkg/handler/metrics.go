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
*/package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/sanselme/helloworld/pkg/host"

	"contrib.go.opencensus.io/exporter/prometheus"
	obsclient "github.com/cloudevents/sdk-go/observability/opencensus/v2/client"
	"go.opencensus.io/examples/exporter"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.opencensus.io/zpages"
)

var MetricsEndpoint = host.Endpoint{Address: "localhost", Port: 9090}

func MetricsServer(
	printExporter *exporter.PrintExporter,
	traceSampler trace.Sampler,
	options prometheus.Options,
) {
	exp, err := prometheus.NewExporter(options)
	if err != nil {
		log.Fatalf("failed to create the stats exporter: %v", err)
	}

	server := http.NewServeMux()

	server.Handle("/metrics", exp)
	zpages.Handle(server, "/debug")

	view.RegisterExporter(exp)
	trace.RegisterExporter(printExporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: traceSampler})

	if err := view.Register(obsclient.LatencyView); err != nil {
		log.Fatalf("failed to register the views: %v", err)
	}

	view.SetReportingPeriod(2 * time.Second)
	log.Fatal(
		"failed metrics ListenAndServe ",
		http.ListenAndServe(MetricsEndpoint.GetURI(), server),
	)
}
