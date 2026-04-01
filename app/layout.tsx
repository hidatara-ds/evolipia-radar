import type { Metadata } from "next";
import "./globals.css";

const BASE_URL = "https://evolipia-radar.vercel.app";

export const metadata: Metadata = {
  metadataBase: new URL(BASE_URL),
  title: "Evolipia Radar — AI Research Intelligence",
  description:
    "Real-time semantic clustering of global AI research signals. Evolipia filters the noise to bring you what actually moves markets.",
  icons: {
    icon: "/assets/icon.webp",
    shortcut: "/assets/icon.webp",
    apple: "/assets/icon.webp",
  },
  openGraph: {
    type: "website",
    url: BASE_URL,
    siteName: "Evolipia Radar",
    title: "Evolipia Radar — AI Research Intelligence",
    description:
      "Real-time semantic clustering of global AI research signals. Predict the future of AI innovation.",
    images: [
      {
        url: "/assets/evolipiaradar.webp",
        width: 1200,
        height: 630,
        alt: "Evolipia Radar — AI Research Intelligence Platform",
      },
    ],
  },
  twitter: {
    card: "summary_large_image",
    site: "@evolipia",
    title: "Evolipia Radar — AI Research Intelligence",
    description:
      "Real-time semantic clustering of global AI research signals. Predict the future of AI innovation.",
    images: ["/assets/evolipiaradar.webp"],
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <head>
        {/* Preload LCP image — mascot is the largest contentful element */}
        <link
          rel="preload"
          as="image"
          href="/assets/maskot1.webp"
          type="image/webp"
        />
      </head>
      <body className="min-h-screen bg-black text-white antialiased font-sans">
        {children}
      </body>
    </html>
  );
}
