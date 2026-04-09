import type { CameraStatus } from "../hooks/useCamera";

export interface CameraControlsProps {
  status: CameraStatus;
  onStart: () => void;
  onStop: () => void;
}

export function CameraControls({
  status,
  onStart,
  onStop,
}: CameraControlsProps) {
  return (
    <div className="flex gap-3 items-center">
      <button
        onClick={onStart}
        disabled={status === "starting" || status === "running"}
        className="px-3 py-2"
      >
        Start
      </button>
      <button
        onClick={onStop}
        disabled={status === "idle" || status === "starting"}
        className="px-3 py-2"
      >
        Stop
      </button>


    </div>
  );
}
