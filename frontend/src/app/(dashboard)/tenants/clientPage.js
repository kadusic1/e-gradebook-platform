"use client";
import { useState } from "react";
import AddButton from "@/app/components/common/AddButton";
import CreateUpdateModal from "@/app/components/modal/CreateUpdateModal";
import ErrorModal from "../../components/modal/ErrorModal";
import DynamicCardParent from "../../components/dynamic_card/DynamicCardParent";
import { FaBuilding, FaChalkboardTeacher } from "react-icons/fa";
import TeachersPerTenantPageClient from "./teachersPerTenantClient";
import { useEffect } from "react";
import Title from "@/app/components/common/Title";
import DynamicTab from "@/app/components/dynamic_card/DynamicTab";
import { FaBookOpen } from "react-icons/fa";
import { FaUserCog } from "react-icons/fa";
import CurriculumAssignmentPageClient from "./curriculumAssignmentPageClient";
import SectionsPageClient from "./sectionsClientPage";
import { FaLayerGroup } from "react-icons/fa";
import TenantSemestersPageClient from "./semestersPageClient";
import { FaRegCalendarAlt } from "react-icons/fa";
import ColorPreviewCard from "@/app/components/previews/color_preview";
import {
  FaRegIdCard,
  FaTable,
  FaPalette,
  FaUserPlus,
  FaBook,
  FaUsers,
  FaUserGraduate,
  FaUserCheck,
  FaCalendarAlt,
  FaChalkboard,
  FaUserTimes,
  FaCalendarDay,
  FaArchive,
} from "react-icons/fa";
import Loading from "@/app/components/common/Loading";
import ClassroomPageClient from "./classroomPageClient";

