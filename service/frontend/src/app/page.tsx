"use client";

import CreateReplButton from "@/components/create-repl-button";
import { Button } from "@/components/ui/button";
import { CodeIcon, RocketIcon, ReloadIcon } from "@radix-ui/react-icons"
import { useMutation, useQueryClient } from "@tanstack/react-query";

export default function Page() {
  return (
    <main className="flex h-full w-full flex-col items-center p-24 justify-center space-y-5">
      <div className="text-2xl italic">Want a clean /home?</div>
      <div className="text-5xl font-bold">Use replme!</div>
      <div>
        Hack together your ideas in development environments,
        use a throwaway shell - all in one place.
      </div>
      <div className="flex flex-row space-x-3 items-center">
        <a href="/signup">
          <Button>
            <RocketIcon className="mr-2 h-4 w-4" /> Register for free
          </Button>
        </a>
        <CreateReplButton />
      </div>
    </main>
  );
}
