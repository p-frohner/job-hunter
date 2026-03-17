import { createTheme } from "@mui/material";

export const theme = createTheme({
	palette: {
		primary: {
			main: "#3d1f0a",
			light: "#6b3a1f",
			dark: "#1e0e04",
			contrastText: "#fff",
		},
		secondary: {
			main: "#000",
		},
	},
	components: {
		MuiPaper: {
			styleOverrides: {
				root: {
					backgroundColor: "#fdf5ee",
					border: "1px secondary solid",
					borderRadius: 16,
				},
			},
		},
		MuiBackdrop: {
			styleOverrides: {
				root: {
					backdropFilter: "blur(8px)",
					backgroundColor: "rgba(0, 0, 0, 0.4)",
				},
			},
		},
		MuiButtonBase: {
			defaultProps: {
				disableRipple: true,
			},
		},
		MuiButton: {
			defaultProps: {
				color: "primary",
			},
			styleOverrides: {
				root: {
					textTransform: "none",
					fontWeight: 600,
					borderRadius: 16,
				},
			},
		},
		MuiCssBaseline: {
			styleOverrides: {
				body: {
					backgroundColor: "#ffe0ba",
				},
			},
		},
		MuiDialog: {
			styleOverrides: {
				paper: {
					borderRadius: 10,
				},
			},
		},
		MuiDialogTitle: {
			styleOverrides: {
				root: {
					textAlign: "center",
					fontWeight: 600,
					fontSize: "1.5rem",
					padding: 24,
				},
			},
		},
		MuiLink: {
			defaultProps: {
				underline: "hover",
				color: "primary",
			},
			styleOverrides: {
				root: {
					fontWeight: 600,
					fontSize: "1.2rem",
				},
			},
		},
		MuiTextField: {
			defaultProps: {
				InputLabelProps: { shrink: true },
				variant: "outlined",
				size: "small",
			},
		},
		MuiInputLabel: {
			styleOverrides: {
				root: {
					position: "relative",
					transform: "none",
					marginBottom: "8px",
					fontSize: "0.875rem",
					fontWeight: 600,
					color: "#333",
					textTransform: "uppercase",
					"&.Mui-focused": {
						color: "#3d1f0a",
					},
				},
			},
		},
		MuiOutlinedInput: {
			styleOverrides: {
				root: {
					borderRadius: 16,
					"& legend": { display: "none" },
					"& .MuiOutlinedInput-notchedOutline": {
						top: "0",
						borderWidth: "3px",
					},
					"&:hover .MuiOutlinedInput-notchedOutline": {
						borderWidth: "3px",
					},
					"&.Mui-focused .MuiOutlinedInput-notchedOutline": {
						borderWidth: "3px",
					},
				},
			},
		},
		MuiTableCell: {
			styleOverrides: {
				head: {
					fontWeight: 600,
				},
			},
		},
	},
});
