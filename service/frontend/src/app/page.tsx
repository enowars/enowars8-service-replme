"use client";

import CreateDevenvButton from "@/components/create-devenv-button";
import CreateReplButton from "@/components/create-repl-button";
import { Button } from "@/components/ui/button";
import { useUserQuery } from "@/hooks/use-user-query";
import { RocketIcon } from "@radix-ui/react-icons"

export default function Page() {
  const userQuery = useUserQuery()
  const isAuthenticatedMode = !userQuery.isStale && userQuery.isSuccess

  return (
    <main className="flex h-screen w-screen flex-col items-center p-24 justify-center space-y-5">
      <div className="text-2xl italic">Want a clean /home?</div>
      <div className="text-5xl font-bold">Use replme!</div>
      <div>
        Hack together your ideas in development environments,
        use a throwaway shell - all in one place.
      </div>
      <div className="flex flex-row space-x-3 items-center">
        {isAuthenticatedMode ?
          <CreateDevenvButton /> :
          <a href="/register">
            <Button>
              <RocketIcon className="mr-2 h-4 w-4" /> REGISTER!
            </Button>
          </a>}
        <CreateReplButton />
      </div>
    </main>
  );
}
