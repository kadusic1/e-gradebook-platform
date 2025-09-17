"use client";
import { useRouter } from "next/navigation";
import { useState } from "react";
import Form from "../../components/form/Form";
import Title from "../../components/common/Title";
import { FaUserPlus } from "react-icons/fa";
import ParentCard from "../../components/common/ParentCard";
import { teacherFields } from "@/app/components/shared_data/teachers_shared";
import CustomLink from "@/app/components/common/CustomLink";
import Text from "@/app/components/common/Text";
import SuccessModal from "@/app/components/modal/SuccessModal";

export default function TeacherRegisterPage() {
  const [errorMessage, setErrorMessage] = useState("");
  const [successMessage, setSuccessMessage] = useState("");

  const handleSubmit = async (data) => {
    try {
      setErrorMessage("");
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/common/register_teacher`,
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
        "Nastavnik je uspješno kreiran! Verifikujte račun putem emaila kako biste se mogli prijaviti.",
      );
    } catch (error) {
      console.error("Error creating teacher:", error);
      setErrorMessage("Failed to create teacher");
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
        <Title icon={FaUserPlus}>Registracija Nastavnika</Title>
        <Form
          fields={[
            ...teacherFields,
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
