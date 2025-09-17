"use client";
import { useState } from "react";
import DynamicCardParent from "../../components/dynamic_card/DynamicCardParent";
import {
  FaLayerGroup,
  FaUsers,
  FaChalkboard,
  FaCalendarDay,
  FaUserTimes,
  FaUserFriends,
  FaBook,
  FaClipboardCheck,
  FaArchive,
  FaAward,
  FaRegListAlt,
} from "react-icons/fa";
import Title from "@/app/components/common/Title";
import SchedulePageClient from "../tenants/schedulePageClient";
import LessonPageClient from "../tenants/lessonPageClient";
import AbsencePageClient from "./absencePageClient";
import { fetchSectionPupilsUtil } from "@/app/util/fetchSectionPupils";
import ErrorMessage from "@/app/components/custom_messages/ErrorMessage";
import PupilsPageClient from "../tenants/pupilsPageClient";
import GradebookPageClient from "@/app/components/client_pages/GradebookPageClient";
import { fetchPupilBehaviourGradeNoPupilID } from "@/app/util/fetchPupilBehaviourGrades";
import {
  handleArchiveConfirmUtil,
  ArchiveConfirmModal,
} from "@/app/util/archive_util";
import ErrorModal from "@/app/components/modal/ErrorModal";
import { CertificatePageClient } from "@/app/components/client_pages/CertificatePage";
import { BehaviourGradeModal } from "@/app/components/modal/BehaviourGradeModal";
import { CompleteGradebookPage } from "@/app/components/client_pages/CompleteGradebookPage";

