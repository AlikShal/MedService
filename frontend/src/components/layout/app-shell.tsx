import type { ReactNode } from "react";
import { Activity, ArrowLeftRight, CalendarHeart, LogOut, ShieldCheck, Stethoscope } from "lucide-react";
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
	const navigation = [
		{ to: "/workspace", label: "Workspace", icon: CalendarHeart },
		{ to: "/operations", label: "Operations", icon: Activity },
	];

	return (
		<div className="min-h-screen bg-[#070d19] text-white">
			<div className="flex min-h-screen">
				<aside className="hidden min-h-screen w-[136px] shrink-0 flex-col border-r border-slate-800 bg-[#0a111f] lg:flex">
					<div className="px-3 py-2">
						<Link to="/" className="flex items-center gap-2">
							<div className="flex h-7 w-7 items-center justify-center rounded-lg bg-lagoon-500/10 text-lagoon-400">
								<Stethoscope className="h-3.5 w-3.5" />
							</div>
							<div>
								<p className="text-[0.68rem] font-extrabold leading-none text-white">MedSync</p>
								<p className="mt-1 text-[0.46rem] font-bold uppercase tracking-[0.26em] text-slate-500">
									Clinical command
								</p>
							</div>
						</Link>
					</div>

					<div className="mx-3 mt-3 rounded-lg border border-slate-800 bg-[#0d1422] p-3">
						<p className="text-[0.5rem] font-extrabold uppercase tracking-[0.22em] text-slate-500">
							Signed in
						</p>
						<h2 className="mt-3 truncate text-[0.72rem] font-extrabold text-white">
							{user?.full_name ?? "Guest observer"}
						</h2>
						<div className="mt-2 inline-flex items-center gap-1 rounded-full bg-lagoon-500/10 px-2 py-0.5 text-[0.5rem] font-extrabold uppercase tracking-[0.15em] text-lagoon-400">
							<ShieldCheck className="h-2.5 w-2.5" />
							{user?.role ?? "preview"}
						</div>
						<p className="mt-3 truncate text-[0.58rem] text-slate-400">
							{user?.email ?? "Sign in to unlock workspace."}
						</p>

						{user ? (
							<button
								type="button"
								onClick={logout}
								className="mt-3 flex h-6 w-full items-center justify-center gap-1 rounded-md border border-slate-800 bg-[#070d19] text-[0.58rem] font-bold text-slate-300 transition hover:text-white"
							>
								<LogOut className="h-3 w-3" />
								Log out
							</button>
						) : (
							<Link
								to="/"
								className="mt-3 flex h-6 w-full items-center justify-center rounded-md border border-slate-800 bg-[#070d19] text-[0.58rem] font-bold text-slate-300 transition hover:text-white"
							>
								Sign in
							</Link>
						)}
					</div>

					<nav className="mx-3 mt-2 grid gap-1">
						{navigation.map(({ to, label, icon: Icon }) => (
							<NavLink
								key={to}
								to={to}
								className={({ isActive }) =>
									cn(
										"flex h-9 items-center gap-2 rounded-md px-2 text-[0.62rem] font-extrabold transition",
										isActive
											? "bg-lagoon-500/10 text-white"
											: "text-slate-400 hover:bg-slate-900/70 hover:text-white",
									)
								}
							>
								<Icon className="h-3.5 w-3.5 text-lagoon-400" />
								{label}
							</NavLink>
						))}
					</nav>

					<div className="mx-3 mb-3 mt-auto rounded-lg border border-slate-800 bg-[#0d1422] p-3">
						<p className="text-[0.5rem] font-extrabold uppercase tracking-[0.22em] text-slate-500">
							Environment
						</p>
						<div className="mt-3 space-y-2 text-[0.55rem] text-slate-400">
							<div className="flex items-center justify-between gap-2">
								<span>Gateway</span>
								<span className="font-bold text-white">Nginx</span>
							</div>
							<div className="flex items-center justify-between gap-2">
								<span>Frontend</span>
								<span className="font-bold text-white">React + Tailwind</span>
							</div>
							<div className="flex items-center justify-between gap-2">
								<span>Monitoring</span>
								<span className="inline-flex items-center gap-1 font-bold text-white">
									<ArrowLeftRight className="h-2.5 w-2.5 text-lagoon-400" />
									Prom/Graf
								</span>
							</div>
						</div>
					</div>
				</aside>

				<main className="min-w-0 flex-1 px-4 py-4 sm:px-6 lg:px-0">
					<div className="mb-5 flex items-center justify-between rounded-lg border border-slate-800 bg-[#0d1422] p-3 lg:hidden">
						<Link to="/" className="flex items-center gap-2">
							<div className="flex h-8 w-8 items-center justify-center rounded-lg bg-lagoon-500/10 text-lagoon-400">
								<Stethoscope className="h-4 w-4" />
							</div>
							<div>
								<p className="text-sm font-extrabold leading-none text-white">MedSync</p>
								<p className="mt-1 text-[0.52rem] font-bold uppercase tracking-[0.25em] text-slate-500">
									Clinical command
								</p>
							</div>
						</Link>
						<div className="flex items-center gap-2">
							<Link to="/workspace" className="rounded-md px-2 py-1 text-xs font-bold text-slate-300">
								Workspace
							</Link>
							<Link to="/operations" className="rounded-md px-2 py-1 text-xs font-bold text-slate-300">
								Ops
							</Link>
							{user ? (
								<button type="button" onClick={logout} className="rounded-md px-2 py-1 text-xs font-bold text-slate-300">
									Log out
								</button>
							) : null}
						</div>
					</div>
					<div className="mx-auto max-w-[660px]">
						<header className="mb-8">
							<p className="lovable-eyebrow text-slate-500">MedSync experience</p>
							<h1 className="mt-1 text-2xl font-extrabold leading-tight tracking-[-0.02em] text-white">
								{title}
							</h1>
							<p className="mt-2 max-w-[520px] text-xs leading-5 text-slate-400">{description}</p>
						</header>

						{children}
					</div>
				</main>
			</div>
		</div>
	);
}
