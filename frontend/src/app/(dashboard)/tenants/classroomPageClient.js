"use client";
import { useState, useEffect } from "react";
import ErrorModal from "../../components/modal/ErrorModal";
import DynamicCardParent from "../../components/dynamic_card/DynamicCardParent";
import { FaChalkboard } from "react-icons/fa";
import BackButton from "@/app/components/common/BackButton";
import CreateUpdateModal from "../../components/modal/CreateUpdateModal";
import AddButton from "@/app/components/common/AddButton";
import ConfirmModal from "@/app/components/modal/ConfirmModal";
import Title from "@/app/components/common/Title";

export default function ClassroomPageClient({ accessToken, tenant, onBack }) {
  const [errorMessage, setErrorMessage] = useState(null);
  const [classrooms, setClassrooms] = useState([]);
  const [showModal, setShowModal] = useState(false);
  const [itemToEdit, setItemToEdit] = useState(null);
  const [itemToDelete, setItemToDelete] = useState(null);

  const handleClassroomCreate = async (data) => {
    try {
      const payload = {
        ...data,
        capacity: data.capacity ? parseInt(data.capacity, 10) : undefined,
      };
      const resp = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/create_classroom/${tenant.id}`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
          body: JSON.stringify(payload),
        },
      );
      if (!resp.ok) {
        const error_message = await resp.text();
        setErrorMessage(error_message);
      } else {
        // Append new data to classrooms
        setClassrooms((prev) => [
          ...(prev || []),
          { ...data, name: `${data.type} ${data.code}` },
        ]);
        setShowModal(false);
      }
    } catch (error) {
      setErrorMessage("Greška pri kreiranju učionice.");
    }
  };

  const handleClassroomEdit = async (data) => {
    try {
      const payload = {
        ...data,
        capacity: data.capacity ? parseInt(data.capacity, 10) : undefined,
      };
      const resp = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/update_classroom/${tenant.id}/${itemToEdit.code}`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
          body: JSON.stringify(payload),
        },
      );
      if (!resp.ok) {
        const error_message = await resp.text();
        setErrorMessage(error_message);
      } else {
        // Find the classroom with itemToEdit.code and update it
        setClassrooms((prev) =>
          prev.map((classroom) =>
            classroom.code === itemToEdit.code
              ? { ...classroom, ...data, name: `${data.type} ${data.code}` }
              : classroom,
          ),
        );
      }
      setItemToEdit(null);
    } catch (error) {
      setErrorMessage("Greška pri uređivanju učionice.");
    }
  };

  const getClassrooms = async () => {
    try {
      const resp = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/get_all_classrooms/${tenant.id}`,
        {
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
        },
      );
      if (!resp.ok) {
        const error_message = await resp.text();
        setErrorMessage(error_message);
      }
      const classrooms = await resp.json();
      setClassrooms(classrooms);
    } catch (e) {
      setErrorMessage("Greška pri učitavanju učionica.");
    }
  };

  const deleteClassroom = async (classroomToDelete) => {
    try {
      const resp = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/delete_classroom/${tenant.id}/${classroomToDelete.code}`,
        {
          method: "DELETE",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
        },
      );
      if (!resp.ok) {
        const error_message = await resp.text();
        setErrorMessage(error_message);
      } else {
        setClassrooms((prev) =>
          prev.filter((classroom) => classroom.code !== classroomToDelete.code),
        );
      }
      setItemToDelete(null);
    } catch (error) {
      setErrorMessage("Greška pri brisanju učionice.");
    }
  };

  const classroom_fields = [
    {
      label: "Broj učionice",
      name: "code",
      placeholder: "Unesite broj učionice (npr. 204 ili 15A)",
    },
    {
      label: "Tip učionice",
      name: "type",
      placeholder:
        "Unesite tip učionice (npr. Učionica, Laboratorija, Sport sala)",
    },
    {
      label: "Kapacitet učionice",
      name: "capacity",
      placeholder: "Unesite kapacitet učionice (npr. 30)",
      type: "positive-number",
    },
  ];

  useEffect(() => {
    getClassrooms();
  }, [accessToken, tenant.id]);

  return (
    <>
      <Title icon={FaChalkboard} colorConfig={tenant.color_config}>
        Učionice
      </Title>
      <div className="flex justify-end mb-4 mr-2">
        {onBack && (
          <div className="mr-2">
            <BackButton onClick={onBack} colorConfig={tenant.color_config} />
          </div>
        )}
        <AddButton
          onClick={() => setShowModal(true)}
          colorConfig={tenant.color_config}
        >
          Dodaj učionicu
        </AddButton>
      </div>
      <div className="flex justify-end mb-4">
        {errorMessage && (
          <ErrorModal
            onClose={() => setErrorMessage(null)}
            colorConfig={tenant.color_config}
          >
            {errorMessage}
          </ErrorModal>
        )}
      </div>
      {showModal && (
        <CreateUpdateModal
          title="Dodaj učionicu"
          fields={classroom_fields}
          onClose={() => setShowModal(false)}
          onSave={handleClassroomCreate}
          colorConfig={tenant.color_config}
        />
      )}
      {itemToEdit && (
        <CreateUpdateModal
          title={`Uredi učionicu`}
          fields={classroom_fields}
          onClose={() => setItemToEdit(null)}
          onSave={handleClassroomEdit}
          initialValues={itemToEdit}
          colorConfig={tenant.color_config}
        />
      )}
      {itemToDelete && (
        <ConfirmModal
          onClose={() => setItemToDelete(null)}
          colorConfig={tenant.color_config}
          onConfirm={() => deleteClassroom(itemToDelete)}
        >
          Da li ste sigurni da želite izbrisati stavku &ldquo;
          {itemToDelete.code}&rdquo; ({itemToDelete.type})?
        </ConfirmModal>
      )}
      <DynamicCardParent
        data={classrooms}
        setData={setClassrooms}
        icon={<FaChalkboard />}
        prefix="učionicu"
        accessToken={accessToken}
        tenantColorConfig={tenant.color_config}
        titleField="name"
        keysToIgnore={["code"]}
        showDelete={true}
        deleteButton={{
          onClick: (data) => setItemToDelete(data),
        }}
        showEdit={true}
        editFields={classroom_fields}
        editButton={{
          onClick: (data) => setItemToEdit(data),
        }}
        mode={tenant.classroom_display}
      />
    </>
  );
}
