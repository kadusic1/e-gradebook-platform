import React from "react";

export default function Modal({ children, onClose, className = "" }) {
  return (
    <div
      className={`fixed inset-0 bg-white/70 backdrop-blur-sm flex items-center justify-center z-50 animate-fadeIn ${className}`}
    >
      <div className="bg-white rounded-xl shadow-2xl p-8 w-full max-w-lg md:max-w-xl relative animate-fadeIn max-h-[80vh] overflow-y-auto">
        <button
          onClick={onClose}
          className="hover:cursor-pointer absolute top-3 right-3 text-gray-400 hover:text-gray-600 text-2xl font-bold focus:outline-none"
          aria-label="Zatvori"
        >
          ×
        </button>
        {children}
      </div>
    </div>
  );
}
