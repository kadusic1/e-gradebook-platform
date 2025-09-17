import ConfirmModal from "../components/modal/ConfirmModal";

export const handleArchiveConfirmUtil = async (
  tenantId,
  sectionId,
  accessToken,
) => {
  try {
    const response = await fetch(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/teacher/archive_section/${tenantId}/${sectionId}`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${accessToken}`,
        },
      },
    );
    if (response.ok) {
      return null;
    } else {
      const errorData = await response.text();
      return errorData;
    }
  } catch (error) {
    console.error(error);
    return null;
  }
};

export const ArchiveConfirmModal = ({
  isOpen,
  selectedSection,
  onClose,
  onConfirm,
  colorConfig,
}) => {
  if (!isOpen || !selectedSection) return null;

  return (
    <ConfirmModal
      title="Potvrda arhiviranja"
      onClose={onClose}
      onConfirm={onConfirm}
      colorConfig={colorConfig}
    >
      Da li ste sigurni da želite arhivirati odjeljenje{" "}
      <span className="font-bold">{selectedSection.name}? </span>
      <span className="font-bold text-red-500">Ova radnja je nepovratna!</span>
    </ConfirmModal>
  );
};
