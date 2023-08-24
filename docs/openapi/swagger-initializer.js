window.onload = function () {
	//<editor-fold desc="Changeable Configuration Block">

	// the following lines will be replaced by docker/configurator, when it runs in a docker-container
	window.ui = SwaggerUIBundle({
		urls: [
			{ url: "cloudevents.swagger.json", name: "cloudevents.swagger.json" },
			{
				url: "v1alpha2/services.swagger.json",
				name: "v1alpha2/services.swagger.json",
			},
			{
				url: "v1alpha2/messages.swagger.json",
				name: "v1alpha2/messages.swagger.json",
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
