import React from "react";

export default function CheckboxInput({ name, ...props }) {
  return (
    <input
      id={name}
      name={name}
      type="checkbox"
      className="mr-2 accent-indigo-500"
      {...props}
    />
  );
}
