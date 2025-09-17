"use client";
import Subtitle from "@/app/components/common/Subtitle";
import Text from "@/app/components/common/Text";
import Title from "@/app/components/common/Title";
import {
  Table,
  TableHeader,
  TableHead,
  TableBody,
  TableCell,
  TableRow,
} from "@/app/components/table/TableComponents";
import BackButton from "../common/BackButton";
import { useEffect, useState, useRef } from "react";
import { formatDateToDDMMYYYY, formatToDate } from "@/app/util/date_util";
import { FaAward } from "react-icons/fa";
import { PDFButton } from "../common/PDFButton";
import { handleDownloadPDF } from "@/app/util/pdf_util";

export const CertificatePageClient = ({
  tenantID,
  sectionID,
  pupilID,
  onBack,
  colorConfig,
  accessToken,
}) => {
  const gradeNames = {
    5: "odličan (5)",
    4: "vrlo dobar (4)",
    3: "dobar (3)",
    2: "dovoljan (2)",
  };

  const fetchCertificateData = async () => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/pupil/certificate/${tenantID}/${sectionID}/${pupilID}`,
        {
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        },
      );

      if (response.ok) {
        const data = await response.json();
        setCertificateData(data);
      }
    } catch (e) {
      console.error(e);
    }
  };

  const [certificateData, setCertificateData] = useState(null);
  const [pdfLoading, setPdfLoading] = useState(false);
  const certificateRef = useRef();

  const downloadCertificatePDF = () => {
    handleDownloadPDF({
      filename: "svjedočanstvo.pdf",
      landscape: false,
      setPdfLoading: setPdfLoading,
      pdfElementRef: certificateRef,
      additionalStyles: `
      .border-gray-500 {
        border: none !important;
        border-color: transparent !important;
      }
      .w-3\\/4 {
        width: 100% !important;
      }
    `,
    });
  };

  useEffect(() => {
    fetchCertificateData();
  }, []);

  if (!certificateData?.passed) {
    return (
      <div>
        <Title icon={FaAward} colorConfig={colorConfig}>
          Svjedočanstvo za razred
        </Title>
        <div className="flex justify-end mb-4">
          <BackButton onClick={onBack} colorConfig={colorConfig} />
        </div>
        <Title
          showLine={false}
          colorConfig={colorConfig}
          className="text-center"
          textSize="text-2xl"
        >
          Svjedočanstvo nije dostupno. Učenik je pao ovaj razred.
        </Title>
      </div>
    );
  }

  return (
    <>
      <Title icon={FaAward} colorConfig={colorConfig}>
        Svjedočanstvo za razred
      </Title>
      <div className="flex justify-end mb-4">
        <BackButton onClick={onBack} colorConfig={colorConfig} />
        <PDFButton
          colorConfig={colorConfig}
          pdfLoading={pdfLoading}
          handleDownloadPDF={downloadCertificatePDF}
          className="ml-2"
        />
      </div>
      <div className="flex justify-center items-center min-h-screen animate-fadeIn">
        <div
          ref={certificateRef}
          className="bg-white shadow-lg border border-gray-500 p-8 text-center w-3/4 h-[28.5cm] relative"
        >
          <div className="mb-4">
            <Subtitle showLine={false} className="uppercase" textSize="text-md">
              Bosna i Hercegovina
            </Subtitle>
            <Subtitle showLine={false} className="uppercase" textSize="text-md">
              Federacija Bosne i Hercegovine
            </Subtitle>
            <Subtitle showLine={false} className="uppercase" textSize="text-md">
              {certificateData?.tenant?.canton_name} Kanton
            </Subtitle>
          </div>
          <Text className="mb-2">
            Školska {certificateData?.section?.year}. godina
          </Text>
          <div className="grid grid-cols-2 gap-2 mb-4 text-left">
            <Text textSize="text-md">
              <span className="font-bold">Naziv škole:</span>{" "}
              {certificateData?.tenant?.tenant_name}
            </Text>
            <Text textSize="text-md" className="text-right">
              <span className="font-bold">Sjedište:</span>{" "}
              {certificateData?.tenant?.tenant_city}
            </Text>
          </div>
          <div className="mb-2">
            <Title showLine={false} className="uppercase" textSize="text-2xl">
              SVJEDODŽBA
            </Title>
            <Subtitle className="uppercase" textSize="text-md" showLine={false}>
              O ZAVRŠENOM {certificateData?.section?.class_code} RAZREDU{" "}
              {certificateData?.tenant?.tenant_type == "primary" ? (
                <>OSNOVNE</>
              ) : (
                <>SREDNJE</>
              )}{" "}
              ŠKOLE
            </Subtitle>
          </div>
          <div className="grid grid-cols-2 gap-4 mt-4 text-left mb-4">
            <Text textSize="text-md">
              <span className="font-bold">
                Prezime i ime{" "}
                {certificateData?.pupil?.gender == "M" ? (
                  <>učenika:</>
                ) : (
                  <>učenice:</>
                )}{" "}
              </span>
              <span>
                {certificateData?.pupil?.name}{" "}
                {certificateData?.pupil?.last_name}
              </span>
            </Text>
            <Text textSize="text-md" className="text-right">
              <span className="font-bold">Ime staratelja: </span>
              <span>{certificateData?.pupil?.guardian_name}</span>
            </Text>
            {certificateData?.pupil?.date_of_birth && (
              <Text textSize="text-md">
                <span className="font-bold">Datum rođenja: </span>
                <span>
                  {formatDateToDDMMYYYY(certificateData?.pupil?.date_of_birth)}
                </span>
              </Text>
            )}
            <Text textSize="text-md" className="text-right">
              <span className="font-bold">Mjesto rođenja: </span>
              <span>{certificateData?.pupil?.place_of_birth}</span>
            </Text>
            {certificateData?.section?.course_name && (
              <Text textSize="text-md">
                <span className="font-bold">Stručno zvanje: </span>
                <span>{certificateData?.section?.course_name}</span>
              </Text>
            )}
          </div>
          <Text className="text-left" textSize="text-sm">
            {certificateData?.pupil?.gender == "M" ? <>Učenik</> : <>Učenica</>}{" "}
            je u naznačenom razredu{" "}
            {certificateData?.pupil?.gender == "M" ? (
              <>postigao</>
            ) : (
              <>postigla</>
            )}{" "}
            slijedeći uspjeh:
          </Text>
          <Table>
            <TableHead>
              <TableRow>
                <TableHeader center bordered>
                  NASTAVNI PREDMETI
                </TableHeader>
                <TableHeader center bordered>
                  OCJENE
                </TableHeader>
                <TableHeader center bordered>
                  NASTAVNI PREDMETI
                </TableHeader>
                <TableHeader center bordered>
                  OCJENE
                </TableHeader>
              </TableRow>
            </TableHead>
            <TableBody>
              {certificateData?.final_grades &&
                certificateData.final_grades
                  .reduce((rows, grade, idx, arr) => {
                    if (idx % 2 === 0) rows.push(arr.slice(idx, idx + 2));
                    return rows;
                  }, [])
                  .map((pair, idx) => (
                    <TableRow key={idx}>
                      <TableCell bordered yPadding="py-2">
                        {pair[0]?.subject_name}
                      </TableCell>
                      <TableCell bordered yPadding="py-2">
                        {gradeNames[pair[0]?.grade]}
                      </TableCell>
                      <TableCell bordered yPadding="py-2">
                        {pair[1]?.subject_name || ""}
                      </TableCell>
                      <TableCell bordered yPadding="py-2">
                        {pair[1] ? gradeNames[pair[1].grade] : ""}
                      </TableCell>
                    </TableRow>
                  ))}
            </TableBody>
          </Table>
          <Text className="text-left mt-2" textSize="text-sm">
            <span className="font-bold">Vladanje: </span>
            <span>{certificateData?.behaviour_grades?.behaviour}</span>
          </Text>
          {certificateData?.graduate_grade && (
            <Text className="text-left mt-2" textSize="text-sm">
              Prema tome{" "}
              {certificateData?.pupil?.gender == "M" ? (
                <> učenik</>
              ) : (
                <> učenica</>
              )}{" "}
              je{" "}
              {certificateData?.pupil?.gender == "M" ? (
                <> završio</>
              ) : (
                <> završila</>
              )}{" "}
              {certificateData?.section?.class_code} razred s uspjehom{" "}
              {gradeNames[certificateData?.graduate_grade]} i prosjekom (
              {Number(certificateData?.average_grade).toFixed(2)}).
            </Text>
          )}
          <div className="absolute bottom-10 left-8 right-8">
            <Text className="mt-4" textSize="text-md">
              <span className="font-bold">Mjesto i datum izdavanja: </span>
              <span>
                {certificateData?.tenant?.tenant_city}, {formatToDate()}
              </span>
            </Text>
            <div className="flex justify-between mt-4">
              <div className="text-left">
                <Text textSize="text-md" className="font-bold">
                  RAZREDNIK
                </Text>
                <Text textSize="text-md">
                  {certificateData?.section?.homeroom_teacher_full_name}
                </Text>
              </div>
              <div className="text-left">
                <Text textSize="text-md" className="font-bold">
                  DIREKTOR ŠKOLE
                </Text>
                <Text textSize="text-md">
                  {certificateData?.tenant?.director_name}
                </Text>
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
};
