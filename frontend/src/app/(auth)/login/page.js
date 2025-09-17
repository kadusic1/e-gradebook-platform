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

export default function LoginPage() {
  const router = useRouter();
  const [error, setError] = useState("");

  const fields = [
    {
      label: "Email",
      name: "username",
      type: "text",
      placeholder: "Unesite email",
    },
    {
      label: "Lozinka",
      name: "password",
      type: "password",
      placeholder: "Unesite lozinku",
    },
  ];

  const handleSubmit = async (form) => {
    setError("");
    if (form.username && form.password) {
      const res = await signIn("credentials", {
        redirect: false,
        username: form.username,
        password: form.password,
      });
      if (res.ok) {
        router.push("/");
        return;
      }
      if (res.status === 401) {
        setError("Neispravno korisničko ime ili lozinka");
        return;
      }
      if (res.error) {
        setError("Greška pri prijavi");
        return;
      }
    } else {
      setError("Unesite korisničko ime i lozinku");
    }
  };

  return (
    <ParentCard>
      <Title icon={FaSignInAlt}>eDnevnik Prijava</Title>
      <Form
        fields={fields}
        onSubmit={handleSubmit}
        showCancel={false}
        submitText="Prijavi se"
        error={error}
        fixedHeight={false}
      />
      <Text className="mt-8 text-center md:text-left">
        Nemate korisnički račun? Registrujte se kao{" "}
        <CustomLink href="/register_teacher">nastavnik</CustomLink> ili{" "}
        <CustomLink href="/register_pupil">učenik</CustomLink>.
      </Text>
      <Text className="mt-4 text-center md:text-left">
        Da li ste roditelj? Prijavite se{" "}
        <CustomLink href="/parent_login">ovdje</CustomLink>.
      </Text>
    </ParentCard>
  );
}
