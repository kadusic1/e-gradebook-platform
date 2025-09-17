import React from "react";

const Text = ({
  children,
  textSize = "text-lg",
  className = "",
  textColor = "text-gray-800",
}) => {
  return (
    <div className={`${textSize} ${textColor} ${className} drop-shadow-sm`}>
      {children}
    </div>
  );
};

export default Text;
