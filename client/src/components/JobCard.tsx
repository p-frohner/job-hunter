import { Card, CardContent, Chip, Link, Skeleton, Stack, Typography } from "@mui/material";
import { memo } from "react";
import type { Job } from "../api";

const SOURCE_LABELS: Record<string, string> = {
	linkedin: "LinkedIn",
	nofluffjobs: "NoFluffJobs",
	professionhu: "Profession.hu",
};

const SOURCE_COLORS: Record<string, string> = {
	linkedin: "#c76b1a",
	nofluffjobs: "#9e3a2b",
	professionhu: "#6b7228",
};

const sourceLabel = (source: string) => SOURCE_LABELS[source] ?? source;
const sourceColor = (source: string) => SOURCE_COLORS[source] ?? "#6b3a1f";

export const JobCard = memo(({ job, index }: { job: Job; index: number }) => {
	const color = job.source ? sourceColor(job.source) : "#6b3a1f";
	return (
		<Card
			elevation={4}
			sx={{
				borderLeft: `32px solid ${color}`,
				animation: "fadeSlideIn 0.3s ease both",
				animationDelay: `${Math.min(index, 10) * 40}ms`,
			}}
		>
			<CardContent sx={{ p: 1.5, "&:last-child": { pb: 1.5 } }}>
				<Stack direction="row" justifyContent="space-between" alignItems="flex-start" gap={1}>
					<Link
						href={job.url}
						target="_blank"
						rel="noopener noreferrer"
						sx={{ fontSize: "0.95rem", lineHeight: 1.3, display: "block" }}
					>
						{job.title}
					</Link>
					{job.source && (
						<Chip
							label={sourceLabel(job.source)}
							size="small"
							sx={{
								flexShrink: 0,
								bgcolor: color,
								color: "#fff",
								fontWeight: 600,
								fontSize: "0.7rem",
								height: 20,
							}}
						/>
					)}
				</Stack>
				<Typography variant="body2" color="text.secondary" sx={{ mt: 0.5, fontSize: "0.8rem" }}>
					{job.company}
					{job.location ? ` · ${job.location}` : ""}
				</Typography>
				{job.snippet && (
					<Typography
						variant="body2"
						color="text.secondary"
						sx={{
							mt: 0.75,
							fontSize: "0.78rem",
							display: "-webkit-box",
							WebkitLineClamp: 2,
							WebkitBoxOrient: "vertical",
							overflow: "hidden",
						}}
					>
						{job.snippet}
					</Typography>
				)}
				{job.posted_at && (
					<Typography variant="caption" color="text.disabled" sx={{ mt: 0.5, display: "block" }}>
						{job.posted_at}
					</Typography>
				)}
			</CardContent>
		</Card>
	);
});

export const SkeletonCard = () => (
	<Card elevation={2}>
		<CardContent sx={{ p: 1.5, "&:last-child": { pb: 1.5 } }}>
			<Skeleton variant="text" width="80%" height={20} />
			<Skeleton variant="text" width="50%" height={16} sx={{ mt: 0.5 }} />
			<Skeleton variant="text" width="100%" height={14} sx={{ mt: 0.75 }} />
			<Skeleton variant="text" width="90%" height={14} />
		</CardContent>
	</Card>
);
