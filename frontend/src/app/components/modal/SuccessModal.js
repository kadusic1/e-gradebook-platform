import React from "react";
import Modal from "./Modal";
import Title from "../common/Title";
import Button from "../common/Button";

const SuccessModal = ({ children, onClose, colorConfig }) => (
  <Modal onClose={onClose}>
    <Title colorConfig={colorConfig}>Uspjeh</Title>
    <div>{children}</div>
    <div className="flex justify-center mt-4">
      <Button onClick={onClose} colorConfig={colorConfig}>
        Zatvori
      </Button>
    </div>
  </Modal>
);

export default SuccessModal;
