import { getServerSession } from "next-auth";
import { authOptions } from "@/app/api/auth/[...nextauth]/route";
import PupilSectionsInvitePageClient from "./clientPage";
import Title from "@/app/components/common/Title";
import { FaRegEnvelope } from "react-icons/fa";

// Server Component: fetch data before render
export default async function PupilInvitesPage() {
  const session = await getServerSession(authOptions);
  const accessToken = session?.accessToken;
  const userID = session?.user?.id;

  // Fetch data from backend API
  const pupil_section_invites_resp = await fetch(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/common/pupil_section_invites/${userID}`,
    {
      cache: "no-store",
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    },
  );
  let pupil_section_invite_data = [];
  if (pupil_section_invites_resp.ok) {
    pupil_section_invite_data = await pupil_section_invites_resp.json();
  }

  return (
    <>
      <Title icon={FaRegEnvelope}>Pozivi u odjeljenja</Title>
      <PupilSectionsInvitePageClient
        initialInvites={pupil_section_invite_data}
        accessToken={accessToken}
      />
    </>
  );
}
