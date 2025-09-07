"use client";
import { useEffect, useState } from "react";
import '../globals.css';

type WeightEntry = {
  id: number;
  date: string;
  weight_kg: number;
};

export default function WeightsPage() {
  const [weights, setWeights] = useState<WeightEntry[]>([]);

  useEffect(() => {
    async function fetchWeights() {
      const res = await fetch("http://localhost:8080/weights");
      const data = await res.json();
      setWeights(data);
    }
    fetchWeights();
  }, []);

  return (
    <main className="min-h-screen bg-black p-8 flex flex-col items-center">
      <h1 className="text-4xl font-extrabold mb-10 text-orange-500 drop-shadow-lg text-center">
        Your Weight Progress
      </h1>
      <h1 className="text-4xl font-bold text-red-600">Test</h1>

      {weights.length === 0 ? (
        <p className="text-gray-400 text-lg animate-pulse">Loading your entries...</p>
      ) : (
        <ul className="w-full max-w-xl space-y-5 mx-auto">
          {weights.map((w, i) => (
            <li
              key={w.id}
              className="group flex justify-between items-center bg-white/10 backdrop-blur-md rounded-2xl shadow-lg border-l-8 border-orange-500 p-6 transition-all duration-300 hover:scale-105 hover:shadow-orange-500/50"
              style={{
                animation: `fadeIn 0.4s ease both`,
                animationDelay: `${i * 0.07}s`,
              }}
              title={`Recorded on ${new Date(w.date).toLocaleDateString()}`}
            >
              <div>
                <span className="text-2xl font-bold text-white drop-shadow-md">{w.weight_kg} kg</span>
                <div className="text-sm text-gray-200 mt-1">
                  {new Date(w.date).toLocaleDateString(undefined, {
                    weekday: "short",
                    year: "numeric",
                    month: "short",
                    day: "numeric",
                  })}
                </div>
              </div>
            </li>
          ))}
        </ul>
      )}

      <style>{`
        @keyframes fadeIn {
          from {
            opacity: 0;
            transform: translateY(24px);
          }
          to {
            opacity: 1;
            transform: none;
          }
        }
      `}</style>
    </main>
  );
}