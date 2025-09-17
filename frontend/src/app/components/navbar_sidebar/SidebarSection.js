import React from "react";

export default function SidebarSection({ title, children }) {
  return (
    <div className="mb-6">
      <div className="text-xs font-semibold text-gray-400 uppercase mb-2 px-2 select-none">
        {title}
      </div>
      <div className="flex flex-col gap-2">{children}</div>
    </div>
  );
}
