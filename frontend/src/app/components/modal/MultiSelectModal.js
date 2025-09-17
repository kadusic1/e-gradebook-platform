import React, { useState } from "react";
import Modal from "./Modal";
import Title from "../common/Title";
import Button from "../common/Button";
import Text from "../common/Text";
import Note from "../common/Note";
import { FaTimes, FaSave } from "react-icons/fa";

export default function MultiSelectModal({
  title = "Odaberi stavke",
  items = [],
  onClose,
  onSave,
  labelField,
  searchPlaceholder = "Pretraži...",
  note,
  keyField = "id",
  showOnlyOnSearch = false,
  colorConfig,
}) {
  const [search, setSearch] = useState("");
  const [selected, setSelected] = useState([]);

  const filteredItems = [
    ...(items || []).filter(
      (item) =>
        item?.[labelField]?.toLowerCase().includes(search.toLowerCase()) ||
        selected.includes(item?.[keyField]),
    ),
  ];

  const toggleItem = (id) => {
    setSelected((prev) =>
      prev.includes(id) ? prev.filter((sid) => sid !== id) : [...prev, id],
    );
  };

  const handleSave = () => {
    onSave(selected);
    onClose();
  };

  return (
    <Modal onClose={onClose}>
      <Title colorConfig={colorConfig}>{title}</Title>
      {note && <Note className="mb-2">{note}</Note>}
      <input
        type="text"
        placeholder={searchPlaceholder}
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        className="w-full mb-3 p-2 border rounded"
      />
      {showOnlyOnSearch && search.trim() === "" ? (
        <Text className="text-gray-500">Unesite pojam za pretragu.</Text>
      ) : (
        <>
          {filteredItems?.map((item) => (
            <label
              key={item?.[keyField]}
              className="flex items-center gap-2 py-1"
            >
              <input
                type="checkbox"
                className="h-5 w-5"
                checked={selected.includes(item?.[keyField])}
                onChange={() => toggleItem(item?.[keyField])}
              />
              <Text className="text-xl">{item?.[labelField]}</Text>
            </label>
          ))}
          {(filteredItems?.length == 0 || !filteredItems) && (
            <Text>Nema rezultata.</Text>
          )}
        </>
      )}
      <div className="flex justify-end gap-2 mt-4">
        <Button
          onClick={onClose}
          color="ternary"
          icon={FaTimes}
          colorConfig={colorConfig}
        >
          Otkaži
        </Button>
        <Button
          disabled={selected?.length === 0}
          onClick={handleSave}
          color="secondary"
          icon={FaSave}
          colorConfig={colorConfig}
        >
          Sačuvaj
        </Button>
      </div>
    </Modal>
  );
}
