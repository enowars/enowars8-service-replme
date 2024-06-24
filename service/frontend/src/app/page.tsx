"use client";

import CreateDevenvButton from "@/components/create-devenv-button";
import CreateReplButton from "@/components/create-repl-button";
import { Button } from "@/components/ui/button";
import { GetUserResponse } from "@/lib/types";
import { RocketIcon } from "@radix-ui/react-icons"
import { useQuery } from "@tanstack/react-query";
import axios from "axios";

export default function Page() {
  const query = useQuery({
    queryKey: ['user'],
    queryFn: () => axios.get<GetUserResponse>(
      process.env.NEXT_PUBLIC_API + "/api/auth/user",
      {
        withCredentials: true
      }
    ),
    staleTime: Infinity,
  })


  const isAuthenticatedMode = !query.isStale && query.isSuccess

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
              <RocketIcon className="mr-2 h-4 w-4" /> Register for free
            </Button>
          </a>}
        <CreateReplButton />
      </div>
    </main>
  );
}
