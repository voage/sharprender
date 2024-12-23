import { cn } from "@/lib/utils";
import "@/styles/globals.css";
import type { AppProps } from "next/app";
import { Karla } from "next/font/google";
import { QueryClient, QueryClientProvider } from "react-query";

const karla = Karla({
  subsets: ["latin"],
  variable: "--font-karla",
});

const queryClient = new QueryClient();

export default function App({ Component, pageProps }: AppProps) {
  return (
    <QueryClientProvider client={queryClient}>
      <main className={cn(karla.variable, "font-karla")}>
        <Component {...pageProps} />
      </main>
    </QueryClientProvider>
  );
}
