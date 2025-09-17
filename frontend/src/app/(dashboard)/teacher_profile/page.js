import { getServerSession } from "next-auth";
import { authOptions } from "@/app/api/auth/[...nextauth]/route";
import GeneralProfileClient from "../pupil_profile/generalDataPageClient";

export default async function PupilProfilePage() {
  const session = await getServerSession(authOptions);
  const accessToken = session?.accessToken;
  const userID = session?.user?.id;

  const teacher_profile_data_resp = await fetch(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/teacher/get_teacher/${userID}`,
    {
      cache: "no-store",
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    },
  );
  let teacher_profile_data = [];
  if (teacher_profile_data_resp.ok) {
    teacher_profile_data = await teacher_profile_data_resp.json();
  }

  return (
    <GeneralProfileClient
        accessToken={accessToken}
        generalData={teacher_profile_data}
        userID={userID}
        mode="teacher"
    />
  );
}
