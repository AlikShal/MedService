import { cn } from "../../lib/cn";
import type { AppointmentStatus } from "../../types/domain";

type StatusBadgeProps = {
	status: AppointmentStatus | "ok" | "degraded" | "down";
	className?: string;
};

const palette: Record<string, string> = {
	new: "border-slate-600 bg-slate-800/60 text-slate-200",
	in_progress: "border-lagoon-500/40 bg-lagoon-500/10 text-lagoon-300",
	done: "border-emerald-400/40 bg-emerald-500/10 text-emerald-300",
	cancelled: "border-rose-400/40 bg-rose-500/10 text-rose-300",
	ok: "border-emerald-400/40 bg-emerald-500/10 text-emerald-300",
	degraded: "border-sun-400/40 bg-sun-500/10 text-sun-300",
	down: "border-rose-400/40 bg-rose-500/10 text-rose-300",
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
