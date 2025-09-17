export async function fetchSectionPupilsUtil({
  tenantId,
  sectionId,
  accessToken,
}) {
  try {
    const response = await fetch(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/teacher/pupils/${tenantId}/${sectionId}`,
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      },
    );
    if (response.ok) {
      const data = await response.json();
      return {
        pupils: data?.pupils || [],
        pendingPupils: data?.pending_pupils || [],
        pupilsForAssignment: data?.pupils_for_assignment || [],
        error: null,
      };
    } else {
      const error_message = await response.text();
      return { error: error_message };
    }
  } catch (error) {
    console.error("Error fetching section pupils:", error);
    return { error: "Failed to fetch section pupils" };
  }
}
