import { getServerSession } from "next-auth";
import { authOptions } from "@/app/api/auth/[...nextauth]/route";
import GlobalAdminPageClient from "./clientPage";

// Server Component: fetch data before render
export default async function GlobalAdminPage() {
  const session = await getServerSession(authOptions);
  const accessToken = session?.accessToken;

  // Fetch data from backend API
  const npp_semesters_response = await fetch(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/superadmin/npp_semesters`,
    {
      cache: "no-store",
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    },
  );
  let npp_semester_data = [];
  if (npp_semesters_response.ok) {
    npp_semester_data = await npp_semesters_response.json();
  }

  const domains_response = await fetch(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/superadmin/all_domains`,
    {
      cache: "no-store",
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    },
  );
  let domains_data = [];
  if (domains_response.ok) {
    domains_data = await domains_response.json();
  }

  return (
    <div>
      <GlobalAdminPageClient
        initialSemesters={npp_semester_data}
        initialDomains={domains_data}
        accessToken={accessToken}
      />
    </div>
  );
}
