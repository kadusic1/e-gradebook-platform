"use client";
import Navbar from "../navbar_sidebar/Navbar";
import MobileSidebar from "../navbar_sidebar/MobileSidebar";
import { useState } from "react";
import { usePathname } from "next/navigation";
import { SessionProvider } from "next-auth/react";
import Footer from "../navbar_sidebar/Footer";
import AIChatEntry from "../modal/AIChatEntry";

export default function AuthenticatedLayout({ session, children }) {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [chatOpen, setChatOpen] = useState(false);
  const pathname = usePathname();

  if (pathname === "/login") {
    return (
      <SessionProvider>
        <div className="flex min-h-screen">
          <main className="main-content-bg flex-1 p-4 md:p-8">{children}</main>
        </div>
      </SessionProvider>
    );
  }

  return (
    <SessionProvider>
      <div className="flex-1 flex flex-col">
        {/* Mobile sidebar overlay */}
        {sidebarOpen && (
          <div
            className="fixed inset-0 z-50 bg-black/40"
            onClick={() => setSidebarOpen(false)}
          >
            <div className="absolute left-0 top-0 h-full w-64 bg-white shadow-lg">
              <MobileSidebar session={session} />
            </div>
          </div>
        )}
        <Navbar session={session} onMenuClick={() => setSidebarOpen(true)} />
        <main className="main-content-bg flex-1 p-4 md:p-8">{children}</main>
        <AIChatEntry
          accessToken={session?.accessToken}
          open={chatOpen}
          setOpen={setChatOpen}
        />
        <Footer />
      </div>
    </SessionProvider>
  );
}
