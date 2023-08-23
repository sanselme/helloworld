window.onload = function () {
	//<editor-fold desc="Changeable Configuration Block">

	// the following lines will be replaced by docker/configurator, when it runs in a docker-container
	window.ui = SwaggerUIBundle({
		urls: [
			{
				url: "v1alpha2/cloudevents.swagger.json",
				name: "v1alpha2/cloudevents.swagger.json",
			},
			{
				url: "v1alpha2/hello_world.swagger.json",
				name: "v1alpha2/hello_world.swagger.json",
			},
			{
				url: "v1alpha1/hello_world.swagger.json",
				name: "v1alpha1/hello_world.swagger.json",
			},
		],
		dom_id: "#swagger-ui",
		deepLinking: true,
		presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
		plugins: [SwaggerUIBundle.plugins.DownloadUrl],
		layout: "StandaloneLayout",
	});

	//</editor-fold>
};
