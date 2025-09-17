import { FaArrowLeft } from "react-icons/fa";
import Button from "./Button";

export default function BackButton({ onClick, colorConfig }) {
  return (
    <Button
      onClick={onClick}
      color="secondary"
      className="flex items-center justify-center gap-2"
      colorConfig={colorConfig}
    >
      <FaArrowLeft className="mr-2" />
      Nazad
    </Button>
  );
}
