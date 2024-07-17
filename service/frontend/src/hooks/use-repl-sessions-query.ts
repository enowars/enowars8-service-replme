"use client";

import { useQuery } from "@tanstack/react-query";
import axios from "axios";

export function useReplSessionsQuery() {
  return useQuery({
    queryKey: ["repl-sessions"],
    queryFn: () =>
      axios.get<string[]>(
        (process.env.NEXT_PUBLIC_API ?? "") + "/api/repl/sessions",
        { withCredentials: true },
      ),
  });
}
