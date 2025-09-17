import { getServerSession } from "next-auth";
import { authOptions } from "../auth/[...nextauth]/route";
import puppeteer from "puppeteer";

export async function POST(req) {
  let browser;
  try {
    const session = await getServerSession(authOptions);
    if (!session) {
      return new Response(JSON.stringify({ error: "Unauthorized" }), {
        status: 401,
      });
    }

    const { htmlContent, landscape, showPageNumbers } = await req.json();

    browser = await puppeteer.launch({ headless: "new" });
    const page = await browser.newPage();

    // Simply inject Tailwind CSS
    await page.setContent(
      `
      <!DOCTYPE html>
      <html>
        <head>
          <meta charset="utf-8" />
          <script src="https://cdn.tailwindcss.com"></script>
          <style>
            .new-page { page-break-before: always; }
          </style>
        </head>
        <body>
          ${htmlContent.replace(/<\/?html>|<\/?head>|<\/?body>/g, "")}
        </body>
      </html>
    `,
      { waitUntil: "networkidle0" },
    );

    const pdfOptions = {
      format: "A4",
      landscape: landscape,
      printBackground: true,
      margin: {
        top: "20px",
        right: "5px",
        bottom: showPageNumbers ? "50px" : "20px",
        left: "5px",
      },
    };

    if (showPageNumbers) {
      pdfOptions.displayHeaderFooter = true;
      pdfOptions.footerTemplate = `
        <div style="width:100%;font-size:14px;color:#888;text-align:center;">
          <span class="pageNumber"></span>
        </div>
      `;
      pdfOptions.headerTemplate = `<div></div>`;
    }

    const pdfBuffer = await page.pdf(pdfOptions);

    await page.close();

    return new Response(pdfBuffer, {
      headers: {
        "Content-Type": "application/pdf",
        "Content-Disposition": "attachment; filename=generated.pdf",
      },
    });
  } catch (error) {
    console.error("PDF generation error:", error);
    return new Response(JSON.stringify({ error: "PDF generation failed" }), {
      status: 500,
    });
  } finally {
    if (browser) {
      await browser.close();
    }
  }
}
