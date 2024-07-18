"use client";

import CreateDevenvButton from "@/components/create-devenv-button";
import CreateReplButton from "@/components/create-repl-button";
import RandomAdBanner from "@/components/random-ad-banner";
import { Button } from "@/components/ui/button";
import { useUserQuery } from "@/hooks/use-user-query";
import { RocketIcon } from "@radix-ui/react-icons";

export default function Page() {
  const userQuery = useUserQuery();
  const isAuthenticatedMode = !userQuery.isStale && userQuery.isSuccess;

  return (
    <main className="flex h-screen w-screen flex-row pl-24 pr-5 py-24 space-x-10">
      <div className="flex h-full w-full flex-col items-center justify-center space-y-5 grow-0">
        <div className="text-2xl italic">Want a clean /home?</div>
        <div className="text-5xl font-bold">Use replme!</div>
        <div>
          Hack together your ideas in development environments, use a throwaway
          shell - all in one place.
        </div>
        <div className="flex flex-row items-center space-x-3">
          <div className="py-1 px-3 border rounded-xl border-black dark:border-white">
            🔐 100% E2EE
          </div>
          <div className="py-1 px-3 border rounded-xl border-black dark:border-white">
            🦾 100% Reliability
          </div>
          <div className="py-1 px-3 border rounded-xl border-black dark:border-white">
            🚀 Blazingly fast
          </div>
        </div>
        <div className="flex flex-row space-x-3 items-center">
          {isAuthenticatedMode ? (
            <CreateDevenvButton />
          ) : (
            <a href="/register">
              <Button>
                <RocketIcon className="mr-2 h-4 w-4" /> REGISTER!
              </Button>
            </a>
          )}
          <CreateReplButton />
        </div>
      </div>

      <div className="relative w-96 h-full">
        <RandomAdBanner />
      </div>
    </main>
  );
}
