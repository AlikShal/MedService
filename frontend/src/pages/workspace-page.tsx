import { useEffect, useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
	Activity,
	CalendarDays,
	Clock3,
	ClipboardPlus,
	HeartHandshake,
	Plus,
	ShieldAlert,
	Stethoscope,
	UsersRound,
} from "lucide-react";
import { toast } from "sonner";
import { AppShell } from "../components/layout/app-shell";
import { PanelHeader } from "../components/shared/panel-header";
import { StatCard } from "../components/shared/stat-card";
import { StatusBadge } from "../components/shared/status-badge";
import { api, ApiError } from "../lib/api";
import { useAuth } from "../providers/auth-provider";
import type { Appointment, AppointmentStatus, PatientProfile } from "../types/domain";

function formatDateTime(value: string) {
	return new Intl.DateTimeFormat(undefined, {
		dateStyle: "medium",
		timeStyle: "short",
	}).format(new Date(value));
}

function datetimeLocalFromDate(date: Date) {
	const offset = date.getTimezoneOffset();
	const adjusted = new Date(date.getTime() - offset * 60_000);
	return adjusted.toISOString().slice(0, 16);
}

function nextDefaultSlot() {
	const base = new Date();
	base.setDate(base.getDate() + 2);
	base.setHours(10, 30, 0, 0);
	return datetimeLocalFromDate(base);
}

function countHealthyServices(statuses: { status: string }[]) {
	return statuses.filter((item) => item.status === "ok").length;
}

function getProgressionButtons(status: AppointmentStatus) {
	switch (status) {
		case "new":
			return [
				{ label: "Start", status: "in_progress" as const },
				{ label: "Cancel", status: "cancelled" as const },
			];
		case "in_progress":
			return [
				{ label: "Complete", status: "done" as const },
				{ label: "Cancel", status: "cancelled" as const },
			];
		default:
			return [];
	}
}

