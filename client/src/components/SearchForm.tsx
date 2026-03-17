import { Box, Button, CircularProgress, Stack, TextField } from "@mui/material";
import { useState } from "react";

type SearchFormProps = {
	onSearch: (query: string, location?: string) => void;
	loading: boolean;
};

export const SearchForm = ({ onSearch, loading }: SearchFormProps) => {
	const [query, setQuery] = useState("frontend");
	const [location, setLocation] = useState("budapest");

	const handleSubmit = (e: React.FormEvent) => {
		e.preventDefault();
		if (!query.trim()) {
			return;
		}
		onSearch(query.trim(), location.trim() || undefined);
	};

	return (
		<Box component="form" onSubmit={handleSubmit} maxWidth="md" margin="0 auto">
			<Stack direction={{ xs: "column", sm: "row" }} spacing={2} alignItems="flex-start">
				<TextField
					size="medium"
					label="Keywords"
					placeholder="e.g. senior engineer"
					value={query}
					onChange={(e) => setQuery(e.target.value)}
					required
					fullWidth
					sx={{ flexGrow: 1, minWidth: 240 }}
				/>
				<TextField
					size="medium"
					label="Location"
					placeholder="Remote"
					value={location}
					onChange={(e) => setLocation(e.target.value)}
					fullWidth
					sx={{ minWidth: 180 }}
				/>
			</Stack>
			<Box height="60px" pt={1}>
				<Button
					type="submit"
					variant="contained"
					disabled={loading || !query.trim()}
					fullWidth
					size="large"
					sx={{
						height: "100%",
					}}
				>
					{loading ? <CircularProgress size={22} color="inherit" /> : "Search"}
				</Button>
			</Box>
		</Box>
	);
};
