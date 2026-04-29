import type { ReactNode } from "react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import { Toaster } from "sonner";
import { LoadingScreen } from "./components/shared/loading-screen";
import { AuthProvider, useAuth } from "./providers/auth-provider";
import { LandingPage } from "./pages/landing-page";
import { OperationsPage } from "./pages/operations-page";
import { WorkspacePage } from "./pages/workspace-page";

const queryClient = new QueryClient({
	defaultOptions: {
		queries: {
			refetchOnWindowFocus: false,
			retry: 1,
		},
	},
});

function ProtectedRoute({ children }: { children: ReactNode }) {
	const { user, isBooting } = useAuth();

	if (isBooting) {
		return <LoadingScreen />;
	}

	if (!user) {
		return <Navigate to="/" replace />;
	}

	return <>{children}</>;
}

function AppRoutes() {
	const { isBooting } = useAuth();

	if (isBooting) {
		return <LoadingScreen />;
	}

	return (
		<Routes>
			<Route path="/" element={<LandingPage />} />
			<Route
				path="/workspace"
				element={
					<ProtectedRoute>
						<WorkspacePage />
					</ProtectedRoute>
				}
			/>
			<Route path="/operations" element={<OperationsPage />} />
			<Route path="*" element={<Navigate to="/" replace />} />
		</Routes>
	);
}

export default function App() {
	return (
		<QueryClientProvider client={queryClient}>
			<AuthProvider>
				<BrowserRouter>
					<AppRoutes />
				</BrowserRouter>
				<Toaster
					position="top-right"
					richColors
					theme="dark"
					toastOptions={{
						style: {
							borderRadius: "18px",
							border: "1px solid rgba(255,255,255,0.08)",
							background: "rgba(8,18,30,0.92)",
						},
					}}
				/>
			</AuthProvider>
		</QueryClientProvider>
	);
}
