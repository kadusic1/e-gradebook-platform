import React from "react";
import { getColor } from "../colors/colors";

export default function SidebarNavItem({
  href,
  icon: Icon,
  children,
  onClick,
  className = "",
  colorConfig,
}) {
  const primaryTextColor = getColor("primary", "text", colorConfig);

  return (
    <a
      href={href}
      className={`flex items-center gap-3 px-3 py-2 rounded-lg text-gray-900 font-medium ${className}`}
      onClick={onClick}
    >
      {Icon && <Icon className={primaryTextColor} />}
      {children}
    </a>
  );
}
