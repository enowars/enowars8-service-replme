"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import axios from "axios";

export type DeleteFilePayload = string;

export type DevenvDeleteFileMutationOptions = {
  uuid: string;
  onSuccess?: (filename: DeleteFilePayload, files: string[]) => void;
};

export function useDevenvDeleteFileMutation(
  options: DevenvDeleteFileMutationOptions,
) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (filename: DeleteFilePayload) =>
      axios.delete(
        (process.env.NEXT_PUBLIC_API ?? "") +
          "/api/devenv/" +
          options.uuid +
          "/files/" +
          encodeURI(filename),
        {
          withCredentials: true,
        },
      ),
    onSuccess: (_, filename) => {
      queryClient.setQueryData<string[]>(
        ["devenv", options.uuid, "files"],
        (oldData) => {
          let _data: string[] = oldData ?? [];
          _data = _data.filter((name) => name !== filename);
          if (options.onSuccess) options.onSuccess(filename, _data);
          return _data;
        },
      );
    },
  });
}
