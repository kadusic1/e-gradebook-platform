import React from "react";
import Subtitle from "../common/Subtitle";
import { FaCopyright } from "react-icons/fa";

export default function Footer() {
  const year = new Date().getFullYear();
  return (
    <footer className="animate-fadeIn py-6 text-center bg-gradient-to-r from-gray-50 to-gray-100 border-t border-gray-200 shadow-lg relative overflow-hidden z-10">
      {/* Subtle background pattern */}
      <div className="absolute inset-0 opacity-5">
        <div className="absolute inset-0 bg-gradient-to-br from-blue-500 to-purple-600"></div>
      </div>

      {/* Main content */}
      <div className="relative">
        <Subtitle
          textSize="text-lg"
          showLine={false}
          icon={FaCopyright}
          textColor="text-gray-700 hover:text-gray-900 transition-colors duration-300"
        >
          Eacon d.o.o. {year}
        </Subtitle>
      </div>
    </footer>
  );
}
