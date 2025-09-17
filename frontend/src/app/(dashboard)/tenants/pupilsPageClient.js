"use client";
import AssignmentPage from "@/app/components/client_pages/AssignmentPage";
import { FaUserFriends } from "react-icons/fa";
import { FaTimes } from "react-icons/fa";
import CreateUpdateModal from "@/app/components/modal/CreateUpdateModal";
import { useState, useEffect } from "react";
import { fetchPupilBehaviourGrade } from "@/app/util/fetchPupilBehaviourGrades";
import { FaClipboardCheck } from "react-icons/fa";
import { fetchSectionPupilsUtil } from "@/app/util/fetchSectionPupils";
import ErrorModal from "@/app/components/modal/ErrorModal";
import { FaAward, FaUserMinus, FaUserTimes } from "react-icons/fa";
import { CertificatePageClient } from "@/app/components/client_pages/CertificatePage";
import { BehaviourGradeModal } from "@/app/components/modal/BehaviourGradeModal";
import ConfirmModal from "@/app/components/modal/ConfirmModal";

export default function PupilsPageClient({
  accessToken,
  onBack,
  tenant,
  section,
  mode,
  pendingMode,
  archived = 0,
}) {
  const [pupils, setPupils] = useState([]);
  const [pupilsForAssignment, setPupilsForAssignment] = useState([]);
  const [pendingPupils, setPendingPupils] = useState([]);
  const [selectedBehaviourGrades, setSelectedBehaviourGrades] = useState(null);
  const [selectedBehaviourGrade, setSelectedBehaviourGrade] = useState(null);
  const [errorMessage, setErrorMessage] = useState();
  const [selectedPupil, setSelectedPupil] = useState(null);
  const [showCertificatePage, setShowCertificatePage] = useState(false);
  const [showUnenrollModal, setShowUnenrollModal] = useState(false);
  const [showBehaviourGradeEdit, setShowBehaviourGradeEdit] = useState(false);

  const fetchSectionPupils = async (tenantId, sectionId) => {
    const result = await fetchSectionPupilsUtil({
      tenantId,
      sectionId,
      accessToken,
    });
    if (result.error) {
      setErrorMessage(result.error);
      setPupils([]);
      setPendingPupils([]);
      setPupilsForAssignment([]);
    } else {
      setPupils(result.pupils);
      setPendingPupils(result.pendingPupils);
      setPupilsForAssignment(result.pupilsForAssignment);
    }
  };

  useEffect(() => {
    fetchSectionPupils(tenant.id, section.id);
  }, []);

  const deletePupilInvite = async (data) => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/teacher/delete_pupil_invite`,
        {
          method: "DELETE",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
          body: JSON.stringify({
            invite_id: data.id,
            pupil_id: data.pupil_id,
            tenant_id: data.tenant_id,
          }),
        },
      );

      if (!response.ok) {
        const errorData = await response.text();
        console.error("Failed to delete invite:", errorData);
      }

      const pupil = await response.json();

      setPupilsForAssignment((prev) => [...prev, pupil]);
      setPendingPupils((prevPupils) =>
        prevPupils.filter((pupil) => pupil.pupil_id !== data.pupil_id),
      );
    } catch (e) {
      console.error(e);
    }
  };

  const config = {
    // API endpoints
    assignEndpoint: `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/teacher/send_section_invite/${tenant.id}/${section.id}`,
    unassignEndpoint: `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/teacher/unassign_pupil_from_section/${tenant.id}/${section.id}`,

    archived: archived,

    // UI configuration
    addButtonText: "Pozovi učenike",
    modalTitle: "Pozovi učenike u odjeljenje",
    searchPlaceholder: "Pretraži email učenika",
    modalNote: "Možete odabrati više učenika za odjeljenje.",

    // Field mappings
    keyField: "id", // Default value, but being explicit
    labelField: "email",
    titleField: "email",
    prefix: "učenika",
    keysToIgnore: ["unenrolled"],

    // Icon
    icon: FaUserFriends,
    getIcon: (item) => {
      return item && item.unenrolled == true ? FaUserTimes : FaUserFriends;
    },

    // Error messages
    assignErrorMessage:
      "Došlo je do greške prilikom pozivanja učenika odjeljenju.",
    unassignErrorMessage:
      "Došlo je do greške prilikom uklanjanja učenika iz odjeljenja.",
    pendingTitle: "Pozvani učenici",
    mainTitle: `${section.name} - Učenici`,
    itemsTitle: "Učenici u odjeljenju",
    showEdit: true,
    editButton: {
      onClick: async (item) => {
        setShowBehaviourGradeEdit(!item?.unenrolled);
        const result = await fetchPupilBehaviourGrade(
          tenant.id,
          section.id,
          item.id,
          accessToken,
        );
        setSelectedBehaviourGrades(result);
      },
      label: "Vladanje",
      icon: FaClipboardCheck,
    },
    archivedEditButton: {
      onClick: (item) => {
        setSelectedPupil(item);
        setShowCertificatePage(true);
      },
      label: "Svjedočanstvo",
      icon: FaAward,
    },
    extraButton: (item) =>
      item?.unenrolled
        ? null
        : {
            onClick: (it) => {
              setSelectedPupil(it);
              setShowUnenrollModal(true);
            },
            label: "Ispiši",
            icon: FaUserMinus,
          },
    getBgColor: (item) => {
      if (item?.unenrolled) {
        return "secondary";
      }
      return "primary";
    },
    getTextColor: (item) => {
      if (item?.unenrolled) {
        return "secondary";
      }
      return "primary";
    },
  };
  const pendingConfig = {
    pendingKeyField: "id",
    pendingTitleField: "pupil_full_name",
    pendingPrefix: "učenika",
    pendingKeysToIgnore: ["pupil_id", "section_id", "tenant_id"],
    deleteInviteButton: {
      label: "Izbriši",
      onClick: deletePupilInvite,
      icon: FaTimes,
    },
    getDeleteInviteMessage: (data) => {
      return (
        <>
          Da li ste sigurni da želite izbrisati poziv u odjeljenje{" "}
          <span className="font-bold">{data["section_name"]}</span> (
          <span className="font-bold">{data["tenant_name"]}</span>) za učenika{" "}
          <span className="font-bold">{data["pupil_full_name"]}</span>.
        </>
      );
    },
  };

  const updatePupilBehaviourGrade = async (behaviour_grade) => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/teacher/update_behaviour_grade/${tenant.id}`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
          body: JSON.stringify(behaviour_grade),
        },
      );

      if (response.ok) {
        const newBehaviourGrade = await response.json();
        setSelectedBehaviourGrades((prev) =>
          prev.map((grade) =>
            grade.id === newBehaviourGrade.id ? newBehaviourGrade : grade,
          ),
        );
      }
    } catch (e) {
      console.error(e);
    }
  };

  const unenrollPupilFromSection = async (pupil) => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/teacher/unenroll_pupil_from_section/${tenant.id}/${section.id}/${pupil.id}`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
        },
      );

      if (response.ok) {
        // Update pupil unenrolled to true
        setPupils((prev) =>
          prev.map((p) => (p.id === pupil.id ? { ...p, unenrolled: true } : p)),
        );

        setPupilsForAssignment((prev) => [...prev, pupil]);
      }
    } catch (e) {
      console.error(e);
    }
  };

  if (selectedPupil && showCertificatePage) {
    return (
      <CertificatePageClient
        tenantID={tenant.id}
        sectionID={section.id}
        pupilID={selectedPupil.id}
        onBack={() => {
          setShowCertificatePage(false);
          setSelectedPupil(null);
        }}
        colorConfig={tenant.color_config}
        accessToken={accessToken}
      />
    );
  }

  return (
    <>
      {errorMessage && (
        <ErrorModal
          onClose={() => setErrorMessage()}
          colorConfig={tenant.color_config}
        >
          {errorMessage}
        </ErrorModal>
      )}
      {selectedBehaviourGrades && (
        <BehaviourGradeModal
          data={selectedBehaviourGrades}
          onClose={() => {
            setSelectedBehaviourGrade(null);
            setSelectedBehaviourGrades(null);
            setShowBehaviourGradeEdit(false);
          }}
          colorConfig={tenant.color_config}
          onEditClick={(item) => setSelectedBehaviourGrade(item)}
          showEdit={showBehaviourGradeEdit}
          tenantID={tenant.id}
          accessToken={accessToken}
        />
      )}
      {selectedBehaviourGrade && (
        <CreateUpdateModal
          title="Uredi vladanje"
          initialValues={{ behaviour: selectedBehaviourGrade.behaviour }}
          colorConfig={tenant.color_config}
          onClose={() => {
            setSelectedBehaviourGrade(null);
          }}
          onSave={(item) => {
            const updatedBehaviourGrade = {
              ...selectedBehaviourGrade,
              behaviour: item.behaviour,
            };
            updatePupilBehaviourGrade(updatedBehaviourGrade);
            setSelectedBehaviourGrade(null);
          }}
          fields={[
            {
              label: "Vladanje",
              name: "behaviour",
              placeholder: "Odaberite vladanje",
              type: "select",
              options: [
                { value: "primjerno", label: "Primjerno" },
                { value: "vrlodobro", label: "Vrlo dobro" },
                { value: "dobro", label: "Dobro" },
                { value: "zadovoljavajuće", label: "Zadovoljavajuće" },
                { value: "loše", label: "Loše" },
              ],
            },
          ]}
        />
      )}
      {selectedPupil && showUnenrollModal && (
        <ConfirmModal
          title="Potvrda ispisa"
          onClose={() => {
            setSelectedPupil(null);
            setShowUnenrollModal(false);
          }}
          onConfirm={() => {
            unenrollPupilFromSection(selectedPupil);
            setSelectedPupil(null);
            setShowUnenrollModal(false);
          }}
          colorConfig={tenant.color_config}
        >
          Da li ste sigurni da želite ispisati učenika{" "}
          <span className="font-bold">
            {selectedPupil.name} {selectedPupil.last_name}
          </span>{" "}
          iz odjeljenja <span className="font-bold">{section.name}?</span>
        </ConfirmModal>
      )}
      <AssignmentPage
        assignedItems={pupils}
        setAssignedItems={setPupils}
        availableItems={pupilsForAssignment}
        setAvailableItems={setPupilsForAssignment}
        tenant={tenant}
        onBack={onBack}
        config={config}
        accessToken={accessToken}
        showOnlyOnSearch={true}
        inviteMode={true}
        pendingItems={pendingPupils}
        setPendingItems={setPendingPupils}
        pendingConfig={pendingConfig}
        colorConfig={tenant.color_config}
        mode={mode}
        pendingMode={pendingMode}
      />
    </>
  );
}
