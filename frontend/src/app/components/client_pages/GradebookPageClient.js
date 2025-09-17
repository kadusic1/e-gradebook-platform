"use client";
import BackButton from "../common/BackButton";
import Title from "../common/Title";
import SelectInput from "../Input/SelectInput";
import { useState, useEffect } from "react";
import EmptyState from "../dynamic_card/EmptyState";
import { GradeInputParent } from "../Input/GradeInput";
import { FaBook } from "react-icons/fa";
import DynamicTab from "../dynamic_card/DynamicTab";
import { formatDateToDDMMYYYY } from "@/app/util/date_util";
import { FaCalendarAlt } from "react-icons/fa";

const PageHeader = ({ section, setShowGradebookPage, colorConfig }) => {
  return (
    <>
      <Title icon={FaBook} colorConfig={colorConfig}>
        Ocjene - {section?.name}
      </Title>
      <div className="flex justify-end mb-4">
        <BackButton
          colorConfig={colorConfig}
          onClick={() => setShowGradebookPage(false)}
        />
      </div>
    </>
  );
};

export default function GradebookPageClient({
  colorConfig,
  setShowGradebookPage,
  accessToken,
  section,
  mode = "teacher",
  archived = 0,
}) {
  const [selectedSubject, setSelectedSubject] = useState("");
  const [teacherSubjects, setTeacherSubjects] = useState();
  const [pupilCountInSection, setPupilCountInSection] = useState(0);
  const [sectionSemesters, setSectionSemesters] = useState([]);
  const [activeTab, setActiveTab] = useState(0);
  const [items, setItems] = useState();

  useEffect(() => {
    setSelectedSubject("");
    setItems(null);
  }, [activeTab]);

  const fetchItems = async (subjectCode) => {
    const semesterCode = sectionSemesters[activeTab]?.semester_code;
    try {
      let urlPath;
      if (mode === "teacher") {
        urlPath = `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/teacher/section_grades_for_subject/${section?.tenant_id}/${section?.id}/${subjectCode}/${semesterCode}`;
      } else {
        urlPath = `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/pupil/section_grades_for_pupil/${section?.tenant_id}/${section?.id}/${semesterCode}`;
      }
      const response = await fetch(urlPath, {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      });
      if (response.ok) {
        const data = await response.json();
        setItems(data);
      }
    } catch (error) {
      console.error("Error fetching grades:", error);
    }
  };

  const handleSubjectChange = (event) => {
    if (mode !== "teacher") return;
    const value = event.target.value;
    if (value === "") {
      setSelectedSubject("");
      return;
    }
    setSelectedSubject(value);
  };

  const fetchGradebookMetadata = async () => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/teacher/gradebook_metadata/${section?.tenant_id}/${section?.id}`,
        {
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        },
      );
      if (!response.ok) {
        throw new Error("Failed to fetch gradebook metadata");
      }
      const data = await response.json();
      if (mode == "teacher") {
        setTeacherSubjects(data?.subjects || []);
        setPupilCountInSection(data.pupil_count);
      }
      setSectionSemesters(data?.section_semesters || []);
    } catch (error) {
      console.error("Error fetching gradebook metadata:", error);
    }
  };

  useEffect(() => {
    if (selectedSubject && mode == "teacher") {
      fetchItems(selectedSubject);
    }
  }, [selectedSubject]);

  useEffect(() => {
    fetchGradebookMetadata();
  }, []);

  useEffect(() => {
    if (mode === "pupil" && sectionSemesters.length > 0) {
      fetchItems(null);
    }
  }, [sectionSemesters, activeTab]);

  if (pupilCountInSection === 0 && mode === "teacher") {
    return (
      <>
        <PageHeader
          section={section}
          setShowGradebookPage={setShowGradebookPage}
          colorConfig={colorConfig}
        />
        <EmptyState message="Ovaj razred nema učenika." className="mt-12" />
      </>
    );
  }

  const GradebookTab = ({ semester }) => {
    return (
      <>
        {selectedSubject || mode === "pupil" ? (
          <div>
            <GradeInputParent
              colorConfig={colorConfig}
              accessToken={accessToken}
              section={section}
              items={items}
              setItems={setItems}
              selectedSubject={selectedSubject}
              mode={mode}
              semester={semester}
              archived={archived}
            />
          </div>
        ) : (
          <EmptyState
            className="mt-8"
            message="Molimo vas da odaberete predmet kako biste vidjeli ocjene."
          />
        )}
      </>
    );
  };

  return (
    <>
      <DynamicTab
        title={`Ocjene - ${section?.name}`}
        titleIcon={FaBook}
        activeTab={activeTab}
        setActiveTab={setActiveTab}
        colorConfig={colorConfig}
        aboveContent={
          <>
            <div className="flex justify-end mb-4">
              <BackButton
                colorConfig={colorConfig}
                onClick={() => setShowGradebookPage(false)}
              />
            </div>
            {mode === "teacher" && (
              <SelectInput
                value={selectedSubject}
                options={teacherSubjects?.map((subject) => ({
                  value: subject.subject_code,
                  label: subject.subject_name,
                }))}
                className="mb-4 border border-gray-500 rounded-md"
                name="subjectSelect"
                placeholder="Odaberite predmet"
                onChange={handleSubjectChange}
              />
            )}
          </>
        }
        childrenTabs={sectionSemesters.map((semester, idx) => ({
          label: `${semester?.semester_name} (${formatDateToDDMMYYYY(semester?.start_date)} - ${formatDateToDDMMYYYY(semester?.end_date)})`,
          content: <GradebookTab key={idx} semester={semester} />,
          icon: FaCalendarAlt,
        }))}
      />
    </>
  );
}
