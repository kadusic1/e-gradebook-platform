"use client";
import React from "react";
import ScheduleTable from "../../components/table/ScheduleTable";
import Title from "@/app/components/common/Title";
import { useState, useEffect } from "react";
import { getScheduleForSection } from "@/app/api/helper/scheduleApi";
import { FaCalendarDay } from "react-icons/fa";

export default function SchedulePageClient({
  colorConfig,
  setShowSchedulePage,
  section,
  tenantID,
  accessToken,
  readOnly = false,
  initialSchedule = [],
  teacherMode = false,
  archived = 0,
}) {
  const [scheduleGroups, setScheduleGroups] = useState(initialSchedule);

  useEffect(() => {
    if (!teacherMode) {
      getScheduleForSection(
        section.id,
        tenantID,
        accessToken,
        setScheduleGroups,
      );
    }
  }, [section, tenantID, accessToken]);

  return (
    <div>
      {!teacherMode && (
        <Title colorConfig={colorConfig} icon={FaCalendarDay}>
          Raspored časova - {section.name}
        </Title>
      )}
      <ScheduleTable
        initialScheduleGroups={scheduleGroups}
        colorConfig={colorConfig}
        className="mt-4"
        section={section}
        accessToken={accessToken}
        setShowSchedulePage={setShowSchedulePage}
        tenantID={tenantID}
        readOnly={readOnly || archived == 1}
        teacherMode={teacherMode}
      />
    </div>
  );
}
