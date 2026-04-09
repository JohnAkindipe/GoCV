"use client";
import { usePathname } from "next/navigation";

export default function Navbar() {
  const pathname = usePathname();

  return (
    <nav className="flex items-center justify-between px-8 py-4 bg-white border-b-2 border-b-[#F5824A]">
      <a href="/" className="text-xl font-bold tracking-tight">iffekt</a>
      <div className="flex gap-6 text-sm font-medium">
        <a href="#" className="hover:text-green-600 transition-colors" style={{ fontFamily: "var(--font-borel)" }}>deepdive</a>
        {pathname !== "/live-effects" && (
          <a href="/live-effects" className="hover:text-green-600 transition-colors" style={{ fontFamily: "var(--font-borel)" }}>live effects</a>
        )}
      </div>
    </nav>
  );
}
