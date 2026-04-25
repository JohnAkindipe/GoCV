"use client";
import { useRef } from "react";
import { useScroll, useTransform, motion } from "framer-motion";

export default function Home() {
  const containerRef = useRef<HTMLDivElement>(null);

  const { scrollYProgress } = useScroll({
    target: containerRef,
    offset: ["start start", "end end"],
  });

  // As scroll goes 0→1, shift the strip left by 2 full screen widths (3 cards)
  const x = useTransform(scrollYProgress, [0, 1], ["0vw", "-200vw"]);

  return (
    <div ref={containerRef} className="relative h-[300vh] bg-[#f0e9d8]">
      <div className="sticky top-0 h-screen overflow-hidden">
        <motion.div style={{ x }} className="flex h-full">
          {[1, 2, 3].map((n) => (
            <div
              key={n}
              className="w-screen h-full flex-shrink-0 flex items-center justify-center"
            >
              <div className="w-80 h-64 bg-white rounded-2xl flex items-center justify-center text-8xl font-bold text-[#f5824a]"
                style={{ boxShadow: "6px 6px 0px #f5824a" }}
              >
                {n}
              </div>
            </div>
          ))}
        </motion.div>
      </div>
    </div>
  );
}


