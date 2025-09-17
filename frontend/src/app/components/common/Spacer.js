import React from "react";

export default function Spacer({ children, className = "" }) {
  return (
    <div className={`flex flex-col gap-3 md:gap-4 w-full ${className}`}>
      {children}
    </div>
  );
}
