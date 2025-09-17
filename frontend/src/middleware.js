import { NextResponse } from "next/server";
import { getToken } from "next-auth/jwt";
import { sidebarPermissions } from "./app/components/navbar_sidebar/sidebarItems";

export async function middleware(request) {
  const accessToken = await getToken({
    req: request,
    secret: process.env.NEXTAUTH_SECRET,
  });
  const isAuth = !!accessToken;
  const accountType = accessToken?.account_type;

  // Allow unauthenticated access to login, register_teacher, and API/auth routes
  if (
    request.nextUrl.pathname.startsWith("/login") ||
    request.nextUrl.pathname.startsWith("/register_teacher") ||
    request.nextUrl.pathname.startsWith("/register_pupil") ||
    request.nextUrl.pathname.startsWith("/verify") ||
    request.nextUrl.pathname.startsWith("/parent_login")
  ) {
    return NextResponse.next();
  }

  if (!isAuth) {
    const loginUrl = new URL("/login", request.url);
    return NextResponse.redirect(loginUrl);
  }

  if (request.nextUrl.pathname === "/") {
    if (accountType == "root") {
      const redirectUrl = new URL("/tenants", request.url);
      return NextResponse.redirect(redirectUrl);
    } else if (accountType == "tenant_admin") {
      const redirectUrl = new URL("/tenant_admin_administration", request.url);
      return NextResponse.redirect(redirectUrl);
    } else if (accountType == "teacher") {
      const redirectUrl = new URL("/teacher_home", request.url);
      return NextResponse.redirect(redirectUrl);
    } else if (accountType == "pupil" || accountType == "parent") {
      const redirectUrl = new URL("/pupil_home", request.url);
      return NextResponse.redirect(redirectUrl);
    } else if (accountType == "parent") {
      // TODO: Later
      const redirectUrl = new URL("/tenants", request.url);
      return NextResponse.redirect(redirectUrl);
    }
  }

  const matchedSidebar = sidebarPermissions.find((item) =>
    request.nextUrl.pathname.startsWith(item.href),
  );
  if (matchedSidebar && !matchedSidebar.account_types.includes(accountType)) {
    return NextResponse.rewrite(new URL("/404", request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    /*
      Exclude:
      - _next (Next.js internals)
      - static files (images, fonts, etc.)
      - favicon.ico
      - api routes
    */
    "/((?!_next/static|_next/image|favicon.ico|api|.*\\.css$|.*\\.js$|.*\\.png$|.*\\.jpg$|.*\\.svg$|.*\\.ico$).*)",
  ],
};
