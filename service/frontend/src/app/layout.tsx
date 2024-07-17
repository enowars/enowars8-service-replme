import "./globals.css";
import "@xterm/xterm/css/xterm.css";

import type { Metadata } from "next";
import { ThemeProvider } from "@/components/theme-provider";
import Navbar from "@/components/navigation-bar";
import localFont from "next/font/local";
import ReactQueryProvider from "@/components/react-query-provider";

const dejaVuFont = localFont({ src: "./DejaVuSansMNerdFontMono-Regular.ttf" });

export const metadata: Metadata = {
  title: "replme",
  description: "Create development environments and REPLs",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={dejaVuFont.className + " h-screen w-screen"}>
        <ReactQueryProvider>
          <ThemeProvider
            attribute="class"
            defaultTheme="system"
            enableSystem
            disableTransitionOnChange
          >
            <Navbar />
            {children}
          </ThemeProvider>
        </ReactQueryProvider>
      </body>
    </html>
  );
}
