"use client";
import { useState, useEffect, useMemo } from "react";
import ErrorModal from "../../components/modal/ErrorModal";
import DynamicCardParent from "../../components/dynamic_card/DynamicCardParent";
import { FaEnvelope, FaTimes, FaCheck, FaBook } from "react-icons/fa";
import BackButton from "@/app/components/common/BackButton";
import Modal from "@/app/components/modal/Modal";
import Title from "@/app/components/common/Title";
import ConfirmModal from "@/app/components/modal/ConfirmModal";
import { formatDateToDDMMYYYY } from "@/app/util/date_util";

export default function TeacherInvitesPageClient({
  accessToken,
  onBack,
  tenant,
  teacherAccountMode = false,
  initialInvites = [],
  setTeacherInviteData,
  colorConfig,
  mode = "table",
  showHeader = true,
}) {
  const [errorMessage, setErrorMessage] = useState(null);
  const [invites, setInvites] = useState(
    teacherAccountMode ? initialInvites : [],
  );
  const [subjects, setSubjects] = useState([]);
  const [inviteToDelete, setInviteToDelete] = useState(null);
  const [showSubjects, setShowSubjects] = useState(false);

  const dateFormattedInvites = useMemo(() => {
    return invites?.map((invite) => ({
      ...invite,
      invite_date: formatDateToDDMMYYYY(invite.invite_date),
    }));
  }, [invites]);

  // Fetch invites for admin mode
  const getTeacherInvites = async (tenantID) => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/teacher_invites/${tenantID}`,
        {
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
        },
      );

      if (!response.ok) {
        throw new Error("Failed to get teacher invites");
      }
      const data = await response.json();
      setInvites(data);
    } catch (error) {
      console.error("Failed to get teacher invite data");
      setErrorMessage("Greška prilikom dohvatanja poziva.");
    }
  };

  // Handler for accepting/declining invites (teacher mode only)
  const handleAction = async (invite, action) => {
    try {
      const res = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/teacher/handle_invite`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
          body: JSON.stringify({
            invite_id: invite.id,
            action: action,
            tenant_id: invite.tenant_id,
          }),
        },
      );
      if (!res.ok) throw new Error(await res.text());

      // Update the status of the handled invite
      const newStatus = action === "accept" ? "accepted" : "declined";
      setInvites((prev) =>
        prev.map((i) =>
          i.id === invite.id && i.tenant_id === invite.tenant_id
            ? { ...i, status: newStatus }
            : i,
        ),
      );
    } catch (err) {
      setErrorMessage(err.message || "Greška prilikom prihvatanja poziva.");
    }
  };

  const handleAccept = (invite) => handleAction(invite, "accept");
  const handleDecline = (invite) => handleAction(invite, "decline");

  // Fetch invites on mount for admin mode
  useEffect(() => {
    if (!teacherAccountMode) {
      getTeacherInvites(tenant.id);
    }
  }, []);

  // Configure action buttons for teacher mode
  const getActionButtons = () => {
    if (!teacherAccountMode) return {};

    return {
      showDelete: (item) => item.status === "pending",
      deleteButton: {
        onClick: handleDecline,
        label: "Odbij",
        icon: FaTimes,
      },
      showEdit: (item) => item.status === "pending",
      editButton: {
        onClick: handleAccept,
        label: "Prihvati",
        icon: FaCheck,
      },
    };
  };

  const deleteTeacherInvite = async (data) => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/delete_teacher_invite`,
        {
          method: "DELETE",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
          body: JSON.stringify({
            invite_id: data.id,
            teacher_id: data.teacher_id,
            tenant_id: data.tenant_id,
          }),
        },
      );

      if (!response.ok) {
        const errorData = await response.text();
        console.error("Failed to delete invite:", errorData);
      }

      setInvites((prevInvites) =>
        prevInvites.filter((invite) => invite.id !== data.id),
      );

      const resp_data = await response.json();
      setTeacherInviteData(resp_data);
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <>
      {showHeader && (
        <Title colorConfig={colorConfig} icon={FaEnvelope}>
          Pozivi nastavnicima za predmete u odjeljenjima
        </Title>
      )}
      {/* Error Modal */}
      {errorMessage && (
        <ErrorModal
          onClose={() => setErrorMessage(null)}
          colorConfig={colorConfig}
        >
          {errorMessage}
        </ErrorModal>
      )}

      {/* Subjects Modal */}
      {showSubjects && (
        <Modal
          onClose={() => {
            setSubjects([]);
            setShowSubjects(false);
          }}
        >
          <Title colorConfig={colorConfig} icon={FaBook}>
            Predmeti za poziv
          </Title>
          <DynamicCardParent
            data={subjects}
            keysToIgnore={["subject_code"]}
            textTitle="Predmet"
            keysToExclude={["subject_name"]}
            icon={FaBook}
            tenantColorConfig={colorConfig}
          />
        </Modal>
      )}

      {/* Back Button (only for admin mode) */}
      {!teacherAccountMode && (
        <div className="flex justify-end mb-4">
          <div className="mr-2">
            <BackButton onClick={onBack} colorConfig={colorConfig} />
          </div>
        </div>
      )}

      {/* Main Content */}
      <DynamicCardParent
        data={dateFormattedInvites}
        setData={setInvites}
        icon={<FaEnvelope />}
        titleField="teacher_full_name"
        prefix="poziv"
        accessToken={accessToken}
        keysToIgnore={[
          "id",
          "teacher_id",
          "section_id",
          "subjects",
          "tenant_id",
        ]}
        showDelete={(item) => item.status === "pending"}
        deleteButton={{
          label: "Izbriši",
          onClick: (item) => {
            setInviteToDelete(item);
          },
          icon: FaTimes,
        }}
        extraButton={{
          label: "Predmeti",
          onClick: (data) => {
            setSubjects(data?.subjects);
            setShowSubjects(true);
          },
          icon: FaBook,
        }}
        mode={mode}
        tenantColorConfig={colorConfig}
        {...getActionButtons()}
        emptyMessage={
          <>
            <div>Trenutno nema poziva.</div>
            <div>Pozovi nastavnika da se pridruži odjeljenju.</div>
          </>
        }
      />
      {inviteToDelete && (
        <ConfirmModal
          title="Brisanje poziva"
          onClose={() => setInviteToDelete(null)}
          onConfirm={() => {
            deleteTeacherInvite(inviteToDelete);
            setInviteToDelete(null);
          }}
          colorConfig={colorConfig}
        >
          Da li ste sigurni da želite izbrisati poziv u odjeljenje{" "}
          <span className="font-bold">
            Odjeljenje {inviteToDelete.section_name} (
            {inviteToDelete.tenant_name})
          </span>{" "}
          za nastavnika{" "}
          <span className="font-bold">{inviteToDelete.teacher_full_name}</span>?
        </ConfirmModal>
      )}
    </>
  );
}
