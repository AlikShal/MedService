import type { ReactNode } from "react";

type StatCardProps = {
	label: string;
	value: string;
	hint: string;
	icon: ReactNode;
	tone?: "lagoon" | "coral" | "sun";
};

const toneMap = {
	lagoon: "from-lagoon-500/20 to-lagoon-400/5 text-lagoon-300",
	coral: "from-coral-500/20 to-coral-400/5 text-coral-200",
	sun: "from-sun-500/20 to-sun-400/5 text-sun-200",
};

export function StatCard({ label, value, hint, icon, tone = "lagoon" }: StatCardProps) {
	return (
		<div
			className={`rounded-[1.75rem] border border-white/10 bg-gradient-to-br ${toneMap[tone]} p-5 shadow-panel backdrop-blur`}
		>
			<div className="mb-5 flex items-start justify-between gap-4">
				<div>
					<p className="text-xs uppercase tracking-[0.28em] text-slate-400">{label}</p>
					<p className="mt-3 font-display text-4xl text-white">{value}</p>
				</div>
				<div className="rounded-2xl border border-white/10 bg-white/5 p-3 text-current">{icon}</div>
			</div>
			<p className="text-sm text-slate-300">{hint}</p>
		</div>
	);
}
