export default function Label({ name, children, className }) {
  return (
    <label
      htmlFor={name}
      className={`block text-sm font-medium text-gray-700 mb-1 ${className}`}
    >
      {children}
    </label>
  );
}
