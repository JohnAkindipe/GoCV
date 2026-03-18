"use client";
import { useState } from "react";
import { useCamera } from "./hooks/useCamera";
import { useWebRTC } from "./hooks/useWebRTC";
import { CameraControls } from "./components/CameraControls";
import { CameraVideo } from "./components/CameraVideo";
import { ErrorMessage } from "./components/ErrorMessage";
import parse from "html-react-parser";

export default function Home() {
  const [playAudio, setPlayAudio] = useState(false);
  const [serverResponse, setServerResponse] = useState<string>("");
  const [isLoading, setIsLoading] = useState(false);
  const [fetchError, setFetchError] = useState<string | null>(null);

  const { videoRef, status, error, stream, startCamera, stopCamera } = useCamera({
    autoStart: true,
    playAudio,
    facingMode: "user",
  });

  const { isConnected, processedImage, asciiText, expectingASCII, startStreaming, stopStreaming, toggleASCIIMode } = useWebRTC();

  const handleStartStreaming = async () => {
    if (!stream || !videoRef.current) {
      setFetchError("No camera stream available");
      return;
    }
    try {
      setFetchError(null);
      await startStreaming(stream, videoRef.current);
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setFetchError(message);
    }
  };

  const testServerEndpoint = async () => {
    setIsLoading(true);
    setFetchError(null);
    setServerResponse("");

    try {
      const response = await fetch("http://161.35.36.3:4000/test");
      
      if (!response.ok) {
        throw new Error(`Server responded with status: ${response.status}`);
      }

      const text = await response.text();
      setServerResponse(text);
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setFetchError(message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <main style={{ padding: 24, display: "grid", gap: 16 }}>
      <h1 style={{ fontSize: 20, fontWeight: 600 }}>Live Camera Feed with Glitch Effect</h1>

      <CameraControls
        status={status}
        playAudio={playAudio}
        onStart={startCamera}
        onStop={stopCamera}
        onToggleAudio={setPlayAudio}
      />

      <ErrorMessage message={error} />

      <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 16 }}>
        <div>
          <h3 style={{ margin: "0 0 8px 0" }}>Live Camera</h3>
          <CameraVideo videoRef={videoRef} />
        </div>
        
        <div>
          <h3 style={{ margin: "0 0 8px 0" }}>
            {expectingASCII ? "ASCII Art Output" : "Glitched Output"}
          </h3>
          {asciiText ? (
            <div className="term-container">
              {parse(asciiText)}
            </div>
          ) : processedImage ? (
            <img 
              src={processedImage} 
              alt="Processed" 
              style={{ 
                width: "100%", 
                maxWidth: 640, 
                border: "2px solid #4CAF50",
                borderRadius: 8 
              }} 
            />
          ) : (
            <div style={{ 
              width: "100%", 
              maxWidth: 640, 
              aspectRatio: "4/3",
              background: "#f0f0f0", 
              display: "flex", 
              alignItems: "center", 
              justifyContent: "center",
              borderRadius: 8,
              color: "#666"
            }}>
              {isConnected ? "Processing next frame..." : "Start streaming to see effect"}
            </div>
          )}
        </div>
      </div>
      <div style={{ display: "grid", gap: 12 }}>
        <div style={{ display: "flex", gap: 8, alignItems: "center", flexWrap: "wrap" }}>
          <button
            onClick={handleStartStreaming}
            disabled={!stream || isConnected}
            style={{ padding: "10px 16px", fontSize: 14 }}
          >
            {isConnected ? "Streaming..." : "Start Glitch Stream (15s intervals)"}
          </button>
          <button
            onClick={stopStreaming}
            disabled={!isConnected}
            style={{ padding: "10px 16px", fontSize: 14 }}
          >
            Stop Stream
          </button>
          <button
            onClick={toggleASCIIMode}
            style={{ 
              padding: "10px 16px", 
              fontSize: 14,
              background: expectingASCII ? "#4CAF50" : "#f0f0f0",
              color: expectingASCII ? "#fff" : "#333",
              border: "1px solid #ccc",
              borderRadius: 4
            }}
          >
            {expectingASCII ? "ASCII Mode: ON" : "ASCII Mode: OFF"}
          </button>
          {isConnected && (
            <span style={{ color: "green", fontSize: 14 }}>
              ✓ Connected - Processing every 15s
            </span>
          )}
        </div>

        <button
          onClick={testServerEndpoint}
          disabled={isLoading}
          style={{ padding: "10px 16px", fontSize: 14 }}
        >
          {isLoading ? "Fetching..." : "Test Server Endpoint"}
        </button>

        {fetchError && (
          <p style={{ color: "crimson", margin: 0 }}>
            Error: {fetchError}
          </p>
        )}

        {serverResponse && (
          <div
            style={{
              padding: 16,
              background: "#f0f0f0",
              borderRadius: 8,
              color: "#333",
            }}
          >
            <strong>Server Response:</strong>
            <p style={{ margin: "8px 0 0", fontFamily: "monospace" }}>
              {serverResponse}
            </p>
          </div>
        )}
      </div>
    </main>
  );
}
