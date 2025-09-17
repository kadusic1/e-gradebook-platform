import Link from "next/link";

export default function Sidebar() {
  return (
    <aside className="w-64 h-full bg-white shadow-lg p-6 flex flex-col gap-4">
      <h2 className="text-2xl font-bold text-gray-800 mb-6 flex items-center gap-2">
        <span className="inline-block w-3 h-3 bg-indigo-500 rounded-full"></span>
        eDnevnik
      </h2>
      <nav className="flex flex-col gap-2">
        <Link
          href="/dashboard"
          className="flex items-center gap-2 px-3 py-2 text-gray-700 hover:bg-gray-100 rounded-lg transition-all duration-200"
        >
          Dashboard
        </Link>
        <Link
          href="/schools"
          className="flex items-center gap-2 px-3 py-2 text-gray-700 hover:bg-gray-100 rounded-lg transition-all duration-200"
        >
          Škole
        </Link>
        <Link
          href="/pupils"
          className="flex items-center gap-2 px-3 py-2 text-gray-700 hover:bg-gray-100 rounded-lg transition-all duration-200"
        >
          Učenici
        </Link>
        <Link
          href="/classes"
          className="flex items-center gap-2 px-3 py-2 text-gray-700 hover:bg-gray-100 rounded-lg transition-all duration-200"
        >
          Razredi
        </Link>
        <Link
          href="/courses"
          className="flex items-center gap-2 px-3 py-2 text-gray-700 hover:bg-gray-100 rounded-lg transition-all duration-200"
        >
          Predmeti
        </Link>
        <Link
          href="/grades"
          className="flex items-center gap-2 px-3 py-2 text-gray-700 hover:bg-gray-100 rounded-lg transition-all duration-200"
        >
          Ocjene
        </Link>
        <Link
          href="/settings"
          className="flex items-center gap-2 px-3 py-2 text-gray-700 hover:bg-gray-100 rounded-lg transition-all duration-200"
        >
          Postavke
        </Link>
        <Link
          href="/profile"
          className="flex items-center gap-2 px-3 py-2 text-gray-700 hover:bg-gray-100 rounded-lg transition-all duration-200"
        >
          Profil
        </Link>
      </nav>
    </aside>
  );
}
