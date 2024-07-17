"use client";

import dynamic from "next/dynamic";

const Terminal = dynamic(() => import("@/components/terminal"), {
  ssr: false,
});

export default function Page({ params }: { params: { slug: string } }) {
  return (
    <main className="w-screen h-screen pt-16 px-2 pb-2">
      <Terminal
        id="terminal"
        path={"/api/repl/" + params.slug}
        catchClose={true}
        className="w-full h-full"
      />
    </main>
  );
}
