import TenantsPageClient from "./clientPage";
import { getServerSession } from "next-auth";
import { authOptions } from "@/app/api/auth/[...nextauth]/route";

// Server Component: fetch data before render
export default async function TenantsPage() {
  const session = await getServerSession(authOptions);
  const accessToken = session?.accessToken;

  // Fetch data from backend API
  const tenants_data_response = await fetch(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/superadmin/tenants`,
    {
      cache: "no-store",
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    },
  );
  let tenant_data = [];
  if (tenants_data_response.ok) {
    tenant_data = await tenants_data_response.json();
  }

  const canton_data_response = await fetch(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/common/cantons`,
    {
      cache: "no-store",
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    },
  );
  let canton_data = [];
  if (canton_data_response.ok) {
    canton_data = await canton_data_response.json();
  }

  return (
    <div>
      <TenantsPageClient
        cantons={canton_data}
        initialTenants={tenant_data}
        accessToken={accessToken}
      />
    </div>
  );
}
