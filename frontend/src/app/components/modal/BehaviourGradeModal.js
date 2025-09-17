"use client";
import Modal from "./Modal";
import Title from "../common/Title";
import DynamicCardParent from "../dynamic_card/DynamicCardParent";
import { FaClipboardCheck, FaHistory } from "react-icons/fa";
import { useState } from "react";
import { formatToFullDateTime } from "@/app/util/date_util";

export const BehaviourGradeModal = ({
  data,
  onClose,
  colorConfig,
  onEditClick,
  showEdit = false,
  tenantID,
  accessToken,
}) => {
  const [historyBehaviour, setHistoryBehaviour] = useState(null);
  const [showHistory, setShowHistory] = useState(false);

  const fetchBehaviourGradeHistory = async (behaviourGradeID) => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/pupil/behaviour_grade_history/${tenantID}/${behaviourGradeID}`,
        {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
        },
      );

      if (response.ok) {
        const data = await response.json();
        setHistoryBehaviour(
          data?.map((item) => ({
            ...item,
            valid_until: formatToFullDateTime(item.valid_until),
          })),
        );
        setShowHistory(true);
      }
    } catch (error) {
      console.error(error);
    }
  };

  if (showHistory) {
    return (
      <Modal
        onClose={() => {
          setHistoryBehaviour(null);
          setShowHistory(false);
        }}
      >
        <Title colorConfig={colorConfig} icon={FaHistory}>
          Historija vladanja
        </Title>
        <DynamicCardParent
          data={historyBehaviour}
          icon={FaHistory}
          titleField="behaviour"
          tenantColorConfig={colorConfig}
          keysToIgnore={["pupil_id", "section_id", "semester_code"]}
          twoColumnsCard={false}
        />
      </Modal>
    );
  }

  return (
    <Modal onClose={onClose}>
      <Title colorConfig={colorConfig} icon={FaClipboardCheck}>
        Vladanje učenika
      </Title>
      <DynamicCardParent
        data={data}
        icon={FaClipboardCheck}
        titleField="behaviour"
        tenantColorConfig={colorConfig}
        keysToIgnore={["pupil_id", "section_id", "semester_code"]}
        twoColumnsCard={false}
        showEdit={showEdit}
        editButton={{
          onClick: onEditClick,
        }}
        extraButton={{
          onClick: (item) => fetchBehaviourGradeHistory(item.id),
          label: "Historija promjena",
          icon: FaHistory,
        }}
      />
    </Modal>
  );
};
