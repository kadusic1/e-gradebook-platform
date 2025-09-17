"use client";
import { CardContainer } from "../dynamic_card/DynamicCard";
import Subtitle from "../common/Subtitle";
import Text from "../common/Text";
import Spacer from "../common/Spacer";
import { getColor } from "../colors/colors";
import {
  FaFileAlt,
  FaMicrophone,
  FaPencilAlt,
  FaTrash,
  FaEdit,
  FaPlus,
  FaLock,
  FaEye,
  FaEyeSlash,
  FaClipboardCheck,
  FaBook,
  FaUser,
  FaUserTimes,
} from "react-icons/fa";
import EmptyStateContainerless from "../dynamic_card/EmptyStateContainerless";
import TextInput from "./TextInput";
import { useMemo, useState } from "react";
import CreateUpdateModal from "../modal/CreateUpdateModal";
import EmptyState from "../dynamic_card/EmptyState";
import {
  formatDateToDDMMYYYY,
  formatToFullDateTime,
} from "@/app/util/date_util";
import ConfirmModal from "../modal/ConfirmModal";
import ErrorModal from "../modal/ErrorModal";
import Modal from "../modal/Modal";

const getInitials = (name) => {
  return (
    name
      .split(" ")
      .map((word) => word.charAt(0).toUpperCase())
      .join(".") + "."
  );
};

const DeleteIcon = ({ onClick }) => {
  return (
    <span className={`inline-flex hover:cursor-pointer mr-2`} onClick={onClick}>
      <FaTrash className="w-6 h-6 text-red-500" />
    </span>
  );
};

const EditIcon = ({ onClick }) => {
  return (
    <span className={`inline-flex hover:cursor-pointer mr-2`} onClick={onClick}>
      <FaEdit className="w-6 h-6 text-green-500" />
    </span>
  );
};

const AddIcon = ({ color, onClick = () => {} }) => {
  return (
    <span className={`inline-flex hover:cursor-pointer mr-2`} onClick={onClick}>
      <FaPlus className={`w-6 h-6 ${color}`} />
    </span>
  );
};

const FinalizeIcon = ({ color, onClick = () => {} }) => {
  return (
    <span className={`inline-flex hover:cursor-pointer mr-2`} onClick={onClick}>
      <FaLock className={`w-6 h-6 ${color}`} />
    </span>
  );
};

const ShowDeletedGradesIcon = ({ color, onClick = () => {} }) => {
  return (
    <span className={`inline-flex hover:cursor-pointer mr-2`} onClick={onClick}>
      <FaEye className={`w-6 h-6 ${color}`} />
    </span>
  );
};

const HideDeletedGradesIcon = ({ color, onClick = () => {} }) => {
  return (
    <span className={`inline-flex hover:cursor-pointer mr-2`} onClick={onClick}>
      <FaEyeSlash className={`w-6 h-6 ${color}`} />
    </span>
  );
};

const gradeTypeConfig = {
  exam: {
    icon: FaFileAlt,
    label: "Kontrolni rad",
    colorKey: "quaternary",
  },
  oral: {
    icon: FaMicrophone,
    label: "Usmeni odgovor",
    colorKey: "secondary",
  },
  written_assignment: {
    icon: FaPencilAlt,
    label: "Pismena zadaća",
    colorKey: "ternary",
  },
  final: {
    icon: FaLock,
    label: "Zaključna ocjena",
    colorKey: "primary",
  },
  behaviour: {
    icon: FaClipboardCheck,
    label: "Vladanje",
    colorKey: "primary",
  },
};

