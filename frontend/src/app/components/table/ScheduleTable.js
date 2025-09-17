import React, { useState, useEffect, useRef } from "react";
import { getColor } from "../colors/colors";
import { FaTimes, FaSave } from "react-icons/fa";
import Button from "../common/Button";
import ConfirmModal from "../modal/ConfirmModal";
import BackButton from "../common/BackButton";
import EmptyState from "../dynamic_card/EmptyState";
import { formatToDateTime } from "@/app/util/date_util";
import {
  Table,
  TableRow,
  TableHead,
  TableHeader,
  TableCell,
  TableBody,
} from "./TableComponents";
import { PDFButton } from "../common/PDFButton";
import { handleDownloadPDF } from "@/app/util/pdf_util";

const WEEKDAYS = ["Ponedjeljak", "Utorak", "Srijeda", "Četvrtak", "Petak"];

export default function ScheduleTable({
  initialScheduleGroups,
  colorConfig,
  className,
  section,
  accessToken,
  setShowSchedulePage,
  tenantID,
  readOnly = false,
  teacherMode = false,
  showButtons = true,
}) {
  const default_schedule = [
    {
      time_period: {
        section_id: section?.id,
        start_time: "08:00",
        end_time: "08:45",
      },
      schedules: [],
    },
    {
      time_period: {
        section_id: section?.id,
        start_time: "08:55",
        end_time: "09:40",
      },
      schedules: [],
    },
    {
      time_period: {
        section_id: section?.id,
        start_time: "09:50",
        end_time: "10:35",
      },
      schedules: [],
    },
    {
      time_period: {
        section_id: section?.id,
        start_time: "10:45",
        end_time: "11:30",
      },
      schedules: [],
    },
    {
      time_period: {
        section_id: section?.id,
        start_time: "11:40",
        end_time: "12:25",
      },
      schedules: [],
    },
    {
      time_period: {
        section_id: section?.id,
        start_time: "12:35",
        end_time: "13:20",
      },
      schedules: [],
    },
  ];

  const [scheduleGroups, setScheduleGroups] = useState(
    initialScheduleGroups && initialScheduleGroups.length > 0
      ? initialScheduleGroups
      : readOnly
        ? []
        : default_schedule,
  );

  // Add this effect:
  useEffect(() => {
    if (initialScheduleGroups && initialScheduleGroups.length > 0) {
      setScheduleGroups(initialScheduleGroups);
    }
  }, [initialScheduleGroups]);
  const [subjects, setSubjects] = useState([]);
  const [classrooms, setClassrooms] = useState([]);
  const [editingTime, setEditingTime] = useState(null);
  const [dragOverCell, setDragOverCell] = useState(null);
  const [showConfirmModal, setShowConfirmModal] = useState(false);
  const [pdfLoading, setPdfLoading] = useState(false);
  const tableRef = useRef();

  const downloadTeacherSchedule = async () => {
    let session = null;
    if (teacherMode) {
      const { getSession } = await import("next-auth/react");
      session = await getSession();
    }
    handleDownloadPDF({
      headerElements: [
        `<h1 class="font-bold text-xl text-center">
        ${
          teacherMode
            ? `Raspored nastave za ${session?.user?.name} ${session?.user?.lastName}`
            : `Raspored nastave za ${section?.name}`
        } - važi od ${formatToDateTime()}
      </h1>`,
      ],
      filename: "raspored.pdf",
      landscape: true,
      setPdfLoading,
      pdfElementRef: tableRef,
    });
  };

  const primaryTextColor = getColor("primary", "text", colorConfig);
  const secondaryTextColor = getColor("secondary", "text", colorConfig);
  const primaryBgColor = getColor("primaryComplement", "bg", colorConfig);

  const getSubjectsForCurriculum = async (curriculumCode) => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/common/get_subjects_for_curriculum/${curriculumCode}`,
        {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
        },
      );
      if (!response.ok) {
        throw new Error("Failed to fetch subjects");
      }
      const data = await response.json();
      setSubjects(data);
    } catch (error) {
      console.error("Error fetching subjects:", error);
    }
  };

  useEffect(() => {
    if (readOnly) return;
    if (section?.curriculum_code) {
      getSubjectsForCurriculum(section.curriculum_code);
    }
  }, [section?.curriculum_code]);

  const getClassrooms = async () => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/get_all_classrooms/${tenantID}`,
        {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
        },
      );
      if (!response.ok) {
        throw new Error("Failed to fetch classrooms");
      }
      const data = await response.json();
      setClassrooms(data);
    } catch (error) {
      console.error("Error fetching classrooms:", error);
    }
  };

  useEffect(() => {
    if (readOnly) return;
    getClassrooms();
  }, [tenantID, accessToken]);

  const handleTimeDoubleClick = (periodIdx, field) => {
    setEditingTime(`${periodIdx}-${field}`);
  };

  const handleTimeChange = (e, periodIdx, field) => {
    const newTime = e.target.value;
    setScheduleGroups((prev) =>
      prev.map((group, idx) =>
        idx === periodIdx
          ? {
              ...group,
              time_period: {
                ...group.time_period,
                [field]: newTime,
              },
            }
          : group,
      ),
    );
  };

  const handleTimeBlur = () => {
    setEditingTime(null);
  };

  const handleTimeKeyPress = (e) => {
    if (e.key === "Enter") {
      setEditingTime(null);
    }
  };

  // Drag and Drop handlers
  const handleDragStart = (e, subject) => {
    e.dataTransfer.setData("application/json", JSON.stringify(subject));
    e.dataTransfer.effectAllowed = "copy";
  };

  // Classroom drag handler
  const handleClassroomDragStart = (e, classroom) => {
    e.dataTransfer.setData(
      "application/json",
      JSON.stringify({ ...classroom, isClassroom: true }),
    );
    e.dataTransfer.effectAllowed = "copy";
  };

  const handleDragOver = (e) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = "copy";
  };

  const handleDragEnter = (e, periodIdx, weekday) => {
    e.preventDefault();
    setDragOverCell(`${periodIdx}-${weekday}`);
  };

  const handleDragLeave = (e) => {
    // Only clear if we're leaving the cell entirely
    if (!e.currentTarget.contains(e.relatedTarget)) {
      setDragOverCell(null);
    }
  };

  const handleDrop = (e, periodIdx, weekday) => {
    e.preventDefault();
    setDragOverCell(null);

    try {
      const subjectData = JSON.parse(
        e.dataTransfer.getData("application/json"),
      );

      // Handle classroom drop
      if (subjectData.isClassroom) {
        // Only allow dropping on cells with existing subjects
        const existingSubject = scheduleGroups[periodIdx].schedules.find(
          (s) => s.weekday === weekday.toLowerCase(),
        );

        if (existingSubject) {
          // Update classroom for existing subject
          setScheduleGroups((prev) =>
            prev.map((group, idx) =>
              idx === periodIdx
                ? {
                    ...group,
                    schedules: group.schedules.map((s) =>
                      s.weekday === weekday.toLowerCase()
                        ? { ...s, classroom_code: subjectData.code }
                        : s,
                    ),
                  }
                : group,
            ),
          );
        }
        return;
      }

      // Check if there's already a subject in this slot
      const existingSubject = scheduleGroups[periodIdx].schedules.find(
        (s) => s.weekday === weekday.toLowerCase(),
      );

      if (existingSubject) {
        // Remove existing subject first
        setScheduleGroups((prev) =>
          prev.map((group, idx) =>
            idx === periodIdx
              ? {
                  ...group,
                  schedules: group.schedules.filter(
                    (s) => s.weekday !== weekday.toLowerCase(),
                  ),
                }
              : group,
          ),
        );
      }

      // Add new subject
      const newScheduleItem = {
        section_id: section?.id,
        subject_code: subjectData.subject_code,
        subject_name: subjectData.subject_name,
        weekday: weekday.toLowerCase(),
        classroom_code: null,
      };

      setScheduleGroups((prev) =>
        prev.map((group, idx) =>
          idx === periodIdx
            ? {
                ...group,
                schedules: [...group.schedules, newScheduleItem],
              }
            : group,
        ),
      );
    } catch (error) {
      console.error("Error dropping subject:", error);
    }
  };

  const handleRemoveSubject = (periodIdx, weekday) => {
    setScheduleGroups((prev) =>
      prev.map((group, idx) =>
        idx === periodIdx
          ? {
              ...group,
              schedules: group.schedules.filter(
                (s) => s.weekday !== weekday.toLowerCase(),
              ),
            }
          : group,
      ),
    );
  };

  const handleAddTimePeriod = () => {
    const newTimePeriod = {
      time_period: {
        section_id: section?.id,
        start_time: "14:00",
        end_time: "14:45",
      },
      schedules: [],
    };

    setScheduleGroups((prev) => [...prev, newTimePeriod]);
  };

  const handleRemoveTimePeriod = (periodIdx) => {
    if (scheduleGroups.length > 1) {
      // Prevent removing all periods
      setScheduleGroups((prev) => prev.filter((_, idx) => idx !== periodIdx));
    }
  };

  // Function to create schedule
  const createSchedule = async (data) => {
    try {
      if (readOnly) return;
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/schedule_create/${tenantID}/${section.id}`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
          body: JSON.stringify(data),
        },
      );
      if (!response.ok) {
        throw new Error("Failed to create schedule");
      }
    } catch (error) {
      console.error("Error creating schedule:", error);
    }
  };

  if (scheduleGroups.length === 0 && readOnly) {
    let emptyMessage;
    if (teacherMode) {
      emptyMessage =
        "Ne predajete ni jedan predmet ili vaša odjeljenja nemaju raspored.";
    } else {
      emptyMessage = "Odjeljenje nema raspored časova.";
    }
    return (
      <>
        {!teacherMode && (
          <div className="flex justify-end">
            <BackButton
              onClick={() => setShowSchedulePage(false)}
              colorConfig={colorConfig}
            />
          </div>
        )}
        <div className="flex flex-col items-center justify-center min-h-[60vh] w-full">
          <EmptyState message={emptyMessage} />
        </div>
      </>
    );
  }

  return (
    <>
      {showConfirmModal && (
        <ConfirmModal
          title="Spremanje Rasporeda"
          colorConfig={colorConfig}
          onClose={() => setShowConfirmModal(false)}
          onConfirm={() => {
            createSchedule(scheduleGroups);
            setShowConfirmModal(false);
          }}
        >
          Da li ste sigurni da želite spremiti ovaj raspored?
        </ConfirmModal>
      )}
      <div className={`shadow-lg ${className}`}>
        <div className="flex justify-end gap-2 mb-4">
          {!teacherMode && showButtons && (
            <BackButton
              onClick={() => setShowSchedulePage(false)}
              colorConfig={colorConfig}
            />
          )}
          {showButtons && (
            <PDFButton
              colorConfig={colorConfig}
              handleDownloadPDF={downloadTeacherSchedule}
              pdfLoading={pdfLoading}
            />
          )}
          {!readOnly && (
            <Button
              onClick={() => setShowConfirmModal(true)}
              color="primary"
              colorConfig={colorConfig}
              icon={FaSave}
            >
              Spremi Raspored
            </Button>
          )}
        </div>
        {/* Classroom drag source */}
        {!readOnly && classrooms?.length > 0 && (
          <div className="grid grid-cols-4 gap-2 mb-4">
            {classrooms?.map((classroom, idx) => (
              <div
                key={idx}
                draggable
                onDragStart={(e) => handleClassroomDragStart(e, classroom)}
                className={`font-semibold text-center rounded py-2 px-3 cursor-grab active:cursor-grabbing select-none ${secondaryTextColor} border border-opacity-20 shadow-sm hover:shadow-md transition-shadow duration-200`}
                title="Povuci učionicu na predmet"
              >
                {classroom.name}
              </div>
            ))}
          </div>
        )}
        {/* Subject drag source */}
        {!readOnly && subjects?.length > 0 && (
          <div className="grid grid-cols-4 gap-2 mb-4">
            {subjects?.map((subj, idx) => (
              <div
                key={idx}
                draggable
                onDragStart={(e) => handleDragStart(e, subj)}
                className={`font-semibold text-center rounded py-2 px-3 cursor-grab active:cursor-grabbing select-none ${primaryTextColor} ${primaryBgColor} border border-opacity-20 shadow-sm hover:shadow-md transition-shadow duration-200`}
                title="Povuci predmet na raspored"
              >
                {subj.subject_name}
              </div>
            ))}
          </div>
        )}
        {!readOnly && (
          <Button
            color="ternary"
            colorConfig={colorConfig}
            className="ml-2"
            onClick={handleAddTimePeriod}
          >
            + Dodaj čas
          </Button>
        )}
        <Table ref={tableRef}>
          <TableHead>
            <TableRow mode="header">
              <TableHeader>#</TableHeader>
              <TableHeader>Vrijeme</TableHeader>
              {WEEKDAYS.map((weekday, idx) => (
                <TableHeader key={idx}>{weekday}</TableHeader>
              ))}
            </TableRow>
          </TableHead>
          <TableBody>
            {scheduleGroups.map((group, periodIdx) => (
              <TableRow key={periodIdx}>
                {/* Period number */}
                <TableCell className="font-semibold">
                  <div className="flex items-center gap-2">
                    {!readOnly && scheduleGroups.length > 1 && (
                      <button
                        onClick={() => handleRemoveTimePeriod(periodIdx)}
                        className="pdf-hide text-red-500 hover:text-red-700 transition-colors hover:cursor-pointer"
                        title="Ukloni čas"
                      >
                        <FaTimes className="w-3 h-3" />
                      </button>
                    )}
                    <span>{periodIdx + 1}.</span>
                  </div>
                </TableCell>
                {/* Time editable */}
                <TableCell>
                  <div className="flex items-center gap-1 w-16">
                    {!readOnly && editingTime === `${periodIdx}-start_time` ? (
                      <input
                        type="time"
                        value={group.time_period.start_time}
                        onChange={(e) =>
                          handleTimeChange(e, periodIdx, "start_time")
                        }
                        onBlur={handleTimeBlur}
                        onKeyDown={handleTimeKeyPress}
                        className="border border-gray-300 rounded px-1 py-0.5 text-sm w-20"
                        autoFocus
                      />
                    ) : (
                      <span
                        onDoubleClick={() =>
                          !readOnly &&
                          handleTimeDoubleClick(periodIdx, "start_time")
                        }
                        className={`${!readOnly ? "cursor-pointer hover:bg-gray-100" : ""} px-1 py-0.5 rounded`}
                        title="Dupli klik za uređivanje"
                      >
                        {group.time_period.start_time?.slice(0, 5)}
                      </span>
                    )}
                    <span>-</span>
                    {!readOnly && editingTime === `${periodIdx}-end_time` ? (
                      <input
                        type="time"
                        value={group.time_period.end_time}
                        onChange={(e) =>
                          handleTimeChange(e, periodIdx, "end_time")
                        }
                        onBlur={handleTimeBlur}
                        onKeyDown={handleTimeKeyPress}
                        className="border border-gray-300 rounded px-1 py-0.5 text-sm w-20"
                        autoFocus
                      />
                    ) : (
                      <span
                        onDoubleClick={() =>
                          !readOnly &&
                          handleTimeDoubleClick(periodIdx, "end_time")
                        }
                        className={`${!readOnly ? "cursor-pointer hover:bg-gray-100" : ""} px-1 py-0.5 rounded`}
                        title="Dupli klik za uređivanje"
                      >
                        {group.time_period.end_time?.slice(0, 5)}
                      </span>
                    )}
                  </div>
                </TableCell>
                {/* Weekday cells */}
                {WEEKDAYS.map((weekday, colIdx) => {
                  const subject = group.schedules.find(
                    (s) => s.weekday === weekday.toLowerCase(),
                  );
                  // If teacher mode use color config from subject
                  const primaryTextColorSubjects = teacherMode
                    ? getColor("primary", "text", subject?.color_config)
                    : primaryTextColor;
                  const secondaryTextColorSubjects = teacherMode
                    ? getColor("secondary", "text", subject?.color_config)
                    : secondaryTextColor;
                  const primaryBgColorSubjects = teacherMode
                    ? getColor("primaryComplement", "bg", subject?.color_config)
                    : primaryBgColor;

                  const cellKey = `${periodIdx}-${weekday}`;
                  const isDragOver = dragOverCell === cellKey;

                  return (
                    <TableCell
                      key={colIdx}
                      className={`min-h-[60px] ${
                        isDragOver
                          ? "bg-blue-100 border-2 border-blue-300 border-dashed"
                          : ""
                      }`}
                      onDragOver={!readOnly ? handleDragOver : undefined}
                      onDragEnter={
                        !readOnly
                          ? (e) => handleDragEnter(e, periodIdx, weekday)
                          : undefined
                      }
                      onDragLeave={!readOnly ? handleDragLeave : undefined}
                      onDrop={
                        !readOnly
                          ? (e) => handleDrop(e, periodIdx, weekday)
                          : undefined
                      }
                    >
                      {!subject ? (
                        <div className="w-24 text-gray-300 text-center py-2 border-2 border-dashed border-gray-200 rounded">
                          prazno
                        </div>
                      ) : (
                        <div
                          className={`inline-block rounded font-semibold text-center py-1 px-2 relative group ${primaryTextColorSubjects} ${primaryBgColorSubjects} border border-opacity-20 shadow-sm`}
                        >
                          {teacherMode && (
                            <div className="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-2 py-1 bg-gray-900 text-white text-xs rounded opacity-0 group-hover:opacity-100 transition-opacity duration-200 whitespace-nowrap z-10">
                              {subject.tenant_name} - {subject.section_name}
                            </div>
                          )}
                          {subject.subject_name}
                          {subject.classroom_code && (
                            <span className={`${secondaryTextColorSubjects}`}>
                              {" "}
                              ({subject.classroom_code})
                            </span>
                          )}
                          {!readOnly && (
                            <button
                              onClick={() =>
                                handleRemoveSubject(periodIdx, weekday)
                              }
                              className="pdf-hide absolute -top-2 -right-2 bg-red-500 hover:bg-red-600 text-white rounded-full w-5 h-5 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-all duration-200 shadow-sm hover:cursor-pointer"
                              title="Ukloni predmet"
                            >
                              <FaTimes className="w-3 h-3" />
                            </button>
                          )}
                        </div>
                      )}
                    </TableCell>
                  );
                })}
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>
    </>
  );
}
