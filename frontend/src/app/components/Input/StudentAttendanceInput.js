import React, { useState, useMemo, useEffect } from "react";
import { useFormContext } from "react-hook-form";
import { FaSearch, FaUserCheck, FaUserTimes } from "react-icons/fa";

export default function StudentAttendanceField({
  name = "pupil_attendance_data",
  label = "Prisustvo učenika",
  students = [],
  ...props
}) {
  const { setValue, watch } = useFormContext();
  const [searchTerm, setSearchTerm] = useState("");

  const [attendances, setAttendances] = useState(
    // Add student name last_name
    watch(name)?.map((item) => {
      return {
        ...item,
        name: students.find((s) => s.id === item.pupil_id)?.name,
        last_name: students.find((s) => s.id === item.pupil_id)?.last_name,
      };
    }) ||
      students.map((student) => {
        return {
          pupil_id: student.id,
          status: "present",
          name: student.name,
          last_name: student.last_name,
        };
      }),
  );

  useEffect(() => {
    setValue(name, attendances);
  }, [attendances]);

  // Filter students based on search term - only show when searching
  const filteredAttendances = useMemo(() => {
    if (!searchTerm.trim()) return [];

    return attendances.filter((attendance) => {
      const fullName =
        `${attendance.name} ${attendance.last_name}`.toLowerCase();
      const searchLower = searchTerm.toLowerCase();
      return (
        fullName.includes(searchLower) ||
        (attendance.name || "").toLowerCase().includes(searchLower) ||
        (attendance.last_name || "").toLowerCase().includes(searchLower)
      );
    });
  }, [searchTerm, attendances]);

  // Handle attendance status change
  const handleAttendanceChange = (pupilId, status) => {
    const updatedAttendances = attendances.map((attendance) =>
      attendance.pupil_id === pupilId ? { ...attendance, status } : attendance,
    );

    setAttendances(updatedAttendances);
  };

  // Count statistics
  const { presentCount, absentCount } = useMemo(() => {
    let present = 0;
    let absent = 0;

    attendances.forEach((attendance) => {
      if (attendance.status === "present") present++;
      else if (
        attendance.status === "absent" ||
        attendance.status === "unexcused" ||
        attendance.status === "excused"
      )
        absent++;
    });

    return {
      presentCount: present,
      absentCount: absent,
    };
  }, [attendances]);

  const totalStudents = students.length;

  return (
    <div className="space-y-4">
      {/* Statistics */}
      <div className="bg-gray-50 rounded-lg p-3 grid grid-cols-3 gap-4 text-sm">
        <div className="text-center">
          <div className="text-gray-600">Ukupno učenika</div>
          <div className="font-semibold text-lg">{totalStudents}</div>
        </div>
        <div className="text-center">
          <div className="text-green-600">Prisutni</div>
          <div className="font-semibold text-lg text-green-700">
            {presentCount}
          </div>
        </div>
        <div className="text-center">
          <div className="text-red-600">Odsutni</div>
          <div className="font-semibold text-lg text-red-700">
            {absentCount}
          </div>
        </div>
      </div>

      {/* Search */}
      <div className="relative">
        <FaSearch className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" />
        <input
          type="text"
          placeholder="Pretražite učenike..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
        />
      </div>

      {/* Student List */}
      {filteredAttendances.length > 0 && (
        <div className="space-y-2 border border-gray-200 rounded-lg p-2">
          {filteredAttendances.map((attendance, idx) => {
            return (
              <div
                key={idx}
                className="flex flex-col sm:flex-row sm:items-center justify-between p-3 bg-white border border-gray-200 rounded-lg hover:bg-gray-50 gap-3"
              >
                <div className="flex-1">
                  <div className="font-medium text-gray-900">
                    {attendance.name} {attendance.last_name}
                  </div>
                </div>

                <div className="flex flex-col sm:flex-row gap-2 w-full sm:w-auto">
                  <button
                    type="button"
                    onClick={() =>
                      handleAttendanceChange(attendance.pupil_id, "present")
                    }
                    className={`hover:cursor-pointer flex items-center justify-center gap-1 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                      attendance.status === "present"
                        ? "bg-green-100 text-green-800 border-2 border-green-300"
                        : "bg-gray-100 text-gray-600 border-2 border-transparent hover:bg-green-50"
                    }`}
                  >
                    <FaUserCheck className="w-3 h-3" />
                    Prisutan
                  </button>

                  <button
                    type="button"
                    onClick={() =>
                      handleAttendanceChange(attendance.pupil_id, "absent")
                    }
                    className={`hover:cursor-pointer flex items-center justify-center gap-1 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                      attendance.status === "absent" ||
                      attendance.status === "unexcused" ||
                      attendance.status === "excused"
                        ? "bg-red-100 text-red-800 border-2 border-red-300"
                        : "bg-gray-100 text-gray-600 border-2 border-transparent hover:bg-red-50"
                    }`}
                  >
                    <FaUserTimes className="w-3 h-3" />
                    Odsutan
                  </button>
                </div>
              </div>
            );
          })}
        </div>
      )}

      {/* Quick Actions */}
      <div className="flex gap-2 pt-2">
        <button
          type="button"
          onClick={() => {
            const allPresentAttendance = attendances.map((attendance) => ({
              ...attendance,
              status: "present",
            }));
            setAttendances(allPresentAttendance);
          }}
          className="hover:cursor-pointer flex items-center gap-1 px-3 py-2 bg-green-100 text-green-800 rounded-lg text-sm font-medium hover:bg-green-200 transition-colors"
        >
          <FaUserCheck className="w-3 h-3" />
          Svi prisutni
        </button>

        <button
          type="button"
          onClick={() => {
            const allAbsentAttendance = attendances.map((attendance) => ({
              ...attendance,
              status: "absent",
            }));
            setAttendances(allAbsentAttendance);
          }}
          className="hover:cursor-pointer flex items-center gap-1 px-3 py-2 bg-red-100 text-red-800 rounded-lg text-sm font-medium hover:bg-red-200 transition-colors"
        >
          <FaUserTimes className="w-3 h-3" />
          Svi odsutni
        </button>
      </div>
    </div>
  );
}