export const GradeType = ({
  gradeObject,
  onEdit,
  onDelete,
  mode,
  finalGrade,
  archived = 0,
  onEditHistoryClick,
  colorConfig,
  behaviourMode = false,
}) => {
  const config = behaviourMode
    ? gradeTypeConfig.behaviour
    : gradeTypeConfig[gradeObject.type];

  const color = getColor(config.colorKey, "text", colorConfig);
  if (!config) return null;
  const Icon = config.icon;

  return (
    <Text
      textSize="lg"
      className={`${gradeObject.is_deleted ? "line-through" : ""} animate-fadeIn`}
    >
      <span className={`inline-flex mr-2 ${color}`}>
        <Icon className="w-6 h-6" />
      </span>
      {mode !== "view" && (
        <span className="font-semibold">{config.label}:</span>
      )}{" "}
      <span className={`${color} font-bold`}>
        {behaviourMode ? gradeObject.behaviour : gradeObject.grade}
      </span>{" "}
      <span className="ml-1 mr-1">
        {behaviourMode
          ? formatToFullDateTime(gradeObject.date).split(" ")[0]
          : formatDateToDDMMYYYY(gradeObject.grade_date)}
      </span>
      <span
        className="ml-1 mr-4 font-bold italic"
        title={gradeObject.signature}
      >
        {behaviourMode
          ? getInitials(gradeObject.behaviour_determined_by_teacher)
          : getInitials(gradeObject.signature)}
      </span>
      {mode === "teacher" &&
        (!finalGrade || gradeObject.type === "final") &&
        archived == 0 &&
        !gradeObject.is_deleted && (
          <>
            {onEdit && <EditIcon onClick={onEdit} />}
            {onDelete && <DeleteIcon onClick={onDelete} />}
          </>
        )}
      {gradeObject.is_edited && (
        <span
          className={`italic ${!gradeObject.valid_until ? "hover:cursor-pointer text-blue-700" : ""}`}
          onClick={onEditHistoryClick}
        >
          {gradeObject.valid_until
            ? `(Važilo do ${formatToFullDateTime(gradeObject.valid_until)})`
            : `(Uređeno)`}
        </span>
      )}
    </Text>
  );
};

const GradeInput = ({
  pupil,
  grades,
  averageGrade,
  colorConfig,
  onAddClick,
  onEditClick,
  onDeleteClick,
  mode,
  subject,
  onFinalClick,
  onFinalEditClick,
  archived = 0,
  getGradeEditHistory,
  ...props
}) => {
  let modeOption = mode;
  if (pupil?.unenrolled) {
    modeOption = "pupil_unenrolled";
  }

  const primaryTextColor = getColor("primary", "text", colorConfig);
  const secondaryTextColor = getColor("secondary", "text", colorConfig);
  const quaternaryTextColor = getColor("quaternary", "text", colorConfig);

  const finalGrade =
    grades?.find((grade) => grade.type === "final" && !grade.is_deleted) ||
    null;
  const [showDeleted, setShowDeleted] = useState(false);

  const filteredGrades = useMemo(() => {
    return showDeleted ? grades : grades?.filter((grade) => !grade.is_deleted);
  }, [showDeleted]);

  const nonDeletedGrades = grades?.filter((grade) => !grade.is_deleted);

  let cellTitleIcon;
  if (modeOption === "teacher") {
    cellTitleIcon = FaUser;
  } else if (modeOption === "pupil_unenrolled") {
    cellTitleIcon = FaUserTimes;
  } else {
    cellTitleIcon = FaBook;
  }

  return (
    <CardContainer className="min-h-64" {...props}>
      <div className="flex justify-between items-center mb-6">
        <Subtitle
          icon={cellTitleIcon}
          colorConfig={colorConfig}
          showLine={false}
          textColor={modeOption === "pupil_unenrolled" ? "text-red-500" : null}
        >
          {mode === "teacher" ? (
            <>
              {pupil.name} {pupil.last_name}
            </>
          ) : (
            subject.subject_name
          )}
          {averageGrade && averageGrade > 0.0 && (
            <span> ({averageGrade.toFixed(2)})</span>
          )}
        </Subtitle>
        <div className="flex gap-2">
          {modeOption === "teacher" && !finalGrade && (
            <>
              {nonDeletedGrades?.length > 0 && (
                <FinalizeIcon color={primaryTextColor} onClick={onFinalClick} />
              )}
              {archived == 0 && (
                <AddIcon color={quaternaryTextColor} onClick={onAddClick} />
              )}
            </>
          )}
          <span
            className="inline-flex hover:cursor-pointer"
            title={
              showDeleted
                ? "Sakrij izbrisane ocjene"
                : "Prikaži izbrisane ocjene"
            }
          >
            {showDeleted ? (
              <ShowDeletedGradesIcon
                color={secondaryTextColor}
                onClick={() => setShowDeleted(false)}
              />
            ) : (
              <HideDeletedGradesIcon
                color={secondaryTextColor}
                onClick={() => setShowDeleted(true)}
              />
            )}
          </span>
        </div>
      </div>
      {filteredGrades?.length > 0 ? (
        <Spacer>
          {finalGrade && (
            <div>
              <GradeType
                gradeObject={finalGrade}
                onEdit={() => onFinalEditClick(finalGrade)}
                onDelete={() => onDeleteClick(finalGrade)}
                mode={modeOption}
                archived={archived}
                onEditHistoryClick={() => getGradeEditHistory(finalGrade.id)}
                colorConfig={colorConfig}
              />
            </div>
          )}
          {filteredGrades.map((grade, idx) => {
            if (grade.type === "final" && !grade.is_deleted) return null;
            const config = gradeTypeConfig[grade.type];
            if (!config) return null;
            return (
              <GradeType
                key={idx}
                gradeObject={grade}
                onEdit={() => onEditClick(grade)}
                onDelete={() => onDeleteClick(grade)}
                mode={modeOption}
                finalGrade={finalGrade}
                archived={archived}
                onEditHistoryClick={() => getGradeEditHistory(grade.id)}
                colorConfig={colorConfig}
              />
            );
          })}
        </Spacer>
      ) : (
        <EmptyStateContainerless message="Nema ocjena za prikaz." />
      )}
    </CardContainer>
  );
};

