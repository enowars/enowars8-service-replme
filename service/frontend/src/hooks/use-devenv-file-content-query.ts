"use client";

import { useQuery } from "@tanstack/react-query";
import axios from "axios";

export type DevenvFileContentQueryOptions = {
  uuid: string;
  filename?: string;
};

export function useDevenvFileContentQuery(
  options: DevenvFileContentQueryOptions,
) {
  return useQuery({
    queryKey: ["devenv", options.uuid, "files", options.filename, "content"],
    queryFn: () =>
      axios
        .get<string>(
          (process.env.NEXT_PUBLIC_API ?? "") +
            "/api/devenv/" +
            options.uuid +
            "/files/" +
            options.filename,
          {
            withCredentials: true,
          },
        )
        .then((data) => {
          return data.data;
        }),
    staleTime: Infinity,
    enabled: Boolean(options.filename),
  });
}
