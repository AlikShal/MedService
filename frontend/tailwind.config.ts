import type { Config } from "tailwindcss";

const config: Config = {
	content: ["./index.html", "./src/**/*.{ts,tsx}"],
	theme: {
		extend: {
			fontFamily: {
				display: ["Inter", "ui-sans-serif", "system-ui", "sans-serif"],
				body: ["Inter", "ui-sans-serif", "system-ui", "sans-serif"],
			},
			colors: {
				midnight: {
					950: "#070d19",
					900: "#0c1422",
					800: "#111c2e",
				},
				lagoon: {
					500: "#14b8a6",
					400: "#22d3c5",
				},
				coral: {
					500: "#F26C4F",
					400: "#FF9878",
				},
				sand: {
					100: "#F6F0E6",
					200: "#EADFCF",
				},
				sun: {
					400: "#F2B75A",
					500: "#E59A2E",
				},
			},
			boxShadow: {
				panel: "0 18px 55px rgba(0, 0, 0, 0.18)",
				glow: "0 0 0 1px rgba(20,184,166,0.12), 0 18px 45px rgba(20,184,166,0.08)",
			},
			backgroundImage: {
				"mesh-radial":
					"linear-gradient(180deg, #070d19 0%, #070d19 100%)",
			},
			keyframes: {
				float: {
					"0%, 100%": { transform: "translateY(0px)" },
					"50%": { transform: "translateY(-12px)" },
				},
				shine: {
					"0%": { backgroundPosition: "0% 50%" },
					"50%": { backgroundPosition: "100% 50%" },
					"100%": { backgroundPosition: "0% 50%" },
				},
				pulseSoft: {
					"0%, 100%": { opacity: "0.6" },
					"50%": { opacity: "1" },
				},
			},
			animation: {
				float: "float 8s ease-in-out infinite",
				shine: "shine 10s ease infinite",
				"pulse-soft": "pulseSoft 3.2s ease-in-out infinite",
			},
		},
	},
	plugins: [],
};

export default config;
