import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Evolipia Radar Dashboard",
  description: "Autonomous Web Intelligence System",
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
