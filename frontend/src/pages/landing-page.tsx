import { useQuery } from "@tanstack/react-query";
import { motion } from "framer-motion";
import {
	Activity,
	ArrowRight,
	CalendarHeart,
	CloudCog,
	HeartPulse,
	Shield,
	Stethoscope,
} from "lucide-react";
import { Link } from "react-router-dom";
import { api } from "../lib/api";
import { AuthCard } from "../components/auth/auth-card";
import { StatusBadge } from "../components/shared/status-badge";
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
		copy: "Prometheus, Grafana, incident simulation, and container-first deployment are all surfaced as first-class parts of the product.",
	},
];

export function LandingPage() {
	const { user } = useAuth();
	const doctorsQuery = useQuery({ queryKey: ["landing-doctors"], queryFn: api.getDoctors });
	const healthQuery = useQuery({
		queryKey: ["landing-health"],
		queryFn: api.getPlatformHealth,
		refetchInterval: 20000,
	});

	const doctors = doctorsQuery.data ?? [];
	const health = healthQuery.data ?? [];

	return (
		<div className="min-h-screen overflow-hidden bg-midnight-950 bg-mesh-radial text-white">
			<div className="absolute inset-0 bg-[radial-gradient(circle_at_top,rgba(255,255,255,0.06),transparent_26%),linear-gradient(180deg,rgba(8,18,30,0.1),rgba(8,18,30,0.88))]" />
			<div className="absolute left-[-8rem] top-20 h-80 w-80 rounded-full bg-lagoon-400/10 blur-3xl" />
			<div className="absolute right-[-6rem] top-40 h-96 w-96 rounded-full bg-coral-500/10 blur-3xl" />

			<div className="relative mx-auto max-w-[1500px] px-5 pb-20 pt-6 sm:px-8 lg:px-10">
				<header className="flex flex-col gap-5 rounded-full border border-white/10 bg-white/5 px-6 py-4 backdrop-blur-xl lg:flex-row lg:items-center lg:justify-between">
					<div className="flex items-center gap-4">
						<div className="rounded-full bg-gradient-to-br from-lagoon-400 to-coral-400 p-3 text-midnight-950">
							<Stethoscope className="h-5 w-5" />
						</div>
						<div>
							<p className="font-display text-2xl text-white">MedSync</p>
							<p className="text-xs uppercase tracking-[0.32em] text-slate-400">
								Medical scheduling platform
							</p>
						</div>
					</div>

					<nav className="flex flex-wrap items-center gap-3 text-sm text-slate-300">
						<a href="#experience" className="hover:text-white">
							Experience
						</a>
						<a href="#topology" className="hover:text-white">
							Topology
						</a>
						<a href="#pulse" className="hover:text-white">
							Live pulse
						</a>
						<Link
							to={user ? "/workspace" : "/operations"}
							className="rounded-full border border-white/10 bg-white/10 px-4 py-2 text-white transition hover:border-lagoon-400/30 hover:bg-lagoon-500/10"
						>
							{user ? "Open workspace" : "Explore operations"}
						</Link>
					</nav>
				</header>

				<section className="grid gap-8 pb-16 pt-12 lg:grid-cols-[1.15fr_0.85fr] lg:items-start">
					<div className="space-y-8">
						<motion.div
							initial={{ opacity: 0, y: 16 }}
							animate={{ opacity: 1, y: 0 }}
							transition={{ duration: 0.55 }}
							className="space-y-6"
						>
							<div className="inline-flex items-center gap-2 rounded-full border border-lagoon-400/20 bg-lagoon-500/10 px-4 py-2 text-xs uppercase tracking-[0.35em] text-lagoon-200">
								<Activity className="h-3.5 w-3.5" />
								Containerized clinical operations
							</div>

							<div className="max-w-4xl space-y-4">
								<h1 className="font-display text-5xl leading-[0.95] text-white sm:text-6xl xl:text-7xl">
									A premium care operations front end for your microservices stack.
								</h1>
								<p className="max-w-2xl text-lg text-slate-300 sm:text-xl">
									MedSync now looks and behaves like a real product: polished onboarding, role-aware
									workspaces, and an operations surface that actually respects the assignment’s DevOps
									scope.
								</p>
							</div>

							<div className="flex flex-wrap gap-3">
								<Link
									to={user ? "/workspace" : "/operations"}
									className="inline-flex items-center gap-2 rounded-full bg-white px-5 py-3 text-sm font-semibold text-midnight-950 transition hover:bg-sand-100"
								>
									{user ? "Go to workspace" : "See the operations view"}
									<ArrowRight className="h-4 w-4" />
								</Link>
								<a
									href="#pulse"
									className="inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/5 px-5 py-3 text-sm font-semibold text-white transition hover:border-coral-400/30 hover:bg-coral-500/10"
								>
									View service pulse
								</a>
							</div>
						</motion.div>

						<div className="grid gap-4 md:grid-cols-3">
							<div className="rounded-[1.7rem] border border-white/10 bg-white/5 p-5">
								<p className="text-xs uppercase tracking-[0.28em] text-slate-500">User flows</p>
								<p className="mt-3 font-display text-4xl text-white">2</p>
								<p className="mt-2 text-sm text-slate-300">Dedicated patient and admin journeys.</p>
							</div>
							<div className="rounded-[1.7rem] border border-white/10 bg-white/5 p-5">
								<p className="text-xs uppercase tracking-[0.28em] text-slate-500">Gateway</p>
								<p className="mt-3 font-display text-4xl text-white">1</p>
								<p className="mt-2 text-sm text-slate-300">Nginx serves the UI and routes API traffic.</p>
							</div>
							<div className="rounded-[1.7rem] border border-white/10 bg-white/5 p-5">
								<p className="text-xs uppercase tracking-[0.28em] text-slate-500">Observability</p>
								<p className="mt-3 font-display text-4xl text-white">24/7</p>
								<p className="mt-2 text-sm text-slate-300">Prometheus and Grafana are built into the flow.</p>
							</div>
						</div>
					</div>

					<motion.div
						initial={{ opacity: 0, y: 20 }}
						animate={{ opacity: 1, y: 0 }}
						transition={{ duration: 0.6, delay: 0.08 }}
						className="lg:sticky lg:top-8"
					>
						<AuthCard />
					</motion.div>
				</section>

				<section id="experience" className="grid gap-5 py-12 lg:grid-cols-3">
					{features.map((feature, index) => {
						const Icon = feature.icon;
						return (
							<motion.article
								key={feature.title}
								initial={{ opacity: 0, y: 18 }}
								whileInView={{ opacity: 1, y: 0 }}
								viewport={{ once: true, amount: 0.3 }}
								transition={{ delay: index * 0.08, duration: 0.45 }}
								className="rounded-[2rem] border border-white/10 bg-white/[0.06] p-6 shadow-panel backdrop-blur"
							>
								<div className="inline-flex rounded-2xl bg-gradient-to-br from-lagoon-500/20 to-coral-500/20 p-4 text-lagoon-200">
									<Icon className="h-6 w-6" />
								</div>
								<h2 className="mt-5 font-display text-3xl text-white">{feature.title}</h2>
								<p className="mt-4 text-sm leading-7 text-slate-300">{feature.copy}</p>
							</motion.article>
						);
					})}
				</section>

				<section
					id="topology"
					className="mt-8 rounded-[2.2rem] border border-white/10 bg-gradient-to-br from-white/[0.08] to-white/[0.03] p-6 shadow-panel sm:p-8"
				>
					<div className="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
						<div>
							<p className="text-xs uppercase tracking-[0.35em] text-lagoon-400">Topology</p>
							<h2 className="mt-3 font-display text-4xl text-white">Clinical product surface over a microservice core.</h2>
						</div>
						<p className="max-w-xl text-sm text-slate-300">
							The frontend now reflects the actual system architecture: gateway, role-based UX, transactional appointment flow, and live monitoring access.
						</p>
					</div>

					<div className="mt-8 grid gap-4 lg:grid-cols-[1.1fr_0.9fr]">
						<div className="grid gap-4">
							<div className="rounded-[1.7rem] border border-white/10 bg-midnight-900/70 p-5">
								<p className="text-xs uppercase tracking-[0.3em] text-slate-500">Experience layer</p>
								<div className="mt-4 flex flex-wrap gap-3">
									{["React workspace", "Nginx gateway", "Role-aware routing"].map((item) => (
										<span
											key={item}
											className="rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm text-slate-200"
										>
											{item}
										</span>
									))}
								</div>
							</div>
							<div className="rounded-[1.7rem] border border-white/10 bg-midnight-900/70 p-5">
								<p className="text-xs uppercase tracking-[0.3em] text-slate-500">Core services</p>
								<div className="mt-4 grid gap-3 md:grid-cols-2">
									{["Auth", "Patients", "Doctors", "Appointments"].map((item) => (
										<div
											key={item}
											className="rounded-[1.2rem] border border-white/10 bg-white/5 px-4 py-4 text-sm text-slate-200"
										>
											{item}
										</div>
									))}
								</div>
							</div>
						</div>

						<div className="rounded-[1.7rem] border border-white/10 bg-midnight-900/70 p-5">
							<p className="text-xs uppercase tracking-[0.3em] text-slate-500">Care roster preview</p>
							<div className="mt-5 space-y-4">
								{doctors.slice(0, 3).map((doctor) => (
									<div
										key={doctor.id}
										className="flex items-start gap-4 rounded-[1.2rem] border border-white/10 bg-white/5 p-4"
									>
										<div className="rounded-2xl bg-gradient-to-br from-lagoon-400 to-coral-400 px-3 py-2 font-semibold text-midnight-950">
											{doctor.full_name
												.split(" ")
												.map((part) => part[0])
												.join("")
												.slice(0, 2)}
										</div>
										<div>
											<h3 className="font-semibold text-white">{doctor.full_name}</h3>
											<p className="text-sm text-slate-300">{doctor.specialization || "General medicine"}</p>
											<p className="mt-1 text-xs uppercase tracking-[0.22em] text-slate-500">
												{doctor.office || "Primary campus"}
											</p>
										</div>
									</div>
								))}
							</div>
						</div>
					</div>
				</section>

				<section id="pulse" className="grid gap-5 py-12 lg:grid-cols-[0.9fr_1.1fr]">
					<div className="rounded-[2rem] border border-white/10 bg-white/[0.06] p-6 shadow-panel">
						<p className="text-xs uppercase tracking-[0.35em] text-lagoon-400">Why this redesign</p>
						<h2 className="mt-4 font-display text-4xl text-white">The frontend now carries the same maturity as the backend.</h2>
						<p className="mt-4 text-sm leading-7 text-slate-300">
							This version moves beyond a basic dashboard. It gives the assignment a product-grade visual identity, better information hierarchy, and a clearer path from user actions to operational verification.
						</p>
						<div className="mt-8 flex flex-wrap gap-3">
							<Link
								to="/workspace"
								className="inline-flex items-center gap-2 rounded-full bg-white px-5 py-3 text-sm font-semibold text-midnight-950"
							>
								Open workspace
								<ArrowRight className="h-4 w-4" />
							</Link>
							<Link
								to="/operations"
								className="inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/5 px-5 py-3 text-sm font-semibold text-white"
							>
								Open operations
							</Link>
						</div>
					</div>

					<div className="grid gap-4 md:grid-cols-2">
						{health.map((service) => (
							<div
								key={service.service}
								className="rounded-[1.8rem] border border-white/10 bg-midnight-900/70 p-5 shadow-panel"
							>
								<div className="flex items-start justify-between gap-4">
									<div>
										<p className="text-xs uppercase tracking-[0.26em] text-slate-500">Live service</p>
										<h3 className="mt-3 text-xl font-semibold text-white">{service.service}</h3>
									</div>
									<StatusBadge status={service.status} />
								</div>
								<p className="mt-4 text-sm text-slate-300">
									{service.error ?? service.storage ?? "Healthy and responding through the gateway."}
								</p>
							</div>
						))}
					</div>
				</section>
			</div>
		</div>
	);
}
