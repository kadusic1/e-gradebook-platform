import React from "react";
import Modal from "./Modal";
import Button from "../common/Button";
import Title from "../common/Title";
import { FaTimes } from "react-icons/fa";
import { FaCheck } from "react-icons/fa";

export default function ConfirmModal({
  title,
  onClose,
  onConfirm,
  children,
  colorConfig,
}) {
  return (
    <Modal onClose={onClose}>
      <Title
        colorConfig={colorConfig}
        className="text-xl font-bold text-gray-800 mb-2"
      >
        {title || "Potvrda brisanja"}
      </Title>
      {children}
      <div className="flex flex-row gap-4 justify-center mt-5">
        <Button
          icon={FaCheck}
          onClick={onConfirm}
          color="secondary"
          className="w-[130px]"
          colorConfig={colorConfig}
        >
          Potvrdi
        </Button>
        <Button
          icon={FaTimes}
          color="ternary"
          onClick={onClose}
          className="w-[130px]"
          colorConfig={colorConfig}
        >
          Odustani
        </Button>
      </div>
    </Modal>
  );
}
