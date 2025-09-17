import React from "react";
import { FaPlus } from "react-icons/fa";
import Button from "./Button";

export default function AddButton({ children, colorConfig, ...props }) {
  return (
    <Button {...props} colorConfig={colorConfig}>
      <div className="flex items-center gap-2">
        <FaPlus />
        {children}
      </div>
    </Button>
  );
}
