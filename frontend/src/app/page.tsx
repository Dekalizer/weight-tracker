"use client";
import Link from "next/link";
import './globals.css';

export default function Home() {
  return (
    <main className="flex flex-col items-center justify-center min-h-screen bg-black">
      <div className="mb-6 text-5xl animate-bounce">⚡️</div>
      <h1 className="text-4xl font-bold mb-6 drop-shadow-lg text-orange">Weight Tracker</h1>
      <nav className="flex gap-4">
        <Link
          href="/weights"
          className="px-6 py-2 rounded-lg bg-orange text-black font-semibold shadow-lg hover:bg-orange/80 transition-all drop-shadow-sm focus:outline-none"
        >
          View Weights
        </Link>
        <Link
          href="/add"
          className="px-6 py-2 rounded-lg bg-green text-black font-semibold shadow-lg hover:bg-green/80 transition-all drop-shadow-sm focus:outline-none"
        >
          Add Entry
        </Link>
      </nav>
    </main>
  );
}