export default function UserHomePageClient({
  initialSections = [],
  accessToken,
  mode = "pupil",
  user,
  archived = 0,
}) {
  const [sections, setSections] = useState(initialSections);
  const [showSchedulePage, setShowSchedulePage] = useState(false);
  const [selectedSection, setSelectedSection] = useState(null);
  const [showLessonPage, setShowLessonPage] = useState(false);
  const [showAbsencePage, setShowAbsencePage] = useState(false);
  const [sectionPupils, setSectionPupils] = useState([]);
  const [pupilsPageSelected, setPupilsPageSelected] = useState(false);
  const [pendingPupils, setPendingPupils] = useState([]);
  const [pupilsForAssignment, setPupilsForAssignment] = useState([]);
  const [errorMessage, setErrorMessage] = useState(null);
  const [selectedTenant, setSelectedTenant] = useState(null);
  const [showSectionAbsencePage, setShowSectionAbsencePage] = useState(false);
  const [showGradebookPage, setShowGradebookPage] = useState(false);
  const [selectedBehaviourGrades, setSelectedBehaviourGrades] = useState(null);
  const [showArchiveModal, setShowArchiveModal] = useState(false);
  const [archiveErrorMessage, setArchiveErrorMessage] = useState(null);
  const [showCertificatePage, setShowCertificatePage] = useState(false);
  const [showCompleteGradebookPage, setShowCompleteGradebookPage] =
    useState(false);

  let emptyMessage;
  if (mode == "pupil") {
    if (archived == 1) {
      emptyMessage = "Trenutno nema arhiviranih odjeljenja koje ste pohađali.";
    } else {
      emptyMessage =
        "Niste ni u jednom odjeljenju. Zatražite poziv od nastavnika i provjerite pozive.";
    }
  } else if (mode == "teacher") {
    if (archived == 1) {
      emptyMessage = "Nema arhiviranih odjeljenja u kojima ste držali nastavu.";
    } else {
      emptyMessage =
        "Ne predajete ni jednom odjeljenju. Zatražite poziv od administratora i provjerite pozive.";
    }
  }

  const fetchSectionPupils = async (tenantId, sectionId) => {
    const result = await fetchSectionPupilsUtil({
      tenantId,
      sectionId,
      accessToken,
      sections,
    });
    if (result.error) {
      setErrorMessage(result.error);
    } else {
      setSectionPupils(result.pupils);
      setPendingPupils(result.pendingPupils);
      setPupilsForAssignment(result.pupilsForAssignment);
      setPupilsPageSelected(true);
    }
  };

  const fetchTenantByID = async (tenantId) => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/tenant/${tenantId}`,
        {
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        },
      );
      if (!response.ok) {
        throw new Error("Failed to fetch tenant");
      }
      const data = await response.json();
      setSelectedTenant(data);
    } catch (error) {
      setErrorMessage(error.message);
    }
  };

  const handleArchiveConfirm = async (sectionId, tenantId) => {
    const errorMsg = await handleArchiveConfirmUtil(
      tenantId,
      sectionId,
      accessToken,
    );
    if (errorMsg) {
      setArchiveErrorMessage(errorMsg || "Failed to archive section");
    } else {
      setSelectedSection(null);
      setSections((prev) => prev.filter((section) => section.id !== sectionId));
    }
  };

  if (pupilsPageSelected && selectedSection && selectedTenant) {
    return (
      <PupilsPageClient
        pupils={sectionPupils}
        setPupils={setSectionPupils}
        pupilsForAssignment={pupilsForAssignment}
        setPupilsForAssignment={setPupilsForAssignment}
        accessToken={accessToken}
        onBack={() => {
          setPupilsPageSelected(false);
          setSelectedSection(null);
        }}
        tenant={selectedTenant}
        section={selectedSection}
        pendingPupils={pendingPupils}
        setPendingPupils={setPendingPupils}
        mode={selectedSection.pupil_display}
        pendingMode={selectedSection.pupil_invite_display}
        archived={archived}
      />
    );
  }

  if (showSchedulePage && selectedSection) {
    return (
      <SchedulePageClient
        colorConfig={selectedSection?.color_config}
        setShowSchedulePage={setShowSchedulePage}
        section={selectedSection}
        tenantID={selectedSection?.tenant_id}
        accessToken={accessToken}
        readOnly={true}
        archived={archived}
      />
    );
  }

  if (showLessonPage && selectedSection) {
    return (
      <LessonPageClient
        section={selectedSection}
        tenantID={selectedSection?.tenant_id}
        accessToken={accessToken}
        colorConfig={selectedSection?.color_config}
        setShowLessonPage={setShowLessonPage}
        mode={selectedSection?.lesson_display}
        archived={archived}
      />
    );
  }

  if (showAbsencePage && selectedSection) {
    return (
      <AbsencePageClient
        section={selectedSection}
        tenantID={selectedSection?.tenant_id}
        pupilID={user?.id}
        accessToken={accessToken}
        colorConfig={selectedSection?.color_config}
        setShowAbsencePage={setShowAbsencePage}
        displayMode={selectedSection?.absence_display}
        archived={archived}
      />
    );
  }

  if (showSectionAbsencePage && selectedSection) {
    return (
      <AbsencePageClient
        section={selectedSection}
        tenantID={selectedSection?.tenant_id}
        accessToken={accessToken}
        colorConfig={selectedSection?.color_config}
        setShowAbsencePage={setShowSectionAbsencePage}
        mode="teacher"
        displayMode={selectedSection?.absence_display}
        archived={archived}
      />
    );
  }

  if (showGradebookPage && selectedSection) {
    return (
      <GradebookPageClient
        colorConfig={selectedSection?.color_config}
        setShowGradebookPage={setShowGradebookPage}
        accessToken={accessToken}
        section={selectedSection}
        mode={mode}
        archived={archived}
      />
    );
  }

  if (showCertificatePage && selectedSection) {
    return (
      <CertificatePageClient
        tenantID={selectedSection.tenant_id}
        sectionID={selectedSection.id}
        pupilID={user.id}
        onBack={() => {
          setShowCertificatePage(false);
          setSelectedSection(null);
        }}
        colorConfig={selectedSection.color_config}
        accessToken={accessToken}
      />
    );
  }

  if (showCompleteGradebookPage && selectedSection) {
    return (
      <CompleteGradebookPage
        colorConfig={selectedSection?.color_config}
        section={selectedSection}
        onBack={() => {
          setShowCompleteGradebookPage(false);
          setSelectedSection(null);
        }}
        accessToken={accessToken}
      />
    );
  }

  return (
    <>
      {archiveErrorMessage && (
        <ErrorModal
          onClose={() => {
            setArchiveErrorMessage(null);
            setSelectedSection(null);
          }}
          colorConfig={selectedSection?.color_config}
        >
          {archiveErrorMessage}
        </ErrorModal>
      )}
      <ArchiveConfirmModal
        isOpen={showArchiveModal}
        selectedSection={selectedSection}
        onClose={() => {
          setSelectedSection(null);
          setShowArchiveModal(false);
        }}
        onConfirm={() => {
          handleArchiveConfirm(selectedSection.id, selectedSection.tenant_id);
          setShowArchiveModal(false);
        }}
        colorConfig={selectedSection?.color_config}
      />
      {selectedBehaviourGrades && mode === "pupil" && selectedSection && (
        <BehaviourGradeModal
          data={selectedBehaviourGrades}
          onClose={() => {
            setSelectedBehaviourGrades(null);
          }}
          colorConfig={selectedSection?.color_config}
          tenantID={selectedSection?.tenant_id}
          accessToken={accessToken}
        />
      )}
      {errorMessage && <ErrorMessage message={errorMessage} />}
      <Title icon={archived == 0 ? FaUsers : FaArchive} outerClassName="mb-8">
        {archived == 0 ? "Moja odjeljenja" : "Arhivirana odjeljenja"}
      </Title>
      <DynamicCardParent
        data={sections}
        setData={setSections}
        icon={archived == 0 ? <FaLayerGroup /> : <FaArchive />}
        titleField="name"
        prefix="odjeljenje"
        accessToken={accessToken}
        keysToIgnore={[
          "section_code",
          "class_code",
          "semester_code",
          "tenant_id",
          "homeroom_teacher_id",
          "curriculum_code",
          "color_config",
          "pupil_display",
          "pupil_invite_display",
          "lesson_display",
          "absence_display",
        ]}
        tenantColorConfig={(item) => item.color_config}
        emptyMessage={emptyMessage}
        extraButton={{
          label: "Raspored časova",
          icon: FaChalkboard,
          onClick: (data) => {
            setShowSchedulePage(true);
            setSelectedSection(data);
          },
        }}
        showEdit={true}
        editButton={{
          onClick: (section) => {
            if (mode === "teacher") {
              setSelectedSection(section);
              setShowLessonPage(true);
            } else {
              setSelectedSection(section);
              setShowAbsencePage(true);
            }
          },
          label: mode === "teacher" ? "Časovi" : "Izostanci",
          icon: mode === "teacher" ? FaCalendarDay : FaUserTimes,
        }}
        extraActions={[
          {
            label: "Učenici",
            icon: FaUserFriends,
            onClick: (data) => {
              setSelectedSection(data);
              fetchTenantByID(data.tenant_id);
              fetchSectionPupils(data.tenant_id, data.id);
            },
            show: (data) => {
              if (data.homeroom_teacher_id === user?.id && mode == "teacher") {
                return true;
              } else {
                return false;
              }
            },
          },
          {
            label: "Izostanci",
            icon: FaUserTimes,
            onClick: (data) => {
              setSelectedSection(data);
              setShowSectionAbsencePage(true);
            },
            show: (data) => {
              if (data.homeroom_teacher_id === user?.id && mode == "teacher") {
                return true;
              } else {
                return false;
              }
            },
          },
          {
            label: "Vladanje",
            icon: FaClipboardCheck,
            onClick: async (data) => {
              const behaviourGrades = await fetchPupilBehaviourGradeNoPupilID(
                data.tenant_id,
                data.id,
                accessToken,
              );
              setSelectedBehaviourGrades(behaviourGrades);
              setSelectedSection(data);
            },
            show: () => mode === "pupil",
          },
          {
            label: "Svjedočanstvo",
            icon: FaAward,
            onClick: (data) => {
              setSelectedSection(data);
              setShowCertificatePage(true);
            },
            show: () => mode === "pupil" && archived == 1,
          },
          {
            label: "Arhiviraj",
            icon: FaArchive,
            onClick: (data) => {
              setSelectedSection(data);
              setShowArchiveModal(true);
            },
            show: (data) => {
              if (
                data.homeroom_teacher_id === user?.id &&
                archived == 0 &&
                mode == "teacher"
              ) {
                return true;
              } else {
                return false;
              }
            },
          },
          {
            label: "Dnevnik",
            icon: FaRegListAlt,
            onClick: (data) => {
              setSelectedSection(data);
              setShowCompleteGradebookPage(true);
            },
            show: (data) => {
              if (data.homeroom_teacher_id === user?.id && mode == "teacher") {
                return true;
              } else {
                return false;
              }
            },
          },
        ]}
        showDelete={true}
        deleteButton={{
          onClick: (section) => {
            setSelectedSection(section);
            setShowGradebookPage(true);
          },
          label: "Ocjene",
          icon: FaBook,
        }}
      />
    </>
  );
}
