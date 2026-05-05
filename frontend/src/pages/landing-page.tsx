import { useQuery } from "@tanstack/react-query";
import { Activity, ArrowRight, CloudCog, HeartPulse, Shield, Stethoscope } from "lucide-react";
import { Link } from "react-router-dom";
import { AuthCard } from "../components/auth/auth-card";
import { StatusBadge } from "../components/shared/status-badge";
import { api } from "../lib/api";
import { useAuth } from "../providers/auth-provider";

const features = [
	{
		icon: Shield,
		title: "Role-aware access",
		copy: "Patients book care without admin friction, while clinic staff retain control over doctors and appointment workflow.",
	},
	{
		icon: HeartPulse,
		title: "Clinical-grade experience",
		copy: "The workspace is organized around care operations, not generic CRUD screens, so the flow feels intentional and production-ready.",
	},
	{
		icon: CloudCog,
		title: "Observability built in",
		copy: "Prometheus, Grafana, incident simulation, and container-first deployment are surfaced as first-class parts of the product.",
	},
];

export function LandingPage() {
	const { user } = useAuth();
	const healthQuery = useQuery({
		queryKey: ["landing-health"],
		queryFn: api.getPlatformHealth,
		refetchInterval: 20000,
	});

	const health = healthQuery.data ?? [];

	return (
		<div className="min-h-screen bg-[#070d19] text-white">
			<div className="mx-auto max-w-[925px] px-6 pb-16 pt-1">
				<header className="flex h-12 items-center justify-between rounded-xl border border-slate-800 bg-[#0d1422] px-4">
					<div className="flex items-center gap-3">
						<div className="flex h-8 w-8 items-center justify-center rounded-lg bg-lagoon-500/10 text-lagoon-400">
							<Stethoscope className="h-4 w-4" />
						</div>
						<div>
							<p className="text-sm font-extrabold leading-none text-white">MedSync</p>
							<p className="mt-1 text-[0.58rem] font-bold uppercase tracking-[0.32em] text-slate-500">
								Medical scheduling platform
							</p>
						</div>
					</div>

					<nav className="hidden items-center gap-8 text-sm font-medium text-slate-400 md:flex">
						<a href="#experience" className="transition hover:text-white">
							Experience
						</a>
						<a href="#topology" className="transition hover:text-white">
							Topology
						</a>
						<a href="#pulse" className="transition hover:text-white">
							Live pulse
						</a>
						{user ? (
							<Link to="/workspace" className="lovable-button px-4 py-2">
								Open workspace
							</Link>
						) : (
							<a href="#auth" className="lovable-button px-4 py-2">
								Open workspace
							</a>
						)}
					</nav>
				</header>

				<section className="grid gap-8 pb-16 pt-12 lg:grid-cols-[1.45fr_0.95fr] lg:items-start">
					<div className="space-y-7">
						<div className="space-y-6">
							<div className="inline-flex items-center gap-2 rounded-full border border-lagoon-500/35 bg-lagoon-500/10 px-4 py-1.5 text-[0.68rem] font-extrabold uppercase tracking-[0.34em] text-lagoon-400">
								<Activity className="h-3.5 w-3.5" />
								Containerized clinical operations
							</div>

							<div className="max-w-[600px] space-y-5">
								<h1 className="text-[2.65rem] font-extrabold leading-[1.02] tracking-[-0.02em] text-white sm:text-[3.35rem]">
									A premium care operations front end{" "}
									<span className="text-lagoon-500">for your microservices stack.</span>
								</h1>
								<p className="max-w-[590px] text-base font-medium leading-7 text-slate-400">
									MedSync looks and behaves like a real product: polished onboarding, role-aware
									workspaces, and an operations surface that respects the assignment&apos;s DevOps
									scope.
								</p>
							</div>

							<div className="flex flex-wrap gap-3">
								{user ? (
									<Link to="/workspace" className="lovable-button min-w-40">
										Go to workspace
										<ArrowRight className="h-4 w-4" />
									</Link>
								) : (
									<a href="#auth" className="lovable-button min-w-40">
										Go to workspace
										<ArrowRight className="h-4 w-4" />
									</a>
								)}
								<a
									href="#pulse"
									className="inline-flex min-w-36 items-center justify-center rounded-md border border-lagoon-500/45 bg-transparent px-5 py-2.5 text-sm font-semibold text-lagoon-400 transition hover:bg-lagoon-500/10"
								>
									View service pulse
								</a>
							</div>
						</div>

						<div className="grid gap-3 sm:grid-cols-3">
							{[
								["User flows", "2", "Dedicated patient and admin journeys."],
								["Gateway", "1", "Nginx serves the UI and routes API traffic."],
								["Observability", "24/7", "Prometheus and Grafana built into the flow."],
							].map(([label, value, copy]) => (
								<div key={label} className="lovable-panel p-4">
									<p className="lovable-eyebrow">{label}</p>
									<p className="mt-3 text-2xl font-extrabold leading-none text-white">{value}</p>
									<p className="mt-2 text-xs leading-5 text-slate-400">{copy}</p>
								</div>
							))}
						</div>
					</div>

					<div id="auth" className="lg:sticky lg:top-8">
						<AuthCard />
					</div>
				</section>

				<section id="experience" className="grid gap-4 py-1 lg:grid-cols-3">
					{features.map((feature) => {
						const Icon = feature.icon;
						return (
							<article key={feature.title} className="lovable-panel p-5">
								<div className="flex h-10 w-10 items-center justify-center rounded-lg bg-lagoon-500/10 text-lagoon-400">
									<Icon className="h-5 w-5" />
								</div>
								<h2 className="mt-5 text-lg font-extrabold text-white">{feature.title}</h2>
								<p className="mt-3 text-sm font-medium leading-6 text-slate-400">{feature.copy}</p>
							</article>
						);
					})}
				</section>

				<section id="topology" className="mt-8 grid gap-4 lg:grid-cols-[0.85fr_1.15fr]">
					<div className="lovable-panel p-5">
						<p className="lovable-eyebrow text-lagoon-400">Topology</p>
						<h2 className="mt-3 text-xl font-extrabold text-white">Microservice product surface</h2>
						<p className="mt-3 text-sm leading-6 text-slate-400">
							The gateway routes a React UI over authentication, patients, doctors, appointments, chat,
							and monitoring.
						</p>
					</div>

					<div className="lovable-panel grid gap-3 p-5 sm:grid-cols-5">
						{["Auth", "Patients", "Doctors", "Appointments", "Chat"].map((item) => (
							<div key={item} className="rounded-lg border border-slate-800 bg-[#080f1d] p-3 text-sm font-bold text-slate-200">
								{item}
							</div>
						))}
					</div>
				</section>

				<section id="pulse" className="grid gap-4 py-8 sm:grid-cols-2 lg:grid-cols-5">
					{health.map((service) => (
						<div key={service.service} className="lovable-panel p-4">
							<div className="flex items-start justify-between gap-3">
								<div>
									<p className="lovable-eyebrow">Live</p>
									<h3 className="mt-2 text-sm font-extrabold text-white">{service.service}</h3>
								</div>
								<StatusBadge status={service.status} className="px-2 py-0.5 text-[0.55rem]" />
							</div>
							<p className="mt-3 text-xs leading-5 text-slate-400">
								{service.error ?? service.storage ?? "Healthy and ready."}
							</p>
						</div>
					))}
				</section>
			</div>
		</div>
	);
}
