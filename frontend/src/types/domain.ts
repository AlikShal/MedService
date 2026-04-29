export type Role = "admin" | "patient";

export interface User {
	id: string;
	full_name: string;
	email: string;
	role: Role;
	created_at: string;
}

export interface AuthPayload {
	token: string;
	user: User;
}

export interface Doctor {
	id: string;
	full_name: string;
	specialization: string;
	email: string;
	office: string;
	created_at: string;
}

export interface PatientProfile {
	id: string;
	user_id: string;
	full_name: string;
	email: string;
	phone: string;
	date_of_birth: string;
	notes: string;
	created_at: string;
	updated_at: string;
}

export type AppointmentStatus = "new" | "in_progress" | "done" | "cancelled";

export interface Appointment {
	id: string;
	title: string;
	description: string;
	doctor_id: string;
	doctor_name: string;
	patient_id: string;
	patient_name: string;
	scheduled_at: string;
	status: AppointmentStatus;
	created_at: string;
	updated_at: string;
}

export interface ServiceHealth {
	service: string;
	status: "ok" | "degraded" | "down";
	storage?: string;
	error?: string;
}
