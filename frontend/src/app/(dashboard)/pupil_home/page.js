import { getServerSession } from "next-auth";
import { authOptions } from "@/app/api/auth/[...nextauth]/route";
import UserHomePageClient from "./clientPage";

// Server Component: fetch data before render
export default async function PupilHomePage({ archived = 0 }) {
  const session = await getServerSession(authOptions);
  const accessToken = session?.accessToken;
  const pupilID = session?.user?.id;

  // Fetch data from backend API
  const pupil_sections_response = await fetch(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/pupil/sections/${pupilID}/${archived}`,
    {
      cache: "no-store",
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    },
  );
  let pupil_sections_data = [];
  if (pupil_sections_response.ok) {
    pupil_sections_data = await pupil_sections_response.json();
  }

  return (
    <>
      <UserHomePageClient
        initialSections={pupil_sections_data}
        accessToken={accessToken}
        mode="pupil"
        user={session?.user}
        archived={archived}
      />
    </>
  );
}
