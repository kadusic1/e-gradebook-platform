import React from "react";

export default function HorizontalSpacer({ children, className = "" }) {
  return (
    <div
      className={`animate-fadeIn flex flex-row gap-2 items-center ${className}`}
    >
      {children}
    </div>
  );
}
