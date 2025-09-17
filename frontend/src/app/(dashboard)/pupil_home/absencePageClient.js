"use client";
import { useState, useEffect, useMemo } from "react";
import BackButton from "@/app/components/common/BackButton";
import DynamicCardParent from "@/app/components/dynamic_card/DynamicCardParent";
import {
  FaUserTimes,
  FaCheck,
  FaTimes,
  FaCalculator,
  FaCheckCircle,
  FaTimesCircle,
} from "react-icons/fa";
import { formatDateToDDMMYYYY } from "@/app/util/date_util";
import Title from "@/app/components/common/Title";
import ConfirmModal from "@/app/components/modal/ConfirmModal";
import Subtitle from "@/app/components/common/Subtitle";
import Spacer from "@/app/components/common/Spacer";

export default function AbsencePageClient({
  accessToken,
  tenantID,
  pupilID,
  section,
  colorConfig,
  setShowAbsencePage,
  mode = "pupil",
  displayMode,
  archived = 0,
}) {
  const [absences, setAbsences] = useState([]);
  const [action, setAction] = useState(null);

  const fetchAbsences = async () => {
    try {
      let urlPath;
      if (mode === "pupil") {
        urlPath = `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/pupil/get_absent_attendances_for_pupil/${tenantID}/${section.id}/${pupilID}`;
      } else {
        urlPath = `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/teacher/get_absent_attendances_for_section/${tenantID}/${section.id}`;
      }
      const response = await fetch(urlPath, {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      });
      if (!response.ok) {
        throw new Error("Failed to fetch absences");
      }
      const data = await response.json();
      setAbsences(data);
    } catch (error) {
      console.error(error);
    }
  };

  const {
    formattedAbsences,
    totalAbsences,
    unexcusedAbsences,
    excusedAbsences,
  } = useMemo(() => {
    if (!absences || absences.length === 0) {
      return {
        formattedAbsences: [],
        totalAbsences: 0,
        unexcusedAbsences: 0,
        excusedAbsences: 0,
      };
    }
    const formatted = absences?.map((absence) => ({
      ...absence,
      date: formatDateToDDMMYYYY(absence.date),
    }));
    const totalCount = formatted.length;
    const unexcusedCount = formatted.filter(
      (absence) => absence.status === "unexcused",
    ).length;
    const excusedCount = formatted.filter(
      (absence) => absence.status === "excused",
    ).length;
    return {
      formattedAbsences: formatted,
      totalAbsences: totalCount,
      unexcusedAbsences: unexcusedCount,
      excusedAbsences: excusedCount,
    };
  }, [absences]);

  useEffect(() => {
    fetchAbsences();
  }, []);

  const handleAttendanceAction = async (action) => {
    // Remove pupil_name from action
    action = {
      type: action.type,
      pupil_id: action.pupil_id,
      lesson_id: action.lesson_id,
    };
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/teacher/handle_attendance_action/${tenantID}`,
        {
          method: "POST",
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
          body: JSON.stringify(action),
        },
      );
      if (response.ok) {
        // Find the absence with pupil_id and lesson_id and set its status
        // to type
        setAbsences((prev) =>
          prev.map((absence) =>
            absence.pupil_id === action.pupil_id &&
            absence.lesson_id === action.lesson_id
              ? { ...absence, status: action.type }
              : absence,
          ),
        );
      }
    } catch (error) {
      console.error("Error handling attendance action:", error);
    }
  };

  return (
    <>
      <Title colorConfig={colorConfig} icon={FaUserTimes}>
        {mode === "pupil" ? "Moji izostanci" : "Izostanci"} - {section?.name}
      </Title>
      <div className="flex justify-end mb-4">
        <BackButton
          colorConfig={colorConfig}
          onClick={() => setShowAbsencePage(false)}
        />
      </div>
      <Spacer className="mb-8">
        <Subtitle
          icon={FaCalculator}
          colorConfig={colorConfig}
          showLine={false}
        >
          Ukupno: {totalAbsences}
        </Subtitle>
        <Subtitle
          icon={FaCheckCircle}
          colorConfig={colorConfig}
          showLine={false}
        >
          Opravdanih: <span className="text-green-500">{excusedAbsences}</span>
        </Subtitle>
        <Subtitle
          icon={FaTimesCircle}
          colorConfig={colorConfig}
          showLine={false}
        >
          Neopravdanih:{" "}
          <span className="text-red-500">{unexcusedAbsences}</span>
        </Subtitle>
      </Spacer>
      {action && (
        <ConfirmModal
          colorConfig={colorConfig}
          title="Potvrda akcije"
          onClose={() => setAction(null)}
          onConfirm={() => {
            handleAttendanceAction(action);
            setAction(null);
          }}
        >
          Da li ste sigurni da želite označiti izostanak učenika{" "}
          <span className="font-bold">{action.pupil_name}</span> kao{" "}
          <span className="font-bold">
            {action.type === "excused" ? "opravdan" : "neopravdan"}
          </span>
          ?
        </ConfirmModal>
      )}
      <DynamicCardParent
        data={formattedAbsences}
        icon={<FaUserTimes />}
        tenantColorConfig={colorConfig}
        titleField={["subject_name", "date"]}
        keysToIgnore={["pupil_id", "lesson_id"]}
        emptyMessage={
          mode === "pupil"
            ? "Trenutno nemate izostanaka."
            : "Nema izostanaka za ovo odjeljenje."
        }
        accessToken={accessToken}
        showEdit={mode === "teacher" && archived == 0 ? true : false}
        showDelete={mode === "teacher" && archived == 0 ? true : false}
        editButton={{
          label: "Opravdan",
          onClick: (absence) => {
            setAction({
              type: "excused",
              pupil_id: absence.pupil_id,
              lesson_id: absence.lesson_id,
              pupil_name: `${absence.name} ${absence.last_name}`,
            });
          },
          icon: FaCheck,
        }}
        deleteButton={{
          label: "Neopravdan",
          onClick: (absence) => {
            setAction({
              type: "unexcused",
              pupil_id: absence.pupil_id,
              lesson_id: absence.lesson_id,
              pupil_name: `${absence.name} ${absence.last_name}`,
            });
          },
          icon: FaTimes,
        }}
        mode={displayMode}
      />
    </>
  );
}
