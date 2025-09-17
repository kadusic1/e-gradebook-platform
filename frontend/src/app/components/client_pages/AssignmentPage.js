"use client";
import { useState } from "react";
import AddButton from "@/app/components/common/AddButton";
import ErrorModal from "@/app/components/modal/ErrorModal";
import DynamicCardParent from "@/app/components/dynamic_card/DynamicCardParent";
import MultiSelectModal from "@/app/components/modal/MultiSelectModal";
import BackButton from "@/app/components/common/BackButton";
import SuccessModal from "../modal/SuccessModal";
import { FaHourglassHalf } from "react-icons/fa";
import DynamicTab from "../dynamic_card/DynamicTab";
import ConfirmModal from "../modal/ConfirmModal";
import Title from "../common/Title";

export default function AssignmentPage({
  assignedItems = [],
  setAssignedItems,
  availableItems = [],
  setAvailableItems,
  onBack,
  showOnlyOnSearch = false,
  inviteMode = false,
  pendingItems = [],
  setPendingItems,

  // Configuration props
  config: {
    // API endpoints
    assignEndpoint,
    unassignEndpoint,

    archived = 0,

    // UI configuration
    addButtonText,
    modalTitle,
    searchPlaceholder,
    modalNote,

    // Field mappings
    keyField = "id",
    labelField,
    titleField,
    prefix,
    keysToIgnore = [],

    // Icon and display
    icon,
    getIcon,

    // Error messages
    assignErrorMessage = "Došlo je do greške prilikom dodijeljivanja.",
    unassignErrorMessage = "Došlo je do greške prilikom uklanjanja.",
    orderFields = [],
    mainTitle,
    itemsTitle,
    pendingTitle,
    showEdit,
    editButton,
    archivedEditButton,
    extraButton,
    archivedExtraButton,
    getBgColor,
    getTextColor,
  },
  pendingConfig: {
    pendingKeyField = "id",
    pendingTitleField,
    pendingPrefix,
    pendingKeysToIgnore,
    deleteInviteButton,
    getDeleteInviteMessage,
  } = {},
  accessToken,
  colorConfig,
  mode = "card",
  pendingMode = "table",
}) {
  const [showModal, setShowModal] = useState(false);
  const [errorMessage, setErrorMessage] = useState(null);
  const [successMessage, setSuccessMessage] = useState(null);
  const [activeTab, setActiveTab] = useState(0);
  const [inviteToDelete, setInviteToDelete] = useState(null);

  const assignItems = async (itemIds) => {
    try {
      const response = await fetch(assignEndpoint, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${accessToken}`,
        },
        body: JSON.stringify(itemIds),
      });

      if (response.ok) {
        const itemsToAssign = availableItems.filter((item) =>
          itemIds.includes(item[keyField]),
        );

        if (!inviteMode) {
          setAssignedItems((prev) => [...(prev || []), ...itemsToAssign]);
        } else {
          setSuccessMessage("Učenici su uspješno pozvani u odjeljenje.");
          const pendingItemsToAdd = await response.json();
          setPendingItems((prev) => [...(prev || []), ...pendingItemsToAdd]);
        }
        setAvailableItems((prev) =>
          prev.filter((item) => !itemIds.includes(item[keyField])),
        );
      } else {
        const errorData = await response.json();
        setErrorMessage(errorData.message || assignErrorMessage);
      }
      setShowModal(false);
    } catch (error) {
      console.error("Error assigning items:", error);
      setErrorMessage(assignErrorMessage);
    }
  };

  const handleDeleteConfirm = async (itemId, closeModal) => {
    try {
      const response = await fetch(`${unassignEndpoint}/${itemId}`, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${accessToken}`,
        },
      });

      if (response.ok) {
        const itemToMove = assignedItems.find(
          (item) => item[keyField] === itemId,
        );

        setAssignedItems((prev) =>
          prev.filter((item) => item[keyField] !== itemId),
        );

        if (itemToMove) {
          setAvailableItems((prev) =>
            [...(prev || []), itemToMove].sort((a, b) => {
              // Use orderFields for sorting
              for (const field of orderFields) {
                if (a[field] !== b[field]) {
                  return String(a[field] || "").localeCompare(
                    String(b[field] || ""),
                  );
                }
              }
              return 0;
            }),
          );
        }
      } else {
        const errorData = await response.text();
        setErrorMessage(errorData || unassignErrorMessage);
      }
      closeModal();
    } catch (error) {
      console.error("Error unassigning item:", error);
      setErrorMessage(unassignErrorMessage);
    }
  };

  return (
    <>
      {successMessage && inviteMode && (
        <SuccessModal
          colorConfig={colorConfig}
          onClose={() => setSuccessMessage(null)}
        >
          {successMessage}
        </SuccessModal>
      )}
      <div className="flex justify-end mb-4">
        {!inviteMode && (
          <>
            {onBack && (
              <div className="mr-2">
                <BackButton onClick={onBack} colorConfig={colorConfig} />
              </div>
            )}
            {archived == 0 && (
              <AddButton
                onClick={() => setShowModal(true)}
                colorConfig={colorConfig}
              >
                {addButtonText}
              </AddButton>
            )}
          </>
        )}
        {errorMessage && (
          <ErrorModal
            onClose={() => setErrorMessage(null)}
            colorConfig={colorConfig}
          >
            {errorMessage}
          </ErrorModal>
        )}
        {showModal && (
          <MultiSelectModal
            title={modalTitle}
            items={availableItems}
            onClose={() => setShowModal(false)}
            onSave={assignItems}
            labelField={labelField}
            searchPlaceholder={searchPlaceholder}
            note={modalNote}
            keyField={keyField}
            showOnlyOnSearch={showOnlyOnSearch}
            colorConfig={colorConfig}
          />
        )}
      </div>
      {!inviteMode ? (
        <DynamicCardParent
          data={assignedItems}
          setData={setAssignedItems}
          getIcon={getIcon}
          titleField={titleField}
          prefix={prefix}
          keyField={keyField}
          keysToIgnore={keysToIgnore}
          showDelete={true}
          onDeleteConfirm={handleDeleteConfirm}
          tenantColorConfig={colorConfig}
          mode={mode}
          showEdit={showEdit}
          editButton={archived == 0 ? editButton : archivedEditButton}
          extraButton={archived == 1 ? archivedExtraButton : extraButton}
          getBgColorOption={getBgColor}
          getTextColorOption={getTextColor}
        />
      ) : (
        <>
          <DynamicTab
            title={mainTitle}
            titleIcon={icon}
            activeTab={activeTab}
            setActiveTab={setActiveTab}
            colorConfig={colorConfig}
            childrenTabs={[
              {
                label: itemsTitle,
                content: (
                  <>
                    <Title colorConfig={colorConfig} icon={icon}>
                      {itemsTitle}
                    </Title>
                    <div className="flex justify-end mb-4">
                      <div className="mr-2">
                        <BackButton
                          onClick={onBack}
                          colorConfig={colorConfig}
                        />
                      </div>
                      {archived == 0 && (
                        <AddButton
                          onClick={() => setShowModal(true)}
                          colorConfig={colorConfig}
                        >
                          {addButtonText}
                        </AddButton>
                      )}
                    </div>
                    <DynamicCardParent
                      data={assignedItems}
                      setData={setAssignedItems}
                      getIcon={getIcon}
                      titleField={titleField}
                      prefix={prefix}
                      keyField={keyField}
                      keysToIgnore={keysToIgnore}
                      showDelete={archived == 0}
                      onDeleteConfirm={handleDeleteConfirm}
                      tenantColorConfig={colorConfig}
                      mode={mode}
                      showEdit={showEdit}
                      editButton={
                        archived == 0 ? editButton : archivedEditButton
                      }
                      extraButton={
                        archived == 1 ? archivedExtraButton : extraButton
                      }
                      getBgColorOption={getBgColor}
                      getTextColorOption={getTextColor}
                    />
                  </>
                ),
                icon: icon,
              },
              {
                label: pendingTitle,
                content: (
                  <>
                    <Title colorConfig={colorConfig} icon={FaHourglassHalf}>
                      {pendingTitle}
                    </Title>
                    <div className="flex justify-end mb-4">
                      <div className="mr-2">
                        <BackButton
                          onClick={onBack}
                          colorConfig={colorConfig}
                        />
                      </div>
                      {archived == 0 && (
                        <AddButton
                          onClick={() => setShowModal(true)}
                          colorConfig={colorConfig}
                        >
                          {addButtonText}
                        </AddButton>
                      )}
                    </div>
                    <DynamicCardParent
                      data={pendingItems}
                      icon={<FaHourglassHalf />}
                      titleField={pendingTitleField}
                      prefix={pendingPrefix}
                      keyField={pendingKeyField}
                      keysToIgnore={pendingKeysToIgnore}
                      deleteButton={{
                        ...deleteInviteButton,
                        onClick: (item) => {
                          setInviteToDelete(item);
                        },
                      }}
                      mode={pendingMode}
                      showDelete={(item) =>
                        item.status === "pending" && archived == 0
                      }
                      tenantColorConfig={colorConfig}
                    />
                  </>
                ),
                icon: FaHourglassHalf,
              },
            ]}
          />
          {inviteToDelete && (
            <ConfirmModal
              title="Brisanje poziva"
              onClose={() => setInviteToDelete(null)}
              onConfirm={() => {
                deleteInviteButton.onClick(inviteToDelete);
                setInviteToDelete(null);
              }}
              colorConfig={colorConfig}
            >
              {getDeleteInviteMessage(inviteToDelete)}
            </ConfirmModal>
          )}
        </>
      )}
    </>
  );
}
