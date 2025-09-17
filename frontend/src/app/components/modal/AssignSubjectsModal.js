import { useState, useEffect } from "react";
import Modal from "./Modal";
import Button from "../common/Button";
import CheckboxInput from "../Input/CheckboxInput";
import Label from "../Input/Label";
import {
  FaTimes,
  FaSave,
  FaClock,
  FaBook,
  FaCheckCircle,
  FaUserTie,
} from "react-icons/fa";
import Spacer from "../common/Spacer";
import Title from "../common/Title";
import Text from "../common/Text";
import SelectInput from "../Input/SelectInput";
import SuccessModal from "./SuccessModal";
import Subtitle from "../common/Subtitle";

// SubjectCheckbox: Checkbox for a subject in SectionAssignment
function SubjectCheckbox({ name, checked, onChange, label }) {
  return (
    <label className="flex items-center gap-2 bg-white border rounded px-2 py-1 cursor-pointer">
      <CheckboxInput
        name={name}
        checked={checked}
        onChange={onChange}
        className="accent-blue-500"
      />
      {label}
    </label>
  );
}

// HomeroomCheckbox: Checkbox for homeroom teacher in SectionAssignment
function HomeroomCheckbox({ name, checked, onChange, pending }) {
  return (
    <label className="mt-8 inline-flex items-center gap-2 bg-white border rounded px-2 py-1 cursor-pointer">
      <CheckboxInput
        name={name}
        checked={checked}
        onChange={onChange}
        className="accent-blue-500"
      />
      <div className="flex items-center gap-2">
        <span>Razrednik {pending && "(na čekanju)"}</span>
        {pending ? (
          <FaClock className="text-orange-500 text-sm" />
        ) : (
          <FaUserTie className="text-blue-600 text-sm" />
        )}
      </div>
    </label>
  );
}

// Komponenta za prikaz jednog odjeljenja sa predmetima i razrednikom
function SectionAssignment({
  section,
  availableSubjects,
  assignedSubjects,
  pendingSubjects,
  assignment,
  onAvailableSubjectChange,
  onPendingSubjectChange,
  onAssignedSubjectChange,
  onHomeroomChange,
}) {
  return (
    <div className="border border-gray-200 rounded p-4 bg-gray-50 mb-4">
      <div className="font-semibold text-blue-700 mb-2">{section.name}</div>
      <div>
        {availableSubjects?.length > 0 && (
          <>
            <div className="flex items-center gap-2 mb-2">
              <FaBook className="text-green-600 text-sm" />
              <Label name={`subjects-${section.id}`} className="mb-0">
                Dostupni predmeti
              </Label>
            </div>
            <div className="flex flex-wrap gap-3 mb-4">
              {availableSubjects.map((s) => (
                <SubjectCheckbox
                  key={s.subject_code}
                  name={`subject-${section.id}-${s.subject_code}`}
                  checked={s.checked || false}
                  onChange={(e) =>
                    onAvailableSubjectChange(s.subject_code, e.target.checked)
                  }
                  label={s.subject_name}
                />
              ))}
            </div>
          </>
        )}
        {pendingSubjects?.length > 0 && (
          <>
            <div className="flex items-center gap-2 mb-2">
              <FaClock className="text-orange-500 text-sm" />
              <Label name={`subjects-${section.id}`} className="mb-0">
                Predmeti na čekanju
              </Label>
            </div>
            <div className="flex flex-wrap gap-3 mb-4">
              {pendingSubjects?.map((s) => (
                <SubjectCheckbox
                  key={s.subject_code}
                  name={`subject-${section.id}-${s.subject_code}`}
                  checked={s.checked || false}
                  onChange={(e) =>
                    onPendingSubjectChange(s.subject_code, e.target.checked)
                  }
                  label={s.subject_name}
                />
              ))}
            </div>
          </>
        )}
        {assignedSubjects?.length > 0 && (
          <>
            <div className="flex items-center gap-2 mb-2">
              <FaCheckCircle className="text-blue-600 text-sm" />
              <Label name={`subjects-${section.id}`} className="mb-0">
                Dodijeljeni predmeti
              </Label>
            </div>
            <div className="flex flex-wrap gap-3 mb-4">
              {assignedSubjects?.map((s) => (
                <SubjectCheckbox
                  key={s.subject_code}
                  name={`subject-${section.id}-${s.subject_code}`}
                  checked={s.checked || false}
                  onChange={(e) =>
                    onAssignedSubjectChange(s.subject_code, e.target.checked)
                  }
                  label={s.subject_name}
                />
              ))}
            </div>
          </>
        )}
      </div>
      <HomeroomCheckbox
        name={`homeroom-${section.id}`}
        checked={assignment?.homeroomRequest || false}
        pending={assignment?.pendingHomeroom}
        onChange={(e) => onHomeroomChange(e.target.checked)}
      />
    </div>
  );
}

