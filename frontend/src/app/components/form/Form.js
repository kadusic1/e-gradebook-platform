import React from "react";
import { useForm } from "react-hook-form";
import FormField from "./FormField";
import Spacer from "../common/Spacer";
import Button from "../common/Button";
import ErrorMessage from "../custom_messages/ErrorMessage";
import { FaSave } from "react-icons/fa";
import { FaTimes } from "react-icons/fa";
import { FormProvider } from "react-hook-form";

export default function Form({
  fields,
  initialValues,
  onSubmit,
  children,
  onClose,
  showCancel = true,
  submitText = "Sačuvaj",
  error,
  fixedHeight = true,
  colorConfig,
  showSave = true,
  internalOverflow = true,
}) {
  const methods = useForm({ defaultValues: initialValues, mode: "onSubmit" });
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = methods;

  const onFormSubmit = (data) => {
    onSubmit(data);
  };

  return (
    <FormProvider {...methods}>
      <form
        onSubmit={handleSubmit(onFormSubmit)}
        className={`w-full max-w-lg ${fixedHeight ? "max-h-[60vh]" : ""} flex flex-col`}
      >
        <div className={`flex-1 ${internalOverflow ? "overflow-y-auto" : ""}`}>
          <Spacer>
            {fields.map((field, idx) => {
              if (React.isValidElement(field)) {
                return <React.Fragment key={idx}>{field}</React.Fragment>;
              }
              const isCheckbox = field.type === "checkbox";
              const isAttendance = field.type === "attendance";
              const validationRules = {};
              if (!isCheckbox && !isAttendance && (field?.required ?? true)) {
                validationRules.required = true;
              }
              if (field.type === "email") {
                validationRules.pattern = {
                  value: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
                  message: "Unesite ispravnu email adresu",
                };
              }
              // Integer validation
              if (field.type === "integer") {
                validationRules.pattern = {
                  value: /^-?\d+$/,
                  message: "Unesite cijeli broj",
                };
              }
              // Positive integer validation
              if (field.type === "positive-integer") {
                validationRules.pattern = {
                  value: /^\d+$/,
                  message: "Unesite pozitivan cijeli broj",
                };
                validationRules.validate = (v) =>
                  (v !== "" && /^\d+$/.test(v) && parseInt(v, 10) > 0) ||
                  "Unesite broj veći od 0";
              }
              // Number validation
              if (field.type === "number") {
                validationRules.pattern = {
                  value: /^-?\d*(\.\d+)?$/,
                  message: "Unesite broj",
                };
              }
              // Positive number validation
              if (field.type === "positive-number") {
                validationRules.pattern = {
                  value: /^\d*(\.\d+)?$/,
                  message: "Unesite pozitivan broj",
                };
                validationRules.validate = (v) =>
                  (v !== "" && /^\d*(\.\d+)?$/.test(v) && parseFloat(v) > 0) ||
                  "Unesite broj veći od 0";
              }
              // Phone validation
              if (field.type === "phone") {
                validationRules.pattern = {
                  value: /^\+?[0-9\s\-()]{7,20}$/,
                  message: "Unesite ispravan broj telefona",
                };
              }
              return (
                <FormField
                  key={field.name ?? idx}
                  {...field}
                  label={field.label}
                  name={field.name}
                  type={field.type === "email" ? "text" : field.type}
                  options={field.options}
                  placeholder={field.placeholder}
                  icon={field.icon}
                  // Pass ref and error to FormField for integration
                  {...register(field.name, validationRules)}
                  error={
                    errors[field.name]
                      ? errors[field.name].message || "Popunite ovo polje"
                      : undefined
                  }
                  items={field.items}
                  colorConfig={colorConfig}
                />
              );
            })}
          </Spacer>
          {children}
        </div>
        <Spacer className="mt-4 flex-shrink-0">
          {error && (
            <ErrorMessage className="mt-4 animate-fadeIn">{error}</ErrorMessage>
          )}
          {/* Show only the first field error */}
          {!error && Object.values(errors).length > 0 && (
            <div className="mt-2 space-y-1">
              <ErrorMessage className="animate-fadeIn">
                {Object.values(errors)[0].message || "Popunite sva polja"}
              </ErrorMessage>
            </div>
          )}
          {showSave && (
            <Button type="submit" icon={FaSave} colorConfig={colorConfig}>
              {submitText}
            </Button>
          )}
          {showCancel && (
            <Button
              color="secondary"
              type="button"
              onClick={onClose}
              icon={FaTimes}
              colorConfig={colorConfig}
            >
              Otkaži
            </Button>
          )}
        </Spacer>
      </form>
    </FormProvider>
  );
}
