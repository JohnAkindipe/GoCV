import type { Metadata } from "next";
import { Geist, Geist_Mono, Borel, Rubik_Wet_Paint } from "next/font/google";
import "./globals.css";
import Navbar from "./components/Navbar";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

const borel = Borel({
  variable: "--font-borel",
  subsets: ["latin"],
  weight: "400",
});

const rubikWetPaint = Rubik_Wet_Paint({
  variable: "--font-rubik-wet-paint",
  subsets: ["latin"],
  weight: "400",
});

export const metadata: Metadata = {
  title: "iffekt",
  description: "Live camera effects powered by GoCV",
};

const currentYear = new Date().getFullYear();

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} ${borel.variable} ${rubikWetPaint.variable} antialiased flex flex-col min-h-screen`}
      >
        <Navbar />

        <div className="flex-1">
          {children}
        </div>

        <footer className="px-8 py-4 bg-white border-t-2 border-t-[#F5824A] text-center text-sm text-gray-500">
          &copy; iffekts {currentYear}. All rights reserved.
        </footer>
      </body>
    </html>
  );
}
