import "../globals.css";
import "../main-content.css";

export default function AuthLayout({ children }) {
  return (
    <html lang="bs">
      <body>
        <div className="min-h-screen">{children}</div>
      </body>
    </html>
  );
}
