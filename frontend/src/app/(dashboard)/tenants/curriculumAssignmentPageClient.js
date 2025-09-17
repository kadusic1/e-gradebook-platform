// CurriculumAssignmentPageClient.js
"use client";
import { FaBookOpen } from "react-icons/fa";
import AssignmentPage from "@/app/components/client_pages/AssignmentPage";
import Title from "@/app/components/common/Title";
import { useEffect, useState } from "react";

export default function CurriculumAssignmentPageClient({
  onBack,
  tenant,
  accessToken,
}) {
  const [curriculums, setCurriculums] = useState([]);
  const [curriculumsForAssignment, setCurriculumsForAssignment] = useState([]);

  const fetchCurriculumData = async () => {
    const response = await fetch(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/get_curriculums_for_tenant/${tenant?.id}`,
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      },
    );
    if (response.ok) {
      const curriculums = await response.json();
      setCurriculums(curriculums);
    }
  };

  const fetchCurriculumAssignmentData = async () => {
    const response = await fetch(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/get_curriculums_for_assignment/${tenant?.id}`,
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      },
    );
    if (response.ok) {
      const curriculumsForAssignment = await response.json();
      setCurriculumsForAssignment(curriculumsForAssignment);
    }
  };

  useEffect(() => {
    fetchCurriculumData();
    fetchCurriculumAssignmentData();
  }, []);

  const config = {
    // API endpoints
    assignEndpoint: `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/superadmin/assign_curriculums_to_tenant/${tenant.id}`,
    unassignEndpoint: `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/superadmin/unassign_curriculum_from_tenant/${tenant.id}`,

    // UI configuration
    addButtonText: "Dodaj kurikulum",
    modalTitle: "Odaberi kurikulume",
    searchPlaceholder: "Pretraži naziv kurikuluma",
    modalNote: "Možete odabrati više kurikuluma za školu.",

    // Field mappings
    keyField: "curriculum_code",
    labelField: "curriculum_name",
    titleField: "curriculum_name",
    prefix: "kurikulum",
    keysToIgnore: ["curriculum_code", "npp_code"],

    // Icon
    icon: <FaBookOpen />,

    // Error messages
    assignErrorMessage:
      "Došlo je do greške prilikom dodijeljivanja kurikuluma školi.",
    unassignErrorMessage: "Došlo je do greške prilikom uklanjanja kurikuluma.",
    orderFields: ["npp_name", "class_code"],
  };

  return (
    <>
      <Title icon={FaBookOpen} colorConfig={tenant.color_config}>
        Kurikulumi
      </Title>
      <AssignmentPage
        assignedItems={curriculums}
        setAssignedItems={setCurriculums}
        availableItems={curriculumsForAssignment}
        setAvailableItems={setCurriculumsForAssignment}
        tenant={tenant}
        onBack={onBack}
        config={config}
        accessToken={accessToken}
        colorConfig={tenant.color_config}
        mode={tenant.curriculum_display}
      />
    </>
  );
}
