import { defineConfig } from "orval";

export default defineConfig({
	"job-hunter-api": {
		input: "../openapi.yaml",
		output: {
		target: "./src/api.ts",
			client: "react-query",
			httpClient: "fetch",
			override: {
				fetch: {
					includeHttpResponseReturnType: false,
				},
				mutator: {
					path: "./src/custom-fetch.ts",
					name: "customFetch",
				},
			},
		},
	},
	"job-hunter-zod": {
		input: "../openapi.yaml",
		output: {
			target: "./src/api.zod.ts",
			client: "zod",
		},
	},
});
