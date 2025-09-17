"use client";
import { useState, useMemo } from "react";
import ErrorModal from "../../components/modal/ErrorModal";
import DynamicCardParent from "../../components/dynamic_card/DynamicCardParent";
import { FaRegEnvelope, FaTimes, FaCheck } from "react-icons/fa";
import { formatDateToDDMMYYYY } from "@/app/util/date_util";

export default function PupilSectionsInvitePageClient({
  initialInvites = [],
  accessToken,
}) {
  const [errorMessage, setErrorMessage] = useState(null);
  const [invites, setInvites] = useState(initialInvites);

  const dateFormattedInvites = useMemo(() => {
    return invites?.map((invite) => ({
      ...invite,
      invite_date: formatDateToDDMMYYYY(invite.invite_date),
    }));
  }, [invites]);

  // Handler for accepting an invite
  const handleAction = async (invite, action) => {
    try {
      const res = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/common/respond_to_section_invite`,
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

      // Update the status of the handled invite to action
      let newStatus = action === "accept" ? "accepted" : "declined";
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

  const handleAccept = (invite) => {
    handleAction(invite, "accept");
  };

  const handleDecline = (invite) => {
    handleAction(invite, "decline");
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
      <DynamicCardParent
        data={dateFormattedInvites}
        setData={setInvites}
        icon={<FaRegEnvelope />}
        titleField="section_name"
        prefix="poziv"
        accessToken={accessToken}
        keysToIgnore={[
          "id",
          "pupil_id",
          "section_id",
          "pupil_full_name",
          "showActions",
          "tenant_id",
        ]}
        keyField="id"
        showDelete={(item) => item.status === "pending"}
        deleteButton={{
          onClick: handleDecline,
          label: "Odbij",
          icon: FaTimes,
        }}
        showEdit={(item) => item.status === "pending"}
        editButton={{
          onClick: handleAccept,
          label: "Prihvati",
          icon: FaCheck,
        }}
        mode="table"
      />
    </>
  );
}
