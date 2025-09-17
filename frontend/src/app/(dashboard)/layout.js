import "../globals.css";
import "../main-content.css";
import Sidebar from "../components/navbar_sidebar/Sidebar";
import AuthenticatedLayout from "../components/layout/AuthenticatedLayout";
import { getServerSession } from "next-auth";
import { authOptions } from "@/app/api/auth/[...nextauth]/route";

export default async function DashboardLayout({ children }) {
  const session = await getServerSession(authOptions);
  return (
    <html lang="bs">
      <body>
        <div className="flex min-h-screen">
          {/* Desktop sidebar (server component) */}
          <div className="hidden lg:block">
            <Sidebar />
          </div>
          <AuthenticatedLayout session={session}>
            {children}
          </AuthenticatedLayout>
        </div>
      </body>
    </html>
  );
}
