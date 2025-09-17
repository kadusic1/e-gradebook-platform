export const formatDateToDDMMYYYY = (dateString) => {
  const [year, month, day] = dateString.split("-");
  return `${day}.${month}.${year}`;
};

export const formatToDateTime = (date = new Date()) => {
  const targetDate = date instanceof Date ? date : new Date(date);

  const day = String(targetDate.getDate()).padStart(2, "0");
  const month = String(targetDate.getMonth() + 1).padStart(2, "0");
  const year = targetDate.getFullYear();
  const hours = String(targetDate.getHours()).padStart(2, "0");
  const minutes = String(targetDate.getMinutes()).padStart(2, "0");

  return `${day}.${month}.${year}. ${hours}:${minutes}`;
};

export const formatToDate = (date = new Date()) => {
  const targetDate = date instanceof Date ? date : new Date(date);

  const day = String(targetDate.getDate()).padStart(2, "0");
  const month = String(targetDate.getMonth() + 1).padStart(2, "0");
  const year = targetDate.getFullYear();

  return `${day}.${month}.${year}.`;
};

export const formatToFullDateTime = (dateTimeString) => {
  // Handles "YYYY-MM-DD HH:MM:SS" or "YYYY-MM-DD HH:MM:SS.ssssss"
  const [datePart, timePart] = dateTimeString.split(" ");
  if (!datePart || !timePart) return dateTimeString;

  const [year, month, day] = datePart.split("-");
  let [hour, minute, second] = timePart.split(":");
  const sec = second ? second.split(".")[0] : "0";
  if (minute === "00") {
    minute = "0";
  }
  if (hour == "00") {
    hour = "0";
  }

  return `${day}.${month}.${year}. ${hour}h ${minute}min ${sec}s`;
};
