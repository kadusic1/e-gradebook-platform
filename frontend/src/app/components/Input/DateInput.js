export default function DateInput({
  name,
  placeholder,
  className,
  min,
  max,
  ...props
}) {
  return (
    <input
      id={name}
      name={name}
      type="date"
      placeholder={placeholder}
      min={min}
      max={max}
      className={`w-full px-3 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-200 bg-white text-gray-900 placeholder:text-gray-400 ${className}`}
      {...props}
    />
  );
}
