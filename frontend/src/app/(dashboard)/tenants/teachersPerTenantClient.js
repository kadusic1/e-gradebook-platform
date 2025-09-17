"use client";
import {
  FaChalkboardTeacher,
  FaEnvelope,
  FaClipboardList,
  FaTrash,
} from "react-icons/fa";
import AssignSubjectsModal from "@/app/components/modal/AssignSubjectsModal";
import DynamicCardParent from "@/app/components/dynamic_card/DynamicCardParent";
import { useState, useEffect } from "react";
import AddButton from "@/app/components/common/AddButton";
import ErrorModal from "@/app/components/modal/ErrorModal";
import BackButton from "@/app/components/common/BackButton";
import Button from "@/app/components/common/Button";
import TeacherInvitesPageClient from "./teacherInvitesPageClient";
import ConfirmModal from "@/app/components/modal/ConfirmModal";
import Title from "@/app/components/common/Title";

export default function TeachersPerTenantPageClient({
  tenant,
  onBack,
  accessToken,
}) {
  const [teachers, setTeachers] = useState([]);
  const [teacherInviteData, setTeacherInviteData] = useState([]);
  const [showModal, setShowModal] = useState(false);
  const [errorMessage, setErrorMessage] = useState(null);
  const [showInvitePage, setShowInvitePage] = useState(false);
  const [initialTeacherID, setInitialTeacherID] = useState(null);
  const [teacherToDelete, setTeacherToDelete] = useState(null);

  const fetchTenantTeachers = async () => {
    const response = await fetch(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/teachers_per_tenant/${tenant?.id}`,
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      },
    );
    if (response.ok) {
      const teachers = await response.json();
      setTeachers(teachers);
    }
  };

  const fetchTeacherInviteData = async () => {
    const response = await fetch(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/teacher_invite_data/${tenant?.id}`,
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      },
    );
    if (response.ok) {
      const data = await response.json();
      setTeacherInviteData(data);
    }
  };

  useEffect(() => {
    fetchTenantTeachers();
  }, []);

  useEffect(() => {
    fetchTeacherInviteData();
  }, []);

  const sendTeacherAssignments = async (teacherId, assignments) => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/teacher_section_assignments/${tenant.id}/${teacherId}`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
          body: JSON.stringify(assignments),
        },
      );

      if (!response.ok) {
        throw new Error("Failed to assign sections");
      }
      const data = await response.json();
      return data;
    } catch (error) {
      console.error("Error assigning teacher sections:", error);
      setErrorMessage("Failed to assign sections");
    }
  };

  const handleSaveAssignments = async (data, teacherId) => {
    const updatedAssignments = await sendTeacherAssignments(teacherId, data);

    try {
      if (updatedAssignments) {
        setTeacherInviteData((prevData) => {
          return prevData.map((assignment) => {
            // Find assignments for this teacher and update them
            if (assignment.teacher.id == teacherId) {
              // Find the corresponding updated assignment
              const updatedAssignment = updatedAssignments.invite_data.find(
                (updated) => updated.section.id === assignment.section.id,
              );
              return updatedAssignment || assignment;
            }
            return assignment;
          });
        });
      }

      if (updatedAssignments.teacher_data) {
        setTeachers(updatedAssignments.teacher_data);
      }
    } catch (error) {
      setErrorMessage("Greška prilikom uređivanja predmeta i razredništva.");
    }
  };

  const deleteTeacherForTenant = async (teacherID) => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/delete_teacher_from_tenant/${tenant.id}/${teacherID}`,
        {
          method: "DELETE",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
        },
      );

      if (!response.ok) {
        throw new Error("Failed to delete teacher from tenant");
      }

      setTeachers((prevTeachers) =>
        prevTeachers.filter((t) => t.id !== teacherID),
      );
    } catch (error) {
      console.error("Error deleting teacher from tenant:", error);
      setErrorMessage("Failed to delete teacher from tenant");
    }
  };

  return !showInvitePage ? (
    <>
      <Title colorConfig={tenant.color_config} icon={FaChalkboardTeacher}>
        Profesori
      </Title>
      <div className="flex justify-end mb-4">
        {onBack && (
          <div className="mr-2">
            <BackButton onClick={onBack} colorConfig={tenant.color_config} />
          </div>
        )}
        <Button
          onClick={() => setShowInvitePage(true)}
          icon={FaEnvelope}
          color="quaternary"
          className="mr-2"
          colorConfig={tenant.color_config}
        >
          Pozivi
        </Button>
        <AddButton
          onClick={() => setShowModal(true)}
          colorConfig={tenant.color_config}
        >
          Zaduženja
        </AddButton>
        {errorMessage && (
          <ErrorModal
            onClose={() => setErrorMessage(null)}
            colorConfig={tenant.color_config}
          >
            {errorMessage}
          </ErrorModal>
        )}
        {showModal && (
          <AssignSubjectsModal
            open={showModal}
            onClose={() => {
              setShowModal(false);
              setInitialTeacherID(null);
            }}
            assignmentsData={teacherInviteData}
            onSave={handleSaveAssignments}
            initialTeacherID={initialTeacherID}
            colorConfig={tenant.color_config}
          />
        )}
      </div>
      <DynamicCardParent
        data={teachers}
        setData={setTeachers}
        icon={<FaChalkboardTeacher />}
        titleField={["name", "last_name"]}
        prefix="nastavnika"
        accessToken={accessToken}
        keysToIgnore={["account_type"]}
        showDelete={true}
        deleteButton={{
          label: "Izbriši",
          onClick: (data) => {
            setTeacherToDelete(data);
          },
          icon: FaTrash,
        }}
        extraButton={{
          label: "Zaduženja",
          onClick: (data) => {
            setInitialTeacherID(String(data.id));
            setShowModal(true);
          },
          icon: FaClipboardList,
        }}
        tenantColorConfig={tenant.color_config}
        mode={tenant.teacher_display}
      />
      {teacherToDelete && (
        <ConfirmModal
          onClose={() => setTeacherToDelete(null)}
          onConfirm={() => {
            deleteTeacherForTenant(teacherToDelete.id);
            setTeacherToDelete(null);
          }}
          colorConfig={tenant.color_config}
        >
          Da li ste sigurni da želite izbrisati nastavnika{" "}
          <span className="font-bold">
            {teacherToDelete.name} {teacherToDelete.last_name}
          </span>
          ?
        </ConfirmModal>
      )}
    </>
  ) : (
    <TeacherInvitesPageClient
      accessToken={accessToken}
      onBack={() => setShowInvitePage(false)}
      tenant={tenant}
      setTeacherInviteData={setTeacherInviteData}
      colorConfig={tenant.color_config}
      mode={tenant.teacher_invite_display}
    />
  );
}
