import DynamicCard from "./DynamicCard";
import { useState } from "react";
import ConfirmModal from "../modal/ConfirmModal";
import Text from "../common/Text";
import CreateUpdateModal from "../modal/CreateUpdateModal";
import EmptyState from "./EmptyState";
import ErrorModal from "../modal/ErrorModal";
import DynamicTable from "./DynamicTable";
import Button from "../common/Button";
import { FaChevronLeft, FaChevronRight } from "react-icons/fa";

// Mapper function to map English keys to Bosnian labels
const keyToBosnianLabel = {
  name: "Ime",
  address: "Adresa",
  city: "Grad",
  email: "Email",
  type: "Tip",
  canton: "Kanton",
  canton_code: "Kanton kod",
  canton_name: "Naziv kantona",
  tenant_name: "Naziv institucije",
  phone: "Telefon",
  director_name: "Ime direktora",
  tenant_type: "Tip institucije",
  last_name: "Prezime",
  role: "Uloga",
  curriculum_name: "Naziv kurikuluma",
  class_code: "Razred",
  npp_name: "Tip nastavnog plana i programa",
  year: "Godina",
  homeroom_teacher_email: "Email razrednika",
  homeroom_teacher_full_name: "Ime i prezime razrednika",
  semester_name: "Polugodište",
  gender: "Spol",
  date_of_birth: "Datum rođenja",
  place_of_birth: "Mjesto rođenja",
  guardian_name: "Ime staratelja",
  phone_number: "Broj telefona",
  guardian_number: "Broj telefona staratelja",
  attends_religion: "Pohađa vjeronauku",
  religion: "Vjeronauka",
  end_date: "Datum kraja",
  start_date: "Datum početka",
  status: "Status",
  invite_date: "Datum poziva",
  domain: "Domena",
  section_name: "Naziv odjeljenja",
  pupil_full_name: "Ime i prezime učenika",
  subject_name: "Naziv predmeta",
  teacher_full_name: "Ime i prezime nastavnika",
  homeroom_teacher: "Razrednik",
  teacher_email: "Email nastavnika",
  pupil_email: "Email učenika",
  capacity: "Kapacitet",
  code: "Broj učionice",
  description: "Opis",
  period_number: "Redni broj časa",
  present_pupil_count: "Prisutno učenika",
  absent_pupil_count: "Odsutno učenika",
  behaviour_grade: "Vladanje",
  pupil_behaviour_grade: "Moje vladanje",
  tenant_city: "Grad institucije",
  behaviour_determined_by_teacher: "Odredio nastavnik",
  lesson_posted_by_teacher: "Objavio nastavnik",
  valid_until: "Važilo do",
  behaviour: "Vladanje",
  child_of_martyr: "Dijete šehida",
  father_name: "Ime oca",
  mother_name: "Ime majke",
  parents_rvi: "Roditelji RVI",
  living_condition: "Uslovi života",
  student_dorm: "Studentski dom",
  refugee: "Izbjeglica",
  returnee_from_abroad: "Povratnik iz inostranstva",
  country_of_birth: "Država rođenja",
  country_of_living: "Država prebivališta",
  citizenship: "Državljanstvo",
  ethnicity: "Etnička pripadnost",
  father_occupation: "Stručna sprema oca",
  mother_occupation: "Stručna sprema majke",
  has_no_parents: "Bez roditelja",
  extra_information: "Dodatne informacije",
  child_alone: "Dodatne informacije ukoliko dijete nema roditelja",
  is_commuter: "Vozar",
  commuting_type: "Način putovanja",
  distance_to_school_km: "Udaljenost do škole (km)",
  has_hifz: "Dijete hafiz",
  special_honors: "Posebna priznanja",
  contractions: "Oslovljavanje",
  title: "Titula",
  specialization: "Usmjerenje",
  course_code: "Šifra smjera",
  course_name: "Naziv smjera",
};

const valueBosnianLabel = {
  primary: "Osnovna škola",
  secondary: "Srednja škola",
  superadmin: "Super admin",
  tenant_admin: "Školski admin",
  teacher: "Profesor",
  M: "Muški",
  F: "Ženski",
  true: "Da",
  false: "Ne",
  Catholic: "Katolička",
  Orthodox: "Pravoslavna",
  Jewish: "Jevrejska",
  Other: "Ostalo",
  NotAttendingReligion: "Ne pohađa vjeronauku",
  None: "Nema",
  pending: "Na čekanju",
  accepted: "Prihvaćen",
  declined: "Odbijen",
  tenant_domain: "Institucijska domena",
  global_domain: "Globalna domena",
  absent: "Odsutan",
  present: "Prisutan",
  excused: "Opravdan",
  unexcused: "Neopravdan",
  primjerno: "Primjerno",
  vrlodobro: "Vrlo dobro",
  dobro: "Dobro",
  zadovoljavajuće: "Zadovoljavajuće",
  loše: "Loše",
  null: "Nije postavljeno",
  both_parents: "Oba roditelja",
  one_parent: "Jedan roditelj",
  another_family_or_alone: "Druga porodica ili samostalno",
  institution_for_children_without_parents: "Ustanova za djecu bez roditelja",
  Walking: "Pješke",
  Bike: "Biciklo",
  Car: "Automobil",
  Bus: "Autobus",
  Train: "Voz",
  NotTraveling: "Ne putuje",
  NoOccupation: "Nema zanimanja",
  "<=5km": "Do 5km",
  "5km - 10km": "5km - 10km",
  "10km - 25km": "10km - 25km",
  ">25km": "Preko 25km",
  regular: "Obično",
  religion: "Vjersko",
  musical: "Muzičko",
};

