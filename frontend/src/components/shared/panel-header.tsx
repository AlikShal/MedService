import { cn } from "../../lib/cn";

type PanelHeaderProps = {
	eyebrow: string;
	title: string;
	description?: string;
	className?: string;
};

export function PanelHeader({ eyebrow, title, description, className }: PanelHeaderProps) {
	return (
		<div className={cn("space-y-2", className)}>
			<p className="text-[0.7rem] font-semibold uppercase tracking-[0.35em] text-lagoon-400/90">
				{eyebrow}
			</p>
			<h2 className="font-display text-2xl text-white sm:text-3xl">{title}</h2>
			{description ? <p className="max-w-2xl text-sm text-slate-300">{description}</p> : null}
		</div>
	);
}