// Modal za dodjelu više predmeta nastavniku po više odjeljenja
const AssignSubjectsModal = ({
  open,
  onClose,
  // array of assignment objects as described
  // Example structure:
  //   [
  //     {
  //         "section": {
  //             "id": 1,
  //             "section_code": "B",
  //             "class_code": "I",
  //             "year": "2025",
  //             "tenant_id": 1,
  //             "curriculum_code": "bos_primary_1_zdk",
  //             "name": "Odjeljenje I-B",
  //             "curriculum_name": "Nastavni plan i program za I razred osnovne škole ZDK - Bosanski jezik"
  //         },
  //         "teacher": {
  //             "id": 2,
  //             "name": "Test",
  //             "last_name": "Testovic",
  //             "email": "test@test.ba",
  //             "phone": "0644231333"
  //         },
  //         "all_subjects": [
  //             {
  //                 "subject_code": "BJZ",
  //                 "subject_name": "Bosanski jezik i književnost"
  //             },
  //         ],
  //         "assigned_subjects": [],
  //         "pending_subjects": [],
  //         "available_subjects": [
  //             {
  //                 "subject_code": "BJZ",
  //                 "subject_name": "Bosanski jezik i književnost"
  //             },
  //         ]
  //     },
  // ]
  assignmentsData,
  onSave,
  initialTeacherID,
  colorConfig,
}) => {
  const [selectedTeacherId, setSelectedTeacherId] = useState(
    "" || initialTeacherID,
  );
  // assignmentsState: { [sectionId]: { subjectCodes: [], isHomeroom: false } }
  const [assignmentsState, setAssignmentsState] = useState({});
  const [successMessage, setSuccessMessage] = useState(null);

  useEffect(() => {
    // Reset assignmentsState when selectedTeacherId changes
    setAssignmentsState({});
  }, [selectedTeacherId]);

  // Extract unique teachers from assignmentsData:
  // Use Map to automatically remove duplicate teachers (Map keys must be unique)
  // Convert teacher IDs to strings to ensure consistent key types
  // Get the Map’s values (the unique teacher objects) and turn them into an array
  const teachers = Array.from(
    new Map(
      assignmentsData?.map((a) => [
        String(a.teacher.id),
        { ...a.teacher, id: String(a.teacher.id) },
      ]),
    ).values(),
  );

  // Filter assignments for the selected teacher
  const teacherAssignments = assignmentsData?.filter(
    (a) => String(a.teacher.id) === selectedTeacherId,
  );

  const handleSectionSelect = (sectionId, checked) => {
    if (checked) {
      const sectionAssignment = teacherAssignments.find(
        (assignment) => assignment.section.id === sectionId,
      );

      // Add section to assignmentsState with default values
      setAssignmentsState((prev) => ({
        ...prev,
        [sectionId]: {
          availableSubjects: sectionAssignment?.available_subjects?.map(
            (subject) => ({ ...subject, checked: false }),
          ),
          pendingSubjects: sectionAssignment?.pending_subjects?.map(
            (subject) => ({ ...subject, checked: true }),
          ),
          assignedSubjects: sectionAssignment?.assigned_subjects?.map(
            (subject) => ({ ...subject, checked: true }),
          ),
          isHomeroom: sectionAssignment?.is_homeroom_teacher,
          pendingHomeroom: sectionAssignment?.is_pending_homeroom_teacher,
          homeroomRequest:
            sectionAssignment?.is_homeroom_teacher ||
            sectionAssignment?.is_pending_homeroom_teacher,
          inviteIndexID: sectionAssignment?.invite_index_id,
        },
      }));
    } else {
      // Remove section from assignmentsState
      setAssignmentsState((prev) => {
        const newState = { ...prev };
        delete newState[sectionId];
        return newState;
      });
    }
  };

  const handleAvailableSubjectChange = (sectionId, subjectCode, checked) => {
    setAssignmentsState((prev) => {
      const currentSection = prev[sectionId];

      // Find and update the specific subject in availableSubjects array
      const updatedAvailableSubjects = currentSection.availableSubjects.map(
        (subject) =>
          subject.subject_code === subjectCode
            ? { ...subject, checked: checked }
            : subject,
      );

      return {
        ...prev,
        [sectionId]: {
          ...currentSection,
          availableSubjects: updatedAvailableSubjects,
        },
      };
    });
  };

  const handlePendingSubjectChange = (sectionId, subjectCode, checked) => {
    setAssignmentsState((prev) => {
      const currentSection = prev[sectionId];

      // Find and update the specific subject in pendingSubjects array
      const updatedPendingSubjects = currentSection.pendingSubjects.map(
        (subject) =>
          subject.subject_code === subjectCode
            ? { ...subject, checked: checked }
            : subject,
      );

      return {
        ...prev,
        [sectionId]: {
          ...currentSection,
          pendingSubjects: updatedPendingSubjects,
        },
      };
    });
  };

  const handleAssignedSubjectChange = (sectionId, subjectCode, checked) => {
    setAssignmentsState((prev) => {
      const currentSection = prev[sectionId];

      // Find and update the specific subject in assignedSubjects array
      const updatedAssignedSubjects = currentSection.assignedSubjects.map(
        (subject) =>
          subject.subject_code === subjectCode
            ? { ...subject, checked: checked }
            : subject,
      );

      return {
        ...prev,
        [sectionId]: {
          ...currentSection,
          assignedSubjects: updatedAssignedSubjects,
        },
      };
    });
  };

  const handleHomeroomChange = (sectionId, checked) => {
    setAssignmentsState((prev) => {
      // Get current section state or initialize if it doesn't exist
      const currentSection = prev[sectionId] || {
        subjectCodes: [],
        homeroomRequest: false,
      };

      // Update isHomeroom based on the checkbox state
      return {
        ...prev,
        [sectionId]: { ...currentSection, homeroomRequest: checked },
      };
    });
  };

  const handleSave = () => {
    // Return assignments for selected teacher only
    setSelectedTeacherId("");
    if (onSave) onSave(assignmentsState, selectedTeacherId);
    setSuccessMessage("Predmeti i razredništvo su uspješno dodijeljeni.");
  };

  return (
    <>
      <Modal open={open} onClose={onClose}>
        <Title colorConfig={colorConfig}>
          Dodijeli predmete i razredništvo
        </Title>
        {teachers?.length > 0 ? (
          <>
            <Text className="mb-4 text-justify" textSize="md">
              <span className="font-bold">Za svako odjeljenje možete:</span>
              <ol className="mt-4 mb-4 list-decimal list-inside space-y-2">
                <li>
                  Dodijeliti nove predmete za predavanje (označite checkbox u
                  sekciji <i>„Dostupni predmeti“</i>)
                </li>
                <li>
                  Otkazati pozive koji su na čekanju (poništite checkbox u
                  sekciji <i>„Predmeti na čekanju“</i>)
                </li>
                <li>
                  Ukloniti predmete koje nastavnik već predaje (poništite
                  checkbox u sekciji <i>„Dodijeljeni predmeti“</i>)
                </li>
                <li>Postaviti nastavnika za razrednika</li>
              </ol>
              <p className="text-sm text-gray-600">
                * Napomena: Sekcije (Dostupni predmeti, Predmeti na čekanju,
                Dodijeljeni predmeti) se prikazuju samo ako postoje predmeti za
                njih.
              </p>
            </Text>
            <div className="mb-4">
              <SelectInput
                name="teacher"
                value={selectedTeacherId || ""}
                onChange={(e) => setSelectedTeacherId(e.target.value)}
                options={teachers.map((t) => ({
                  value: t.id,
                  label: t.name + " " + t.last_name,
                }))}
                placeholder="Izaberi nastavnika"
              />
            </div>
          </>
        ) : (
          <Subtitle colorConfig={colorConfig} showLine={false}>
            Za zaduživanje nastavnika potrebno je da škola ima barem jedno
            aktivno odjeljenje.
          </Subtitle>
        )}
        <Spacer>
          {selectedTeacherId && teacherAssignments.length > 0 && (
            <div className="mb-4">
              <Label className="mb-2">Odjeljenja</Label>
              <div className="flex flex-wrap gap-4">
                {teacherAssignments.map((a) => (
                  <label
                    key={a.section.id}
                    className="flex items-center gap-2 bg-white border rounded px-2 py-1 cursor-pointer"
                  >
                    <CheckboxInput
                      name={`section-select-${a.section.id}`}
                      checked={a.section.id in assignmentsState}
                      onChange={(e) =>
                        handleSectionSelect(a.section.id, e.target.checked)
                      }
                      className="accent-blue-500"
                    />
                    {a.section.name}
                  </label>
                ))}
              </div>
            </div>
          )}
          {selectedTeacherId &&
            Object.keys(assignmentsState).map((sectionId) => {
              const a = teacherAssignments.find(
                (x) => x.section.id === parseInt(sectionId),
              );
              if (!a) return null;
              return (
                <SectionAssignment
                  key={a.section.id}
                  section={a.section}
                  availableSubjects={
                    assignmentsState[a.section.id].availableSubjects
                  }
                  pendingSubjects={
                    assignmentsState[a.section.id].pendingSubjects
                  }
                  assignedSubjects={
                    assignmentsState[a.section.id].assignedSubjects
                  }
                  assignment={assignmentsState[a.section.id]}
                  onAvailableSubjectChange={(subjectCode, checked) =>
                    handleAvailableSubjectChange(
                      a.section.id,
                      subjectCode,
                      checked,
                    )
                  }
                  onPendingSubjectChange={(subjectCode, checked) =>
                    handlePendingSubjectChange(
                      a.section.id,
                      subjectCode,
                      checked,
                    )
                  }
                  onAssignedSubjectChange={(subjectCode, checked) =>
                    handleAssignedSubjectChange(
                      a.section.id,
                      subjectCode,
                      checked,
                    )
                  }
                  onHomeroomChange={(checked) =>
                    handleHomeroomChange(a.section.id, checked)
                  }
                />
              );
            })}
          <div className="flex justify-end gap-3 mt-6">
            <Button
              colorConfig={colorConfig}
              icon={FaTimes}
              onClick={onClose}
              color="ternary"
            >
              Završi
            </Button>
            {teachers.length > 0 && (
              <Button
                icon={FaSave}
                onClick={handleSave}
                color="secondary"
                disabled={
                  !selectedTeacherId ||
                  Object.keys(assignmentsState).length === 0
                }
                colorConfig={colorConfig}
              >
                Sačuvaj
              </Button>
            )}
          </div>
        </Spacer>
      </Modal>
      {successMessage && (
        <SuccessModal
          colorConfig={colorConfig}
          onClose={() => setSuccessMessage(null)}
        >
          {successMessage}
        </SuccessModal>
      )}
    </>
  );
};

export default AssignSubjectsModal;
