import React from "react";
import Button from "../common/Button";
import { FaEdit, FaChevronRight } from "react-icons/fa";
import { FaTrash, FaPlus } from "react-icons/fa";
import { getColor } from "../colors/colors";
import { useState } from "react";
import Subtitle from "../common/Subtitle";

// Reusable card container
export function CardContainer({
  children,
  minWidth = "min-w-[200px]",
  className = "",
  ...props
}) {
  return (
    <div
      className={`animate-fadeIn bg-white rounded-xl shadow-lg p-4 md:p-6 flex flex-col gap-2 border border-gray-100 transition hover:shadow-xl hover:-translate-y-1 hover:bg-gray-50 duration-200 w-full ${minWidth} ${className}`}
      {...props}
    >
      {children}
    </div>
  );
}

// CardHeader with colored background
function CardHeader({ icon, title, colorConfig, bgColorOption = "primary" }) {
  const headerBgColor = getColor(bgColorOption, "bg", colorConfig);
  const headerTextColor = getColor("primaryComplement", "text", colorConfig);
  return (
    <div
      className={`flex items-center gap-3 mb-2 px-4 py-2 rounded-t-xl ${headerBgColor} ${headerTextColor} shadow-sm`}
    >
      {icon && (
        <span className="inline-flex items-center justify-center w-8 h-8 text-xl">
          {typeof icon === "function" ? React.createElement(icon) : icon}
        </span>
      )}
      <span className="font-bold text-xl leading-8">{title}</span>
    </div>
  );
}

export default function DynamicCard({
  data,
  icon,
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
  mapKeyToBosnian = (key) => key,
  mapValueToBosnian = (value) => value,
  textTitle,
  keysToExclude = [],
  tenantColorConfig,
  getValueColor,
  extraActions = [], // Array of { label, icon, onClick }
  className = "",
  twoColumnsLg = false,
  bgColorOption,
  textColorOption = "primary",
}) {
  const [showMore, setShowMore] = useState(false);
  const colorConfigID =
    typeof tenantColorConfig === "function"
      ? tenantColorConfig(data)
      : tenantColorConfig;

  const shouldShowEdit =
    typeof showEdit === "function" ? showEdit(data) : showEdit;
  const shouldShowDelete =
    typeof showDelete === "function" ? showDelete(data) : showDelete;

  const primaryTextColor = getColor(textColorOption, "text", colorConfigID);

  return (
    <CardContainer className={className}>
      <CardHeader
        icon={icon}
        title={textTitle || getTitle(data, titleField)}
        size="lg"
        colorConfig={colorConfigID}
        bgColorOption={bgColorOption}
      />
      <div
        className={`mt-4 px-4 pb-2 gap-3 ${
          twoColumnsLg ? "grid grid-cols-1 lg:grid-cols-2" : "grid grid-cols-1"
        }`}
      >
        {Object.entries(data).map(([key, value], idx) => {
          if (
            key === "id" ||
            (Array.isArray(titleField)
              ? titleField.includes(key)
              : key === titleField) ||
            keysToIgnore.includes(key)
          )
            return null;
          return (
            <div key={idx} className={`text-base font-semibold text-gray-700`}>
              {!keysToExclude.includes(key) && <>{mapKeyToBosnian(key)}: </>}
              <span
                className={`font-bold ${getValueColor ? getValueColor(value) || primaryTextColor : primaryTextColor}`}
              >
                {String(mapValueToBosnian(value))}
              </span>
            </div>
          );
        })}
      </div>
      <div className="flex flex-col xl:flex-row xl:justify-end gap-2 xl:gap-4 mt-4 items-center">
        {extraButton?.label && extraButton?.onClick && (
          <Button
            color="quaternary"
            onClick={() => extraButton.onClick(data)}
            icon={extraButton?.icon}
            className="w-5/6 xl:w-auto"
            colorConfig={colorConfigID}
          >
            {extraButton.label}
          </Button>
        )}
        {shouldShowEdit && (
          <Button
            color="secondary"
            className="w-5/6 xl:w-auto"
            onClick={() =>
              editButton?.onClick ? editButton.onClick(data) : onEdit(data)
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
            className="w-5/6 xl:w-auto"
            onClick={() =>
              deleteButton?.onClick
                ? deleteButton.onClick(data)
                : onConfirm(data)
            }
            icon={deleteButton?.icon || FaTrash}
            colorConfig={colorConfigID}
          >
            {deleteButton?.label || "Izbriši"}
          </Button>
        )}
        {(() => {
          const filteredActions = extraActions.filter(
            // If action has a show property, filter based on it
            // else show all actions
            (action) => action.show === undefined || action.show(data),
          );
          if (filteredActions.length === 0) return null;
          return (
            <div className="relative flex items-center">
              <span
                className="inline-flex items-center justify-center w-8 h-8 rounded-full hover:bg-gray-200 cursor-pointer transition"
                onClick={() => setShowMore((v) => !v)}
                title="Više opcija"
              >
                <FaChevronRight className="w-5 h-5 text-gray-600" />
              </span>
              {showMore && (
                <div className="p-2 absolute bottom-full left-1/2 -translate-x-1/2 -translate-y-1/3 mb-2 md:top-auto md:right-0 md:left-auto md:bottom-auto md:-translate-x-10 md:mb-0 z-50 min-w-[250px] flex flex-col gap-2 bg-white rounded shadow-md">
                  <Subtitle
                    colorConfig={colorConfigID}
                    showLine={false}
                    icon={FaPlus}
                    boxClassname="mb-2"
                  >
                    Dodatno
                  </Subtitle>
                  {filteredActions.map((action, idx) => (
                    <Button
                      key={idx}
                      color={
                        idx % 4 === 0
                          ? "primary"
                          : idx % 4 === 1
                            ? "secondary"
                            : idx % 4 === 2
                              ? "ternary"
                              : "quaternary"
                      }
                      icon={action.icon}
                      onClick={() => {
                        action.onClick(data);
                        setShowMore(false);
                      }}
                      className="w-full justify-start"
                      colorConfig={colorConfigID}
                    >
                      {action.label}
                    </Button>
                  ))}
                </div>
              )}
            </div>
          );
        })()}
      </div>
    </CardContainer>
  );
}
