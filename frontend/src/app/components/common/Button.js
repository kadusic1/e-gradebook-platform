import React from "react";
import { getColor } from "../colors/colors";

export default function Button({
  children,
  className = "",
  color = "primary",
  onClick,
  icon: Icon,
  colorConfig,
  type = "button",
  iconClassName = "",
  ...props
}) {
  let bgColor = getColor(color, "bg", colorConfig);
  let textColor = getColor("primaryComplement", "text", colorConfig);

  return (
    <button
      type={type}
      className={`animate-fadeIn hover:cursor-pointer font-semibold py-2 px-6 rounded-lg shadow-sm transition duration-200 focus:outline-none focus:ring-2 focus:ring-indigo-200 ${bgColor} ${textColor} ${className}`}
      onClick={onClick}
      {...props}
    >
      <div className="flex items-center gap-2 justify-center">
        {Icon && <Icon className={`${textColor} ${iconClassName}`} />}
        {children}
      </div>
    </button>
  );
}
