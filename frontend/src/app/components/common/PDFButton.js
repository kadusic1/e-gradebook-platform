import Button from "./Button";
import { FaDownload, FaSpinner } from "react-icons/fa";

export const PDFButton = ({
  colorConfig,
  handleDownloadPDF,
  pdfLoading,
  className,
}) => {
  return (
    <Button
      onClick={handleDownloadPDF}
      colorConfig={colorConfig}
      color="ternary"
      className={className}
      icon={!pdfLoading ? FaDownload : FaSpinner}
      iconClassName={pdfLoading ? "animate-spin" : ""}
    >
      {!pdfLoading ? "Preuzmi PDF" : "Generisanje"}
    </Button>
  );
};
