import { dirname } from "path";
import { fileURLToPath } from "url";
import { FlatCompat } from "@eslint/eslintrc";

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const compat = new FlatCompat({
  baseDirectory: __dirname,
});

const eslintConfig = [
  ...compat.extends("next/core-web-vitals"),
  {
    rules: {
      // Turn off warning rules
      "react-hooks/exhaustive-deps": "off",

      // Keep error rules as errors
      "react/display-name": "error",
      "react/no-children-prop": "error",
    }
  }
];

export default eslintConfig;
