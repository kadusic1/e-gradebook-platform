"use client";
import { useState } from "react";
import ErrorModal from "../../components/modal/ErrorModal";
import DynamicCardParent from "../../components/dynamic_card/DynamicCardParent";
import {
  FaRegCalendarAlt,
  FaGlobeEurope,
  FaLink,
  FaTrash,
} from "react-icons/fa";
import DynamicTab from "@/app/components/dynamic_card/DynamicTab";
import CreateUpdateModal from "@/app/components/modal/CreateUpdateModal";
import AddButton from "@/app/components/common/AddButton";
import ConfirmModal from "@/app/components/modal/ConfirmModal";
import Title from "@/app/components/common/Title";

export default function GlobalAdminPageClient({
  initialSemesters = [],
  initialDomains = [],
  accessToken,
}) {
  const [errorMessage, setErrorMessage] = useState(null);
  const [semesters, setSemesters] = useState(initialSemesters);
  const [domains, setDomains] = useState(initialDomains);
  const [activeTab, setActiveTab] = useState(0);
  const [showDomainModal, setShowDomainModal] = useState(false);
  const [domainToDelete, setDomainToDelete] = useState(null);

  const semester_edit_fields = [
    {
      label: "Datum početka",
      name: "start_date",
      placeholder: "Unesite datum početka",
      type: "date",
    },
    {
      label: "Datum kraja",
      name: "end_date",
      placeholder: "Unesite datum kraja",
      type: "date",
    },
  ];

  async function handleSemesterEditSave(
    newData,
    setData,
    closeModal,
    setError,
  ) {
    try {
      const resp = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/superadmin/npp_semesters_update`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
          body: JSON.stringify({
            npp_code: newData.npp_code,
            semester_code: newData.semester_code,
            start_date: newData.start_date,
            end_date: newData.end_date,
          }),
        },
      );
      if (!resp.ok) {
        const error_message = await resp.text();
        setError(error_message);
      } else {
        const updatedSemester = await resp.json();
        setData((prev) =>
          prev.map((item) =>
            item.npp_code === updatedSemester.npp_code &&
            item.semester_code === updatedSemester.semester_code
              ? updatedSemester
              : item,
          ),
        );
        closeModal();
      }
    } catch (e) {
      setError("Greška pri snimanju.");
    }
  }

  const domain_fields = [
    {
      label: "Domena",
      name: "domain",
      type: "text",
      placeholder: "Unesite domenu (npr. gmail.com)",
    },
  ];

  const handleGlobalDomainCreate = async (data) => {
    try {
      const resp = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/superadmin/domain_create`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
          body: JSON.stringify(data),
        },
      );
      if (!resp.ok) {
        const error_message = await resp.text();
        setErrorMessage(error_message);
      } else {
        setDomains((prev) => [
          ...(prev || []),
          {
            domain: data.domain,
            type: "global_domain",
          },
        ]);
        setShowDomainModal(false);
      }
    } catch (e) {
      console.error("Error creating domain:", e);
      setErrorMessage("Greška pri snimanju.");
    }
  };

  const handleGlobalDomainDelete = async (domain) => {
    try {
      const resp = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/superadmin/domains_delete`,
        {
          method: "DELETE",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
          body: JSON.stringify({ domain }),
        },
      );
      if (!resp.ok) {
        const error_message = await resp.text();
        setErrorMessage(error_message);
      } else {
        setDomains((prev) => prev.filter((item) => item.domain !== domain));
        setDomainToDelete(null);
      }
    } catch (e) {
      console.error("Error deleting domain:", e);
      setErrorMessage("Greška pri brisanju domene.");
    }
  };

  return (
    <>
      <div className="flex justify-end mb-4">
        {errorMessage && (
          <ErrorModal onClose={() => setErrorMessage(null)}>
            {errorMessage}
          </ErrorModal>
        )}
      </div>
      <DynamicTab
        title="Globalna administracija"
        titleIcon={FaGlobeEurope}
        activeTab={activeTab}
        setActiveTab={setActiveTab}
        childrenTabs={[
          {
            label: "Polugodišta",
            content: (
              <>
                <Title icon={FaRegCalendarAlt}>Polugodišta</Title>
                <DynamicCardParent
                  data={semesters}
                  setData={setSemesters}
                  icon={<FaRegCalendarAlt />}
                  titleField="full_name"
                  prefix="polugodište"
                  editFields={semester_edit_fields}
                  editUrl={`${process.env.NEXT_PUBLIC_API_BASE_URL}/api/superadmin/npp_semesters_update`}
                  showEdit={true}
                  accessToken={accessToken}
                  keysToIgnore={["npp_code", "semester_code"]}
                  keyField="npp_code"
                  onEditSave={handleSemesterEditSave}
                />
              </>
            ),
            icon: FaRegCalendarAlt,
          },
          {
            label: "Domene",
            content: (
              <>
                <Title icon={FaLink}>Domene</Title>
                {domainToDelete && (
                  <ConfirmModal
                    onClose={() => setDomainToDelete(null)}
                    onConfirm={() => {
                      handleGlobalDomainDelete(domainToDelete);
                      setDomainToDelete(null);
                    }}
                  >
                    Da li ste sigurni da želite izbrisati domenu &ldquo;
                    {domainToDelete}?&rdquo;?
                  </ConfirmModal>
                )}
                <div className="flex justify-end mb-4">
                  <AddButton onClick={() => setShowDomainModal(true)}>
                    Dodaj domenu
                  </AddButton>
                </div>
                {showDomainModal && (
                  <CreateUpdateModal
                    title="Dodaj domenu"
                    fields={domain_fields}
                    onSave={handleGlobalDomainCreate}
                    onClose={() => setShowDomainModal(false)}
                    show={showDomainModal}
                  />
                )}
                <DynamicCardParent
                  data={domains}
                  setData={setDomains}
                  icon={<FaLink />}
                  prefix="domenu"
                  titleField="domain"
                  accessToken={accessToken}
                  keyField="domain"
                  mode="table"
                  showDelete={(item) => item.type === "global_domain"}
                  deleteButton={{
                    label: "Izbriši",
                    onClick: (data) => setDomainToDelete(data.domain),
                    icon: FaTrash,
                  }}
                />
              </>
            ),
            icon: FaLink,
          },
        ]}
      />
    </>
  );
}
