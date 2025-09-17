import PupilProfileClient from "./clientPage";
import { getServerSession } from "next-auth";
import { authOptions } from "@/app/api/auth/[...nextauth]/route";

export default async function PupilProfilePage() {
  const session = await getServerSession(authOptions);
  const accessToken = session?.accessToken;
  const userID = session?.user?.id;

  const general_data_resp = await fetch(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/pupil/profile`,
    {
      cache: "no-store",
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    },
  );
  let general_data = [];
  if (general_data_resp.ok) {
    general_data = await general_data_resp.json();
  }

  const stats_data_resp = await fetch(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/pupil/statistics_fields/${userID}`,
    {
      cache: "no-store",
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    },
  );
  let stats_data = [];
  if (stats_data_resp.ok) {
    stats_data = await stats_data_resp.json();
  }

  return (
    <PupilProfileClient
      userID={userID}
      accessToken={accessToken}
      generalData={general_data}
      statsData={stats_data}
    />
  );
}
