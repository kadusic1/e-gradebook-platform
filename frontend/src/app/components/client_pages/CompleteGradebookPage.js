import BackButton from "../common/BackButton";
import Title from "../common/Title";
import { PDFButton } from "../common/PDFButton";
import { useRef, useEffect, useState } from "react";
import {
  FaRegListAlt,
  FaUserFriends,
  FaCalendarDay,
  FaClock,
  FaBook,
  FaUser,
  FaFileAlt,
  FaCalendarAlt,
  FaCheckCircle,
  FaTimesCircle,
  FaUserTimes,
  FaClipboardList,
  FaHourglassHalf,
} from "react-icons/fa";
import {
  Table,
  TableBody,
  TableHead,
  TableCell,
  TableHeader,
  TableRow,
} from "../table/TableComponents";
import { handleDownloadPDF } from "@/app/util/pdf_util";
import EmptyStateContainerless from "../dynamic_card/EmptyStateContainerless";
import ScheduleTable from "../table/ScheduleTable";
import {
  formatDateToDDMMYYYY,
  formatToFullDateTime,
} from "@/app/util/date_util";
import Subtitle from "../common/Subtitle";
import {
  mapValueToBosnian,
  getValueColor,
} from "../dynamic_card/DynamicCardParent";
import { GradeType } from "../Input/GradeInput";

