import { getServerSession } from "next-auth";
import { authOptions } from "@/app/api/auth/[...nextauth]/route";
import Title from "@/app/components/common/Title";
import { FaRegEnvelope } from "react-icons/fa";
import TeacherInvitesPageClient from "../tenants/teacherInvitesPageClient";

// Server Component: fetch data before render
export default async function TeacherInvitesPage() {
  const session = await getServerSession(authOptions);
  const accessToken = session?.accessToken;
  const userID = session?.user?.id;

  // Fetch data from backend API
  const teacher_section_invites_resp = await fetch(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/teacher/invites/${userID}`,
    {
      cache: "no-store",
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    },
  );
  let teacher_invite_data = [];
  if (teacher_section_invites_resp.ok) {
    teacher_invite_data = await teacher_section_invites_resp.json();
  }

  return (
    <>
      <Title icon={FaRegEnvelope}>Pozivi u odjeljenja</Title>
      <TeacherInvitesPageClient
        initialInvites={teacher_invite_data}
        accessToken={accessToken}
        teacherAccountMode={true}
        showHeader={false}
      />
    </>
  );
}
