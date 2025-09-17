"use client";
import DynamicCard from "@/app/components/dynamic_card/DynamicCard";
import {
  mapValueToBosnian,
  mapKeyToBosnian,
  getValueColor,
  getTitle,
} from "@/app/components/dynamic_card/DynamicCardParent";
import { FaChartBar } from "react-icons/fa";
import { useState } from "react";
import { pupilStatisticsFields } from "@/app/components/shared_data/pupils_shared";
import CreateUpdateModal from "@/app/components/modal/CreateUpdateModal";

export default function PupilStatisticsDataClient({
  accessToken,
  statsData,
  pupilID,
}) {
  const [statsProfileData, setStatsProfileData] = useState(statsData);
  const [selectdPupilData, setSelectedPupilData] = useState(null);

  const updateStatusPupilData = async (data) => {
    try {
      const resp = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/pupil/update_statistics_fields/${pupilID}`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
          body: JSON.stringify(data),
        },
      );
      if (resp.ok) {
        const updatedStats = await resp.json();
        setStatsProfileData(updatedStats);
      }
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <>
      {selectdPupilData && (
        <CreateUpdateModal
          title="Uredi statističke podatke"
          fields={pupilStatisticsFields}
          onClose={() => setSelectedPupilData(null)}
          onSave={(updatedData) => {
            updateStatusPupilData(updatedData);
            setSelectedPupilData(null);
          }}
          initialValues={selectdPupilData}
        />
      )}
      <DynamicCard
        data={statsProfileData || {}}
        titleField=""
        textTitle="Statistika"
        className="mt-8"
        mapValueToBosnian={mapValueToBosnian}
        mapKeyToBosnian={mapKeyToBosnian}
        getValueColor={getValueColor}
        icon={FaChartBar}
        getTitle={getTitle}
        showEdit={true}
        editButton={{
          onClick: (data) => {
            setSelectedPupilData(data);
          },
        }}
        twoColumnsLg={true}
      />
    </>
  );
}