export const CompleteGradebookPage = ({
  onBack,
  colorConfig,
  section,
  accessToken,
}) => {
  const pdfRef = useRef(null);
  const [gradebookData, setGradebookData] = useState([]);
  const [pdfLoading, setPdfLoading] = useState(false);

  const getCompleteGradeBookData = async () => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/teacher/complete_gradebook/${section.tenant_id}/${section.id}`,
        {
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        },
      );
      if (response.ok) {
        const data = await response.json();
        setGradebookData(data);
      }
    } catch (error) {
      console.error("Error fetching complete gradebook data:", error);
    }
  };

  useEffect(() => {
    getCompleteGradeBookData();
  }, []);

  const downloadCompleteGradebookPDF = () => {
    handleDownloadPDF({
      filename: "dnevnik.pdf",
      landscape: true,
      showPageNumbers: true,
      setPdfLoading: setPdfLoading,
      pdfElementRef: pdfRef,
      additionalStyles: `
        body {
          padding: 5px;
        }
      `,
    });
  };

  const chunkArray = (arr, size) =>
    Array.from({ length: Math.ceil(arr.length / size) }, (_, i) =>
      arr.slice(i * size, i * size + size),
    );

  const SUBJECT_CHUNK_SIZE = 4;

  return (
    <>
      <Title colorConfig={colorConfig} icon={FaRegListAlt}>
        {section.name} - Dnevnik
      </Title>
      <div className="flex justify-end mb-4">
        <BackButton colorConfig={colorConfig} onClick={onBack} />
        <PDFButton
          colorConfig={colorConfig}
          pdfLoading={pdfLoading}
          handleDownloadPDF={downloadCompleteGradebookPDF}
          className="ml-2"
        />
      </div>
      <div ref={pdfRef}>
        <div className="new-page">
          <Title icon={FaUserFriends} colorConfig={colorConfig}>
            Učenici
          </Title>
          {gradebookData?.pupils?.length > 0 ? (
            <Table>
              <TableHead>
                <TableRow>
                  <TableHeader>R.B.</TableHeader>
                  <TableHeader>Prezime i ime učenika</TableHeader>
                  <TableHeader>Ime staratelja</TableHeader>
                  <TableHeader>Telefon</TableHeader>
                  <TableHeader>Email</TableHeader>
                  <TableHeader>Vjeronauka</TableHeader>
                  <TableHeader>Vozar</TableHeader>
                </TableRow>
              </TableHead>
              <TableBody>
                {gradebookData?.pupils?.map((pupil, idx) => (
                  <TableRow key={idx}>
                    <TableCell>{idx + 1}</TableCell>
                    <TableCell className="flex items-center">
                      {pupil.unenrolled ? (
                        <FaUserTimes className="mr-2 text-red-500" size={12} />
                      ) : (
                        <FaUser className="mr-2 text-green-500" size={12} />
                      )}
                      {pupil.last_name} {pupil.name}
                    </TableCell>
                    <TableCell>{pupil.guardian_name}</TableCell>
                    <TableCell>{pupil.phone_number}</TableCell>
                    <TableCell>{pupil.email}</TableCell>
                    <TableCell>{mapValueToBosnian(pupil.religion)}</TableCell>
                    <TableCell className={getValueColor(pupil.is_commuter)}>
                      {pupil.is_commuter}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          ) : (
            <EmptyStateContainerless
              message="Nema učenika u ovom odjeljenju."
              className="mt-12"
            />
          )}
        </div>
        <div className="new-page">
          <Title icon={FaCalendarDay} colorConfig={colorConfig}>
            Raspored časova
          </Title>
          {gradebookData?.schedule_history?.length > 0 ? (
            gradebookData?.schedule_history?.map((schedule, idx) => (
              <div key={idx} className="mb-5">
                <Subtitle
                  colorConfig={colorConfig}
                  showLine={false}
                  textSize="text-lg"
                  icon={FaClock}
                >
                  Raspored kreiran:{" "}
                  {formatToFullDateTime(schedule[0].created_at)}
                </Subtitle>
                <ScheduleTable
                  colorConfig={colorConfig}
                  readOnly={true}
                  initialScheduleGroups={schedule}
                  showButtons={false}
                />
              </div>
            ))
          ) : (
            <EmptyStateContainerless
              message="Nema podataka o rasporedu."
              className="mt-12"
            />
          )}
        </div>
        <div className="new-page">
          <Title icon={FaBook} colorConfig={colorConfig}>
            Ocjene
          </Title>

          {gradebookData?.grade_data?.map((grade_data, pupilIdx) => (
            <div key={pupilIdx} className="mb-4">
              <Subtitle
                icon={grade_data?.pupil_unenrolled ? FaUserTimes : FaUser}
                colorConfig={colorConfig}
                showLine={false}
                textColor={grade_data?.pupil_unenrolled ? "text-red-500" : null}
              >
                {pupilIdx + 1}. {grade_data.pupil_name}
              </Subtitle>

              {chunkArray(
                gradebookData?.subjects || [],
                SUBJECT_CHUNK_SIZE,
              ).map((subjectGroup, groupIdx, allGroups) => (
                <Table key={groupIdx} minWidth="">
                  <TableHead>
                    <TableRow>
                      <TableHeader bordered>Polugodište</TableHeader>
                      {subjectGroup.map((subject, idx) => (
                        <TableHeader key={idx} bordered>
                          {subject.subject_name}
                        </TableHeader>
                      ))}
                      {/* Only show behaviour in the last table */}
                      {groupIdx === allGroups.length - 1 && (
                        <TableHeader bordered>Vladanje</TableHeader>
                      )}
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {grade_data.grades_for_semester.map(
                      (semester, semesterIdx) => (
                        <TableRow key={semesterIdx}>
                          <TableCell bordered>
                            <div className="font-medium text-blue-700">
                              {semester.semester_name}
                            </div>
                          </TableCell>
                          {subjectGroup.map((_, subjectIdx) => {
                            const globalSubjectIdx =
                              groupIdx * SUBJECT_CHUNK_SIZE + subjectIdx;
                            const subjectGrade =
                              semester.subject_grades[globalSubjectIdx];
                            return (
                              <TableCell key={subjectIdx} bordered>
                                {subjectGrade?.grades?.map(
                                  (grade, gradeIdx) => (
                                    <GradeType
                                      key={gradeIdx}
                                      gradeObject={grade}
                                      colorConfig={colorConfig}
                                      mode="view"
                                    />
                                  ),
                                )}
                              </TableCell>
                            );
                          })}
                          {groupIdx === allGroups.length - 1 && (
                            <TableCell bordered>
                              {semester?.behaviour_grades?.map(
                                (behaviourGrade, idx) => (
                                  <GradeType
                                    key={idx}
                                    gradeObject={behaviourGrade}
                                    colorConfig={colorConfig}
                                    mode="view"
                                    behaviourMode={true}
                                  />
                                ),
                              )}
                            </TableCell>
                          )}
                        </TableRow>
                      ),
                    )}
                  </TableBody>
                </Table>
              ))}
            </div>
          ))}
        </div>
        {gradebookData?.lessons?.length == 0 ? (
          <div className="new-page">
            <Title icon={FaFileAlt} colorConfig={colorConfig}>
              Časovi
            </Title>
            <EmptyStateContainerless
              message="Nema podataka o časovima."
              className="mt-12"
            />
          </div>
        ) : (
          gradebookData?.lessons?.map((lessonWeekGroup, idx) => {
            return (
              <div className="new-page" key={idx}>
                <Title icon={FaFileAlt} colorConfig={colorConfig}>
                  Časovi {lessonWeekGroup.week}
                </Title>

                <Subtitle
                  icon={FaCalendarAlt}
                  className="mb-2"
                  colorConfig={colorConfig}
                  showLine={false}
                >
                  Planirano časova: {lessonWeekGroup.total_lessons_in_week}
                </Subtitle>
                <Subtitle
                  icon={FaCheckCircle}
                  className="mb-2"
                  colorConfig={colorConfig}
                  showLine={false}
                >
                  Održanih časova: {lessonWeekGroup.held_lessons_in_week}
                </Subtitle>
                <Subtitle
                  icon={FaTimesCircle}
                  className="mb-6"
                  colorConfig={colorConfig}
                  showLine={false}
                >
                  Neodržanih časova: {lessonWeekGroup.unheld_lessons_in_week}
                </Subtitle>

                {lessonWeekGroup?.lesson_date_group?.map(
                  (lessonDateGroup, lessonDateGroupIdx) => {
                    return (
                      <div key={lessonDateGroupIdx} className="mb-6">
                        <Subtitle
                          icon={FaCalendarAlt}
                          showLine={false}
                          colorConfig={colorConfig}
                          className="mb-3"
                        >
                          {formatDateToDDMMYYYY(lessonDateGroup.date)} (Broj
                          održanih časova:{" "}
                          {lessonDateGroup.lessons_for_date.length})
                        </Subtitle>

                        {/* Lessons Table for this date */}
                        <Table>
                          <TableHead>
                            <TableRow>
                              <TableHeader>Predmet</TableHeader>
                              <TableHeader>Opis časa</TableHeader>
                              <TableHeader>Nastavnik</TableHeader>
                            </TableRow>
                          </TableHead>
                          <TableBody>
                            {lessonDateGroup.lessons_for_date.map(
                              (lesson, lessonIdx) => (
                                <TableRow key={lessonIdx}>
                                  <TableCell>
                                    <div className="font-medium text-blue-700">
                                      {lesson.subject_name}
                                    </div>
                                  </TableCell>
                                  <TableCell>
                                    <div className="max-w-md">
                                      {lesson.description || "Nema opisa"}
                                    </div>
                                  </TableCell>
                                  <TableCell>
                                    <div className="flex items-center">
                                      <FaUser
                                        className="mr-2 text-gray-500"
                                        size={12}
                                      />
                                      {lesson.lesson_posted_by_teacher}
                                    </div>
                                  </TableCell>
                                </TableRow>
                              ),
                            )}
                          </TableBody>
                        </Table>
                      </div>
                    );
                  },
                )}
              </div>
            );
          })
        )}

        {gradebookData?.absences?.length == 0 ? (
          <div className="new-page">
            <Title icon={FaUserTimes} colorConfig={colorConfig}>
              Izostanci
            </Title>
            <EmptyStateContainerless
              message="Nema izostanaka."
              className="mt-12"
            />
          </div>
        ) : (
          gradebookData?.absences?.map((week, weekIdx) => (
            <div key={weekIdx} className="new-page">
              <Title icon={FaUserTimes} colorConfig={colorConfig}>
                Izostanci {week.week}
              </Title>

              <div className="flex gap-5 mb-6">
                <Subtitle
                  icon={FaClipboardList}
                  colorConfig={colorConfig}
                  showLine={false}
                >
                  Ukupno: {week.total}
                </Subtitle>
                <Subtitle
                  icon={FaCheckCircle}
                  colorConfig={colorConfig}
                  showLine={false}
                  textColor="text-green-500"
                >
                  Opravdanih: {week.excused_count}
                </Subtitle>
                <Subtitle
                  icon={FaTimesCircle}
                  colorConfig={colorConfig}
                  showLine={false}
                  textColor="text-red-500"
                >
                  Neopravdanih: {week.unexcused_count}
                </Subtitle>
                <Subtitle
                  icon={FaHourglassHalf}
                  colorConfig={colorConfig}
                  showLine={false}
                  textColor="text-orange-500"
                >
                  Neriješenih: {week.pending_count}
                </Subtitle>
              </div>

              {week?.days?.map((day, dayIdx) => (
                <div key={dayIdx} className="mb-4">
                  <div className="pl-4 border-l-4 border-blue-300 bg-blue-50 p-3 rounded mb-4">
                    <Subtitle
                      icon={FaCalendarAlt}
                      showLine={false}
                      colorConfig={colorConfig}
                      className="mb-3"
                    >
                      {formatDateToDDMMYYYY(day.date)}
                    </Subtitle>

                    <div className="flex gap-5 mb-6">
                      <Subtitle
                        icon={FaClipboardList}
                        colorConfig={colorConfig}
                        showLine={false}
                      >
                        Ukupno: {day.total}
                      </Subtitle>
                      <Subtitle
                        icon={FaCheckCircle}
                        colorConfig={colorConfig}
                        showLine={false}
                        textColor="text-green-500"
                      >
                        Opravdanih: {day.excused_count}
                      </Subtitle>
                      <Subtitle
                        icon={FaTimesCircle}
                        colorConfig={colorConfig}
                        showLine={false}
                        textColor="text-red-500"
                      >
                        Neopravdanih: {day.unexcused_count}
                      </Subtitle>
                      <Subtitle
                        icon={FaHourglassHalf}
                        colorConfig={colorConfig}
                        showLine={false}
                        textColor="text-orange-500"
                      >
                        Neriješenih: {day.pending_count}
                      </Subtitle>
                    </div>
                  </div>

                  {day?.subjects?.map((subject, subjIdx) => (
                    <div key={subjIdx} className="mb-4">
                      <div className="ml-4 pl-4 border-l-4 border-green-300 bg-green-50 p-3 rounded mb-4">
                        <Subtitle
                          icon={FaBook}
                          showLine={false}
                          colorConfig={colorConfig}
                          className="mb-3"
                        >
                          {subject.subject_name}
                        </Subtitle>

                        <div className="flex gap-5 mb-6">
                          <Subtitle
                            icon={FaClipboardList}
                            colorConfig={colorConfig}
                            showLine={false}
                          >
                            Ukupno: {subject.total}
                          </Subtitle>
                          <Subtitle
                            icon={FaCheckCircle}
                            colorConfig={colorConfig}
                            showLine={false}
                            textColor="text-green-500"
                          >
                            Opravdanih: {subject.excused_count}
                          </Subtitle>
                          <Subtitle
                            icon={FaTimesCircle}
                            colorConfig={colorConfig}
                            showLine={false}
                            textColor="text-red-500"
                          >
                            Neopravdanih: {subject.unexcused_count}
                          </Subtitle>
                          <Subtitle
                            icon={FaHourglassHalf}
                            colorConfig={colorConfig}
                            showLine={false}
                            textColor="text-orange-500"
                          >
                            Neriješenih: {subject.pending_count}
                          </Subtitle>
                        </div>
                      </div>

                      <Table className="ml-4">
                        <TableHead>
                          <TableRow>
                            <TableHeader>Učenik</TableHeader>
                            <TableHeader>Status</TableHeader>
                            <TableHeader>Redni broj časa</TableHeader>
                          </TableRow>
                        </TableHead>
                        <TableBody>
                          {subject.absences.map((a, idx) => (
                            <TableRow key={idx}>
                              <TableCell>
                                {a.name} {a.last_name}
                              </TableCell>
                              <TableCell className={getValueColor(a.status)}>
                                {mapValueToBosnian(a.status)}
                              </TableCell>
                              <TableCell>{a.period_number}</TableCell>
                            </TableRow>
                          ))}
                        </TableBody>
                      </Table>
                    </div>
                  ))}
                </div>
              ))}
            </div>
          ))
        )}
      </div>
    </>
  );
};
