type NotificationHandler = (message: string, severity: "error" | "success") => void;

let notificationHandler: NotificationHandler | null = null;

export const setNotificationHandler = (handler: NotificationHandler) => {
	notificationHandler = handler;
};

export const showGlobalNotification = (
	message: string,
	severity: "error" | "success" = "error",
) => {
	if (notificationHandler) {
		notificationHandler(message, severity);
	} else {
		console.warn("Notification handler not initialized yet:", message);
	}
};
