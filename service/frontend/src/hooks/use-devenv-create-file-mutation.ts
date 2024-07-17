"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import axios from "axios";

export type CreateFilePayload = {
  name: string;
};

export type DevenvCreateFileMutationOptions = {
  uuid: string;
  onSuccess?: (file: CreateFilePayload) => void;
};

export function useDevenvCreateFileMutation(
  options: DevenvCreateFileMutationOptions,
) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (file: CreateFilePayload) =>
      axios.post(
        (process.env.NEXT_PUBLIC_API ?? "") +
          "/api/devenv/" +
          options.uuid +
          "/files",
        file,
        {
          withCredentials: true,
        },
      ),
    onSuccess: (_, file) => {
      queryClient.setQueryData<string[]>(
        ["devenv", options.uuid, "files"],
        (oldData) => {
          let _data: string[] = oldData ?? [];
          if (_data.includes(file.name)) return _data;
          if (options.onSuccess) options.onSuccess(file);
          return [..._data, file.name].sort();
        },
      );
    },
  });
}
