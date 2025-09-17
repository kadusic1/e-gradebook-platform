// Helper for text color
const getTextClass = ({ color, shade }) => {
  if (color === "white" || color === "black") {
    return `text-${color}`;
  }
  return `text-${color}-${shade}`;
};

// Helper for background color including hover/focus
const getBgClassSet = ({ color, base, hover, focus }) => {
  if (color === "white" || color === "black") {
    return `bg-${color}`; // hover/focus variants don't apply to solid colors
  }
  return `bg-${color}-${base} hover:bg-${color}-${hover} focus:bg-${color}-${focus}`;
};

// Factory function to generate color maps from any config
export const createColorConfig = (config) => {
  const textColorMap = {
    primary: getTextClass({
      color: config.baseColor.color,
      shade: config.baseColor.base,
    }),
    primaryComplement: getTextClass({
      color: config.baseComplement.color,
      shade: config.baseComplement?.shade,
    }),
    secondary: getTextClass({
      color: config.buttonColor1.color,
      shade: config.buttonColor1.base,
    }),
    ternary: getTextClass({
      color: config.buttonColor2.color,
      shade: config.buttonColor2.base,
    }),
    quaternary: getTextClass({
      color: config.buttonColor3.color,
      shade: config.buttonColor3.base,
    }),
  };

  const bgColorMap = {
    primary: getBgClassSet(config.baseColor),
    secondary: getBgClassSet(config.buttonColor1),
    ternary: getBgClassSet(config.buttonColor2),
    quaternary: getBgClassSet(config.buttonColor3),
  };

  return {
    textColorMap,
    bgColorMap,
  };
};
