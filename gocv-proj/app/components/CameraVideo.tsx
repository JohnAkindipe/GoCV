import type { CSSProperties } from "react";

export interface CameraVideoProps {
  videoRef: React.RefObject<HTMLVideoElement | null>;
  style?: CSSProperties;
}

export function CameraVideo({ videoRef, style }: CameraVideoProps) {
  return (
    <video
      ref={videoRef}
      playsInline
      style={{
        width: "min(900px, 100%)",
        aspectRatio: "16 / 9",
        background: "#111",
        borderRadius: 12,
        ...style,
      }}
    />
  );
}
