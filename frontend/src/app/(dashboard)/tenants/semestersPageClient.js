"use client";
import { useState, useMemo, useEffect } from "react";
import ErrorModal from "../../components/modal/ErrorModal";
import DynamicCardParent from "../../components/dynamic_card/DynamicCardParent";
import { FaRegCalendarAlt } from "react-icons/fa";
import BackButton from "@/app/components/common/BackButton";
import Title from "@/app/components/common/Title";
import { formatDateToDDMMYYYY } from "@/app/util/date_util";

export default function TenantSemestersPageClient({
  accessToken,
  tenant,
  onBack,
}) {
  const [tenantSemesters, setTenantSemesters] = useState([]);
  const [errorMessage, setErrorMessage] = useState(null);

  const fetchTenantSemesters = async () => {
    const response = await fetch(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/tenant_semesters/${tenant?.id}`,
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      },
    );
    if (response.ok) {
      const semesters = await response.json();
      setTenantSemesters(semesters);
    }
  };

  useEffect(() => {
    fetchTenantSemesters();
  }, []);

  const dateFormattedSemesters = useMemo(() => {
    return tenantSemesters?.map((semester) => ({
      ...semester,
      start_date: formatDateToDDMMYYYY(semester.start_date),
      end_date: formatDateToDDMMYYYY(semester.end_date),
    }));
  }, [tenantSemesters]);

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
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/tenant_semesters_update`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
          body: JSON.stringify({
            tenant_id: newData.tenant_id,
            semester_code: newData.semester_code,
            start_date: newData.start_date,
            end_date: newData.end_date,
            npp_code: newData.npp_code,
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

  return (
    <>
      <Title icon={FaRegCalendarAlt} colorConfig={tenant.color_config}>
        Polugodišta
      </Title>
      <div className="flex justify-end mb-4 mr-2">
        {onBack && (
          <div className="mr-2">
            <BackButton onClick={onBack} colorConfig={tenant.color_config} />
          </div>
        )}
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
      <DynamicCardParent
        data={dateFormattedSemesters}
        setData={setTenantSemesters}
        icon={<FaRegCalendarAlt />}
        titleField="full_name"
        prefix="polugodište"
        editFields={semester_edit_fields}
        editUrl={`${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/tenant_semesters_update`}
        showEdit={true}
        accessToken={accessToken}
        keysToIgnore={["npp_code", "semester_code", "tenant_id"]}
        keyField="npp_code"
        onEditSave={handleSemesterEditSave}
        tenantColorConfig={tenant.color_config}
        mode={tenant.semester_display}
      />
    </>
  );
}
