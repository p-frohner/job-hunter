import DeleteForeverIcon from "@mui/icons-material/DeleteForever";
import { Box, Card, CardContent, IconButton, Link, Skeleton, Stack, Typography } from "@mui/material";
import { AnimatePresence, animate, motion, useMotionValue, useTransform } from "framer-motion";
import { memo, useRef, useState } from "react";
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

const SHORT_SWIPE = 40; // px: below this → snap back
const LONG_SWIPE = 120; // px: above this → auto-delete
const REVEAL_X = -40; // px: snap-to-reveal position

export const JobCard = memo(
	({ job, index, onDelete }: { job: Job; index: number; onDelete: () => void }) => {
		const color = job.source ? sourceColor(job.source) : "#6b3a1f";
		const [isLeaving, setIsLeaving] = useState(false);
		const dragOffsetX = useRef(0);
		const x = useMotionValue(0);
		const opacity = useTransform(x, [0, REVEAL_X], [0, 1]);
		// Fade out the left edge of the content as it's dragged left
		const maskImage = useTransform(
			x,
			[0, -24],
			[
				"linear-gradient(to right, rgba(0,0,0,1) 0px, rgba(0,0,0,1) 0px, rgba(0,0,0,1) 100%)",
				"linear-gradient(to right, rgba(0,0,0,0) 0px, rgba(0,0,0,1) 30px, rgba(0,0,0,1) 100%)",
			],
		);

		const triggerDelete = () => {
			animate(x, -500, {
				duration: 0.2,
				ease: "easeOut",
				onComplete: () => setIsLeaving(true),
			});
		};

		return (
			<AnimatePresence onExitComplete={onDelete}>
			{!isLeaving && (
			<motion.div
				style={{ overflow: "hidden" }}
				initial={{ opacity: 0, y: 12 }}
				animate={{ opacity: 1, y: 0 }}
				exit={{ height: 0, opacity: 0, y: 0, transition: { duration: 0.2, ease: "easeOut" } }}
				transition={{ duration: 0.3, ease: "easeOut", delay: Math.min(index, 10) * 0.04 }}
			>
				<Card elevation={0} sx={{ display: "flex", overflow: "hidden" }}>
					<Box
						sx={{
							width: 32,
							backgroundColor: color,
							display: "flex",
							alignItems: "center",
							justifyContent: "center",
							flexShrink: 0,
							py: 1,
						}}
					>
						{job.source && (
							<Typography
								sx={{
									writingMode: "vertical-lr",
									transform: "rotate(180deg)",
									color: "#fff",
									fontWeight: 700,
									fontSize: "0.65rem",
									letterSpacing: "0.05em",
									userSelect: "none",
									textTransform: "uppercase",
								}}
							>
								{sourceLabel(job.source)}
							</Typography>
						)}
					</Box>
					<CardContent
						sx={{ flex: 1, p: 1.5, pr: 0, "&:last-child": { pb: 1.5 }, position: "relative" }}
					>
						<motion.div style={{ maskImage, WebkitMaskImage: maskImage, overflow: "hidden" }}>
							<Stack direction="row" position="relative">
								<motion.div
									drag="x"
									style={{ x, width: "calc(100% - 25px)" }}
									dragConstraints={{ left: -300, right: 0 }}
									dragElastic={0.2}
									onDragEnd={(_, info) => {
										const dist = Math.abs(info.offset.x);
										dragOffsetX.current = dist;
										if (info.offset.x < 0) {
											if (dist >= LONG_SWIPE) {
												triggerDelete();
											} else if (dist >= SHORT_SWIPE) {
												animate(x, REVEAL_X, { type: "spring", stiffness: 300, damping: 30 });
											} else {
												animate(x, 0, { type: "spring", stiffness: 300, damping: 30 });
											}
										} else {
											animate(x, 0, { type: "spring", stiffness: 300, damping: 30 });
										}
									}}
								>
									<Link
										href={job.url}
										target="_blank"
										rel="noopener noreferrer"
										draggable={false}
										onClick={(e) => {
											// Prevent click on drag events
											if (dragOffsetX.current > 2) {
												e.preventDefault();
												dragOffsetX.current = 0;
											}
										}}
										sx={{ fontSize: "0.95rem", lineHeight: 1.3, display: "block" }}
									>
										{job.title}
									</Link>
									<Typography
										variant="body2"
										color="text.secondary"
										sx={{ mt: 0.5, fontSize: "0.8rem" }}
									>
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
										<Typography
											variant="caption"
											color="text.disabled"
											sx={{ mt: 0.5, display: "block" }}
										>
											{job.posted_at}
										</Typography>
									)}
								</motion.div>
								<motion.div
									style={{
										opacity,
										position: "relative",
										width: 40,
										right: "8px",
										display: "flex",
										alignItems: "center",
									}}
								>
									<IconButton
										aria-label="dismiss job"
										size="small"
										onClick={triggerDelete}
										sx={{ color: "text.secondary" }}
									>
										<DeleteForeverIcon fontSize="small" />
									</IconButton>
								</motion.div>
							</Stack>
						</motion.div>
					</CardContent>
				</Card>
			</motion.div>
			)}
			</AnimatePresence>
		);
	},
);

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
