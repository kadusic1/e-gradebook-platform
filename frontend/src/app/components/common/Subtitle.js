import React from "react";
import { getColor } from "../colors/colors";

export default function Subtitle({
  children,
  icon: Icon,
  className = "",
  boxClassname = "",
  showLine = true,
  colorConfig,
  textSize = "text-2xl",
  bgColor,
  textColor,
  textColorOption = "primary",
  bgColorOption = "primary",
} = {}) {
  let primaryTextColor;
  if (textColor) {
    primaryTextColor = textColor;
  } else {
    primaryTextColor = getColor(textColorOption, "text", colorConfig);
  }

  let primaryBgColor;
  if (bgColor) {
    primaryBgColor = bgColor;
  } else {
    primaryBgColor = getColor(bgColorOption, "bg", colorConfig);
  }

  return (
    <div className={`animate-fadeIn ${boxClassname}`}>
      <h2
        className={`${textSize} font-bold text-gray-800 drop-shadow-sm ${className}`}
      >
        {Icon && (
          <Icon
            className={`inline-block mr-2 ${primaryTextColor} align-middle`}
          />
        )}{" "}
        {children}
      </h2>
      {showLine && (
        <div className={`h-0.5 w-12 ${primaryBgColor} rounded`}></div>
      )}
    </div>
  );
}
