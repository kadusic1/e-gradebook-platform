"use client";
import { FaUser, FaIdCard, FaChartBar } from "react-icons/fa";
import DynamicTab from "@/app/components/dynamic_card/DynamicTab";
import { useState } from "react";
import GeneralProfileClient from "./generalDataPageClient";
import PupilStatisticsDataClient from "./statsDataPageClient";

export default function PupilProfileClient({
  userID,
  accessToken,
  generalData,
  statsData,
}) {
  const [activeTab, setActiveTab] = useState(0);

  return (
    <>
      <DynamicTab
        title="Moj profil"
        titleIcon={FaUser}
        activeTab={activeTab}
        setActiveTab={setActiveTab}
        childrenTabs={[
          {
            label: "Osnovni podaci",
            content: (
              <GeneralProfileClient
                accessToken={accessToken}
                generalData={generalData}
              />
            ),
            icon: FaIdCard,
          },
          {
            label: "Statistički podaci",
            content: (
              <PupilStatisticsDataClient
                accessToken={accessToken}
                statsData={statsData}
                pupilID={userID}
              />
            ),
            icon: FaChartBar,
          },
        ]}
      />
    </>
  );
}
