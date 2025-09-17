"use client";
import { useState, useEffect } from "react";
import AddButton from "@/app/components/common/AddButton";
import CreateUpdateModal from "@/app/components/modal/CreateUpdateModal";
import ErrorModal from "../../components/modal/ErrorModal";
import DynamicCardParent from "../../components/dynamic_card/DynamicCardParent";
import {
  FaLayerGroup,
  FaUserFriends,
  FaCalendarDay,
  FaFileAlt,
  FaUserTimes,
  FaBook,
  FaArchive,
  FaRegListAlt,
} from "react-icons/fa";
import BackButton from "@/app/components/common/BackButton";
import SchedulePageClient from "./schedulePageClient";
import LessonPageClient from "./lessonPageClient";
import AbsencePageClient from "../pupil_home/absencePageClient";
import Title from "@/app/components/common/Title";
import GradebookPageClient from "@/app/components/client_pages/GradebookPageClient";
import PupilsPageClient from "./pupilsPageClient";
import {
  handleArchiveConfirmUtil,
  ArchiveConfirmModal,
} from "@/app/util/archive_util";
import { CompleteGradebookPage } from "@/app/components/client_pages/CompleteGradebookPage";

// Archived can be 0 or 1. 0 means show only non-archived sections.
// 1 means show only archied sections.
export default function SectionsPageClient({
  accessToken,
  onBack,
  tenant,
  archived = 0,
}) {
  const [showModal, setShowModal] = useState(false);
  const [errorMessage, setErrorMessage] = useState(null);
  const [showSchedulePage, setShowSchedulePage] = useState(false);
  const [selectedSection, setSelectedSection] = useState(null);
  const [showLessonPage, setShowLessonPage] = useState(false);
  const [showAbsencePage, setShowAbsencePage] = useState(false);
  const [showGradebookPage, setShowGradebookPage] = useState(false);
  const [showArchiveModal, setShowArchiveModal] = useState(false);
  const [sections, setSections] = useState([]);
  const [sectionsMetadata, setSectionsMetadata] = useState();
  const [pupilsPageSelected, setPupilsPageSelected] = useState(false);
  const [showCompleteGradebookPage, setShowCompleteGradebookPage] =
    useState(false);

  const fetchSectionMetadata = async () => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/section_creation_metadata/${tenant?.id}`,
        {
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        },
      );
      if (response.ok) {
        const metadata = await response.json();
        setSectionsMetadata(metadata);
      }
    } catch (e) {
      console.error(e);
    }
  };

  const fetchSections = async () => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/sections/${tenant?.id}/${archived}`,
        {
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        },
      );
      if (response.ok) {
        const sections = await response.json();
        setSections(sections);
      }
    } catch (e) {
      console.error(e);
    }
  };

  useEffect(() => {
    if (archived == 0) {
      fetchSectionMetadata();
    }
    fetchSections();
    setSelectedSection(null);
    setPupilsPageSelected(false);
    setShowSchedulePage(false);
    setShowLessonPage(false);
    setShowLessonPage(false);
    setShowGradebookPage(false);
  }, [archived]);

  const sectionFields =
    sectionsMetadata?.curriculums?.length > 0
      ? [
          {
            label: "Kod odjeljenja",
            name: "section_code",
            placeholder: "Unesite kod odjeljenja (npr. B ili 2)",
          },
          {
            label: "Razred",
            name: "class_code",
            placeholder: "Odaberite razred",
            type: "select",
            options: sectionsMetadata?.classes?.map((cls) => ({
              value: cls.class_code,
              label: cls.class_code,
            })),
          },
          {
            label: "Godina",
            name: "year",
            placeholder: "Unesite godinu (npr. 2025/2026)",
          },
          {
            label: "Nastavni plan i program",
            name: "curriculum_code",
            placeholder: "Odaberite NPP",
            type: "select",
            options: sectionsMetadata?.curriculums?.map((curriculum) => ({
              value: curriculum.curriculum_code,
              label: curriculum.curriculum_name,
            })),
          },
        ]
      : [
          {
            label: (
              <span>
                Dodavanje odjeljenja nije moguće dok ne definišete barem jedan
                kurikulum za instituciju. Molimo vas da prvo dodate odgovarajuće
                kurikulume.
              </span>
            ),
            type: "subtitle",
          },
        ];

  const createSection = async (data) => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/section_create/${tenant.id}`,
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
      }
      const newSection = await response.json();
      setSections((prev) => [...(prev || []), newSection]);
      setShowModal(false);
    } catch (error) {
      console.error("Error creating section:", error);
      setErrorMessage("Failed to create section");
    }
  };

  const handleDeleteConfirm = async (sectionId, closeDeleteModal) => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/section/${tenant.id}/${sectionId}`,
        {
          method: "DELETE",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
        },
      );
      if (response.ok) {
        setSections((prev) =>
          prev.filter((section) => section.id !== sectionId),
        );
      } else {
        const errorData = await response.json();
        setErrorMessage(errorData.message || "Failed to delete section");
      }
      closeDeleteModal();
    } catch (error) {
      console.error("Error deleting section:", error);
      setErrorMessage("Failed to delete section");
    }
  };

  const handleArchiveConfirm = async (sectionId, tenantId) => {
    const errorMsg = await handleArchiveConfirmUtil(
      tenantId,
      sectionId,
      accessToken,
    );
    if (errorMsg) {
      setErrorMessage(errorMsg || "Failed to archive section");
    } else {
      setSections((prev) => prev.filter((section) => section.id !== sectionId));
    }
  };

  if (pupilsPageSelected && selectedSection) {
    return (
      <PupilsPageClient
        accessToken={accessToken}
        onBack={() => {
          setSelectedSection(null);
          setPupilsPageSelected(false);
        }}
        tenant={tenant}
        section={selectedSection}
        mode={tenant.pupil_display}
        pendingMode={tenant.pupil_invite_display}
        archived={archived}
      />
    );
  }

  if (showSchedulePage && selectedSection) {
    return (
      <SchedulePageClient
        colorConfig={tenant.color_config}
        setShowSchedulePage={setShowSchedulePage}
        section={selectedSection}
        tenantID={tenant.id}
        accessToken={accessToken}
        archived={archived}
      />
    );
  }

  if (showLessonPage && selectedSection) {
    return (
      <LessonPageClient
        section={selectedSection}
        tenantID={tenant.id}
        accessToken={accessToken}
        colorConfig={tenant.color_config}
        setShowLessonPage={setShowLessonPage}
        mode={tenant.lesson_display}
        archived={archived}
      />
    );
  }

  if (showAbsencePage && selectedSection) {
    return (
      <AbsencePageClient
        section={selectedSection}
        tenantID={tenant.id}
        accessToken={accessToken}
        colorConfig={tenant.color_config}
        setShowAbsencePage={setShowAbsencePage}
        mode="teacher"
        displayMode={tenant.absence_display}
        archived={archived}
      />
    );
  }

  if (showGradebookPage && selectedSection) {
    return (
      <GradebookPageClient
        colorConfig={tenant.color_config}
        setShowGradebookPage={setShowGradebookPage}
        accessToken={accessToken}
        section={selectedSection}
        archived={archived}
      />
    );
  }

  if (showCompleteGradebookPage && selectedSection) {
    return (
      <CompleteGradebookPage
        section={selectedSection}
        colorConfig={tenant.color_config}
        accessToken={accessToken}
        onBack={() => {
          setShowCompleteGradebookPage(false);
          setSelectedSection(null);
        }}
      />
    );
  }

  return (
    <>
      <ArchiveConfirmModal
        isOpen={showArchiveModal}
        selectedSection={selectedSection}
        onClose={() => {
          setSelectedSection(null);
          setShowArchiveModal(false);
        }}
        onConfirm={() => {
          handleArchiveConfirm(selectedSection.id, tenant.id);
          setSelectedSection(null);
          setShowArchiveModal(false);
        }}
        colorConfig={tenant.color_config}
      />
      <Title
        icon={archived == 0 ? FaLayerGroup : FaArchive}
        colorConfig={tenant.color_config}
      >
        {archived == 0 ? "Odjeljenja" : "Arhivirana odjeljenja"}
      </Title>
      <div className="flex justify-end mb-4">
        {onBack && (
          <div className="mr-2">
            <BackButton onClick={onBack} colorConfig={tenant.color_config} />
          </div>
        )}
        {archived == 0 && (
          <AddButton
            onClick={() => setShowModal(true)}
            colorConfig={tenant.color_config}
          >
            Dodaj odjeljenje
          </AddButton>
        )}
        {errorMessage && (
          <ErrorModal
            onClose={() => setErrorMessage(null)}
            colorConfig={tenant.color_config}
          >
            {errorMessage}
          </ErrorModal>
        )}
        {showModal && archived == 0 && (
          <CreateUpdateModal
            title="Dodaj odjeljenje"
            fields={sectionFields}
            onClose={() => setShowModal(false)}
            onSave={createSection}
            colorConfig={tenant.color_config}
            showSave={sectionsMetadata?.curriculums?.length > 0}
            internalFormOverflow={sectionsMetadata?.curriculums?.length > 0}
          />
        )}
      </div>
      <DynamicCardParent
        data={sections}
        setData={setSections}
        icon={archived == 0 ? <FaLayerGroup /> : <FaArchive />}
        titleField="name"
        prefix="odjeljenje"
        onDeleteConfirm={handleDeleteConfirm}
        editFields={sectionFields.filter(
          (field) => field.name !== "curriculum_code",
        )}
        editUrl={`${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/section`}
        showEdit={true}
        showDelete={archived == 0}
        accessToken={accessToken}
        keysToIgnore={[
          "section_code",
          "class_code",
          "semester_code",
          "tenant_id",
          "homeroom_teacher_id",
          "curriculum_code",
        ]}
        extraButton={{
          label: "Učenici",
          onClick: (section) => {
            setSelectedSection(section);
            setPupilsPageSelected(true);
          },
          icon: FaUserFriends,
        }}
        tenantColorConfig={tenant.color_config}
        mode={tenant.section_display}
        extraActions={[
          {
            label: "Raspored",
            icon: FaCalendarDay,
            onClick: (section) => {
              setShowSchedulePage(true);
              setSelectedSection(section);
            },
          },
          {
            label: "Časovi",
            icon: FaFileAlt,
            onClick: (section) => {
              setShowLessonPage(true);
              setSelectedSection(section);
            },
          },
          {
            label: "Izostanci",
            icon: FaUserTimes,
            onClick: (section) => {
              setShowAbsencePage(true);
              setSelectedSection(section);
            },
          },
          {
            label: "Ocjene",
            icon: FaBook,
            onClick: (section) => {
              setShowGradebookPage(true);
              setSelectedSection(section);
            },
          },
          {
            label: "Arhiviraj",
            icon: FaArchive,
            onClick: (section) => {
              setSelectedSection(section);
              setShowArchiveModal(true);
            },
            show: () => archived == 0,
          },
          {
            label: "Dnevnik",
            onClick: (section) => {
              setShowCompleteGradebookPage(true);
              setSelectedSection(section);
            },
            icon: FaRegListAlt,
          },
        ]}
      />
    </>
  );
}
