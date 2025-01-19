import { cn } from "@/lib/utils";
import "@/styles/globals.css";
import type { AppProps } from "next/app";
import { Karla } from "next/font/google";
import { ClerkProvider } from "@clerk/nextjs";

const karla = Karla({
  subsets: ["latin"],
  variable: "--font-karla",
});

export default function App({ Component, pageProps }: AppProps) {
  return (
    <ClerkProvider {...pageProps}>
      <main className={cn(karla.variable, "font-karla")}>
        <Component {...pageProps} />
      </main>
    </ClerkProvider>
  );
}
