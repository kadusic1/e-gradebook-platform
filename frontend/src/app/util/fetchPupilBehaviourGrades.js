export const fetchPupilBehaviourGrade = async (
  tenant_id,
  section_id,
  pupil_id,
  accessToken,
) => {
  try {
    const response = await fetch(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/pupil/behaviour_grades/${tenant_id}/${section_id}/${pupil_id}`,
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      },
    );
    if (response.ok) {
      const data = await response.json();
      return data;
    }
  } catch (e) {
    console.error(e);
  }
};

export const fetchPupilBehaviourGradeNoPupilID = async (
  tenant_id,
  section_id,
  accessToken,
) => {
  try {
    const response = await fetch(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/pupil/behaviour_grades/${tenant_id}/${section_id}`,
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      },
    );
    if (response.ok) {
      const data = await response.json();
      return data;
    }
  } catch (e) {
    console.error(e);
  }
};
