import React from "react";

export default function NavbarTextItem({ children, icon, className = "" }) {
  return (
    <div
      className={`flex text-lg items-center gap-2 font-semibold text-gray-800 ${className}`}
    >
      {icon && <span className="text-xl">{icon}</span>}
      {children}
    </div>
  );
}
