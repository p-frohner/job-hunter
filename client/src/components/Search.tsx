import Masonry from "@mui/lab/Masonry";
import { Box, Typography } from "@mui/material";
import { useTransition } from "react";
import { useSearchStream } from "../hooks/useSearchStream";
import { JobCard, SkeletonCard } from "./JobCard";
import { SearchForm } from "./SearchForm";

export const Search = () => {
	const { search, jobs, loading, error, removeJob } = useSearchStream();
	const [isPending, startTransition] = useTransition();

	const handleSearch = (query: string, location?: string) => {
		startTransition(() => search(query, location));
	};

	return (
		<Box padding={3}>
			<SearchForm onSearch={handleSearch} loading={isPending || loading} />

			{error && (
				<Typography color="error" marginTop={2}>
					{error}
				</Typography>
			)}

			{(jobs.length > 0 || loading) && (
				<Box marginTop={3} sx={{ opacity: isPending ? 0.4 : 1, transition: "opacity 0.15s" }}>
					{jobs.length > 0 && (
						<Typography variant="subtitle2" color="text.secondary" gutterBottom>
							{jobs.length} result{jobs.length !== 1 ? "s" : ""}
							{loading && " — still searching…"}
						</Typography>
					)}
					<Masonry columns={{ xs: 1, sm: 2, md: 3 }} spacing={2}>
						{jobs.map((job, index) => (
							<JobCard key={`${job.source}-${job.id}`} job={job} index={index} onDelete={() => removeJob(job.id, job.source, job.url)} />
						))}
						{loading &&
							jobs.length === 0 &&
							["s0", "s1", "s2", "s3", "s4", "s5"].map((k) => <SkeletonCard key={k} />)}
					</Masonry>
				</Box>
			)}
		</Box>
	);
};
