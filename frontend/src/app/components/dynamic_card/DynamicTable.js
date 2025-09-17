import React from "react";
import Button from "../common/Button";
import { FaEdit, FaTrash, FaChevronRight, FaPlus } from "react-icons/fa";
import Subtitle from "../common/Subtitle";
import { getColor } from "../colors/colors";
import { useState } from "react";
import { createPortal } from "react-dom";

export default function DynamicTable({
  data = [],
  titleField = "name",
  onConfirm,
  onEdit,
  showEdit,
  showDelete,
  extraButton,
  keysToIgnore = [],
  getTitle,
  deleteButton,
  editButton,
  keyField = "id",
  mapKeyToBosnian = (key) => key,
  mapValueToBosnian = (value) => value,
  tenantColorConfig,
  getValueColor,
  extraActions = [], // Array of { label, icon, onClick }
}) {
  const [openExtraIdx, setOpenExtraIdx] = useState(null);
  const [dropdownPosition, setDropdownPosition] = useState({ top: 0, left: 0 });

  const colorConfigID =
    typeof tenantColorConfig === "function"
      ? tenantColorConfig(data)
      : tenantColorConfig;
  if (!data || data.length === 0) return null;

  // Collect all keys to show (excluding ignored, id, and titleField)
  const allKeys = Object.keys(data[0] || {}).filter(
    (key) =>
      key !== keyField &&
      key !== "id" &&
      (Array.isArray(titleField)
        ? !titleField.includes(key)
        : key !== titleField) &&
      !keysToIgnore.includes(key),
  );

  const hasActions =
    (extraButton?.label && extraButton?.onClick) || showEdit || showDelete;

  const primaryComplementTextColor = getColor(
    "primaryComplement",
    "text",
    colorConfigID,
  );
  const primaryTextColor = getColor("primary", "text", colorConfigID);
  const headerBgColor = getColor("primary", "bg", colorConfigID);

  const handleExtraActionsClick = (event, idx) => {
    if (openExtraIdx === idx) {
      setOpenExtraIdx(null);
      return;
    }

    const rect = event.currentTarget.getBoundingClientRect();
    setDropdownPosition({
      top: rect.bottom + window.scrollY + 10,
      left: rect.right - 200 + window.scrollX, // 200px is dropdown width
    });
    setOpenExtraIdx(idx);
  };

  return (
    <>
      <div className="overflow-x-auto rounded-xl shadow-lg bg-white border border-gray-100 min-w-[700px]">
        <table className="min-w-full divide-y divide-gray-200 table-fixed">
          <thead className={`${headerBgColor}`}>
            <tr>
              <th
                className={`px-4 py-3 text-left text-xs font-bold ${primaryComplementTextColor} uppercase tracking-wider`}
              >
                {Array.isArray(titleField)
                  ? titleField.map(mapKeyToBosnian).join(" ")
                  : mapKeyToBosnian(titleField)}
              </th>
              {allKeys.map((key) => (
                <th
                  key={key}
                  className={`px-4 py-3 text-left text-xs font-bold ${primaryComplementTextColor} uppercase tracking-wide`}
                >
                  {mapKeyToBosnian(key)}
                </th>
              ))}
              {hasActions && (
                <th
                  className={`px-4 py-3 text-center text-xs font-bold ${primaryComplementTextColor} uppercase tracking-wider`}
                ></th>
              )}
              {extraActions.length > 0 && (
                <th
                  className={`px-4 py-3 text-center text-xs font-bold ${primaryComplementTextColor} uppercase tracking-wider`}
                ></th>
              )}
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-100">
            {data.map((item, idx) => {
              const shouldShowEdit =
                typeof showEdit === "function" ? showEdit(item) : showEdit;
              const shouldShowDelete =
                typeof showDelete === "function"
                  ? showDelete(item)
                  : showDelete;

              return (
                <tr key={idx} className="hover:bg-gray-50">
                  <td className={`px-4 py-3 font-semibold ${primaryTextColor}`}>
                    <div
                      className={`${hasActions ? "min-h-16" : "min-h-12"} flex items-center`}
                    >
                      {getTitle(item, titleField)}
                    </div>
                  </td>
                  {allKeys.map((key) => (
                    <td key={key} className="px-4 py-3 text-gray-700">
                      <div
                        className={`${hasActions ? "min-h-16" : "min-h-12"} flex items-center`}
                      >
                        <span
                          className={
                            getValueColor
                              ? getValueColor(item[key]) || "text-gray-700"
                              : "text-gray-700"
                          }
                        >
                          {String(mapValueToBosnian(item[key]))}
                        </span>
                      </div>
                    </td>
                  ))}
                  {hasActions && (
                    <td className="px-4 py-3">
                      <div
                        className={`${hasActions ? "min-h-16" : "min-h-12"} flex flex-row gap-2 justify-center items-center`}
                      >
                        {extraButton?.label && extraButton?.onClick && (
                          <Button
                            color="quaternary"
                            onClick={() => extraButton.onClick(item)}
                            icon={extraButton?.icon}
                            colorConfig={colorConfigID}
                          >
                            {extraButton.label}
                          </Button>
                        )}
                        {shouldShowEdit && (
                          <Button
                            color="secondary"
                            onClick={() =>
                              editButton?.onClick
                                ? editButton.onClick(item)
                                : onEdit(item)
                            }
                            icon={editButton?.icon || FaEdit}
                            colorConfig={colorConfigID}
                          >
                            {editButton?.label || "Uredi"}
                          </Button>
                        )}
                        {shouldShowDelete && (
                          <Button
                            color="ternary"
                            onClick={() =>
                              deleteButton?.onClick
                                ? deleteButton.onClick(item)
                                : onConfirm(item)
                            }
                            icon={deleteButton?.icon || FaTrash}
                            colorConfig={colorConfigID}
                          >
                            {deleteButton?.label || "Izbriši"}
                          </Button>
                        )}
                      </div>
                    </td>
                  )}
                  {extraActions.length > 0 && ( // ADD
                    <td className="px-4 py-3 relative">
                      <span
                        className="inline-flex items-center justify-center w-8 h-8 rounded-full hover:bg-gray-200 cursor-pointer transition"
                        onClick={(e) => handleExtraActionsClick(e, idx)}
                        title="Više opcija"
                      >
                        <FaChevronRight className="w-5 h-5 text-gray-600" />
                      </span>
                    </td>
                  )}
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
      {/* Portal for dropdown outside table */}
      {openExtraIdx !== null &&
        createPortal(
          <>
            {/* Backdrop to close dropdown */}
            <div
              className="fixed inset-0 z-40"
              onClick={() => setOpenExtraIdx(null)}
            />
            <div
              className="p-2 absolute z-50 min-w-[200px] flex flex-col gap-2 bg-white rounded shadow-lg border border-gray-200"
              style={{
                top: dropdownPosition.top,
                left: dropdownPosition.left,
              }}
            >
              <Subtitle colorConfig={colorConfigID} showLine={false} icon={FaPlus} boxClassname="mb-2">
                Dodatno
              </Subtitle>
              {extraActions.map((action, aIdx) => (
                <Button
                  key={aIdx}
                  color={
                    aIdx % 4 === 0
                      ? "primary"
                      : aIdx % 4 === 1
                        ? "secondary"
                        : aIdx % 4 === 2
                          ? "ternary"
                          : "quaternary"
                  }
                  icon={action.icon}
                  onClick={() => {
                    action.onClick(data[openExtraIdx]);
                    setOpenExtraIdx(null);
                  }}
                  className="w-full justify-start"
                  colorConfig={colorConfigID}
                >
                  {action.label}
                </Button>
              ))}
            </div>
          </>,
          document.body,
        )}
    </>
  );
}
