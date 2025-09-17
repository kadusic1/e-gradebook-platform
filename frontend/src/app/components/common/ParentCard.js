import React from "react";

export default function ParentCard({ children, className }) {
  return (
    <div
      className={`main-content-bg rounded-xl p-6 md:p-10 w-full max-w-xs md:max-w-md lg:max-w-xl xl:max-w-2xl mx-auto mt-24 shadow-lg flex flex-col items-center animate-fadeIn ${className}`}
    >
      {children}
    </div>
  );
}
