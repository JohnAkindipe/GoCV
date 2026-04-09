export interface ErrorMessageProps {
  message: string | null;
}

export function ErrorMessage({ message }: ErrorMessageProps) {
  if (!message) return null;

  return <p className="text-red-600">Camera error: {message}</p>;
}
