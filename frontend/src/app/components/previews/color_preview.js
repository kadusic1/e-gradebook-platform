import { useFormContext } from "react-hook-form";
import DynamicCard from "../dynamic_card/DynamicCard";
import { FaFlask } from "react-icons/fa";

export default function ColorPreviewCard() {
  const { watch } = useFormContext();
  const color_choice = watch("color_config");

  const sample_card_data = {
    Ime: "Test",
    Prezime: "Testovic",
    Telefon: "+387 061 123 456",
  };

  return (
    <div className={`rounded-lg p-4`}>
      {color_choice && (
        <DynamicCard
          data={sample_card_data}
          textTitle="Pregled boja"
          tenantColorConfig={color_choice}
          showDelete={true}
          deleteButton={{
            label: "Izbriši",
            onClick: () => {},
          }}
          showEdit={true}
          editButton={{
            label: "Uredi",
            onClick: () => {},
          }}
          extraButton={{
            label: "Dodatno",
            onClick: () => {},
            icon: FaFlask,
          }}
        />
      )}
    </div>
  );
}
