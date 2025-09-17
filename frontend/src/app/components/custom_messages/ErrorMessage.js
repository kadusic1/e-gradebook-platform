import { FaExclamationCircle } from "react-icons/fa";
import { formatMessage } from "../modal/ErrorModal";

export default function ErrorMessage({ children, className = "" }) {
  if (!children) return null;

  const formattedMessage = formatMessage(children);

  return (
    <div
      className={`flex justify-center items-center gap-2 bg-red-100 border border-red-400 text-red-700 px-4 py-2 rounded mb-2 animate-fade-in ${className}`}
    >
      <FaExclamationCircle className="text-red-500" />
      <span className="font-medium">{formattedMessage}</span>
    </div>
  );
}
