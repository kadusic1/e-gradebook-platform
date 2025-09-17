import { getServerSession } from "next-auth";
import { authOptions } from "@/app/api/auth/[...nextauth]/route";
import SidebarHeader from "./SidebarHeader";
import SidebarNavItem from "./SidebarNavItem";
import SidebarSection from "./SidebarSection";
import { sidebarSections } from "./sidebarItems";

export default async function Sidebar({ className = "" }) {
  const session = await getServerSession(authOptions);
  const userAccountType = session?.user?.account_type;

  return (
    <aside
      className={`h-screen md:sticky top-0 w-64 max-w-full bg-white border-r border-gray-200 shadow-lg p-4 md:p-6 flex flex-col ${className}`}
    >
      <SidebarHeader />
      <nav className="flex-1 overflow-y-auto">
        {sidebarSections.map((section) => {
          const filteredItems = section.items.filter(
            (item) =>
              !item.account_types ||
              (userAccountType && item.account_types.includes(userAccountType)),
          );
          if (filteredItems.length === 0) return null;
          return (
            <SidebarSection key={section.title} title={section.title}>
              {filteredItems.map((item) => (
                <SidebarNavItem
                  key={item.href}
                  href={item.href}
                  icon={item.icon}
                >
                  {item.label}
                </SidebarNavItem>
              ))}
            </SidebarSection>
          );
        })}
      </nav>
    </aside>
  );
}
