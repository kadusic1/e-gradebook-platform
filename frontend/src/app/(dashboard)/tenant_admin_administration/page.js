import TenantsPageClient from "../tenants/clientPage";
import { getServerSession } from "next-auth";
import { authOptions } from "@/app/api/auth/[...nextauth]/route";

// Server Component: fetch data before render
export default async function TenantAdminAdministrationPage() {
  const session = await getServerSession(authOptions);
  const accessToken = session?.accessToken;

  const tenant_id = session?.user?.tenant_id;

  // Fetch data from backend API
  const tenants_data_response = await fetch(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/tenant_admin/tenant/${tenant_id}`,
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
  // Turn tenant data into an array
  tenant_data = Array.isArray(tenant_data) ? tenant_data : [tenant_data];

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
        tenantAdminMode={true}
      />
    </div>
  );
}
