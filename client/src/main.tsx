import { MutationCache, QueryCache, QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { createRouter, RouterProvider } from "@tanstack/react-router";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { NotificationSnackbar } from "./components/NotificationSnackbar";
import { routeTree } from "./routeTree.gen"; // This file is auto-generated
import { showGlobalNotification } from "./utils/notificationHandler";

const router = createRouter({ routeTree, defaultViewTransition: true });

const queryClient = new QueryClient({
	defaultOptions: {
		queries: {
			staleTime: 1000 * 60 * 5, // Consider data fresh for 5 minutes
		},
	},
	queryCache: new QueryCache({
		onError: (error) => {
			const msg = error?.message || "An unexpected error occurred";
			showGlobalNotification(msg, "error");
		},
	}),
	mutationCache: new MutationCache({
		onError: (error) => {
			const msg = error?.message || "An unexpected error occurred";
			showGlobalNotification(msg, "error");
		},
	}),
});

const root = document.getElementById("root");

if (!root) {
	throw new Error("Failed to find the root element");
}

createRoot(root).render(
	<StrictMode>
		<QueryClientProvider client={queryClient}>
			<NotificationSnackbar />
			<RouterProvider router={router} />
			<ReactQueryDevtools initialIsOpen={false} />
		</QueryClientProvider>
	</StrictMode>,
);
