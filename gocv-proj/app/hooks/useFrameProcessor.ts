import { useRef, useState, useCallback } from 'react';

export type ProcessingMode = 'glitch' | 'blur' | 'sketch' | 'emboss' | 'wave-ripple' | 'pixelate';

export function useFrameProcessor() {
  const [isProcessing, setIsProcessing] = useState(false);
  const [processedImage, setProcessedImage] = useState<string | null>(null);
  const [mode, setMode] = useState<ProcessingMode>('glitch');
  const [snapshotLabel, setSnapshotLabel] = useState<string | null>(null);
  const intervalRef = useRef<NodeJS.Timeout | null>(null);
  const modeRef = useRef<ProcessingMode>('glitch');
  const countRef = useRef<number>(20);

  const updateMode = useCallback((newMode: ProcessingMode) => {
    modeRef.current = newMode;
    setMode(newMode);
  }, []);

  const captureAndSend = useCallback(async (videoElement: HTMLVideoElement) => {
    const canvas = document.createElement('canvas');
    canvas.width = videoElement.videoWidth;
    canvas.height = videoElement.videoHeight;

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    ctx.drawImage(videoElement, 0, 0);

    const blob = await new Promise<Blob | null>((resolve) =>
      canvas.toBlob(resolve, 'image/jpeg', 0.8)
    );
    if (!blob) {
      console.error('Failed to create blob from canvas');
      return;
    }

    const currentMode = modeRef.current;
    console.log(`Sending frame: ${blob.size} bytes, endpoint: /${currentMode}`);

    try {
      const response = await fetch(`http://localhost:4000/${currentMode}`, {
        method: 'POST',
        headers: { 'Content-Type': 'image/jpeg' },
        body: blob,
      });

      if (!response.ok) {
        const errorText = await response.text();
        console.error(`Server error ${response.status}: ${errorText}`);
        return;
      }

      const responseBlob = await response.blob();
      const imageUrl = URL.createObjectURL(responseBlob);
      setProcessedImage(imageUrl);
    } catch (err) {
      console.error('Failed to send frame:', err);
    }
  }, []);

  const startProcessing = useCallback((videoElement: HTMLVideoElement) => {
    setIsProcessing(true);
    setSnapshotLabel('Taking snapshot...');
    captureAndSend(videoElement).then(() => {
      countRef.current = 20;
      setSnapshotLabel('Next snapshot in 20s');
    });

    intervalRef.current = setInterval(() => {
      countRef.current -= 1;
      if (countRef.current <= 0) {
        setSnapshotLabel('Taking snapshot...');
        captureAndSend(videoElement).then(() => {
          countRef.current = 20;
          setSnapshotLabel('Next snapshot in 20s');
        });
      } else {
        setSnapshotLabel(`Next snapshot in ${countRef.current}s`);
      }
    }, 1000);
  }, [captureAndSend]);

  const stopProcessing = useCallback(() => {
    if (intervalRef.current) {
      clearInterval(intervalRef.current);
      intervalRef.current = null;
    }
    setIsProcessing(false);
    setProcessedImage(null);
    setSnapshotLabel(null);
  }, []);

  return { isProcessing, processedImage, mode, setMode: updateMode, startProcessing, stopProcessing, snapshotLabel };
}
