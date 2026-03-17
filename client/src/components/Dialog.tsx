import { Grow, Dialog as MuiDialog } from "@mui/material";
import { type MouseEvent, type ReactNode, useCallback, useState } from "react";

type Props = {
	open: boolean;
	onClose: () => void;
	anchorEl?: HTMLElement | null;
	children: ReactNode;
};

/**
 * A dialog component with a grow animation that scales from an anchor element.
 * Use with `useDialog` hook for easy state management.
 *
 * @example
 * ```tsx
 * const { isOpen, anchorEl, open, close } = useDialog();
 *
 * <Button onClick={open}>Open Dialog</Button>
 * <Dialog open={isOpen} anchorEl={anchorEl} onClose={close}>
 *   <DialogTitle>Title</DialogTitle>
 *   <DialogContent>Content here</DialogContent>
 * </Dialog>
 * ```
 */
export const Dialog = ({ open, onClose, anchorEl, children }: Props) => {
	const getTransformOrigin = () => {
		if (!anchorEl) {
			return undefined;
		}
		const rect = anchorEl.getBoundingClientRect();

		return `${rect.left + rect.width / 2}px ${rect.top + rect.height}px`;
	};

	return (
		<MuiDialog
			disableRestoreFocus
			open={open}
			fullWidth
			onClose={onClose}
			slots={{ transition: Grow }}
			slotProps={{
				transition: { style: { transformOrigin: getTransformOrigin() }, timeout: 300 },
				backdrop: { sx: { transition: "opacity 300ms ease-in-out !important" } },
			}}
		>
			{children}
		</MuiDialog>
	);
};

export const useDialog = () => {
	const [isOpen, setIsOpen] = useState(false);
	const [anchorEl, setAnchorEl] = useState<HTMLElement | null>(null);

	const open = useCallback((e: MouseEvent<HTMLElement>) => {
		setAnchorEl(e.currentTarget);
		setIsOpen(true);
	}, []);

	const close = useCallback(() => {
		setIsOpen(false);
	}, []);

	return { isOpen, anchorEl, open, close };
};