export const GradeInputParent = ({
  items,
  accessToken,
  colorConfig,
  section,
  setItems,
  selectedSubject,
  mode = "teacher",
  semester,
  archived = 0,
}) => {
  const gradeFields = [
    {
      label: "Tip ocjene",
      name: "type",
      type: "select",
      placeholder: "Odaberite tip ocjene",
      options: [
        { value: "exam", label: "Kontrolni rad" },
        { value: "oral", label: "Usmeni odgovor" },
        { value: "written_assignment", label: "Pismena zadaća" },
      ],
    },
    {
      label: "Ocjena",
      name: "grade",
      type: "positive-number",
      placeholder: "Unesite ocjenu",
    },
    {
      label: "Datum ocjene",
      name: "grade_date",
      type: "date",
      min: semester.start_date,
      max: semester.end_date,
    },
  ];

  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showFinalCreateModal, setShowFinalCreateModal] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [showFinalEditModal, setShowFinalEditModal] = useState(false);
  const [selectedPupil, setSelectedPupil] = useState(null);
  const [selectedGrade, setSelectedGrade] = useState(null);
  const [searchTerm, setSearchTerm] = useState("");
  const [errorMessage, setErrorMessage] = useState("");
  const [historyEditGrades, setHistoryEditGrades] = useState(null);
  const filteredItems = useMemo(() => {
    if (searchTerm === "" || searchTerm === null || searchTerm === undefined) {
      return items;
    }
    if (mode === "teacher") {
      return items?.filter(
        (grade) =>
          grade?.pupil?.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
          grade?.pupil?.last_name
            .toLowerCase()
            .includes(searchTerm.toLowerCase()),
      );
    } else {
      return items?.filter((grade) =>
        grade?.subject?.subject_name
          .toLowerCase()
          .includes(searchTerm.toLowerCase()),
      );
    }
  }, [searchTerm, items]);

  const handleGradeSave = async (gradeData, selectedPupil) => {
    if (mode !== "teacher") return;
    const gradePayload = {
      type: gradeData.type,
      pupil_id: selectedPupil.id,
      section_id: section.id,
      subject_code: selectedSubject,
      grade: parseInt(gradeData.grade),
      grade_date: gradeData.grade_date,
      semester_code: semester.semester_code,
    };
    setSelectedPupil(null);
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/teacher/create_grade/${section.tenant_id}`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
          body: JSON.stringify(gradePayload),
        },
      );

      if (response.ok) {
        const createdGrade = await response.json();
        setItems((prevGrades) =>
          prevGrades.map((grade) =>
            grade.pupil.id === createdGrade.pupil.id ? createdGrade : grade,
          ),
        );
      } else {
        const errorData = await response.text();
        setErrorMessage(errorData);
      }
    } catch (error) {
      console.error("Error saving grade:", error);
    }
  };

  const handleGradeEdit = async (gradeData, selectedPupil) => {
    if (mode !== "teacher") return;
    const gradePayload = {
      id: gradeData.id,
      type: gradeData.type,
      pupil_id: selectedPupil.id,
      section_id: section.id,
      subject_code: selectedSubject,
      grade: parseInt(gradeData.grade),
      grade_date: gradeData.grade_date,
      semester_code: semester.semester_code,
    };
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/teacher/update_grade/${section.tenant_id}`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
          body: JSON.stringify(gradePayload),
        },
      );

      if (response.ok) {
        const updatedGradeData = await response.json();
        setItems((prevGrades) =>
          prevGrades.map((grade) =>
            grade.pupil.id === updatedGradeData.pupil.id
              ? updatedGradeData
              : grade,
          ),
        );
      } else {
        const errorData = await response.text();
        setErrorMessage(errorData);
      }
    } catch (error) {
      console.error("Error saving grade:", error);
    }
  };

  const handleGradeDelete = async (grade) => {
    if (mode !== "teacher") return;
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/teacher/delete_grade/${section.tenant_id}`,
        {
          method: "DELETE",
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
          body: JSON.stringify(grade),
        },
      );
      const deletedGradeItems = await response.json();
      if (response.ok) {
        setItems((prevGrades) =>
          prevGrades.map((grade) =>
            grade.pupil.id === deletedGradeItems.pupil.id
              ? deletedGradeItems
              : grade,
          ),
        );
      }
    } catch (error) {
      console.error("Error deleting grade:", error);
    }
  };

  const getGradeEditHistory = async (gradeID) => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/pupil/grade_edit_history/${section.tenant_id}/${gradeID}`,
        {
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
        },
      );

      if (response.ok) {
        const grades = await response.json();
        setHistoryEditGrades(grades);
      }
    } catch (error) {
      console.error("Error getting grade history:", error);
    }
  };

  return (
    <>
      {historyEditGrades && (
        <Modal
          colorConfig={colorConfig}
          onClose={() => setHistoryEditGrades(null)}
        >
          <Subtitle colorConfig={colorConfig}>Historija ocjene</Subtitle>
          <Spacer className="mt-6">
            {historyEditGrades.length === 0 ? (
              <EmptyStateContainerless message="Nema izmjena za ovu ocjenu." />
            ) : (
              historyEditGrades.map((grade, idx) => {
                // Use GradeType for all other types
                const config = gradeTypeConfig[grade.type];
                if (!config) return null;
                return (
                  <GradeType
                    key={idx}
                    gradeObject={grade}
                    mode={mode}
                    archived={archived}
                    colorConfig={colorConfig}
                  />
                );
              })
            )}
          </Spacer>
        </Modal>
      )}
      {errorMessage && mode === "teacher" && (
        <ErrorModal
          onClose={() => setErrorMessage(null)}
          colorConfig={colorConfig}
        >
          {errorMessage}
        </ErrorModal>
      )}
      {showFinalCreateModal && mode === "teacher" && archived == 0 && (
        <CreateUpdateModal
          title="Zaključivanje ocjene"
          fields={gradeFields?.filter((field) => field.name !== "type")}
          onClose={() => {
            setShowFinalCreateModal(false);
            setSelectedPupil(null);
          }}
          onSave={(data) => {
            const finalGrade = {
              ...data,
              type: "final",
            };
            handleGradeSave(finalGrade, selectedPupil);
          }}
          colorConfig={colorConfig}
        />
      )}
      {showCreateModal && mode === "teacher" && archived == 0 && (
        <CreateUpdateModal
          title="Unos ocjene"
          fields={gradeFields}
          onClose={() => {
            setShowCreateModal(false);
            setSelectedPupil(null);
          }}
          onSave={(data) => {
            handleGradeSave(data, selectedPupil);
          }}
          colorConfig={colorConfig}
        />
      )}
      {showEditModal && mode === "teacher" && archived == 0 && (
        <CreateUpdateModal
          title="Uređivanje ocjene"
          fields={gradeFields}
          onClose={() => {
            setShowEditModal(false);
            setSelectedPupil(null);
            setSelectedGrade(null);
          }}
          onSave={(data) => {
            handleGradeEdit(data, selectedPupil);
            setSelectedPupil(null);
            setSelectedGrade(null);
          }}
          colorConfig={colorConfig}
          initialValues={selectedGrade}
        />
      )}
      {showFinalEditModal && mode === "teacher" && archived == 0 && (
        <CreateUpdateModal
          title="Uređivanje ocjene"
          fields={gradeFields?.filter((field) => field.name !== "type")}
          onClose={() => {
            setShowFinalEditModal(false);
            setSelectedPupil(null);
            setSelectedGrade(null);
          }}
          onSave={(data) => {
            handleGradeEdit(data, selectedPupil);
            setSelectedPupil(null);
            setSelectedGrade(null);
          }}
          colorConfig={colorConfig}
          initialValues={selectedGrade}
        />
      )}
      {showDeleteModal && mode === "teacher" && archived == 0 && (
        <ConfirmModal
          title="Brisanje ocjene"
          onClose={() => {
            setShowDeleteModal(false);
            setSelectedGrade(null);
            setSelectedPupil(null);
          }}
          onConfirm={() => {
            handleGradeDelete(selectedGrade);
            setShowDeleteModal(false);
            setSelectedGrade(null);
            setSelectedPupil(null);
          }}
          colorConfig={colorConfig}
        >
          Da li ste sigurni da želite izbrisati ocjenu{" "}
          <span className="font-bold">
            {selectedGrade?.grade} (
            {formatDateToDDMMYYYY(selectedGrade?.grade_date)})
          </span>{" "}
          iz predmeta{" "}
          <span className="font-bold">{selectedGrade?.subject_name}</span> za
          učenika{" "}
          <span className="font-bold">
            {selectedPupil?.name} {selectedPupil?.last_name}
          </span>
          ?
        </ConfirmModal>
      )}
      <TextInput
        name="search"
        placeholder={
          mode === "teacher"
            ? "Pretraži učenike po imenu i/ili prezimenu"
            : "Pretraži predmete"
        }
        className="mb-4 border border-gray-500 rounded-md"
        onChange={(e) => setSearchTerm(e.target.value)}
      />
      <Spacer>
        {filteredItems?.length === 0 && (
          <EmptyState
            className="mt-8"
            message={
              mode === "teacher"
                ? "Nema učenika koji odgovaraju pretrazi."
                : "Nema predmeta koji odgovaraju pretrazi."
            }
          />
        )}
        <div className="grid grid-cols-1 xl:grid-cols-2 gap-4 md:gap-6">
          {filteredItems?.map((item, idx) => (
            <GradeInput
              key={idx}
              colorConfig={colorConfig}
              pupil={item?.pupil}
              averageGrade={item?.average_grade}
              grades={item?.grades}
              onAddClick={() => {
                if (mode !== "teacher") return;
                setShowCreateModal(true);
                setSelectedPupil(item?.pupil);
              }}
              onDeleteClick={(selectedGrade) => {
                if (mode !== "teacher") return;
                setSelectedPupil(item?.pupil);
                setShowDeleteModal(true);
                setSelectedGrade(selectedGrade);
              }}
              onEditClick={(selectedGrade) => {
                if (mode !== "teacher") return;
                setSelectedPupil(item?.pupil);
                setShowEditModal(true);
                setSelectedGrade(selectedGrade);
              }}
              onFinalClick={() => {
                if (mode !== "teacher") return;
                setShowFinalCreateModal(true);
                setSelectedPupil(item?.pupil);
              }}
              onFinalEditClick={(selectedGrade) => {
                if (mode !== "teacher") return;
                setSelectedPupil(item?.pupil);
                setShowFinalEditModal(true);
                setSelectedGrade(selectedGrade);
              }}
              mode={mode}
              subject={item?.subject}
              archived={archived}
              getGradeEditHistory={getGradeEditHistory}
            />
          ))}
        </div>
      </Spacer>
    </>
  );
};
