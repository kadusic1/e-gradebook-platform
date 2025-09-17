"use client";

import Link from "next/link";

const CustomLink = ({ children, href, className = "" }) => {
  return (
    <Link href={href} className={`text-blue-600 underline ${className}`}>
      {children}
    </Link>
  );
};

export default CustomLink;
