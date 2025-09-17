import React from "react";

export default function Paragraph({ children, className = "" }) {
  return (
    <p className={`text-lg text-gray-700 mb-6 ${className}`}>{children}</p>
  );
}
