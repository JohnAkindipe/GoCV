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
    <div style={{ display: "flex", gap: 12, alignItems: "center" }}>
      <button
        onClick={onStart}
        disabled={status === "starting" || status === "running"}
        style={{ padding: "8px 12px" }}
      >
        Start
      </button>
      <button
        onClick={onStop}
        disabled={status === "idle" || status === "starting"}
        style={{ padding: "8px 12px" }}
      >
        Stop
      </button>

      <span style={{ opacity: 0.8 }}>
        Status: <strong>{status}</strong>
      </span>
    </div>
  );
}
