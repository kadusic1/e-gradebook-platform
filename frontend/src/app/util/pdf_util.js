export const handleDownloadPDF = async (config = {}) => {
  const {
    headerElements = [],
    filename = "schedule.pdf",
    landscape = false,
    showPageNumbers = false,
    additionalStyles = "",
    setPdfLoading,
    pdfElementRef,
  } = config;

  if (!pdfElementRef.current) return;

  try {
    setPdfLoading(true);

    const pdfElementHtml = pdfElementRef.current.outerHTML;
    const headerHtml = headerElements.join("\n");

    const htmlContent = `
      <html>
        <head>
          <meta charset="utf-8" />
          <style>
            .pdf-hide {
              display: none !important;
            }
            ${additionalStyles}
          </style>
        </head>
        <body>
          ${headerHtml}
          ${pdfElementHtml}
        </body>
      </html>
    `;

    const res = await fetch("/api/generate-pdf", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ htmlContent, landscape, showPageNumbers }),
    });

    if (res.ok) {
      const blob = await res.blob();
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = filename;
      document.body.appendChild(a);
      a.click();
      URL.revokeObjectURL(url);
      a.remove();
    } else {
      alert("PDF download failed!");
    }
  } finally {
    setPdfLoading(false);
  }
};
