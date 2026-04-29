import {
	createContext,
	useContext,
	useEffect,
	useRef,
	useState,
	type PropsWithChildren,
} from "react";
import { api } from "../lib/api";
import type { AuthPayload, User } from "../types/domain";

const TOKEN_KEY = "medsync_token";

type AuthContextValue = {
	token: string;
	user: User | null;
	isBooting: boolean;
	login: (payload: { email: string; password: string }) => Promise<AuthPayload>;
	register: (payload: { full_name: string; email: string; password: string }) => Promise<AuthPayload>;
	logout: () => void;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: PropsWithChildren) {
	const [token, setToken] = useState<string>(() => localStorage.getItem(TOKEN_KEY) ?? "");
	const [user, setUser] = useState<User | null>(null);
	const [isBooting, setIsBooting] = useState(true);
	const bootRef = useRef(false);

	useEffect(() => {
		let active = true;

		async function bootstrap() {
			if (!token) {
				if (active) {
					setUser(null);
					setIsBooting(false);
				}
				return;
			}

			try {
				const { user: nextUser } = await api.me(token);
				if (active) {
					setUser(nextUser);
				}
			} catch {
				localStorage.removeItem(TOKEN_KEY);
				if (active) {
					setToken("");
					setUser(null);
				}
			} finally {
				if (active) {
					setIsBooting(false);
				}
			}
		}

		if (!bootRef.current || token) {
			bootRef.current = true;
			void bootstrap();
		}

		return () => {
			active = false;
		};
	}, [token]);

	const persistSession = (payload: AuthPayload) => {
		localStorage.setItem(TOKEN_KEY, payload.token);
		setToken(payload.token);
		setUser(payload.user);
		return payload;
	};

	const value: AuthContextValue = {
		token,
		user,
		isBooting,
		async login(payload) {
			return persistSession(await api.login(payload));
		},
		async register(payload) {
			return persistSession(await api.register(payload));
		},
		logout() {
			localStorage.removeItem(TOKEN_KEY);
			setToken("");
			setUser(null);
		},
	};

	return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
	const context = useContext(AuthContext);
	if (!context) {
		throw new Error("useAuth must be used within AuthProvider");
	}
	return context;
}
