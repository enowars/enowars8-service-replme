"use client";

import { Devenv } from "@/lib/types";
import { useQuery } from "@tanstack/react-query";
import axios from "axios";

export function useDevenvsQuery() {
  return useQuery({
    queryKey: ["devenvs"],
    queryFn: () =>
      axios.get<Devenv[]>((process.env.NEXT_PUBLIC_API ?? "") + "/api/devenv", {
        withCredentials: true,
      }),
  });
}
