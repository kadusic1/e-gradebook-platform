"use client";
import { useState } from "react";
import Form from "../../components/form/Form";
import Title from "../../components/common/Title";
import { FaUserPlus } from "react-icons/fa";
import ParentCard from "../../components/common/ParentCard";
import { pupilFields } from "@/app/components/shared_data/pupils_shared";
import CustomLink from "@/app/components/common/CustomLink";
import Text from "@/app/components/common/Text";
import SuccessModal from "@/app/components/modal/SuccessModal";

export default function PupilRegisterPage() {
  const [errorMessage, setErrorMessage] = useState("");
  const [successMessage, setSuccessMessage] = useState("");

  const handleSubmit = async (data) => {
    try {
      setErrorMessage("");
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/common/register_pupil`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(data),
        },
      );
      if (!response.ok) {
        const resp = await response.text();
        setErrorMessage(resp);
        return;
      }
      setSuccessMessage(
        "Učenik je uspješno kreiran! Verifikujte račun putem emaila kako biste se mogli prijaviti.",
      );
    } catch (error) {
      console.error("Error creating pupil:", error);
      setErrorMessage("Failed to create pupil");
    }
  };

  return (
    <>
      {successMessage && (
        <SuccessModal onClose={() => setSuccessMessage("")}>
          {successMessage}
        </SuccessModal>
      )}
      <ParentCard>
        <Title icon={FaUserPlus}>Registracija Učenika</Title>
        <Form
          fields={[
            ...pupilFields,
            {
              label: "Lozinka",
              name: "password",
              type: "password",
              placeholder: "Unesite lozinku",
            },
          ]}
          onSubmit={handleSubmit}
          showCancel={false}
          submitText="Registruj se"
          error={errorMessage}
          fixedHeight={false}
        />
        <Text className="mt-8 text-center md:text-left">
          Već imate korisnički račun? Prijavite se{" "}
          <CustomLink href="/login">ovdje</CustomLink>.
        </Text>
      </ParentCard>
    </>
  );
}
