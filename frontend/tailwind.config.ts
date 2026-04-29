import type { Config } from "tailwindcss";

const config: Config = {
	content: ["./index.html", "./src/**/*.{ts,tsx}"],
	theme: {
		extend: {
			fontFamily: {
				display: ["Fraunces", "serif"],
				body: ["Space Grotesk", "sans-serif"],
			},
			colors: {
				midnight: {
					950: "#08121E",
					900: "#102338",
					800: "#183755",
				},
				lagoon: {
					500: "#0EA5A3",
					400: "#38C6BE",
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
				panel: "0 24px 70px rgba(8, 18, 30, 0.16)",
				glow: "0 0 0 1px rgba(255,255,255,0.06), 0 24px 60px rgba(7, 14, 24, 0.18)",
			},
			backgroundImage: {
				"mesh-radial":
					"radial-gradient(circle at top left, rgba(14,165,163,0.16), transparent 30%), radial-gradient(circle at 82% 12%, rgba(242,108,79,0.14), transparent 28%), radial-gradient(circle at 50% 100%, rgba(229,154,46,0.12), transparent 32%)",
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
