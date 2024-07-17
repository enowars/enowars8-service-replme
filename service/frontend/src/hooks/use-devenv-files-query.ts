"use client";

import { useQuery } from "@tanstack/react-query";
import axios from "axios";

export type DevenvFilesQueryOptions = {
  uuid: string;
  callback?: (files: string[]) => void;
};

export function useDevenvFilesQuery(options: DevenvFilesQueryOptions) {
  return useQuery({
    queryKey: ["devenv", options.uuid, "files"],
    queryFn: () =>
      axios
        .get<string[]>(
          (process.env.NEXT_PUBLIC_API ?? "") +
            "/api/devenv/" +
            options.uuid +
            "/files",
          {
            withCredentials: true,
          },
        )
        .then((data) => {
          const files = data.data.sort();
          if (options.callback) options.callback(files);
          return files;
        }),
  });
}
