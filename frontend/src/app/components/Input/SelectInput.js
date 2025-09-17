import React from "react";

export default function SelectInput({
  name,
  options,
  placeholder,
  className,
  ...props
}) {
  return (
    <select
      id={name}
      name={name}
      className={`w-full px-3 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-200 bg-white text-gray-900 ${className}`}
      {...props}
    >
      {!options || options?.length === 0 ? (
        <option key="-1" value="">
          {"Nema dostupnih opcija"}
        </option>
      ) : (
        <>
          <option key="-1" value="">
            {placeholder || "Odaberite opciju"}
          </option>
          {options &&
            options.map((option, idx) => (
              <option
                key={idx}
                value={option.value !== undefined ? option.value : option}
              >
                {option.label !== undefined ? option.label : option}
              </option>
            ))}
        </>
      )}
    </select>
  );
}
