import CreateUpdateModal from "./CreateUpdateModal";

export default function PasswordChangeModal({
  setErrorMessage,
  setShowPasswordChange,
  accessToken,
  setSuccessMessage,
}) {
  const passwordFields = [
    {
      name: "current_password",
      label: "Trenutna lozinka",
      type: "password",
      placeholder: "Unesite trenutnu lozinku",
    },
    {
      name: "new_password",
      label: "Nova lozinka",
      type: "password",
      placeholder: "Unesite novu lozinku",
    },
    {
      name: "confirm_password",
      label: "Potvrdi novu lozinku",
      type: "password",
      placeholder: "Potvrdite novu lozinku",
    },
  ];

  const changeAccountPassword = async (data) => {
    try {
      const resp = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/common/change_password`,
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
        const errorMsg = await resp.text();
        setErrorMessage(errorMsg);
      } else {
        setSuccessMessage("Lozinka uspješno promijenjena");
      }
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <CreateUpdateModal
      title="Promjena lozinke"
      fields={passwordFields}
      onClose={() => setShowPasswordChange(false)}
      onSave={(data) => {
        changeAccountPassword(data);
        setShowPasswordChange(false);
      }}
    />
  );
}
