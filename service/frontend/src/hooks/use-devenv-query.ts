"use client";

import { Devenv } from "@/lib/types";
import { useQuery } from "@tanstack/react-query";
import axios from "axios";

export type DevenvQueryOptions = {
  uuid: string;
};

export function useDevenvQuery(options: DevenvQueryOptions) {
  return useQuery({
    queryKey: ["devenv", options.uuid],
    queryFn: () =>
      axios
        .get<Devenv>(
          (process.env.NEXT_PUBLIC_API ?? "") + "/api/devenv/" + options.uuid,
          {
            withCredentials: true,
          },
        )
        .then((data) => {
          return data.data;
        }),
    staleTime: Infinity,
  });
}
