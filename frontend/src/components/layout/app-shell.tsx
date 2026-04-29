import type { ReactNode } from "react";
import {
	Activity,
	ArrowLeftRight,
	CalendarHeart,
	LogOut,
	ShieldCheck,
	Stethoscope,
} from "lucide-react";
import { Link, NavLink } from "react-router-dom";
import { cn } from "../../lib/cn";
import { useAuth } from "../../providers/auth-provider";

type AppShellProps = {
	title: string;
	description: string;
	children: ReactNode;
};

export function AppShell({ title, description, children }: AppShellProps) {
	const { user, logout } = useAuth();
	const navigation = user
		? [
				{ to: "/workspace", label: "Workspace", icon: CalendarHeart },
				{ to: "/operations", label: "Operations", icon: Activity },
			]
		: [
				{ to: "/", label: "Home", icon: CalendarHeart },
				{ to: "/operations", label: "Operations", icon: Activity },
			];

	return (
		<div className="min-h-screen bg-midnight-950 bg-mesh-radial text-white">
			<div className="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.045)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.045)_1px,transparent_1px)] bg-[size:120px_120px] opacity-[0.04]" />
			<div className="relative mx-auto flex min-h-screen max-w-[1600px] flex-col lg:flex-row">
				<aside className="border-b border-white/10 bg-white/5 px-5 py-6 backdrop-blur-xl lg:min-h-screen lg:w-[300px] lg:border-b-0 lg:border-r">
					<div className="flex items-start justify-between gap-4 lg:block">
						<div className="space-y-4">
							<div className="inline-flex items-center gap-3 rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm">
								<div className="rounded-full bg-gradient-to-br from-lagoon-400 to-coral-400 p-2 text-midnight-950">
									<Stethoscope className="h-4 w-4" />
								</div>
								<div>
									<p className="font-semibold tracking-wide text-white">MedSync</p>
									<p className="text-xs uppercase tracking-[0.28em] text-slate-400">
										Clinical command
									</p>
								</div>
							</div>

							<div className="rounded-[1.6rem] border border-white/10 bg-gradient-to-br from-white/10 to-white/5 p-5 shadow-panel">
								<p className="text-xs uppercase tracking-[0.3em] text-slate-400">Signed in</p>
								<h2 className="mt-3 font-display text-2xl text-white">
									{user?.full_name ?? "Guest observer"}
								</h2>
								<div className="mt-4 inline-flex items-center gap-2 rounded-full border border-lagoon-400/20 bg-lagoon-500/10 px-3 py-1 text-xs uppercase tracking-[0.24em] text-lagoon-200">
									<ShieldCheck className="h-3.5 w-3.5" />
									{user?.role ?? "preview"}
								</div>
								<p className="mt-4 text-sm text-slate-300">
									{user?.email ?? "Sign in to unlock the full clinical workspace."}
								</p>
							</div>
						</div>

						{user ? (
							<button
								type="button"
								onClick={logout}
								className="inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm text-slate-200 transition hover:border-coral-400/40 hover:bg-coral-500/10 hover:text-white"
							>
								<LogOut className="h-4 w-4" />
								Log out
							</button>
						) : (
							<Link
								to="/"
								className="inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm text-slate-200 transition hover:border-lagoon-400/30 hover:bg-lagoon-500/10 hover:text-white"
							>
								Return to sign in
							</Link>
						)}
					</div>

					<nav className="mt-8 grid gap-3">
						{navigation.map(({ to, label, icon: Icon }) => (
							<NavLink
								key={to}
								to={to}
								className={({ isActive }) =>
									cn(
										"flex items-center gap-3 rounded-2xl border px-4 py-3 text-sm transition",
										isActive
											? "border-lagoon-400/30 bg-lagoon-500/[0.15] text-white shadow-glow"
											: "border-white/10 bg-white/5 text-slate-300 hover:border-white/20 hover:bg-white/10 hover:text-white",
									)
								}
							>
								<Icon className="h-4 w-4" />
								{label}
							</NavLink>
						))}
					</nav>

					<div className="mt-8 rounded-[1.5rem] border border-white/10 bg-white/5 p-5">
						<p className="text-xs uppercase tracking-[0.28em] text-slate-500">Environment</p>
						<div className="mt-4 space-y-3 text-sm text-slate-300">
							<div className="flex items-center justify-between">
								<span>Gateway</span>
								<span className="inline-flex items-center gap-2 text-lagoon-200">
									<ArrowLeftRight className="h-4 w-4" />
									Nginx
								</span>
							</div>
							<div className="flex items-center justify-between">
								<span>Frontend</span>
								<span>React + Tailwind</span>
							</div>
							<div className="flex items-center justify-between">
								<span>Monitoring</span>
								<span>Prometheus / Grafana</span>
							</div>
						</div>
					</div>
				</aside>

				<main className="relative flex-1 px-5 py-6 sm:px-8 lg:px-10">
					<header className="mb-8 flex flex-col gap-5 border-b border-white/10 pb-6 md:flex-row md:items-end md:justify-between">
						<div className="space-y-3">
							<p className="text-xs uppercase tracking-[0.35em] text-lagoon-400">
								MedSync experience
							</p>
							<h1 className="font-display text-4xl text-white sm:text-5xl">{title}</h1>
							<p className="max-w-3xl text-sm text-slate-300 sm:text-base">{description}</p>
						</div>

						<div className="flex flex-wrap gap-3 text-sm">
							<a
								href="http://localhost:3000"
								target="_blank"
								rel="noreferrer"
								className="rounded-full border border-white/10 bg-white/5 px-4 py-2 text-slate-200 transition hover:border-lagoon-400/30 hover:bg-lagoon-500/10 hover:text-white"
							>
								Grafana
							</a>
							<a
								href="http://localhost:9090"
								target="_blank"
								rel="noreferrer"
								className="rounded-full border border-white/10 bg-white/5 px-4 py-2 text-slate-200 transition hover:border-sun-400/30 hover:bg-sun-500/10 hover:text-white"
							>
								Prometheus
							</a>
						</div>
					</header>

					{children}
				</main>
			</div>
		</div>
	);
}
