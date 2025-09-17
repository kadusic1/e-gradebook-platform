import Text from "../common/Text";
import { FaBoxOpen } from "react-icons/fa";
import { CardContainer } from "./DynamicCard";

export default function EmptyState({
  message = "Nema podataka za prikaz.",
  className = "",
}) {
  return (
    <CardContainer
      minWidth="max-w-lg mx-auto text-center"
      className={className}
    >
      <FaBoxOpen className="text-6xl text-gray-400 mb-4 mx-auto" />
      <Text className="text-center font-bold text-xl">{message}</Text>
    </CardContainer>
  );
}
