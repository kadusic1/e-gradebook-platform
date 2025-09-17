import React from "react";
import { FaBookOpen, FaUserCircle } from "react-icons/fa";
import { signOut } from "next-auth/react";
import Button from "../common/Button";
import HorizontalSpacer from "../common/HorizontalSpacer";
import NavbarTextItem from "./NavbarTextItem";

export default function Navbar({ session, onMenuClick }) {
  const accountTypeMapping = {
    root: "Superadmin",
    tenant_admin: "Administrator",
    teacher: "Nastavnik",
    pupil: "Učenik",
    parent: "Roditelj",
  };

  return (
    <nav className="animate-fadeIn w-full flex items-center justify-between px-4 md:px-8 py-4 bg-gray-100 border-b border-gray-200 shadow-sm z-10">
      <div className="flex items-center gap-3">
        {/* Mobile menu button */}
        <button
          className="lg:hidden mr-2 text-2xl text-gray-700"
          onClick={onMenuClick}
          aria-label="Otvori meni"
        >
          <svg width="28" height="28" fill="none" viewBox="0 0 24 24">
            <path
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              d="M4 6h16M4 12h16M4 18h16"
            />
          </svg>
        </button>
        <FaBookOpen className="text-gray-700 text-2xl" />
      </div>
      <HorizontalSpacer>
        <NavbarTextItem icon={<FaUserCircle />} className="text-gray-700">
          {session?.user?.name}
          <span className="hidden sm:inline">
            ({accountTypeMapping[session?.user?.account_type]})
          </span>
        </NavbarTextItem>
        <Button onClick={() => signOut()} className="ml-2">
          Odjava
        </Button>
      </HorizontalSpacer>
    </nav>
  );
}
