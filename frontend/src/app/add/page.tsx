"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import '../globals.css';

export default function AddPage() {
  const [weight, setWeight] = useState("");
  const [date, setDate] = useState("");
  const router = useRouter();

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();

    try {
      const res = await fetch("http://localhost:8080/add", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          weight: parseFloat(weight),
          date: date || new Date().toISOString().split("T")[0], // default today
        }),
      });

      if (!res.ok) {
        console.error("Failed to submit:", res.status);
        return;
      }

      router.push("/weights");
    } catch (err) {
      console.error("Network error:", err);
    }
  }

  return (
    <main className="flex flex-col items-center min-h-screen justify-center bg-black">
      <form
        onSubmit={handleSubmit}
        className="bg-white/10 backdrop-blur-lg p-8 rounded-2xl shadow-lg flex flex-col gap-4 w-full max-w-sm"
      >
        <h1 className="text-3xl mb-5 text-green font-bold">Add New Entry</h1>
        <input
          type="number"
          step="0.1"
          value={weight}
          onChange={e => setWeight(e.target.value)}
          className="border border-green focus:border-orange focus:ring-2 focus:ring-orange/70 rounded-lg p-3 bg-black text-white placeholder-white/60 outline-none transition"
          placeholder="Enter weight (kg)"
          required
        />
        <input
          type="date"
          value={date}
          onChange={e => setDate(e.target.value)}
          className="border border-orange focus:border-green focus:ring-2 focus:ring-green/70 rounded-lg p-3 bg-black text-white placeholder-white/60 outline-none transition"
        />
        <button
          type="submit"
          className="px-6 py-3 rounded-lg bg-orange text-black font-bold hover:bg-orange/90 transition shadow-lg"
        >
          Save
        </button>
      </form>
    </main>
  );
}
