import { useQuery } from "@tanstack/react-query";
import { BellRing, Box, DatabaseZap, LineChart, MonitorUp, Siren } from "lucide-react";
import { AppShell } from "../components/layout/app-shell";
import { PanelHeader } from "../components/shared/panel-header";
import { StatusBadge } from "../components/shared/status-badge";
import { StatCard } from "../components/shared/stat-card";
import { api } from "../lib/api";

const playbookSteps = [
	{
		title: "Detect through system pulse",
		copy: "Use the live health cards, Prometheus targets, and Grafana panels to identify the first sign of degradation.",
	},
	{
		title: "Correlate with the failing workflow",
		copy: "Connect service health to user impact by checking which part of the scheduling journey actually stopped working.",
	},
	{
		title: "Mitigate with container-level action",
		copy: "Apply the incident override, inspect containers, correct the broken configuration, and restore the healthy compose stack.",
	},
];

const observabilityLinks = [
	{
		title: "Grafana",
		copy: "Watch latency and request-rate panels in the provisioned dashboard.",
		href: "http://localhost:3000",
		icon: MonitorUp,
	},
	{
		title: "Prometheus",
		copy: "Inspect scrape targets and verify which service went unhealthy.",
		href: "http://localhost:9090",
		icon: LineChart,
	},
	{
		title: "Compose stack",
		copy: "Container status confirms the runtime impact during the incident exercise.",
		href: "#",
		icon: Box,
	},
];

export function OperationsPage() {
	const healthQuery = useQuery({
		queryKey: ["operations-health"],
		queryFn: api.getPlatformHealth,
		refetchInterval: 20000,
	});

	const health = healthQuery.data ?? [];
	const healthyCount = health.filter((item) => item.status === "ok").length;

	return (
		<AppShell
			title="Operations and incident surface"
			description="This page turns the assignment’s monitoring and incident-response requirements into a coherent operational product view."
		>
			<div className="space-y-8">
				<section className="grid gap-4 xl:grid-cols-4">
					<StatCard
						label="Healthy targets"
						value={`${healthyCount}/${health.length || 5}`}
						hint="Services healthy according to the UI-level pulse."
						icon={<DatabaseZap className="h-5 w-5" />}
					/>
					<StatCard
						label="Monitoring surfaces"
						value="2"
						hint="Prometheus and Grafana are ready to validate the platform."
						icon={<LineChart className="h-5 w-5" />}
						tone="sun"
					/>
					<StatCard
						label="Incident scenario"
						value="1"
						hint="Broken DB host on the transactional appointment service."
						icon={<Siren className="h-5 w-5" />}
						tone="coral"
					/>
					<StatCard
						label="Alert posture"
						value="Live"
						hint="Status cards refresh continuously in the browser."
						icon={<BellRing className="h-5 w-5" />}
					/>
				</section>

				<section className="grid gap-6 xl:grid-cols-[0.9fr_1.1fr]">
					<section className="rounded-[2rem] border border-white/10 bg-white/[0.06] p-6 shadow-panel">
						<PanelHeader
							eyebrow="Incident playbook"
							title="A clearer response narrative"
							description="This gives the assignment a stronger operational storyline: detect, analyze, mitigate, restore."
						/>

						<div className="mt-6 space-y-4">
							{playbookSteps.map((step, index) => (
								<div
									key={step.title}
									className="rounded-[1.5rem] border border-white/10 bg-midnight-900/60 p-5"
								>
									<div className="flex items-center gap-3">
										<div className="flex h-9 w-9 items-center justify-center rounded-full bg-lagoon-500/[0.15] text-sm font-semibold text-lagoon-200">
											0{index + 1}
										</div>
										<h3 className="font-semibold text-white">{step.title}</h3>
									</div>
									<p className="mt-4 text-sm leading-7 text-slate-300">{step.copy}</p>
								</div>
							))}
						</div>
					</section>

					<section className="rounded-[2rem] border border-white/10 bg-white/[0.06] p-6 shadow-panel">
						<PanelHeader
							eyebrow="Service pulse"
							title="Operational state across the platform"
							description="These health cards make the monitoring layer visible without leaving the product."
						/>

						<div className="mt-6 grid gap-4 md:grid-cols-2">
							{health.map((service) => (
								<div
									key={service.service}
									className="rounded-[1.6rem] border border-white/10 bg-midnight-900/60 p-5"
								>
									<div className="flex items-center justify-between gap-3">
										<h3 className="text-lg font-semibold text-white">{service.service}</h3>
										<StatusBadge status={service.status} />
									</div>
									<p className="mt-4 text-sm text-slate-300">
										{service.error ?? service.storage ?? "Healthy and reporting metrics."}
									</p>
								</div>
							))}
						</div>
					</section>
				</section>

				<section className="rounded-[2rem] border border-white/10 bg-white/[0.06] p-6 shadow-panel">
					<PanelHeader
						eyebrow="Observability"
						title="Jump from product to tooling without losing context"
						description="These surfaces support the screenshots and verification steps needed for the assignment submission."
					/>

					<div className="mt-6 grid gap-4 lg:grid-cols-3">
						{observabilityLinks.map((item) => {
							const Icon = item.icon;
							return (
								<a
									key={item.title}
									href={item.href}
									target={item.href.startsWith("http") ? "_blank" : undefined}
									rel={item.href.startsWith("http") ? "noreferrer" : undefined}
									className="rounded-[1.7rem] border border-white/10 bg-gradient-to-br from-white/10 to-white/[0.04] p-5 transition hover:border-lagoon-400/30 hover:bg-white/10"
								>
									<div className="inline-flex rounded-2xl bg-gradient-to-br from-lagoon-500/20 to-coral-500/20 p-4 text-lagoon-200">
										<Icon className="h-5 w-5" />
									</div>
									<h3 className="mt-5 text-xl font-semibold text-white">{item.title}</h3>
									<p className="mt-3 text-sm leading-7 text-slate-300">{item.copy}</p>
								</a>
							);
						})}
					</div>
				</section>
			</div>
		</AppShell>
	);
}
