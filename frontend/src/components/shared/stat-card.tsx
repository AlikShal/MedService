import type { ReactNode } from "react";

type StatCardProps = {
	label: string;
	value: string;
	hint: string;
	icon: ReactNode;
	tone?: "lagoon" | "coral" | "sun";
};

const toneMap = {
	lagoon: "text-lagoon-400",
	coral: "text-coral-400",
	sun: "text-sun-400",
};

export function StatCard({ label, value, hint, icon, tone = "lagoon" }: StatCardProps) {
	return (
		<div className="lovable-panel p-5">
			<div className="mb-5 flex items-start justify-between gap-4">
				<div>
					<p className="lovable-eyebrow">{label}</p>
					<p className="mt-3 text-3xl font-bold leading-none text-white">{value}</p>
				</div>
				<div className={`rounded-md bg-lagoon-500/10 p-2 ${toneMap[tone]}`}>{icon}</div>
			</div>
			<p className="text-xs leading-5 text-slate-400">{hint}</p>
		</div>
	);
}
