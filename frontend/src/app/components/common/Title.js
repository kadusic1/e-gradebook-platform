import React from "react";
import { getColor } from "../colors/colors";

export default function Title({
  children,
  icon: Icon,
  className = "",
  colorConfig,
  outerClassName = "",
  showLine = true,
  textSize = "text-4xl",
}) {
  const primaryTextColor = getColor("primary", "text", colorConfig);
  const primaryBgColor = getColor("primary", "bg", colorConfig);

  return (
    <div className={`animate-fadeIn ${outerClassName}`}>
      <h1
        className={`${textSize} font-extrabold text-gray-800 drop-shadow-md ${className} ${showLine ? "mb-2" : "mb-5"}`}
      >
        {Icon && (
          <Icon
            className={`inline-block mr-2 ${primaryTextColor} align-middle`}
          />
        )}{" "}
        {children}
      </h1>
      {showLine && (
        <div className={`h-1 w-16 ${primaryBgColor} rounded mb-6`}></div>
      )}
    </div>
  );
}