export function mapKeyToBosnian(key) {
  return keyToBosnianLabel[key] || key;
}

export function mapValueToBosnian(value) {
  return valueBosnianLabel[value] || value;
}

export const getValueColor = (value) => {
  const colorMap = {
    pending: "text-yellow-600",
    accepted: "text-green-600",
    declined: "text-red-500",
    true: "text-green-600",
    false: "text-red-500",
    absent: "text-yellow-600",
    excused: "text-green-600",
    unexcused: "text-red-500",
    Da: "text-green-600",
    Ne: "text-red-500",
  };
  return colorMap[value] || null;
};

// Support titleField as string or array
export const getTitle = (item, title) => {
  let finalTitle;
  if (Array.isArray(title)) {
    finalTitle = title
      .map((field) => item[field])
      .filter(Boolean)
      .join(" ");
  } else {
    finalTitle = item[title] || "";
  }
  return mapValueToBosnian(finalTitle);
};

/**
 * DynamicCardParent renders a grid of DynamicCard components.
 * @param {Array} data - Array of data objects for each card
 * @param {string} [className] - Additional classes for the grid container
 * @returns JSX.Element
 */
export default function DynamicCardParent({
  data,
  setData,
  icon,
  titleField,
  className = "",
  prefix = "objekat",
  deleteUrl,
  editFields,
  editUrl,
  showEdit,
  showDelete,
  extraButton,
  onDeleteConfirm, // If we want to override the default delete behavior
  keyField = "id", // Default key field for identifying items
  keysToIgnore = [], // Fields to ignore in the card display
  keysToExclude = [], // Show value but do not show key
  accessToken,
  onEditSave,
  editButton,
  deleteButton,
  showEmptyState = true,
  mode = "card", // card or table mode
  textTitle, // Fixed text title
  tenantColorConfig,
  emptyMessage,
  extraActions = [], // Array of { label, icon, onClick }
  twoColumnsCard = true,
  getIcon,
  getBgColorOption,
  getTextColorOption,
  itemsPerPage = 6,
  paginationEnabled = true,
}) {
  const [currentPage, setCurrentPage] = useState(1);
  const [deleteModal, setDeleteModal] = useState({ open: false, item: null });
  const [editModal, setEditModal] = useState({ open: false, item: null });
  const [errorMessage, setErrorMessage] = useState();

  const isPaginationEnabled =
    paginationEnabled !== false && data?.length > itemsPerPage;
  const totalPages = Math.ceil((data?.length || 0) / itemsPerPage);
  const startIndex = (currentPage - 1) * itemsPerPage;
  const paginatedData = isPaginationEnabled
    ? data?.slice(startIndex, startIndex + itemsPerPage)
    : data;

  const openDeleteModal = (item) => setDeleteModal({ open: true, item });
  const closeDeleteModal = () => setDeleteModal({ open: false, item: null });

  const openEditModal = (item) => {
    setEditModal({ open: true, item });
  };
  const closeEditModal = () => setEditModal({ open: false, item: null });

  const defaultHandleDeleteConfirm = async () => {
    try {
      const resp = await fetch(`${deleteUrl}/${deleteModal.item.id}`, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${accessToken}`,
        },
      });
      if (!resp.ok) {
        const error_message = await resp.text();
        setErrorMessage(error_message);
      } else {
        setData((prevData) =>
          prevData.filter((item) => item.id !== deleteModal.item.id),
        );
      }
    } catch (error) {
      console.error("Error deleting item:", error);
    }
    closeDeleteModal();
  };

  const handleDeleteConfirm = onDeleteConfirm
    ? () => onDeleteConfirm(deleteModal?.item?.[keyField], closeDeleteModal)
    : defaultHandleDeleteConfirm;

  const defaultHandleEditSave = async (newData) => {
    try {
      const resp = await fetch(`${editUrl}/${newData?.[keyField]}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${accessToken}`,
        },
        body: JSON.stringify(newData),
      });
      if (!resp.ok) {
        const error_message = await resp.text();
        setErrorMessage(error_message);
      } else {
        const updatedItem = await resp.json();
        setData((prevData) =>
          prevData.map((item) =>
            item.id === updatedItem.id ? updatedItem : item,
          ),
        );
      }
    } catch (error) {
      console.error("Error updating item:", error);
    }
    closeEditModal();
  };

  const handleEditSave = onEditSave
    ? (newData) => onEditSave(newData, setData, closeEditModal, setErrorMessage)
    : defaultHandleEditSave;

  return (
    <>
      {mode == "card" ? (
        <div
          className={`grid grid-cols-1 ${twoColumnsCard ? "lg:grid-cols-2" : ""} gap-4 md:gap-6 ${className}`}
        >
          {!data || data?.length == 0
            ? showEmptyState && (
                <div className="col-span-full min-h-[60vh] flex items-center justify-center">
                  <EmptyState
                    message={
                      emptyMessage || `Trenutno nema podataka za prikaz.`
                    }
                  />
                </div>
              )
            : paginatedData?.map((item, idx) => {
                const itemShowEdit =
                  typeof showEdit === "function" ? showEdit(item) : showEdit;
                const itemShowDelete =
                  typeof showDelete === "function"
                    ? showDelete(item)
                    : showDelete;
                const itemEditButton =
                  typeof editButton === "function"
                    ? editButton(item)
                    : editButton;
                const itemExtraButton =
                  typeof extraButton === "function"
                    ? extraButton(item)
                    : extraButton;
                return (
                  <DynamicCard
                    key={idx}
                    data={item}
                    icon={getIcon ? getIcon(item) : icon}
                    titleField={titleField}
                    onConfirm={openDeleteModal}
                    onEdit={openEditModal}
                    showEdit={itemShowEdit}
                    showDelete={itemShowDelete}
                    extraButton={itemExtraButton}
                    keysToIgnore={keysToIgnore}
                    getTitle={getTitle}
                    editButton={itemEditButton}
                    deleteButton={deleteButton}
                    mapKeyToBosnian={mapKeyToBosnian}
                    mapValueToBosnian={mapValueToBosnian}
                    textTitle={textTitle}
                    keysToExclude={keysToExclude}
                    tenantColorConfig={tenantColorConfig}
                    getValueColor={getValueColor}
                    extraActions={extraActions}
                    bgColorOption={
                      getBgColorOption ? getBgColorOption(item) : "primary"
                    }
                    textColorOption={
                      getTextColorOption ? getTextColorOption(item) : "primary"
                    }
                  />
                );
              })}
        </div>
      ) : (
        <>
          {(!data || data?.length === 0) && showEmptyState ? (
            <div className="min-h-[60vh] flex items-center justify-center">
              <EmptyState
                message={emptyMessage || `Trenutno nema podataka za prikaz.`}
              />
            </div>
          ) : (
            <DynamicTable
              data={paginatedData}
              titleField={titleField}
              onConfirm={openDeleteModal}
              onEdit={openEditModal}
              showEdit={showEdit}
              showDelete={showDelete}
              extraButton={extraButton}
              keysToIgnore={keysToIgnore}
              getTitle={getTitle}
              editButton={editButton}
              deleteButton={deleteButton}
              keyField={keyField}
              mapKeyToBosnian={mapKeyToBosnian}
              mapValueToBosnian={mapValueToBosnian}
              tenantColorConfig={tenantColorConfig}
              getValueColor={getValueColor}
              extraActions={extraActions}
            />
          )}
        </>
      )}

      {isPaginationEnabled && (
        <div className="flex justify-center items-center gap-4 mt-6">
          <Button
            colorConfig={tenantColorConfig}
            icon={FaChevronLeft}
            onClick={() => {
              setCurrentPage((p) => Math.max(1, p - 1));
            }}
            disabled={currentPage === 1}
          ></Button>
          <Text>
            {currentPage}/{totalPages}
          </Text>
          <Button
            colorConfig={tenantColorConfig}
            icon={FaChevronRight}
            onClick={() => setCurrentPage((p) => Math.min(totalPages, p + 1))}
            disabled={currentPage === totalPages}
          ></Button>
        </div>
      )}

      {deleteModal.open && (
        <ConfirmModal
          onClose={closeDeleteModal}
          onConfirm={handleDeleteConfirm}
          colorConfig={tenantColorConfig}
        >
          <Text className="text-center font-semibold">
            {`Da li ste sigurni da želite izbrisati ${prefix} "${getTitle(deleteModal.item, titleField)}"?`}
          </Text>
        </ConfirmModal>
      )}

      {editModal.open && (
        <CreateUpdateModal
          title={`Uredi ${prefix}`}
          fields={editFields}
          onClose={closeEditModal}
          onSave={handleEditSave}
          initialValues={editModal.item}
          colorConfig={tenantColorConfig}
        />
      )}
      {errorMessage && (
        <ErrorModal onClose={() => setErrorMessage(null)}>
          {errorMessage}
        </ErrorModal>
      )}
    </>
  );
}
