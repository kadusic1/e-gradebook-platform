"use client";
import DynamicCard from "@/app/components/dynamic_card/DynamicCard";
import {
  mapValueToBosnian,
  mapKeyToBosnian,
  getValueColor,
  getTitle,
} from "@/app/components/dynamic_card/DynamicCardParent";
import { FaIdCard, FaKey, FaCopy, FaCheck } from "react-icons/fa";
import { useState } from "react";
import { pupilFields } from "@/app/components/shared_data/pupils_shared";
import CreateUpdateModal from "@/app/components/modal/CreateUpdateModal";
import PasswordChangeModal from "@/app/components/modal/PasswordChangeModal";
import ErrorModal from "@/app/components/modal/ErrorModal";
import SuccessModal from "@/app/components/modal/SuccessModal";
import { teacherFields } from "@/app/components/shared_data/teachers_shared";
import Title from "@/app/components/common/Title";

export default function GeneralProfileClient({
  accessToken,
  generalData,
  userID,
  mode = "pupil",
}) {
  const [generalProfileData, setGeneralProfileData] = useState(generalData);
  const [selectedData, setSelectedData] = useState(null);
  const [showPasswordChange, setShowPasswordChange] = useState(false);
  const [errorMessage, setErrorMessage] = useState();
  const [successMessage, setSuccessMessage] = useState();
  const [copiedCode, setCopiedCode] = useState(false);

  const updateGeneralData = async (data) => {
    try {
      let urlPath;
      if (mode === "pupil") {
        urlPath = `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/pupil/update_general_data`;
      } else {
        urlPath = `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/teacher/${userID}`;
      }
      const resp = await fetch(urlPath, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${accessToken}`,
        },
        body: JSON.stringify(data),
      });
      if (resp.ok) {
        const updatedData = await resp.json();
        setGeneralProfileData(updatedData);
      } else {
        const errorMessage = await resp.text();
        setErrorMessage(errorMessage);
      }
    } catch (error) {
      console.error(error);
    }
  };

  const copyParentAccessCode = async () => {
    if (mode !== "pupil") return;
    try {
      if (generalProfileData?.parent_access_code) {
        await navigator.clipboard.writeText(
          generalProfileData.parent_access_code,
        );
        setCopiedCode(true);
        // Reset the copied state after 2 seconds
        setTimeout(() => {
          setCopiedCode(false);
        }, 2000);
      }
    } catch (error) {
      setErrorMessage("Greška pri kopiranju koda");
    }
  };

  return (
    <>
      {selectedData && (
        <CreateUpdateModal
          title="Uredi osnovne podatke"
          fields={mode === "pupil" ? pupilFields : teacherFields}
          onClose={() => setSelectedData(null)}
          onSave={(updatedData) => {
            updateGeneralData(updatedData);
            setSelectedData(null);
          }}
          initialValues={selectedData}
        />
      )}
      {showPasswordChange && (
        <PasswordChangeModal
          setShowPasswordChange={setShowPasswordChange}
          setErrorMessage={setErrorMessage}
          accessToken={accessToken}
          setSuccessMessage={setSuccessMessage}
        />
      )}
      {errorMessage && (
        <ErrorModal onClose={() => setErrorMessage(null)}>
          {errorMessage}
        </ErrorModal>
      )}
      {successMessage && (
        <SuccessModal onClose={() => setSuccessMessage(null)}>
          {successMessage}
        </SuccessModal>
      )}
      {mode === "teacher" && <Title icon={FaIdCard}>Moj profil</Title>}
      <DynamicCard
        data={generalProfileData || {}}
        titleField={["name", "last_name"]}
        className="mt-8"
        mapValueToBosnian={mapValueToBosnian}
        mapKeyToBosnian={mapKeyToBosnian}
        getValueColor={getValueColor}
        icon={FaIdCard}
        getTitle={getTitle}
        showEdit={true}
        editButton={{
          onClick: (data) => {
            setSelectedData(data);
          },
        }}
        extraButton={{
          onClick: () => setShowPasswordChange(true),
          label: "Promijeni lozinku",
          icon: FaKey,
        }}
        twoColumnsLg={true}
        keysToIgnore={["unenrolled", "parent_access_code"]}
        showDelete={mode === "pupil"}
        deleteButton={{
          onClick: copyParentAccessCode,
          label: copiedCode ? "Kopirano!" : "Kod za roditelje",
          icon: copiedCode ? FaCheck : FaCopy,
        }}
      />
    </>
  );
}
