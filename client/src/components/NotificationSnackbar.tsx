import { Alert, Snackbar } from "@mui/material";
import { useEffect, useState } from "react";

import { setNotificationHandler } from "../utils/notificationHandler";

export const NotificationSnackbar = () => {
	const [open, setOpen] = useState(false);
	const [config, setConfig] = useState({ message: "", severity: "error" as "error" | "success" });

	useEffect(() => {
		setNotificationHandler((message, severity) => {
			setConfig({ message, severity });
			setOpen(true);
		});
	}, []);

	return (
		<Snackbar
			open={open}
			autoHideDuration={5000}
			onClose={() => setOpen(false)}
			anchorOrigin={{ vertical: "bottom", horizontal: "center" }}
		>
			<Alert
				severity={config.severity}
				variant="filled"
				sx={{
					width: "100%",
					bgcolor: config.severity === "error" ? "error.main" : "warning.main",
					backdropFilter: "blur(8px)",
				}}
			>
				{config.message}
			</Alert>
		</Snackbar>
	);
};
