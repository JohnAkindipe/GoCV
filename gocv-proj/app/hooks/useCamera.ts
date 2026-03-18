import { useCallback, useEffect, useRef, useState } from "react";

export type CameraStatus = "idle" | "starting" | "running" | "error";

export interface UseCameraOptions {
  autoStart?: boolean;
  playAudio?: boolean;
  facingMode?: "user" | "environment";
}

export interface UseCameraReturn {
  videoRef: React.RefObject<HTMLVideoElement | null>;
  status: CameraStatus;
  error: string | null;
  stream: MediaStream | null;
  startCamera: () => Promise<void>;
  stopCamera: () => void;
}

export function useCamera(options: UseCameraOptions = {}): UseCameraReturn {
  const { autoStart = true, playAudio = false, facingMode = "user" } = options;

  const videoRef = useRef<HTMLVideoElement | null>(null);
  const streamRef = useRef<MediaStream | null>(null);

  const [status, setStatus] = useState<CameraStatus>("idle");
  const [error, setError] = useState<string | null>(null);

  const stopCamera = useCallback(() => {
    const stream = streamRef.current;
    if (stream) {
      for (const track of stream.getTracks()) {
        track.stop();
      }
    }
    streamRef.current = null;

    const video = videoRef.current;
    if (video) {
      video.srcObject = null;
    }

    setStatus("idle");
  }, []);

  const startCamera = useCallback(async () => {
    setError(null);
    setStatus("starting");

    try {
      if (!navigator.mediaDevices?.getUserMedia) {
        throw new Error("getUserMedia is not supported in this browser.");
      }

      stopCamera();

      const stream = await navigator.mediaDevices.getUserMedia({
        video: { facingMode },
        audio: {
          echoCancellation: true,
          noiseSuppression: true,
          autoGainControl: true,
        },
      });

      streamRef.current = stream;
      const video = videoRef.current;
      if (!video) {
        throw new Error("Video element not ready.");
      }

      video.srcObject = stream;

      // Keep muted by default so autoplay works and to avoid feedback.
      video.muted = !playAudio;
      await video.play();
      setStatus("running");
    } catch (err) {
      stopCamera();
      const message = err instanceof Error ? err.message : String(err);
      setError(message);
      setStatus("error");
    }
  }, [playAudio, facingMode, stopCamera]);

  useEffect(() => {
    const video = videoRef.current;
    if (!video) return;

    video.muted = !playAudio;
    // If the browser blocks unmuting without a gesture, the toggle still flips,
    // but audio may not start until the user presses Start.
    void video.play().catch(() => {});
  }, [playAudio]);

  useEffect(() => {
    if (autoStart) {
      void startCamera();
    }
    return () => stopCamera();
  }, [autoStart, startCamera, stopCamera]);

  return {
    videoRef,
    status,
    error,
    stream: streamRef.current,
    startCamera,
    stopCamera,
  };
}
