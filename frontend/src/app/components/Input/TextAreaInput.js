import React from "react";

export default function TextareaInput({
  name,
  placeholder,
  className,
  ...props
}) {
  return (
    <textarea
      id={name}
      name={name}
      placeholder={placeholder}
      className={`w-full px-3 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-200 bg-white text-gray-900 placeholder:text-gray-400 ${className}`}
      {...props}
    />
  );
}
