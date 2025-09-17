import React from "react";
import Modal from "./Modal";
import Form from "../form/Form";
import Title from "../common/Title";
import Note from "../common/Note";

export default function CreateUpdateModal({
  title,
  fields,
  onClose,
  onSave,
  initialValues = {},
  noteMessage,
  fixedHeight = false,
  colorConfig,
  showSave = true,
  internalFormOverflow = true,
}) {
  const handleSubmit = (data) => {
    onSave(data);
    onClose();
  };

  return (
    <Modal onClose={onClose} overlayClassName="bg-white/70 backdrop-blur-sm">
      <Title colorConfig={colorConfig}>{title}</Title>
      {noteMessage && <Note className="mb-4">({noteMessage})</Note>}
      <Form
        fields={fields}
        onSubmit={handleSubmit}
        onClose={onClose}
        initialValues={initialValues}
        fixedHeight={fixedHeight}
        colorConfig={colorConfig}
        showSave={showSave}
        internalOverflow={internalFormOverflow}
      ></Form>
    </Modal>
  );
}
