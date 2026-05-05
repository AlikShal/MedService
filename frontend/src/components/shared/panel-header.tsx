import { cn } from "../../lib/cn";

type PanelHeaderProps = {
	eyebrow: string;
	title: string;
	description?: string;
	className?: string;
};

export function PanelHeader({ eyebrow, title, description, className }: PanelHeaderProps) {
	return (
		<div className={cn("space-y-1.5", className)}>
			<p className="lovable-eyebrow">{eyebrow}</p>
			<h2 className="text-base font-bold leading-tight text-white">{title}</h2>
			{description ? <p className="max-w-2xl text-xs leading-5 text-slate-400">{description}</p> : null}
		</div>
	);
}