export function WorkspacePage() {
	const queryClient = useQueryClient();
	const { user, token } = useAuth();
	const isAdmin = user?.role === "admin";
	const [profileForm, setProfileForm] = useState({
		full_name: user?.full_name ?? "",
		phone: "",
		date_of_birth: "",
		notes: "",
	});
	const [doctorForm, setDoctorForm] = useState({
		full_name: "",
		specialization: "",
		email: "",
		office: "",
	});
	const [appointmentForm, setAppointmentForm] = useState({
		title: "",
		description: "",
		doctor_id: "",
		scheduled_at: nextDefaultSlot(),
	});

	const doctorsQuery = useQuery({
		queryKey: ["doctors"],
		queryFn: api.getDoctors,
	});

	const profileQuery = useQuery({
		queryKey: ["profile", user?.id],
		enabled: Boolean(token && user?.role === "patient"),
		queryFn: async () => {
			try {
				return await api.getProfile(token!);
			} catch (error) {
				if (error instanceof ApiError && error.status === 404) {
					return null;
				}
				throw error;
			}
		},
	});

	const appointmentsQuery = useQuery({
		queryKey: ["appointments", user?.id],
		enabled: Boolean(token),
		queryFn: () => api.getAppointments(token!),
	});

	const healthQuery = useQuery({
		queryKey: ["workspace-health"],
		queryFn: api.getPlatformHealth,
		refetchInterval: 20000,
	});

	useEffect(() => {
		if (doctorsQuery.data?.length && !appointmentForm.doctor_id) {
			setAppointmentForm((current) => ({
				...current,
				doctor_id: doctorsQuery.data?.[0]?.id ?? "",
			}));
		}
	}, [appointmentForm.doctor_id, doctorsQuery.data]);

	useEffect(() => {
		const profile = profileQuery.data;
		if (!profile) {
			setProfileForm((current) => ({ ...current, full_name: user?.full_name ?? current.full_name }));
			return;
		}

		setProfileForm({
			full_name: profile.full_name,
			phone: profile.phone,
			date_of_birth: profile.date_of_birth,
			notes: profile.notes,
		});
	}, [profileQuery.data, user?.full_name]);

	const profileMutation = useMutation({
		mutationFn: (payload: typeof profileForm) => api.upsertProfile(payload, token!),
		onSuccess: () => {
			toast.success("Patient profile saved.");
			void queryClient.invalidateQueries({ queryKey: ["profile", user?.id] });
		},
		onError: (error) => {
			toast.error(error instanceof Error ? error.message : "Unable to save patient profile.");
		},
	});

	const doctorMutation = useMutation({
		mutationFn: (payload: typeof doctorForm) => api.createDoctor(payload, token!),
		onSuccess: () => {
			toast.success("Doctor added to the clinical roster.");
			setDoctorForm({ full_name: "", specialization: "", email: "", office: "" });
			void queryClient.invalidateQueries({ queryKey: ["doctors"] });
		},
		onError: (error) => {
			toast.error(error instanceof Error ? error.message : "Unable to create doctor.");
		},
	});

	const appointmentMutation = useMutation({
		mutationFn: (payload: typeof appointmentForm) =>
			api.createAppointment(
				{
					...payload,
					scheduled_at: new Date(payload.scheduled_at).toISOString(),
				},
				token!,
			),
		onSuccess: () => {
			toast.success("Appointment booked successfully.");
			setAppointmentForm({
				title: "",
				description: "",
				doctor_id: doctorsQuery.data?.[0]?.id ?? "",
				scheduled_at: nextDefaultSlot(),
			});
			void queryClient.invalidateQueries({ queryKey: ["appointments", user?.id] });
		},
		onError: (error) => {
			toast.error(error instanceof Error ? error.message : "Unable to create appointment.");
		},
	});

	const statusMutation = useMutation({
		mutationFn: ({ id, status }: { id: string; status: AppointmentStatus }) =>
			api.updateAppointmentStatus(id, status, token!),
		onSuccess: () => {
			toast.success("Appointment status updated.");
			void queryClient.invalidateQueries({ queryKey: ["appointments", user?.id] });
		},
		onError: (error) => {
			toast.error(error instanceof Error ? error.message : "Unable to update appointment status.");
		},
	});

	const doctors = doctorsQuery.data ?? [];
	const appointments = appointmentsQuery.data ?? [];
	const health = healthQuery.data ?? [];
	const healthyCount = countHealthyServices(health);
	const nextAppointment = appointments[0];

	const appointmentColumns = useMemo(() => {
		const map: Record<AppointmentStatus, Appointment[]> = {
			new: [],
			in_progress: [],
			done: [],
			cancelled: [],
		};

		for (const appointment of appointments) {
			map[appointment.status].push(appointment);
		}

		return map;
	}, [appointments]);

	const profileReady = Boolean(profileQuery.data);

	return (
		<AppShell
			title={isAdmin ? "Clinical command workspace" : "Your care workspace"}
			description={
				isAdmin
					? "Manage the doctor roster, move appointments through their lifecycle, and keep a constant eye on service health."
					: "Everything a patient needs lives here: care profile, doctor discovery, and appointment scheduling with live operational awareness."
			}
		>
			<div className="space-y-8">
				<section className="grid gap-4 xl:grid-cols-4">
					<StatCard
						label="Doctor roster"
						value={String(doctors.length)}
						hint="Active clinicians available for scheduling."
						icon={<Stethoscope className="h-5 w-5" />}
					/>
					<StatCard
						label="Appointments"
						value={String(appointments.length)}
						hint={isAdmin ? "Total workload across the platform." : "Your scheduled and historical visits."}
						icon={<CalendarDays className="h-5 w-5" />}
						tone="coral"
					/>
					<StatCard
						label="Service health"
						value={`${healthyCount}/${health.length || 4}`}
						hint="Healthy backend targets visible through the UI."
						icon={<Activity className="h-5 w-5" />}
						tone="sun"
					/>
					<StatCard
						label={isAdmin ? "Workflow state" : "Profile readiness"}
						value={isAdmin ? appointmentColumns.new.length.toString() : profileReady ? "Ready" : "Draft"}
						hint={
							isAdmin
								? "Appointments still waiting for clinical action."
								: "Profile completion unlocks the booking flow."
						}
						icon={isAdmin ? <ClipboardPlus className="h-5 w-5" /> : <HeartHandshake className="h-5 w-5" />}
						tone={isAdmin ? "lagoon" : "coral"}
					/>
				</section>

				<section className="grid gap-6 xl:grid-cols-[1.15fr_0.85fr]">
					<div className="space-y-6">
						{isAdmin ? (
							<div className="grid gap-6 lg:grid-cols-[0.46fr_0.54fr]">
								<section className="rounded-[2rem] border border-white/10 bg-white/[0.06] p-6 shadow-panel">
									<PanelHeader
										eyebrow="Roster control"
										title="Add a doctor with context."
										description="Expand the clinical roster with office and specialization data so the experience feels operationally complete."
									/>

									<form
										className="mt-6 space-y-4"
										onSubmit={(event) => {
											event.preventDefault();
											doctorMutation.mutate(doctorForm);
										}}
									>
										<div className="grid gap-4 md:grid-cols-2">
											<input
												className="field"
												placeholder="Full name"
												value={doctorForm.full_name}
												onChange={(event) =>
													setDoctorForm((current) => ({ ...current, full_name: event.target.value }))
												}
												required
											/>
											<input
												className="field"
												placeholder="Specialization"
												value={doctorForm.specialization}
												onChange={(event) =>
													setDoctorForm((current) => ({
														...current,
														specialization: event.target.value,
													}))
												}
											/>
											<input
												className="field"
												type="email"
												placeholder="Doctor email"
												value={doctorForm.email}
												onChange={(event) =>
													setDoctorForm((current) => ({ ...current, email: event.target.value }))
												}
												required
											/>
											<input
												className="field"
												placeholder="Office / room"
												value={doctorForm.office}
												onChange={(event) =>
													setDoctorForm((current) => ({ ...current, office: event.target.value }))
												}
											/>
										</div>

										<button
											type="submit"
											disabled={doctorMutation.isPending}
											className="inline-flex items-center gap-2 rounded-full bg-white px-5 py-3 text-sm font-semibold text-midnight-950"
										>
											<Plus className="h-4 w-4" />
											{doctorMutation.isPending ? "Adding doctor..." : "Add to roster"}
										</button>
									</form>
								</section>

								<section className="rounded-[2rem] border border-white/10 bg-white/[0.06] p-6 shadow-panel">
									<PanelHeader
										eyebrow="Command view"
										title="Live appointment pressure"
										description="Use the workflow board to move bookings from intake to completion."
									/>

									<div className="mt-6 grid gap-4 md:grid-cols-2">
										{health.map((service) => (
											<div
												key={service.service}
												className="rounded-[1.4rem] border border-white/10 bg-midnight-900/60 p-4"
											>
												<div className="flex items-center justify-between gap-3">
													<h3 className="text-sm font-semibold text-white">{service.service}</h3>
													<StatusBadge status={service.status} />
												</div>
												<p className="mt-3 text-sm text-slate-300">
													{service.error ?? service.storage ?? "Healthy and available."}
												</p>
											</div>
										))}
									</div>
								</section>
							</div>
						) : (
							<div className="grid gap-6 lg:grid-cols-[0.48fr_0.52fr]">
								<section className="rounded-[2rem] border border-white/10 bg-white/[0.06] p-6 shadow-panel">
									<PanelHeader
										eyebrow="Patient profile"
										title="Shape a care-ready identity"
										description="We use this profile to attach appointments to a real patient record and make the experience feel grounded."
									/>

									<form
										className="mt-6 space-y-4"
										onSubmit={(event) => {
											event.preventDefault();
											profileMutation.mutate(profileForm);
										}}
									>
										<div className="grid gap-4 md:grid-cols-2">
											<input
												className="field"
												placeholder="Full name"
												value={profileForm.full_name}
												onChange={(event) =>
													setProfileForm((current) => ({ ...current, full_name: event.target.value }))
												}
												required
											/>
											<input
												className="field"
												placeholder="Phone"
												value={profileForm.phone}
												onChange={(event) =>
													setProfileForm((current) => ({ ...current, phone: event.target.value }))
												}
											/>
											<input
												className="field"
												type="date"
												value={profileForm.date_of_birth}
												onChange={(event) =>
													setProfileForm((current) => ({
														...current,
														date_of_birth: event.target.value,
													}))
												}
											/>
											<div className="rounded-[1.2rem] border border-white/10 bg-midnight-900/60 px-4 py-3 text-sm text-slate-200">
												<p className="text-xs uppercase tracking-[0.24em] text-slate-500">Profile state</p>
												<p className="mt-3 font-semibold text-white">{profileReady ? "Care-ready" : "Needs completion"}</p>
											</div>
										</div>
										<textarea
											className="field min-h-32 resize-none"
											placeholder="Allergies, visit preferences, recurring symptoms, or reminders"
											value={profileForm.notes}
											onChange={(event) =>
												setProfileForm((current) => ({ ...current, notes: event.target.value }))
											}
										/>
										<button
											type="submit"
											disabled={profileMutation.isPending}
											className="inline-flex items-center gap-2 rounded-full bg-white px-5 py-3 text-sm font-semibold text-midnight-950"
										>
											{profileMutation.isPending ? "Saving profile..." : "Save patient profile"}
										</button>
									</form>
								</section>

								<section className="rounded-[2rem] border border-white/10 bg-white/[0.06] p-6 shadow-panel">
									<PanelHeader
										eyebrow="Booking studio"
										title="Schedule a polished appointment flow"
										description="The booking flow only unlocks once the patient profile is ready, which mirrors the backend validation."
									/>

									<form
										className="mt-6 space-y-4"
										onSubmit={(event) => {
											event.preventDefault();
											appointmentMutation.mutate(appointmentForm);
										}}
									>
										<input
											className="field"
											placeholder="Appointment title"
											value={appointmentForm.title}
											onChange={(event) =>
												setAppointmentForm((current) => ({ ...current, title: event.target.value }))
											}
											disabled={!profileReady}
											required
										/>
										<div className="grid gap-4 md:grid-cols-2">
											<select
												className="field"
												value={appointmentForm.doctor_id}
												onChange={(event) =>
													setAppointmentForm((current) => ({ ...current, doctor_id: event.target.value }))
												}
												disabled={!profileReady}
											>
												{doctors.map((doctor) => (
													<option key={doctor.id} value={doctor.id}>
														{doctor.full_name} · {doctor.specialization || "General medicine"}
													</option>
												))}
											</select>
											<input
												className="field"
												type="datetime-local"
												value={appointmentForm.scheduled_at}
												onChange={(event) =>
													setAppointmentForm((current) => ({
														...current,
														scheduled_at: event.target.value,
													}))
												}
												disabled={!profileReady}
												required
											/>
										</div>
										<textarea
											className="field min-h-28 resize-none"
											placeholder="Reason for visit, symptoms, history, or desired outcome"
											value={appointmentForm.description}
											onChange={(event) =>
												setAppointmentForm((current) => ({
													...current,
													description: event.target.value,
												}))
											}
											disabled={!profileReady}
										/>
										<button
											type="submit"
											disabled={appointmentMutation.isPending || !profileReady || !doctors.length}
											className="inline-flex items-center gap-2 rounded-full bg-gradient-to-r from-lagoon-400 to-coral-500 px-5 py-3 text-sm font-semibold text-midnight-950 disabled:opacity-50"
										>
											{appointmentMutation.isPending ? "Booking visit..." : "Book appointment"}
										</button>
									</form>
								</section>
							</div>
						)}

						<section className="rounded-[2rem] border border-white/10 bg-white/[0.06] p-6 shadow-panel">
							<PanelHeader
								eyebrow="Doctor directory"
								title="A roster that feels curated, not dumped."
								description="Strong visual hierarchy helps the doctor list feel like an actual product surface instead of a raw API table."
							/>

							<div className="mt-6 grid gap-4 md:grid-cols-2 xl:grid-cols-3">
								{doctors.map((doctor, index) => (
									<div
										key={doctor.id}
										className="rounded-[1.7rem] border border-white/10 bg-gradient-to-br from-white/10 to-white/[0.04] p-5 shadow-panel"
									>
										<div className="flex items-start justify-between gap-3">
											<div>
												<div className="mb-4 inline-flex rounded-2xl bg-gradient-to-br from-lagoon-400 to-coral-400 px-3 py-2 text-sm font-semibold text-midnight-950">
													{doctor.full_name
														.split(" ")
														.map((part) => part[0])
														.join("")
														.slice(0, 2)}
												</div>
												<h3 className="text-lg font-semibold text-white">{doctor.full_name}</h3>
												<p className="mt-1 text-sm text-slate-300">
													{doctor.specialization || "General medicine"}
												</p>
											</div>
											<span className="rounded-full border border-white/10 bg-white/5 px-3 py-1 text-[0.7rem] uppercase tracking-[0.22em] text-slate-400">
												#{index + 1}
											</span>
										</div>
										<div className="mt-6 space-y-2 text-sm text-slate-300">
											<p>{doctor.email}</p>
											<p>{doctor.office || "Primary campus"}</p>
										</div>
									</div>
								))}
							</div>
						</section>
					</div>

					<div className="space-y-6">
						<section className="rounded-[2rem] border border-white/10 bg-white/[0.06] p-6 shadow-panel">
							<PanelHeader
								eyebrow="Priority feed"
								title={nextAppointment ? "What happens next" : "No appointments yet"}
								description={
									nextAppointment
										? "A polished summary of the next relevant visit keeps the workspace grounded in the user’s real flow."
										: "Create an appointment or move workflow states to start populating this feed."
								}
							/>

							{nextAppointment ? (
								<div className="mt-6 rounded-[1.8rem] border border-lagoon-400/20 bg-gradient-to-br from-lagoon-500/10 to-coral-500/10 p-5">
									<div className="flex flex-wrap items-center justify-between gap-3">
										<StatusBadge status={nextAppointment.status} />
										<div className="inline-flex items-center gap-2 text-sm text-slate-200">
											<Clock3 className="h-4 w-4" />
											{formatDateTime(nextAppointment.scheduled_at)}
										</div>
									</div>
									<h3 className="mt-4 font-display text-3xl text-white">{nextAppointment.title}</h3>
									<p className="mt-3 text-sm text-slate-300">
										{nextAppointment.description || "No additional clinical note was captured for this visit."}
									</p>
									<div className="mt-6 grid gap-3 text-sm text-slate-200">
										<div className="rounded-[1.2rem] border border-white/10 bg-white/5 px-4 py-3">
											<strong className="mr-2 text-white">Doctor:</strong>
											{nextAppointment.doctor_name}
										</div>
										<div className="rounded-[1.2rem] border border-white/10 bg-white/5 px-4 py-3">
											<strong className="mr-2 text-white">Patient:</strong>
											{nextAppointment.patient_name}
										</div>
									</div>
								</div>
							) : (
								<div className="mt-6 rounded-[1.8rem] border border-dashed border-white/[0.15] bg-midnight-900/60 p-6 text-sm text-slate-400">
									Nothing scheduled yet. Once the system books a visit, this area highlights the next important care interaction.
								</div>
							)}
						</section>

						<section className="rounded-[2rem] border border-white/10 bg-white/[0.06] p-6 shadow-panel">
							<PanelHeader
								eyebrow="Service pulse"
								title="Operational confidence at a glance"
								description="These cards keep production concerns visible without kicking the user out to another tool."
							/>

							<div className="mt-6 grid gap-4">
								{health.map((service) => (
									<div
										key={service.service}
										className="rounded-[1.4rem] border border-white/10 bg-midnight-900/60 p-4"
									>
										<div className="flex items-center justify-between gap-3">
											<h3 className="text-sm font-semibold text-white">{service.service}</h3>
											<StatusBadge status={service.status} />
										</div>
										<p className="mt-3 text-sm text-slate-300">
											{service.error ?? service.storage ?? "Healthy and ready."}
										</p>
									</div>
								))}
							</div>
						</section>
					</div>
				</section>

				<section className="rounded-[2rem] border border-white/10 bg-white/[0.06] p-6 shadow-panel">
					<PanelHeader
						eyebrow={isAdmin ? "Workflow board" : "Appointment stream"}
						title={
							isAdmin
								? "Move the system through real operational states"
								: "A cleaner timeline of your medical scheduling activity"
						}
						description={
							isAdmin
								? "The board groups workload by state so staff can progress care instead of scanning a flat list."
								: "Patients get a crisp, reassuring view of everything already scheduled."
						}
					/>

					{isAdmin ? (
						<div className="mt-6 grid gap-4 xl:grid-cols-4">
							{(["new", "in_progress", "done", "cancelled"] as AppointmentStatus[]).map((status) => (
								<div
									key={status}
									className="rounded-[1.6rem] border border-white/10 bg-midnight-900/60 p-4"
								>
									<div className="mb-4 flex items-center justify-between gap-3">
										<h3 className="text-sm font-semibold uppercase tracking-[0.24em] text-slate-400">
											{status.replaceAll("_", " ")}
										</h3>
										<StatusBadge status={status} />
									</div>
									<div className="space-y-3">
										{appointmentColumns[status].length ? (
											appointmentColumns[status].map((appointment) => (
												<div
													key={appointment.id}
													className="rounded-[1.3rem] border border-white/10 bg-white/5 p-4"
												>
													<h4 className="font-semibold text-white">{appointment.title}</h4>
													<p className="mt-2 text-sm text-slate-300">
														{appointment.patient_name} · {appointment.doctor_name}
													</p>
													<p className="mt-2 text-xs uppercase tracking-[0.22em] text-slate-500">
														{formatDateTime(appointment.scheduled_at)}
													</p>

													<div className="mt-4 flex flex-wrap gap-2">
														{getProgressionButtons(appointment.status).map((button) => (
															<button
																key={button.status}
																type="button"
																onClick={() =>
																	statusMutation.mutate({
																		id: appointment.id,
																		status: button.status,
																	})
																}
																className="rounded-full border border-white/10 bg-white/10 px-3 py-1.5 text-xs font-semibold uppercase tracking-[0.18em] text-white transition hover:border-lagoon-400/30 hover:bg-lagoon-500/10"
															>
																{button.label}
															</button>
														))}
													</div>
												</div>
											))
										) : (
											<div className="rounded-[1.3rem] border border-dashed border-white/10 px-4 py-6 text-center text-sm text-slate-500">
												No appointments in this state.
											</div>
										)}
									</div>
								</div>
							))}
						</div>
					) : (
						<div className="mt-6 grid gap-4">
							{appointments.length ? (
								appointments.map((appointment) => (
									<div
										key={appointment.id}
										className="rounded-[1.6rem] border border-white/10 bg-midnight-900/60 p-5"
									>
										<div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
											<div className="space-y-3">
												<div className="flex flex-wrap items-center gap-3">
													<StatusBadge status={appointment.status} />
													<span className="inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/5 px-3 py-1 text-xs uppercase tracking-[0.22em] text-slate-400">
														<Clock3 className="h-3.5 w-3.5" />
														{formatDateTime(appointment.scheduled_at)}
													</span>
												</div>
												<h3 className="font-display text-3xl text-white">{appointment.title}</h3>
												<p className="max-w-3xl text-sm text-slate-300">
													{appointment.description || "No additional appointment note captured."}
												</p>
											</div>
											<div className="grid gap-3 text-sm text-slate-300">
												<div className="rounded-[1.1rem] border border-white/10 bg-white/5 px-4 py-3">
													<strong className="mr-2 text-white">Doctor</strong>
													{appointment.doctor_name}
												</div>
												<div className="rounded-[1.1rem] border border-white/10 bg-white/5 px-4 py-3">
													<strong className="mr-2 text-white">Patient</strong>
													{appointment.patient_name}
												</div>
											</div>
										</div>
									</div>
								))
							) : (
								<div className="rounded-[1.6rem] border border-dashed border-white/10 bg-midnight-900/60 px-4 py-8 text-center text-sm text-slate-500">
									No appointments yet. Once you book a visit, the timeline will appear here.
								</div>
							)}
						</div>
					)}
				</section>

				{isAdmin ? (
					<section className="rounded-[2rem] border border-coral-400/[0.15] bg-coral-500/5 p-6 shadow-panel">
						<div className="flex items-start gap-4">
							<div className="rounded-2xl bg-coral-500/[0.15] p-3 text-coral-200">
								<ShieldAlert className="h-5 w-5" />
							</div>
							<div>
								<h2 className="font-display text-3xl text-white">Incident-ready by design</h2>
								<p className="mt-3 max-w-3xl text-sm leading-7 text-slate-300">
									The admin workspace is intentionally operational. It complements the assignment’s incident simulation by making degraded service states and scheduling disruption visible where decisions actually happen.
								</p>
							</div>
						</div>
					</section>
				) : null}
			</div>
		</AppShell>
	);
}

type AppointmentCardProps = {
	appointment: Appointment;
};

function AppointmentCard({ appointment }: AppointmentCardProps) {
	return (
		<div className="rounded-[1.4rem] border border-white/10 bg-white/5 p-4">
			<h4 className="font-semibold text-white">{appointment.title}</h4>
			<p className="mt-2 text-sm text-slate-300">{appointment.description}</p>
		</div>
	);
}
