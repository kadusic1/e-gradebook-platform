"use client";
import { useState } from "react";
import AddButton from "@/app/components/common/AddButton";
import CreateUpdateModal from "@/app/components/modal/CreateUpdateModal";
import ErrorModal from "../../components/modal/ErrorModal";
import DynamicCardParent from "../../components/dynamic_card/DynamicCardParent";
import { FaLayerGroup, FaBuilding } from "react-icons/fa";
import { pupilFields } from "@/app/components/shared_data/pupils_shared";
import Title from "@/app/components/common/Title";
import Modal from "@/app/components/modal/Modal";

export default function PupilAccountsPageClient({
  pupils,
  setPupils,
  accessToken,
}) {
  const [showModal, setShowModal] = useState(false);
  const [errorMessage, setErrorMessage] = useState(null);
  const [pupilTenants, setPupilTenants] = useState(null);

  const createPupil = async (data) => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/pupils`,
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
      const newPupil = await response.json();
      setPupils((prev) => [...(prev || []), newPupil]);
      setShowModal(false);
    } catch (error) {
      console.error("Error creating pupil:", error);
      setErrorMessage("Failed to create pupil");
    }
  };

  const getPupilTenants = async (pupil_id) => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/tenants_for_pupil/${pupil_id}`,
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
      setPupilTenants(data || []);
    } catch (error) {
      console.error(error);
      setErrorMessage(error);
    }
  };

  return (
    <>
      <Title icon={FaLayerGroup}>Učenički nalozi</Title>
      <div className="flex justify-end mb-4">
        <AddButton onClick={() => setShowModal(true)}>Dodaj učenika</AddButton>
        {errorMessage && (
          <ErrorModal onClose={() => setErrorMessage(null)}>
            {errorMessage}
          </ErrorModal>
        )}
        {showModal && (
          <CreateUpdateModal
            title="Dodaj učenika"
            fields={[
              ...pupilFields,
              {
                label: "Lozinka",
                name: "password",
                type: "password",
                placeholder: "Unesite lozinku",
              },
            ]}
            onClose={() => setShowModal(false)}
            onSave={(data) => {
              createPupil(data);
            }}
          />
        )}
      </div>
      {pupilTenants && (
        <Modal
          onClose={() => {
            setPupilTenants(null);
          }}
        >
          <Title icon={FaBuilding}>Institucije za učenika</Title>
          <DynamicCardParent
            data={pupilTenants}
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
        data={pupils}
        setData={setPupils}
        icon={<FaLayerGroup />}
        titleField={["name", "last_name"]}
        prefix="učenika"
        deleteUrl={`${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/pupils`}
        editFields={pupilFields}
        editUrl={`${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/pupils`}
        showEdit={true}
        showDelete={true}
        accessToken={accessToken}
        extraButton={{
          label: "Institucije",
          onClick: (pupil) => {
            getPupilTenants(pupil?.id);
          },
          icon: FaBuilding,
        }}
        keysToIgnore={["parent_access_code", "unenrolled"]}
      />
    </>
  );
}
