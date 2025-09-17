import { getServerSession } from "next-auth";
import { authOptions } from "@/app/api/auth/[...nextauth]/route";
import SchedulePageClient from "../tenants/schedulePageClient";
import Title from "@/app/components/common/Title";
import { FaChalkboard } from "react-icons/fa";

// Server Component: fetch data before render
export default async function TeacherSchedulePage() {
  const session = await getServerSession(authOptions);
  const accessToken = session?.accessToken;
  const teacherID = session?.user?.id;

  // Fetch data from backend API
  const teacher_schedule_response = await fetch(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/teacher/schedule/${teacherID}`,
    {
      cache: "no-store",
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    },
  );
  let teacher_schedule_data = [];
  if (teacher_schedule_response.ok) {
    teacher_schedule_data = await teacher_schedule_response.json();
  }

  return (
    <div>
      <Title icon={FaChalkboard}>Moj raspored</Title>
      <SchedulePageClient
        accessToken={accessToken}
        readOnly={true}
        initialSchedule={teacher_schedule_data}
        teacherMode={true}
      />
    </div>
  );
}
