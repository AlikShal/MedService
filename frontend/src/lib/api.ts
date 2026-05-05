import type {
	Appointment,
	AppointmentStatus,
	AuthPayload,
	ChatMessage,
	ChatThread,
	Doctor,
	PatientProfile,
	ServiceHealth,
	User,
} from "../types/domain";

export class ApiError extends Error {
	status: number;
	data: unknown;

	constructor(message: string, status: number, data: unknown) {
		super(message);
		this.name = "ApiError";
		this.status = status;
		this.data = data;
	}
}

type RequestOptions = {
	method?: string;
	body?: unknown;
	token?: string;
	signal?: AbortSignal;
};

async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
	const headers = new Headers();
	headers.set("Content-Type", "application/json");

	if (options.token) {
		headers.set("Authorization", `Bearer ${options.token}`);
	}

	const response = await fetch(path, {
		method: options.method ?? "GET",
		headers,
		body: options.body ? JSON.stringify(options.body) : undefined,
		signal: options.signal,
	});

	const text = await response.text();
	const data = text ? JSON.parse(text) : null;

	if (!response.ok) {
		throw new ApiError(
			typeof data?.error === "string" ? data.error : `Request failed with ${response.status}`,
			response.status,
			data,
		);
	}

	return data as T;
}

export const api = {
	login(payload: { email: string; password: string }) {
		return request<AuthPayload>("/api/auth/login", { method: "POST", body: payload });
	},

	register(payload: { full_name: string; email: string; password: string }) {
		return request<AuthPayload>("/api/auth/register", { method: "POST", body: payload });
	},

	me(token: string) {
		return request<{ user: User }>("/api/auth/me", { token });
	},

	getDoctors() {
		return request<Doctor[]>("/api/doctors");
	},

	createDoctor(
		payload: { full_name: string; specialization: string; email: string; office: string },
		token: string,
	) {
		return request<Doctor>("/api/doctors", { method: "POST", body: payload, token });
	},

	getProfile(token: string) {
		return request<PatientProfile>("/api/patients/me", { token });
	},

	upsertProfile(
		payload: { full_name: string; phone: string; date_of_birth: string; notes: string },
		token: string,
	) {
		return request<PatientProfile>("/api/patients/me", { method: "PUT", body: payload, token });
	},

	getAppointments(token: string) {
		return request<Appointment[]>("/api/appointments", { token });
	},

	createAppointment(
		payload: { title: string; description: string; doctor_id: string; scheduled_at: string },
		token: string,
	) {
		return request<Appointment>("/api/appointments", { method: "POST", body: payload, token });
	},

	updateAppointmentStatus(id: string, status: AppointmentStatus, token: string) {
		return request<{ message: string }>(`/api/appointments/${id}/status`, {
			method: "PATCH",
			body: { status },
			token,
		});
	},

	getChatThreads(token: string) {
		return request<ChatThread[]>("/api/chat/threads", { token });
	},

	createChatThread(payload: { appointment_id: string; subject?: string }, token: string) {
		return request<ChatThread>("/api/chat/threads", { method: "POST", body: payload, token });
	},

	getChatMessages(threadId: string, token: string) {
		return request<ChatMessage[]>(`/api/chat/threads/${threadId}/messages`, { token });
	},

	sendChatMessage(threadId: string, body: string, token: string) {
		return request<ChatMessage>(`/api/chat/threads/${threadId}/messages`, {
			method: "POST",
			body: { body },
			token,
		});
	},

	async getPlatformHealth(): Promise<ServiceHealth[]> {
		const endpoints = [
			"/api/health/auth",
			"/api/health/patient",
			"/api/health/doctor",
			"/api/health/appointment",
			"/api/health/chat",
		];

		const settled = await Promise.allSettled(
			endpoints.map(async (path) => request<ServiceHealth>(path)),
		);

		return settled.map((result, index) => {
			if (result.status === "fulfilled") {
				return result.value;
			}

			const fallbackName = endpoints[index].split("/").at(-1) ?? "service";
			return {
				service: `${fallbackName}-service`,
				status: "down",
				error: result.reason instanceof Error ? result.reason.message : "Unreachable",
			} satisfies ServiceHealth;
		});
	},
};
