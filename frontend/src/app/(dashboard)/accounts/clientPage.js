"use client";

import DynamicTab from "@/app/components/dynamic_card/DynamicTab";
import TeacherAccountsPageClient from "./teacherAccountsPageClient";
import PupilAccountsPageClient from "./pupilAccountsPageClient";
import { FaUser, FaChalkboardTeacher } from "react-icons/fa";
import { useState } from "react";

export default function AccountsPageClient({
  initialTeachers = [],
  initialPupils = [],
  accessToken,
}) {
  const [activeTab, setActiveTab] = useState(0);
  const [teachers, setTeachers] = useState(initialTeachers);
  const [pupils, setPupils] = useState(initialPupils);

  return (
    <>
      <DynamicTab
        title={"Korisnički nalozi"}
        titleIcon={FaUser}
        activeTab={activeTab}
        setActiveTab={setActiveTab}
        childrenTabs={[
          {
            label: "Profesori",
            content: (
              <TeacherAccountsPageClient
                teachers={teachers}
                setTeachers={setTeachers}
                accessToken={accessToken}
              />
            ),
            icon: FaChalkboardTeacher,
          },
          {
            label: "Učenici",
            content: (
              <PupilAccountsPageClient
                pupils={pupils}
                setPupils={setPupils}
                accessToken={accessToken}
              />
            ),
            icon: FaUser,
          },
        ]}
      />
    </>
  );
}
