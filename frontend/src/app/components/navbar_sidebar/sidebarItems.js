import {
  FaChalkboardTeacher,
  FaBuilding,
  FaUserShield,
  FaRegEnvelope,
  FaUsers,
  FaChalkboard,
  FaArchive,
  FaUser,
} from "react-icons/fa";

export const sidebarSections = [
  {
    title: "Opcije",
    items: [
      {
        href: "/tenants",
        icon: FaBuilding,
        label: "Institucije",
        account_types: ["root"],
      },
      {
        href: "/tenant_admin_administration",
        icon: FaBuilding,
        label: "Administracija",
        account_types: ["tenant_admin"],
      },
      {
        href: "/accounts",
        icon: FaChalkboardTeacher,
        label: "Korisnički nalozi",
        account_types: ["root", "tenant_admin"],
      },
      {
        href: "/global_admin",
        icon: FaUserShield,
        label: "Administracija",
        account_types: ["root"],
      },
      {
        href: "/pupil_home",
        icon: FaUsers,
        label: "Odjeljenja",
        account_types: ["pupil"],
      },
      {
        href: "/pupil_archived_sections",
        icon: FaArchive,
        label: "Arhivirana odjeljenja",
        account_types: ["pupil"],
      },
      {
        href: "/pupil_invites",
        icon: FaRegEnvelope,
        label: "Pozivi",
        account_types: ["pupil"],
      },
      {
        href: "/teacher_home",
        icon: FaUsers,
        label: "Odjeljenja",
        account_types: ["teacher"],
      },
      {
        href: "/teacher_archived_sections",
        icon: FaArchive,
        label: "Arhivirana odjeljenja",
        account_types: ["teacher"],
      },
      {
        href: "/teacher_invites",
        icon: FaRegEnvelope,
        label: "Pozivi",
        account_types: ["teacher"],
      },
      {
        href: "/teacher_schedule",
        icon: FaChalkboard,
        label: "Moj raspored",
        account_types: ["teacher"],
      },
      {
        href: "/pupil_profile",
        icon: FaUser,
        label: "Moj profil",
        account_types: ["pupil"],
      },
      {
        href: "/teacher_profile",
        icon: FaUser,
        label: "Moj profil",
        account_types: ["root", "tenant_admin", "teacher"],
      },
    ],
  },
];

export const sidebarPermissions = sidebarSections
  .flatMap((section) => section.items)
  .map((item) => ({
    href: item.href,
    account_types: item.account_types,
  }));
