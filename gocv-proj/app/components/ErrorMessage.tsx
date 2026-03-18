export interface ErrorMessageProps {
  message: string | null;
}

export function ErrorMessage({ message }: ErrorMessageProps) {
  if (!message) return null;

  return <p style={{ color: "crimson" }}>Camera error: {message}</p>;
}
