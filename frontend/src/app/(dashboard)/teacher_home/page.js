import { getServerSession } from "next-auth";
import { authOptions } from "@/app/api/auth/[...nextauth]/route";
import UserHomePageClient from "../pupil_home/clientPage";

// Server Component: fetch data before render
export default async function TeacherHomePage({ archived = 0 }) {
  const session = await getServerSession(authOptions);
  const accessToken = session?.accessToken;
  const teacherID = session?.user?.id;

  // Fetch data from backend API
  const teacher_sections_response = await fetch(
    // 0 in archived means that we get non-archived sections
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/teacher/sections/${teacherID}/${archived}`,
    {
      cache: "no-store",
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    },
  );
  let teacher_sections_data = [];
  if (teacher_sections_response.ok) {
    teacher_sections_data = await teacher_sections_response.json();
  }

  return (
    <>
      <UserHomePageClient
        initialSections={teacher_sections_data}
        accessToken={accessToken}
        mode="teacher"
        user={session?.user}
        archived={archived}
      />
    </>
  );
}
