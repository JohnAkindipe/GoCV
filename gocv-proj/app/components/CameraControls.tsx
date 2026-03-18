import type { CameraStatus } from "../hooks/useCamera";

export interface CameraControlsProps {
  status: CameraStatus;
  playAudio: boolean;
  onStart: () => void;
  onStop: () => void;
  onToggleAudio: (checked: boolean) => void;
}

export function CameraControls({
  status,
  playAudio,
  onStart,
  onStop,
  onToggleAudio,
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

      <label style={{ display: "flex", gap: 8, alignItems: "center" }}>
        <input
          type="checkbox"
          checked={playAudio}
          onChange={(e) => onToggleAudio(e.target.checked)}
        />
        Play audio
      </label>

      <span style={{ opacity: 0.8 }}>
        Status: <strong>{status}</strong>
      </span>
    </div>
  );
}
