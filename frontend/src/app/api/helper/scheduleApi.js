// utils/api/scheduleApi.js
export const getScheduleForSection = async (
  sectionID,
  tenantId,
  accessToken,
  setSchedule,
) => {
  try {
    const response = await fetch(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/pupil/schedule/${tenantId}/${sectionID}`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${accessToken}`,
        },
      },
    );

    if (!response.ok) {
      throw new Error("Failed to fetch schedule");
    }

    const data = await response.json();
    setSchedule(data);
  } catch (error) {
    console.error("Error fetching schedule:", error);
    throw error; // Re-throw to let the component handle the error
  }
};
