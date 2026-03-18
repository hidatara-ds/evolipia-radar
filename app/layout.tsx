import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Evolipia Radar",
  description: "AI Research Intelligence - Your source for the latest AI/ML news and research",
  icons: {
    icon: "/assets/icon.png",
    shortcut: "/assets/icon.png",
    apple: "/assets/icon.png",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="min-h-screen bg-black text-white antialiased font-sans">
        {children}
      </body>
    </html>
  );
}
