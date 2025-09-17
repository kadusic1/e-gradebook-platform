import React, { useState } from "react";
import Button from "../common/Button";
import HorizontalSpacer from "../common/HorizontalSpacer";
import Title from "../common/Title";
import Spacer from "../common/Spacer";
import { FaBars } from "react-icons/fa";
import SidebarNavItem from "../navbar_sidebar/SidebarNavItem";
import { getColor } from "../colors/colors";

// Custom Tab Button Component
const TabButton = ({
  children,
  isActive,
  onClick,
  icon: Icon,
  className = "",
  tenantColorConfig,
}) => {
  const activeButtonBgColor = getColor("primary", "bg", tenantColorConfig);
  const activeButtonTextColor = getColor(
    "primaryComplement",
    "text",
    tenantColorConfig,
  );

  const baseClasses =
    "font-semibold py-3 px-5 rounded-lg transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-indigo-200 hover:cursor-pointer";
  const activeClasses = `${activeButtonBgColor} ${activeButtonTextColor} shadow-lg transform scale-105 ring-2 ring-indigo-300 ring-opacity-50`;
  const inactiveClasses =
    "bg-white hover:bg-gray-50 text-gray-700 border border-gray-200 hover:border-gray-300 hover:shadow-sm";

  return (
    <button
      className={`${baseClasses} ${isActive ? activeClasses : inactiveClasses} ${className}`}
      onClick={onClick}
    >
      <div className="flex items-center gap-2 justify-center">
        {Icon && <Icon className={isActive ? "text-white" : "text-gray-600"} />}
        {children}
      </div>
    </button>
  );
};

export default function DynamicTab({
  childrenTabs,
  title,
  titleIcon,
  className = "",
  activeTab,
  setActiveTab,
  colorConfig,
  aboveContent,
}) {
  const [menuOpen, setMenuOpen] = useState(false);

  // Handle menu open/close for mobile
  const handleMenuToggle = () => setMenuOpen((open) => !open);
  const handleTabSelect = (idx) => {
    setActiveTab(idx);
    setMenuOpen(false);
  };

  return (
    <>
      <Spacer className={`flex flex-col ${className}`}>
        {/* Header div */}
        <div className="bg-gradient-to-r from-gray-100 via-gray-200 to-gray-300 rounded-xl shadow-lg p-4 mb-3">
          <div className="mb-3">
            <Title icon={titleIcon} colorConfig={colorConfig} showLine={false}>
              {title}
            </Title>
          </div>

          {/* Mobile: Dropdown menu */}
          <div className="block md:hidden relative">
            <Button
              className="flex justify-center items-center gap-2"
              onClick={handleMenuToggle}
              color="secondary"
              aria-label="Open tab menu"
              colorConfig={colorConfig}
            >
              Opcije
              <FaBars className="w-5 h-5" />
            </Button>
            {menuOpen && (
              <div className="absolute z-50 w-48 bg-white border rounded shadow-lg mt-1">
                {childrenTabs.map((tab, idx) => (
                  <SidebarNavItem
                    key={tab.label}
                    icon={tab?.icon}
                    className={`${activeTab === idx ? "bg-indigo-50 text-indigo-900" : ""}`}
                    onClick={() => handleTabSelect(idx)}
                    colorConfig={colorConfig}
                  >
                    {tab.label}
                  </SidebarNavItem>
                ))}
              </div>
            )}
          </div>

          {/* Desktop: Tab buttons */}
          <HorizontalSpacer className="hidden md:flex gap-2 flex-wrap">
            {childrenTabs.map((tab, idx) => (
              <TabButton
                key={idx}
                isActive={activeTab === idx}
                onClick={() => setActiveTab(idx)}
                icon={tab?.icon}
                tenantColorConfig={colorConfig}
              >
                {tab.label}
              </TabButton>
            ))}
          </HorizontalSpacer>
        </div>

        <div>{aboveContent}</div>

        {/* Content */}
        <div>{childrenTabs[activeTab]?.content}</div>
      </Spacer>
    </>
  );
}
