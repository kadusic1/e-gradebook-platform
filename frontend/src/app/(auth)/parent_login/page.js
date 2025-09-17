"use client";
import { useRouter } from "next/navigation";
import { useState } from "react";
import Form from "../../components/form/Form";
import Title from "../../components/common/Title";
import { FaSignInAlt } from "react-icons/fa";
import ParentCard from "../../components/common/ParentCard";
import { signIn } from "next-auth/react";
import Text from "@/app/components/common/Text";
import CustomLink from "@/app/components/common/CustomLink";

export default function ParentLoginPage() {
  const router = useRouter();
  const [error, setError] = useState("");

  const fields = [
    {
      label: "Roditeljski kod",
      name: "parent_access_code",
      type: "text",
      placeholder: "Unesite roditeljski kod",
    },
  ];

  const handleSubmit = async (form) => {
    setError("");
    if (form.parent_access_code) {
      const res = await signIn("parent-access", {
        redirect: false,
        parent_access_code: form.parent_access_code,
      });
      
      if (res.ok) {
        router.push("/");
        return;
      }
      if (res.status === 401) {
        setError("Neispravan roditeljski kod");
        return;
      }
      if (res.error) {
        setError("Greška pri prijavi");
        return;
      }
    } else {
      setError("Unesite roditeljski kod");
    }
  };

  return (
    <ParentCard>
      <Title icon={FaSignInAlt}>eDnevnik Prijava za roditelje</Title>
      <Form
        fields={fields}
        onSubmit={handleSubmit}
        showCancel={false}
        submitText="Prijavi se"
        error={error}
        fixedHeight={false}
      />
      <Text className="mt-8 text-center md:text-left">
        Roditeljski kod možete zatražiti od vašeg djeteta.{" "}
        Vrati me na{" "}
        <CustomLink href="/login">prijavu</CustomLink>.
      </Text>
    </ParentCard>
  );
}