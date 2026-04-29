import { cn } from "../../lib/cn";
import type { AppointmentStatus } from "../../types/domain";

type StatusBadgeProps = {
	status: AppointmentStatus | "ok" | "degraded" | "down";
	className?: string;
};

const palette: Record<string, string> = {
	new: "border-sun-400/30 bg-sun-500/[0.15] text-sun-200",
	in_progress: "border-lagoon-400/30 bg-lagoon-500/[0.15] text-lagoon-100",
	done: "border-emerald-400/30 bg-emerald-500/[0.15] text-emerald-100",
	cancelled: "border-coral-400/30 bg-coral-500/[0.15] text-coral-100",
	ok: "border-emerald-400/30 bg-emerald-500/[0.15] text-emerald-100",
	degraded: "border-sun-400/30 bg-sun-500/[0.15] text-sun-200",
	down: "border-coral-400/30 bg-coral-500/[0.15] text-coral-100",
};

export function StatusBadge({ status, className }: StatusBadgeProps) {
	return (
		<span
			className={cn(
				"inline-flex items-center rounded-full border px-3 py-1 text-[0.72rem] font-semibold uppercase tracking-[0.22em]",
				palette[status],
				className,
			)}
		>
			{String(status).replaceAll("_", " ")}
		</span>
	);
}
