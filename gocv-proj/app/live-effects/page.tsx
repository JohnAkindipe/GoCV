"use client";
import { useCamera } from "../hooks/useCamera";
import { useFrameProcessor } from "../hooks/useFrameProcessor";
import { CameraVideo } from "../components/CameraVideo";
import { ErrorMessage } from "../components/ErrorMessage";

export default function LiveEffects() {
  const { videoRef, status, error, startCamera, stopCamera } = useCamera({
    autoStart: true,
    facingMode: "user",
  });

  const { isProcessing, processedImage, mode, setMode, startProcessing, stopProcessing, snapshotLabel } = useFrameProcessor();

  const handleStartProcessing = () => {
    if (!videoRef.current) return;
    startProcessing(videoRef.current);
  };

  return (
    <main className="p-6 grid gap-4 min-h-screen bg-[#f0e9d8]">
        <div className="mt-8 flex flex-col items-center gap-3">
          <h1
            className="text-4xl font-normal text-center tracking-tight heading-fade"
          >
            Check out cool effects below
          </h1>
          <div className="ornament-divider">
            <span>— ✦ —</span>
          </div>
        </div>

        <ErrorMessage message={error} />

        <div className="flex flex-wrap justify-center gap-4">
          <div style={{ minWidth: 300, maxWidth: "70vw", flex: 1 }}>
            <h2 className="mb-2 text-left font-bold" style={{ textDecoration: "underline wavy #f5824a", textUnderlineOffset: "4px" }}>Camera feed</h2>
            <CameraVideo videoRef={videoRef} />
          </div>

          <div style={{ minWidth: 300, maxWidth: "70vw", flex: 1 }}>
            <h2 className="mb-2 text-left font-bold" style={{ textDecoration: "underline wavy #f5824a", textUnderlineOffset: "4px" }}>Effect</h2>
            {processedImage ? (
              <img
                src={processedImage}
                alt="Processed"
                className="w-full max-w-[900px] aspect-video object-cover rounded-xl"
              />
            ) : (
              <div className="w-full max-w-[900px] aspect-video bg-gray-100 flex items-center justify-center rounded-xl text-gray-500">
                {isProcessing ? "Applying Effects..." : "Effect displayed here"}
              </div>
            )}
          </div>
        </div>
        <div className="flex gap-2 items-center justify-center flex-wrap">
          <button
            onClick={status === "running" ? stopCamera : startCamera}
            disabled={status === "starting"}
            className="px-4 py-2.5 text-sm camera-btn"
          >
            {status === "running" ? "Stop Camera" : "Start Camera"}
          </button>
          <button
            onClick={isProcessing ? stopProcessing : handleStartProcessing}
            disabled={status !== "running"}
            className="px-4 py-2.5 text-sm processing-btn"
          >
            {isProcessing ? "Stop Effects" : "Begin Effects"}
          </button>
          {snapshotLabel && (
            <span className="text-sm text-gray-600 italic">{snapshotLabel}</span>
          )}
        </div>
        <div className="flex flex-row justify-center items-center gap-2 mb-5">
          <label className="text-sm text-gray-700 underline">Click to select effect</label>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="20"
            height="20"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
            className="text-gray-700"
          >
            <line x1="5" y1="12" x2="19" y2="12" />
            <polyline points="12 5 19 12 12 19" />
          </svg>
          <select
            value={mode}
            onChange={(e) => setMode(e.target.value as import("../hooks/useFrameProcessor").ProcessingMode)}
            className="effect-select"
          >
            <option value="glitch">Glitch</option>
            <option value="blur">Blur</option>
            <option value="sketch">Sketch</option>
            <option value="emboss">Emboss</option>
            <option value="wave-ripple">Wave Ripple</option>
            <option value="pixelate">Pixelate</option>
          </select>
        </div>

    </main>
  );
}
