export interface CameraVideoProps {
  videoRef: React.RefObject<HTMLVideoElement | null>;
  className?: string;
}

export function CameraVideo({ videoRef, className }: CameraVideoProps) {
  return (
    <video
      ref={videoRef}
      playsInline
      className={`w-full max-w-[900px] aspect-video bg-[#111] rounded-xl ${className ?? ""}`}
    />
  );
}
