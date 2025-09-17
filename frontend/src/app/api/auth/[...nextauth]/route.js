import NextAuth from "next-auth";
import CredentialsProvider from "next-auth/providers/credentials";
import { jwtDecode } from "jwt-decode";

export const authOptions = {
  providers: [
    CredentialsProvider({
      name: "Credentials",
      credentials: {
        username: { label: "username", type: "text" },
        password: { label: "password", type: "password" },
      },
      async authorize(credentials) {
        try {
          const res = await fetch(
            `${process.env.NEXT_PUBLIC_API_BASE_URL}/login`,
            {
              method: "POST",
              headers: { "Content-Type": "application/json" },
              body: JSON.stringify({
                email: credentials.username,
                password: credentials.password,
              }),
            },
          );
          if (!res.ok) return null;
          const token = await res.text(); // Get JWT token as plain text
          let claims = {};
          try {
            claims = jwtDecode(token);
          } catch (e) {
            console.error("JWT decode failed:", e);
            return null; // Prevent login if decode fails
          }
          return { ...claims, token };
        } catch (error) {
          return null;
        }
      },
    }),
    CredentialsProvider({
      id: "parent-access",
      name: "Parent Access",
      credentials: {
        parent_access_code: { label: "Parent Access Code", type: "text" },
      },
      async authorize(credentials) {
        try {
          const res = await fetch(
            `${process.env.NEXT_PUBLIC_API_BASE_URL}/parent-login`,
            {
              method: "POST",
              headers: { "Content-Type": "application/json" },
              body: JSON.stringify({
                parent_access_code: credentials.parent_access_code,
              }),
            },
          );
          if (!res.ok) return null;
          const token = await res.text();
          let claims = {};
          try {
            claims = jwtDecode(token);
          } catch (e) {
            console.error("JWT decode failed:", e);
            return null;
          }
          return { ...claims, token, loginType: "parent" };
        } catch (error) {
          return null;
        }
      },
    }),
  ],
  session: {
    strategy: "jwt",
  },
  callbacks: {
    async jwt({ token, user }) {
      if (user) {
        token.accessToken = user.token;
        token.id = user.id;
        token.name = user.name;
        token.lastName = user.last_name;
        token.email = user.email;
        token.phone = user.phone;
        token.school = user.school;
        token.account_type = user.account_type;
        token.account_id = user.account_id;
        if (user.account_type === "tenant_admin") {
          token.tenant_id = user.tenant_id;
        }
      }
      return token;
    },
    async session({ session, token }) {
      session.accessToken = token.accessToken;
      session.user.id = token.id;
      session.user.name = token.name;
      session.user.lastName = token.lastName;
      session.user.email = token.email;
      session.user.phone = token.phone;
      session.user.school = token.school;
      session.user.account_type = token.account_type;
      session.user.account_id = token.account_id;
      if (token.account_type === "tenant_admin") {
        session.user.tenant_id = token.tenant_id;
      }
      return session;
    },
  },
  pages: {
    signIn: "/login",
  },
};

const handler = NextAuth(authOptions);
export { handler as GET, handler as POST };
