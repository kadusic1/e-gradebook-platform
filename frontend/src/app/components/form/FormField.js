import React from "react";
import Label from "../Input/Label";
import TextInput from "../Input/TextInput";
import DateInput from "../Input/DateInput";
import CheckboxInput from "../Input/CheckboxInput";
import ReadonlyInput from "../Input/ReadonlyInput";
import SelectInput from "../Input/SelectInput";
import PasswordInput from "../Input/PasswordInput";
import Subtitle from "../common/Subtitle";
import DisplayChoiceInput from "../Input/DisplayChoiceInput";
import StudentAttendanceField from "../Input/StudentAttendanceInput";
import { getColor } from "../colors/colors";
import TextareaInput from "../Input/TextAreaInput";

export default function FormField({
  label,
  name,
  type = "text",
  options,
  placeholder,
  icon: Icon,
  items, // If form field renders multiple items
  colorConfig,
  ...props
}) {
  const primaryTextColor = getColor("primary", "text", colorConfig);

  return (
    <div>
      {type != "subtitle" && (
        <Label name={name}>
          {Icon && (
            <Icon
              className={`inline-block mr-2 ${primaryTextColor} align-middle`}
            />
          )}{" "}
          {label}
        </Label>
      )}
      {type === "checkbox" ? (
        <CheckboxInput name={name} {...props} />
      ) : type === "select" ? (
        <SelectInput
          name={name}
          options={options}
          placeholder={placeholder}
          {...props}
        />
      ) : type === "readonly" ? (
        <ReadonlyInput name={name} options={options} {...props} />
      ) : type === "date" ? (
        <DateInput name={name} placeholder={placeholder} {...props} />
      ) : type === "password" ? (
        <PasswordInput name={name} placeholder={placeholder} {...props} />
      ) : type === "subtitle" ? (
        <Subtitle icon={Icon} showLine={false} colorConfig={colorConfig}>
          {label}
        </Subtitle>
      ) : type === "display-choice" ? (
        <DisplayChoiceInput name={name} options={options} {...props} />
      ) : type === "attendance" ? (
        <StudentAttendanceField
          name={name}
          label={label}
          students={items || []}
          {...props}
        />
      ) : type === "text-area" ? (
        <TextareaInput name={name} placeholder={placeholder} {...props} />
      ) : (
        <TextInput name={name} placeholder={placeholder} {...props} />
      )}
    </div>
  );
}
