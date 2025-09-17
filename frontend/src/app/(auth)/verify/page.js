"use client";
import ParentCard from "../../components/common/ParentCard";
import Text from "@/app/components/common/Text";
import Title from "@/app/components/common/Title";
import CustomLink from "@/app/components/common/CustomLink";
import { FaCheckCircle } from "react-icons/fa";
import { useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";
import { Suspense } from "react";

function VerifyContent() {
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);
  const searchParams = useSearchParams();
  const token = searchParams.get("token") || "";

  const verifyAccount = async (token) => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/common/verify_account?token=${token}`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
        },
      );

      if (!response.ok) {
        setError("Greška pri verifikaciji računa.");
      }
    } catch (error) {
      console.error("Error verifying account:", error);
      setError("Greška pri verifikaciji računa.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (token) {
      verifyAccount(token);
    } else {
      setError("Token za verifikaciju nije pronađen.");
      setLoading(false);
    }
  }, [token]);

  if (loading) {
    return (
      <ParentCard className="border-2 border-black">
        <Title>Verifikacija u toku...</Title>
        <Text className="text-center mt-4 font-bold">Molimo sačekajte...</Text>
      </ParentCard>
    );
  }

  if (error) {
    return (
      <ParentCard className="border-2 border-black">
        <Title className="text-red-600">Greška</Title>
        <Text className="text-center mt-4 font-bold text-red-600">{error}</Text>
      </ParentCard>
    );
  }

  return (
    <ParentCard className="border-2 border-black">
      <Title icon={FaCheckCircle}>eDnevnik Verifikacija</Title>
      <Text className="text-center mt-4 font-bold">
        Uspješno ste verifikovali svoj račun. Sada se možete{" "}
        <CustomLink href="/login">prijaviti</CustomLink>.
      </Text>
    </ParentCard>
  );
}

export default function Verify() {
  return (
    <Suspense
      fallback={
        <ParentCard className="border-2 border-black">
          <Title>Učitavanje...</Title>
          <Text className="text-center mt-4 font-bold">
            Molimo sačekajte...
          </Text>
        </ParentCard>
      }
    >
      <VerifyContent />
    </Suspense>
  );
}
