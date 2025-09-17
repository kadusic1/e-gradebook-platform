export default function Note({ children, className = "" }) {
  return (
    <div className={`text-gray-500 text-sm font-semibold ${className}`}>
      {children}
    </div>
  );
}
