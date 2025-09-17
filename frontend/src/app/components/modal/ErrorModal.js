import React from "react";
import Modal from "./Modal";
import Title from "../common/Title";
import Button from "../common/Button";

export function formatMessage(message) {
  if (typeof message !== "string") return message;

  // Trim whitespace first, then capitalize
  let formatted = message.trim();
  formatted = formatted.charAt(0).toUpperCase() + formatted.slice(1);

  // Check if it ends with a dot after trimming
  if (!formatted.endsWith(".")) {
    formatted += ".";
  }

  return formatted;
}

const ErrorModal = ({ children, onClose, colorConfig }) => (
  <Modal onClose={onClose}>
    <Title colorConfig={colorConfig}>Greška</Title>
    <div>{formatMessage(children)}</div>
    <div className="flex justify-center mt-4">
      <Button onClick={onClose} colorConfig={colorConfig}>
        Zatvori
      </Button>
    </div>
  </Modal>
);

export default ErrorModal;
