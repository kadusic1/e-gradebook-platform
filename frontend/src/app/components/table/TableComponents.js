import { forwardRef } from "react";

export const Table = forwardRef(
  ({ className, minWidth = "min-w-full", children, ...props }, ref) => {
    return (
      <table
        ref={ref}
        className={`${minWidth} bg-white mt-4 border border-gray-400 ${className}`}
        {...props}
      >
        {children}
      </table>
    );
  },
);
Table.displayName = "Table";

export const TableRow = ({ mode = "body", children, className, ...props }) => {
  return (
    <tr
      className={`${mode == "body" ? "hover:bg-gray-50" : "hover:bg-gray-300"} transition duration-150 ${className}`}
      {...props}
    >
      {children}
    </tr>
  );
};

export const TableHead = ({ children, className, ...props }) => {
  return (
    <thead className={`bg-gray-100 ${className}`} {...props}>
      {children}
    </thead>
  );
};

export const TableHeader = ({
  children,
  className,
  center = false,
  bordered = false,
  ...props
}) => {
  return (
    <th
      className={`px-4 py-3 ${center ? "text-center" : "text-left"} text-xs font-medium text-gray-500 uppercase tracking-wider border-b border-gray-200 ${bordered ? "border border-gray-400" : ""} ${className}`}
      {...props}
    >
      {children}
    </th>
  );
};

export const TableCell = ({
  children,
  className,
  bordered = false,
  yPadding = "py-4",
  ...props
}) => {
  return (
    <td
      className={`px-6 ${yPadding} whitespace-nowrap text-sm text-gray-700 ${bordered ? "border border-gray-400" : ""} ${className}`}
      {...props}
    >
      {children}
    </td>
  );
};

export const TableBody = ({ children, className, ...props }) => {
  return (
    <tbody className={`divide-y divide-gray-200 ${className}`} {...props}>
      {children}
    </tbody>
  );
};