export default function TenantsPageClient({
  cantons,
  initialTenants = [],
  accessToken,
  tenantAdminMode = false,
}) {
  const [showModal, setShowModal] = useState(false);
  const [errorMessage, setErrorMessage] = useState(null);
  const [tenants, setTenants] = useState(initialTenants);
  const [expandedTenant, setExpandedTenant] = useState(null);
  const [activeTab, setActiveTab] = useState(0);

  // Initialize expandedTenant for tenant admin mode
  useEffect(() => {
    if (tenantAdminMode && initialTenants.length > 0 && !expandedTenant) {
      setExpandedTenant(initialTenants[0]);
    }
  }, [tenantAdminMode, initialTenants, expandedTenant]);

  const tenantFields = [
    {
      label: "Generalije",
      type: "subtitle",
      icon: FaBuilding,
    },
    {
      label: "Naziv institucije",
      name: "tenant_name",
      placeholder: "Unesite naziv institucije",
    },
    {
      label: "Grad institucije",
      name: "tenant_city",
      placeholder: "Unesite grad institucije",
    },
    {
      label: "Kanton",
      name: "canton_code",
      placeholder: "Odaberite kanton",
      type: "select",
      options: cantons?.map((canton) => ({
        value: canton.canton_code,
        label: `${canton.canton_name}`,
      })),
    },
    { label: "Adresa", name: "address", placeholder: "Unesite adresu" },
    {
      label: "Telefon",
      name: "phone",
      placeholder: "Unesite telefon +387",
      type: "phone",
    },
    {
      label: "Email",
      name: "email",
      type: "email",
      placeholder: "Unesite email",
    },
    {
      label: "Domena",
      name: "domain",
      type: "text",
      placeholder: "Unesite domenu institucije (npr. skola.ba)",
      required: false,
    },
    {
      label: "Ime direktora",
      name: "director_name",
      placeholder: "Unesite ime direktora",
    },
    {
      label: "Tip institucije",
      name: "tenant_type",
      placeholder: "Odaberite tip institucije",
      type: "select",
      options: [
        { value: "primary", label: "Osnovna škola" },
        { value: "secondary", label: "Srednja škola" },
      ],
    },
    {
      label: "Usmjerenje institucije",
      name: "specialization",
      placeholder: "Odaberite usmjerenje institucije",
      type: "select",
      options: [
        { value: "regular", label: "Obično" },
        { value: "religion", label: "Vjersko" },
        { value: "musical", label: "Muzičko" },
      ],
    },
    {
      label: "Administrator institucije",
      type: "subtitle",
      icon: FaUserCog,
    },
    {
      label: "Ime",
      name: "teacher_name",
      placeholder: "Unesite ime",
    },
    {
      label: "Prezime",
      name: "teacher_last_name",
      placeholder: "Unesite prezime",
    },
    {
      label: "Email",
      name: "teacher_email",
      type: "email",
      placeholder: "Unesite email",
    },
    {
      label: "Telefon",
      name: "teacher_phone",
      type: "phone",
      placeholder: "Unesite telefon +387",
    },
    {
      label: "Oslovljavanje",
      name: "teacher_contractions",
      placeholder: "Unesite oslovljavanje (npr. Gdin.)",
    },
    {
      label: "Titula",
      name: "teacher_title",
      placeholder: "Unesite titulu (npr. Prof.)",
    },
    {
      label: "Lozinka",
      name: "teacher_password",
      type: "password",
      placeholder: "Unesite lozinku",
    },
    {
      label: "Konfiguracija boja",
      type: "subtitle",
      icon: FaPalette,
    },
    {
      label: "Paleta boja",
      name: "color_config",
      placeholder: "Odaberite paletu boja",
      type: "select",
      options: [
        { value: 0, label: "Klasična indigo" },
        { value: 1, label: "Plavo-zelena" },
        { value: 2, label: "Neutralna siva" },
        { value: 3, label: "Tamno siva" },
      ],
    },
    <ColorPreviewCard key="colorConfigPreview" />,
    {
      label: "Konfiguracija dizajna",
      type: "subtitle",
      icon: FaRegIdCard,
    },
    {
      label: "Prikaz profesora",
      name: "teacher_display",
      type: "display-choice",
      icon: FaChalkboardTeacher,
      options: [
        { value: "card", label: "Kartice", icon: FaRegIdCard },
        { value: "table", label: "Tabela", icon: FaTable },
      ],
    },
    {
      label: "Prikaz poziva profesora",
      name: "teacher_invite_display",
      type: "display-choice",
      icon: FaUserPlus,
      options: [
        { value: "card", label: "Kartice", icon: FaRegIdCard },
        { value: "table", label: "Tabela", icon: FaTable },
      ],
    },
    {
      label: "Prikaz kurikuluma",
      name: "curriculum_display",
      type: "display-choice",
      icon: FaBook,
      options: [
        { value: "card", label: "Kartice", icon: FaRegIdCard },
        { value: "table", label: "Tabela", icon: FaTable },
      ],
    },
    {
      label: "Prikaz odjeljenja",
      name: "section_display",
      type: "display-choice",
      icon: FaUsers,
      options: [
        { value: "card", label: "Kartice", icon: FaRegIdCard },
        { value: "table", label: "Tabela", icon: FaTable },
      ],
    },
    {
      label: "Prikaz učenika",
      name: "pupil_display",
      type: "display-choice",
      icon: FaUserGraduate,
      options: [
        { value: "card", label: "Kartice", icon: FaRegIdCard },
        { value: "table", label: "Tabela", icon: FaTable },
      ],
    },
    {
      label: "Prikaz poziva učenika",
      name: "pupil_invite_display",
      type: "display-choice",
      icon: FaUserCheck,
      options: [
        { value: "card", label: "Kartice", icon: FaRegIdCard },
        { value: "table", label: "Tabela", icon: FaTable },
      ],
    },
    {
      label: "Prikaz polugodišta",
      name: "semester_display",
      type: "display-choice",
      icon: FaCalendarAlt,
      options: [
        { value: "card", label: "Kartice", icon: FaRegIdCard },
        { value: "table", label: "Tabela", icon: FaTable },
      ],
    },
    {
      label: "Prikaz časova",
      name: "lesson_display",
      type: "display-choice",
      icon: FaCalendarDay,
      options: [
        { value: "card", label: "Kartice", icon: FaRegIdCard },
        { value: "table", label: "Tabela", icon: FaTable },
      ],
    },
    {
      label: "Prikaz izostanaka",
      name: "absence_display",
      type: "display-choice",
      icon: FaUserTimes,
      options: [
        { value: "card", label: "Kartice", icon: FaRegIdCard },
        { value: "table", label: "Tabela", icon: FaTable },
      ],
    },
    {
      label: "Prikaz učionica",
      name: "classroom_display",
      type: "display-choice",
      icon: FaChalkboard,
      options: [
        { value: "card", label: "Kartice", icon: FaRegIdCard },
        { value: "table", label: "Tabela", icon: FaTable },
      ],
    },
  ];

  const createTenant = async (data) => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/superadmin/tenant`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
          body: JSON.stringify(data),
        },
      );
      if (!response.ok) {
        const error_message = await response.text();
        setErrorMessage(error_message);
      } else {
        const newTenant = await response.json();
        setTenants((prev) => [...(prev || []), newTenant]);
        setShowModal(false);
      }
    } catch (error) {
      console.error("Error creating tenant:", error);
      setErrorMessage("Failed to create tenant");
    }
  };

  const handleShowAdministration = (tenant) => {
    setExpandedTenant(tenant);
  };

  return (
    <>
      {expandedTenant ? (
        <>
          <DynamicTab
            title={`${expandedTenant.tenant_name}`}
            titleIcon={FaBuilding}
            activeTab={activeTab}
            setActiveTab={setActiveTab}
            colorConfig={expandedTenant.color_config}
            childrenTabs={[
              {
                label: "Profesori",
                content: (
                  <TeachersPerTenantPageClient
                    tenant={expandedTenant}
                    onBack={
                      tenantAdminMode ? null : () => setExpandedTenant(null)
                    }
                    accessToken={accessToken}
                  />
                ),
                icon: FaChalkboardTeacher,
              },
              {
                label: "Kurikulumi",
                content: (
                  <CurriculumAssignmentPageClient
                    onBack={
                      tenantAdminMode ? null : () => setExpandedTenant(null)
                    }
                    tenant={expandedTenant}
                    accessToken={accessToken}
                  />
                ),
                icon: FaBookOpen,
              },
              {
                label: "Odjeljenja",
                content: (
                  <SectionsPageClient
                    accessToken={accessToken}
                    onBack={
                      tenantAdminMode ? null : () => setExpandedTenant(null)
                    }
                    tenant={expandedTenant}
                    archived={0}
                  />
                ),
                icon: FaLayerGroup,
              },
              {
                label: "Arhivirana odjeljenja",
                content: (
                  <SectionsPageClient
                    accessToken={accessToken}
                    onBack={
                      tenantAdminMode ? null : () => setExpandedTenant(null)
                    }
                    tenant={expandedTenant}
                    archived={1}
                  />
                ),
                icon: FaArchive,
              },
              {
                label: "Polugodišta",
                content: (
                  <TenantSemestersPageClient
                    accessToken={accessToken}
                    tenant={expandedTenant}
                    onBack={
                      tenantAdminMode ? null : () => setExpandedTenant(null)
                    }
                  />
                ),
                icon: FaRegCalendarAlt,
              },
              {
                label: "Učionice",
                content: (
                  <ClassroomPageClient
                    accessToken={accessToken}
                    tenant={expandedTenant}
                    onBack={
                      tenantAdminMode ? null : () => setExpandedTenant(null)
                    }
                  />
                ),
                icon: FaChalkboard,
              },
            ]}
          />
        </>
      ) : !tenantAdminMode ? (
        <>
          <Title icon={FaBuilding}>Institucije</Title>
          <div className="flex justify-end mb-4">
            <AddButton onClick={() => setShowModal(true)}>
              Dodaj instituciju
            </AddButton>
            {errorMessage && (
              <ErrorModal onClose={() => setErrorMessage(null)}>
                {errorMessage}
              </ErrorModal>
            )}
            {showModal && (
              <CreateUpdateModal
                title="Dodaj instituciju"
                fields={tenantFields}
                onClose={() => setShowModal(false)}
                onSave={createTenant}
              />
            )}
          </div>
          <DynamicCardParent
            data={tenants}
            setData={setTenants}
            icon={<FaBuilding />}
            titleField="tenant_name"
            prefix="instituciju"
            deleteUrl={`${process.env.NEXT_PUBLIC_API_BASE_URL}/api/superadmin/tenant`}
            editFields={tenantFields.filter(
              (field) =>
                field.name !== "teacher_password" &&
                field.name !== "tenant_type" &&
                field.name !== "canton_code",
            )}
            editUrl={`${process.env.NEXT_PUBLIC_API_BASE_URL}/api/superadmin/tenant`}
            showEdit={true}
            showDelete={true}
            extraButton={{
              label: "Administracija",
              onClick: (tenant) => handleShowAdministration(tenant),
              icon: FaUserCog,
            }}
            keysToIgnore={[
              "color_config",
              "teacher_display",
              "teacher_invite_display",
              "pupil_display",
              "pupil_invite_display",
              "section_display",
              "curriculum_display",
              "semester_display",
              "teacher_name",
              "teacher_last_name",
              "teacher_email",
              "teacher_phone",
              "teacher_id",
              "lesson_display",
              "absence_display",
              "grade_display",
              "classroom_display",
              "teacher_contractions",
              "teacher_title",
            ]}
            accessToken={accessToken}
            tenantColorConfig={(item) => item.color_config}
          />
        </>
      ) : (
        <Loading />
      )}
    </>
  );
}
