import Text from "../common/Text";
import { FaBoxOpen } from "react-icons/fa";

export default function EmptyStateContainerless({
  message = "Nema podataka za prikaz.",
  className = "",
}) {
  return (
    <div className={`max-w-lg mx-auto text-center ${className}`}>
      <FaBoxOpen className="text-6xl text-gray-400 mb-4 mx-auto" />
      <Text className="text-center font-bold text-xl">{message}</Text>
    </div>
  );
}
