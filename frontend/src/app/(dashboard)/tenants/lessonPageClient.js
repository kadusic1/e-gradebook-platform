"use client";
import Title from "@/app/components/common/Title";
import { useState, useEffect, useMemo } from "react";
import {
  FaFileAlt,
  FaUsers,
  FaCalendarDay,
  FaInfoCircle,
} from "react-icons/fa";
import BackButton from "@/app/components/common/BackButton";
import DynamicCardParent from "@/app/components/dynamic_card/DynamicCardParent";
import CreateUpdateModal from "@/app/components/modal/CreateUpdateModal";
import AddButton from "@/app/components/common/AddButton";
import { formatDateToDDMMYYYY } from "@/app/util/date_util";
import ConfirmModal from "@/app/components/modal/ConfirmModal";
import { getScheduleForSection } from "@/app/api/helper/scheduleApi";

export default function LessonPageClient({
  section,
  tenantID,
  mode,
  accessToken,
  colorConfig,
  setShowLessonPage,
  archived = 0,
}) {
  const [lessons, setLessons] = useState([]);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [itemToUpdate, setItemToUpdate] = useState(null);
  const [pupils, setPupils] = useState([]);
  const [subjects, setSubjects] = useState([]);
  const [lessonToDelete, setLessonToDelete] = useState(null);
  const [scheduleGroups, setScheduleGroups] = useState(null);

  useEffect(() => {
    getScheduleForSection(section.id, tenantID, accessToken, setScheduleGroups);
  }, [section, tenantID, accessToken]);

  const lessonDisplay = useMemo(() => {
    return (
      lessons?.map((item) => {
        return {
          id: item.lesson_data.id,
          description: item.lesson_data.description,
          date: formatDateToDDMMYYYY(item.lesson_data.date),
          period_number: item.lesson_data.period_number,
          subject_code: item.lesson_data.subject_code,
          subject_name: item.lesson_data.subject_name,
          present_pupil_count: item?.pupil_attendance_data?.filter(
            (attendance) => attendance?.status === "present",
          ).length,
          absent_pupil_count: item?.pupil_attendance_data?.filter(
            (attendance) =>
              attendance?.status === "absent" ||
              attendance?.status === "unexcused" ||
              attendance?.status === "excused",
          ).length,
          lesson_posted_by_teacher: item.lesson_data.lesson_posted_by_teacher,
        };
      }) || []
    );
  }, [lessons]);

  const fetchClasses = async () => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/teacher/get_lessons_for_section/${tenantID}/${section.id}`,
        {
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        },
      );
      if (response.ok) {
        const data = await response.json();
        setLessons(data?.lessons || []);
        setPupils(data?.pupils || []);
        setSubjects(data?.subjects || []);
      }
    } catch (error) {
      console.error("Error fetching tenant children data:", error);
    }
  };

  useEffect(() => {
    fetchClasses();
  }, []);

  let lessonFields;

  if (
    (!pupils?.length || pupils?.length == 0) &&
    (!scheduleGroups?.length || scheduleGroups?.length == 0)
  ) {
    lessonFields = [
      {
        label:
          "Ne možete dodati čas. Odjeljenje nema učenika ni raspored časova.",
        type: "subtitle",
      },
    ];
  } else if (!pupils?.length || pupils?.length == 0) {
    lessonFields = [
      {
        label: "Ne možete dodati čas. Odjeljenje nema učenika.",
        type: "subtitle",
      },
    ];
  } else if (!scheduleGroups?.length || scheduleGroups?.length == 0) {
    lessonFields = [
      {
        label: "Ne možete dodati čas. Odjeljenje nema raspored časova.",
        type: "subtitle",
      },
    ];
  } else {
    lessonFields = [
      {
        label: "Detalji časa",
        type: "subtitle",
        icon: FaInfoCircle,
      },
      {
        label: "Predmet",
        name: "lesson_data.subject_code",
        placeholder: "Odaberite predmet",
        type: "select",
        options: subjects?.map((subject) => ({
          value: subject.subject_code,
          label: `${subject.subject_name}`,
        })),
      },
      {
        label: "Tema časa",
        name: "lesson_data.description",
        placeholder:
          "Unesite temu časa (npr. sabiranje i oduzimanje razlomaka)",
        type: "text-area",
      },
      {
        label: "Datum časa",
        name: "lesson_data.date",
        placeholder: "Unesite datum časa",
        type: "date",
      },
      {
        label: "Redni broj časa",
        name: "lesson_data.period_number",
        placeholder:
          "Npr. ako imate 2 časa fizike na isti dan, unesite 1 ili 2",
        type: "positive-number",
      },
      {
        label: "Izostanci",
        type: "subtitle",
        icon: FaUsers,
      },
      {
        label: "Prisustvo učenika",
        name: "pupil_attendance_data",
        type: "attendance",
        items: pupils,
      },
    ];
  }

  const createLesson = async (data) => {
    try {
      data.lesson_data.section_id = section.id;
      data.lesson_data.period_number = parseInt(data.lesson_data.period_number);
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/teacher/create_lesson/${tenantID}`,
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
        const resp = await response.text();
        setErrorMessage(resp);
        return;
      } else {
        const new_lesson = await response.json();
        setLessons((prev) => [...(prev || []), new_lesson]);
      }
    } catch (error) {
      console.error("Error creating lesson:", error);
      setErrorMessage("Failed to create lesson");
    }
  };

  const updateLesson = async (data) => {
    try {
      data.lesson_data.section_id = section.id;
      data.lesson_data.period_number = parseInt(data.lesson_data.period_number);
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/teacher/update_lesson/${tenantID}/${data.lesson_data.id}`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
          body: JSON.stringify(data),
        },
      );
      if (!response.ok) {
        const resp = await response.text();
        setErrorMessage(resp);
        return;
      } else {
        const updated_lesson = await response.json();
        setLessons((prev) =>
          prev.map((lesson) =>
            lesson.lesson_data.id === updated_lesson.lesson_data.id
              ? updated_lesson
              : lesson,
          ),
        );
      }
    } catch (error) {
      console.error("Error updating lesson:", error);
      setErrorMessage("Failed to update lesson");
    }
  };

  const deleteLesson = async (data) => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/teacher/delete_lesson/${tenantID}/${data.id}`,
        {
          method: "DELETE",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
        },
      );
      if (!response.ok) {
        const resp = await response.text();
        setErrorMessage(resp);
        return;
      } else {
        setLessons((prev) =>
          prev.filter((lesson) => lesson.lesson_data.id !== data.id),
        );
      }
    } catch (error) {
      console.error("Error deleting lesson:", error);
      setErrorMessage("Failed to delete lesson");
    }
  };

  return (
    <>
      <Title icon={FaFileAlt} colorConfig={colorConfig}>
        Časovi - {section.name}
      </Title>
      {showCreateModal && (
        <CreateUpdateModal
          title="Dodaj čas"
          onSave={(data) => {
            createLesson(data);
          }}
          colorConfig={colorConfig}
          fields={lessonFields}
          onClose={() => setShowCreateModal(false)}
          internalFormOverflow={
            pupils?.length > 0 && scheduleGroups?.length > 0 ? true : false
          }
          showSave={
            pupils?.length > 0 && scheduleGroups?.length > 0 ? true : false
          }
        />
      )}
      {itemToUpdate && (
        <CreateUpdateModal
          title="Uredi čas"
          onSave={(data) => {
            updateLesson(data);
          }}
          colorConfig={colorConfig}
          fields={lessonFields}
          onClose={() => setItemToUpdate(null)}
          initialValues={itemToUpdate}
          internalFormOverflow={
            pupils?.length > 0 && scheduleGroups?.length > 0 ? true : false
          }
          showSave={
            pupils?.length > 0 && scheduleGroups?.length > 0 ? true : false
          }
        />
      )}
      {lessonToDelete && (
        <ConfirmModal
          onConfirm={() => {
            deleteLesson(lessonToDelete);
            setLessonToDelete(null);
          }}
          onClose={() => setLessonToDelete(null)}
          colorConfig={colorConfig}
        >
          Da li ste sigurni da želite izbrisati{" "}
          {lessonToDelete.period_number + ". "}
          čas{" "}
          <span className="font-bold">
            {lessonToDelete.subject_name} (
            {formatDateToDDMMYYYY(lessonToDelete.date)})
          </span>
          ?
        </ConfirmModal>
      )}
      <div className="flex justify-end mb-4">
        <div className="mr-2">
          <BackButton
            colorConfig={colorConfig}
            onClick={() => setShowLessonPage(false)}
          />
        </div>
        {archived == 0 && (
          <AddButton
            onClick={() => setShowCreateModal(true)}
            colorConfig={colorConfig}
          >
            Dodaj čas
          </AddButton>
        )}
      </div>
      <DynamicCardParent
        data={lessonDisplay}
        icon={FaCalendarDay}
        titleField={["subject_name", "date"]}
        prefix="čas"
        showEdit={archived == 0}
        showDelete={archived == 0}
        accessToken={accessToken}
        tenantColorConfig={colorConfig}
        editFields={lessonFields}
        keysToIgnore={["subject_code", "id"]}
        editButton={{
          onClick: (item) => {
            const lessonToUpdate = lessons.find(
              (lesson) => lesson.lesson_data.id === item.id,
            );
            setItemToUpdate(lessonToUpdate);
          },
        }}
        deleteButton={{
          onClick: (item) => {
            const lessonToDelete = lessons.find(
              (lesson) => lesson.lesson_data.id === item.id,
            );
            setLessonToDelete(lessonToDelete.lesson_data);
          },
        }}
        mode={mode}
      />
    </>
  );
}
