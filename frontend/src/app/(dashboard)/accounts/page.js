import { FaChalkboardTeacher } from "react-icons/fa";
import Title from "../../components/common/Title";
import AccountsPageClient from "./clientPage";
import { getServerSession } from "next-auth";
import { authOptions } from "@/app/api/auth/[...nextauth]/route";

// Server Component: fetch data before render
export default async function AccountsPage() {
  const session = await getServerSession(authOptions);
  const accessToken = session?.accessToken;

  // Fetch data from backend API
  const teachers_data_response = await fetch(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/teachers`,
    {
      cache: "no-store",
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    },
  );
  let teacher_data = [];
  if (teachers_data_response.ok) {
    teacher_data = await teachers_data_response.json();
  }

  const pupils_data_response = await fetch(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/pupil_accounts`,
    {
      cache: "no-store",
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    },
  );
  let pupils_data = [];
  if (pupils_data_response.ok) {
    pupils_data = await pupils_data_response.json();
  }

  return (
    <div>
      <AccountsPageClient
        initialTeachers={teacher_data}
        initialPupils={pupils_data}
        accessToken={accessToken}
      />
    </div>
  );
}
