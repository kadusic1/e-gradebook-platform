import React from "react";
import { useFormContext, useController } from "react-hook-form";

export default function DisplayChoiceInput({
  name,
  options,
  icon: Icon,
  ...props
}) {
  const { control } = useFormContext();
  const {
    field: { value, onChange },
  } = useController({ name, control });

  return (
    <div className="flex gap-4">
      {options.map((opt) => (
        <label
          key={opt.value}
          className={`flex items-center gap-2 cursor-pointer border rounded px-3 py-2 transition
            ${value === opt.value ? "border-blue-500 bg-blue-100 shadow" : "border-gray-300 bg-white"}
          `}
        >
          <input
            type="radio"
            name={name}
            value={opt.value}
            checked={value === opt.value}
            onChange={() => onChange(opt.value)}
            {...props}
            className="accent-blue-500"
          />
          {opt.icon && (
            <opt.icon
              className={
                value === opt.value ? "text-blue-600" : "text-gray-400"
              }
            />
          )}
          <span
            className={value === opt.value ? "font-semibold text-blue-700" : ""}
          >
            {opt.label}
          </span>
        </label>
      ))}
    </div>
  );
}
