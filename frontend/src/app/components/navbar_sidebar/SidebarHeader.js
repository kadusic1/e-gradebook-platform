import React from "react";

export default function SidebarHeader() {
  return (
    <div className="bg-gray-100 rounded-lg px-3 py-2 mb-8 flex items-center gap-2">
      <h2 className="text-xl font-bold text-gray-800 tracking-tight flex items-center gap-2 mb-0">
        <span className="inline-block w-2 h-2 bg-indigo-500 rounded-full"></span>
        eDnevnik
      </h2>
    </div>
  );
}
