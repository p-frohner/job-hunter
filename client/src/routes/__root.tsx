import { Box, Container, CssBaseline, Divider, GlobalStyles, ThemeProvider } from "@mui/material";
import { createRootRoute, Outlet } from "@tanstack/react-router";
import { ErrorBoundary } from "react-error-boundary";

import logo from "/logo.png";
import { theme } from "../themeProvider";

export const Route = createRootRoute({
	component: () => {
		return (
			<ErrorBoundary FallbackComponent={({ error }) => <pre>{error.message}</pre>}>
				<CssBaseline />
				<ThemeProvider theme={theme}>
					<GlobalStyles
						styles={{
							body: {
								backgroundColor: "#ffe0ba",
							},
						}}
					/>
					<Container maxWidth="lg">
						<Box>
							<Box sx={{ backgroundColor: "primary.paper", viewTransitionName: "main-content" }}>
								<Box textAlign="center" sx={{ backgroundColor: "#ffe0ba" }}>
									<Logo />
								</Box>
								<Divider
									sx={{
										height: 8,
										borderBottom: "none",
										backgroundColor: "primary.dark",
									}}
								/>
								<Outlet />
							</Box>
						</Box>
					</Container>
				</ThemeProvider>
			</ErrorBoundary>
		);
	},
});

const Logo = () => <Box component="img" height={{ xs: 200, sm: 320 }} src={logo} margin={1} />;
