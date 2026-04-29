import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { ArrowRight, ShieldCheck } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";
import { ApiError } from "../../lib/api";
import { useAuth } from "../../providers/auth-provider";
import { cn } from "../../lib/cn";

type Mode = "login" | "register";

const initialLogin = { email: "", password: "" };
const initialRegister = { full_name: "", email: "", password: "" };

export function AuthCard() {
	const navigate = useNavigate();
	const { login, register } = useAuth();
	const [mode, setMode] = useState<Mode>("login");
	const [loginForm, setLoginForm] = useState(initialLogin);
	const [registerForm, setRegisterForm] = useState(initialRegister);

	const mutation = useMutation({
		mutationFn: async () => {
			if (mode === "login") {
				return login(loginForm);
			}
			return register(registerForm);
		},
		onSuccess: (payload) => {
			toast.success(mode === "login" ? "Welcome back." : "Account created successfully.");
			navigate("/workspace");
			if (mode === "login") {
				setLoginForm(initialLogin);
			} else {
				setRegisterForm(initialRegister);
			}
			if (payload.user.role === "patient") {
				toast.info("Patient workspace is ready. Complete the profile to start booking care.");
			}
		},
		onError: (error) => {
			const message = error instanceof ApiError ? error.message : "Unable to complete authentication";
			toast.error(message);
		},
	});

	return (
		<div className="overflow-hidden rounded-[2rem] border border-white/10 bg-white/[0.08] shadow-panel backdrop-blur-2xl">
			<div className="border-b border-white/10 px-6 py-5">
				<div className="inline-flex items-center gap-2 rounded-full border border-lagoon-400/20 bg-lagoon-500/10 px-3 py-1 text-xs uppercase tracking-[0.22em] text-lagoon-200">
					<ShieldCheck className="h-3.5 w-3.5" />
					Secure access
				</div>
				<h3 className="mt-4 font-display text-3xl text-white">Step into the control center</h3>
				<p className="mt-2 text-sm text-slate-300">
					Patients can onboard and book visits. Admins can orchestrate doctors and appointments.
				</p>
			</div>

			<div className="px-6 pt-5">
				<div className="grid grid-cols-2 gap-2 rounded-full border border-white/10 bg-white/5 p-1">
					{(["login", "register"] as Mode[]).map((candidate) => (
						<button
							key={candidate}
							type="button"
							onClick={() => setMode(candidate)}
							className={cn(
								"rounded-full px-4 py-2 text-sm font-medium capitalize transition",
								mode === candidate
									? "bg-white text-midnight-950"
									: "text-slate-300 hover:text-white",
							)}
						>
							{candidate}
						</button>
					))}
				</div>
			</div>

			<form
				className="space-y-4 px-6 py-6"
				onSubmit={(event) => {
					event.preventDefault();
					mutation.mutate();
				}}
			>
				{mode === "register" ? (
					<div className="space-y-2">
						<label className="text-xs uppercase tracking-[0.25em] text-slate-400">Full name</label>
						<input
							value={registerForm.full_name}
							onChange={(event) =>
								setRegisterForm((current) => ({ ...current, full_name: event.target.value }))
							}
							required
							placeholder="Aruzhan Saparova"
							className="field"
						/>
					</div>
				) : null}

				<div className="space-y-2">
					<label className="text-xs uppercase tracking-[0.25em] text-slate-400">Email</label>
					<input
						type="email"
						value={mode === "login" ? loginForm.email : registerForm.email}
						onChange={(event) => {
							const email = event.target.value;
							if (mode === "login") {
								setLoginForm((current) => ({ ...current, email }));
								return;
							}
							setRegisterForm((current) => ({ ...current, email }));
						}}
						required
						placeholder="you@clinic.local"
						className="field"
					/>
				</div>

				<div className="space-y-2">
					<label className="text-xs uppercase tracking-[0.25em] text-slate-400">Password</label>
					<input
						type="password"
						value={mode === "login" ? loginForm.password : registerForm.password}
						onChange={(event) => {
							const password = event.target.value;
							if (mode === "login") {
								setLoginForm((current) => ({ ...current, password }));
								return;
							}
							setRegisterForm((current) => ({ ...current, password }));
						}}
						required
						minLength={mode === "register" ? 6 : undefined}
						placeholder={mode === "login" ? "Enter your password" : "At least 6 characters"}
						className="field"
					/>
				</div>

				<button
					type="submit"
					disabled={mutation.isPending}
					className="group inline-flex w-full items-center justify-center gap-2 rounded-full bg-gradient-to-r from-lagoon-400 via-lagoon-500 to-coral-500 bg-[length:200%_200%] px-5 py-3 text-sm font-semibold text-midnight-950 transition hover:animate-shine disabled:cursor-not-allowed disabled:opacity-60"
				>
					{mutation.isPending ? "Processing..." : mode === "login" ? "Open workspace" : "Create patient account"}
					<ArrowRight className="h-4 w-4 transition group-hover:translate-x-1" />
				</button>

				<div className="rounded-[1.4rem] border border-white/10 bg-midnight-900/60 p-4">
					<p className="text-xs uppercase tracking-[0.22em] text-slate-500">Demo admin access</p>
					<div className="mt-3 flex flex-wrap gap-2 text-sm text-slate-200">
						<span className="rounded-full border border-white/10 bg-white/5 px-3 py-1">
							admin@medsync.local
						</span>
						<span className="rounded-full border border-white/10 bg-white/5 px-3 py-1">admin123</span>
					</div>
				</div>
			</form>
		</div>
	);
}
