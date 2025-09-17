import React from "react";

export default function ReadonlyInput({ name, options = {}, ...props }) {
  if (options?.value && options?.label) {
    return (
      <>
        <input
          id={name}
          name={name}
          type="text"
          readOnly
          value={options.label}
          className="w-full px-3 py-2 border border-gray-200 rounded-lg bg-gray-100 text-gray-600"
        />
        <input type="hidden" name={name} value={options.value} {...props} />
      </>
    );
  }
  return (
    <input
      id={name}
      name={name}
      type="text"
      readOnly
      value="Nema podataka"
      className="w-full px-3 py-2 border border-gray-200 rounded-lg bg-gray-100 text-gray-600"
    />
  );
}
