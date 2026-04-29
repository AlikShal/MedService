export function LoadingScreen() {
	return (
		<div className="relative flex min-h-screen items-center justify-center overflow-hidden bg-midnight-950 text-white">
			<div className="absolute inset-0 bg-mesh-radial" />
			<div className="relative flex flex-col items-center gap-6">
				<div className="flex items-center gap-3 rounded-full border border-white/10 bg-white/5 px-5 py-3 text-sm uppercase tracking-[0.35em] text-lagoon-300">
					<div className="h-3 w-3 animate-pulse rounded-full bg-lagoon-400" />
					Loading MedSync
				</div>
				<div className="h-32 w-32 animate-float rounded-full bg-gradient-to-br from-lagoon-400/30 to-coral-400/20 blur-3xl" />
			</div>
		</div>
	);
}
