"use client";
import { useState } from "react";
import AddButton from "@/app/components/common/AddButton";
import CreateUpdateModal from "@/app/components/modal/CreateUpdateModal";
import ErrorModal from "../../components/modal/ErrorModal";
import DynamicCardParent from "../../components/dynamic_card/DynamicCardParent";
import { FaBuilding, FaChalkboardTeacher } from "react-icons/fa";
import { teacherFields } from "@/app/components/shared_data/teachers_shared";
import Title from "@/app/components/common/Title";
import Modal from "@/app/components/modal/Modal";

export default function TeacherAccountsPageClient({
  teachers,
  setTeachers,
  accessToken,
}) {
  const [showModal, setShowModal] = useState(false);
  const [errorMessage, setErrorMessage] = useState(null);
  const [teacherTenants, setTeacherTenants] = useState(null);

  const createTeacher = async (data) => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/teacher`,
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
        const resp = await response.text();
        setErrorMessage(resp);
        return;
      }
      const newTeacher = await response.json();
      setTeachers((prev) => [...(prev || []), newTeacher]);
      setShowModal(false);
    } catch (error) {
      console.error("Error creating teacher:", error);
      setErrorMessage("Failed to create teacher");
    }
  };

  const getTeacherTenants = async (teacher_id) => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/tenants_for_teacher/${teacher_id}`,
        {
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
        },
      );
      if (!response.ok) {
        const resp = await response.text();
        setErrorMessage(resp);
        return;
      }
      const data = await response.json();
      setTeacherTenants(data || []);
    } catch (error) {
      console.error(error);
      setErrorMessage(error);
    }
  };

  return (
    <>
      <Title icon={FaChalkboardTeacher}>Profesorski nalozi</Title>
      <div className="flex justify-end mb-4">
        <AddButton onClick={() => setShowModal(true)}>
          Dodaj profesora
        </AddButton>
        {errorMessage && (
          <ErrorModal onClose={() => setErrorMessage(null)}>
            {errorMessage}
          </ErrorModal>
        )}
        {showModal && (
          <CreateUpdateModal
            title="Dodaj profesora"
            fields={[
              ...teacherFields,
              {
                label: "Lozinka",
                name: "password",
                type: "password",
                placeholder: "Unesite lozinku",
              },
            ]}
            onClose={() => setShowModal(false)}
            onSave={createTeacher}
          />
        )}
      </div>
      {teacherTenants && (
        <Modal
          onClose={() => {
            setTeacherTenants(null);
          }}
        >
          <Title icon={FaBuilding}>Institucije za nastavnika</Title>
          <DynamicCardParent
            data={teacherTenants}
            titleField="tenant_name"
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
            ]}
            icon={FaBuilding}
            tenantColorConfig={(item) => item.color_config}
            twoColumnsCard={false}
          />
        </Modal>
      )}
      <DynamicCardParent
        data={teachers}
        setData={setTeachers}
        icon={<FaBuilding />}
        titleField="email"
        prefix="profesora"
        deleteUrl={`${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/teacher`}
        editFields={teacherFields}
        editUrl={`${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/teacher`}
        showEdit={true}
        showDelete={true}
        accessToken={accessToken}
        extraButton={{
          label: "Institucije",
          onClick: (teacher) => {
            getTeacherTenants(teacher?.id);
          },
          icon: FaBuilding,
        }}
      />
    </>
  );
}
