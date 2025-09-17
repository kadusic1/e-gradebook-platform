import { createColorConfig } from "./colorConfig";
import { tenantColorConfigs } from "./defaultTenantConfig";

export const getColor = (variant, type, configID = 0) => {
  if (
    configID < 0 ||
    configID >= tenantColorConfigs.length ||
    !tenantColorConfigs[configID]
  ) {
    configID = 0;
  }

  let config = tenantColorConfigs[configID];

  let processedConfig = createColorConfig(config);

  if (type === "text") {
    return processedConfig.textColorMap[variant];
  }

  if (type === "bg") {
    return processedConfig.bgColorMap[variant];
  }
};
