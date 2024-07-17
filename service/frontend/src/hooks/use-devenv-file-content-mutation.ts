"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import axios from "axios";
import { useDevenvStateMutation } from "./use-devenv-state-mutation";

export type DevenvFileContentMutationOptions = {
  uuid: string;
  filename?: string;
};

export function useDevenvFileContentMutation(
  options: DevenvFileContentMutationOptions,
) {
  const queryClient = useQueryClient();

  const devenvStateMutation = useDevenvStateMutation({
    uuid: options.uuid,
  });

  return useMutation({
    mutationKey: ["devenv", options.uuid, "files", options.filename, "content"],
    mutationFn: (value: string) =>
      axios.post(
        (process.env.NEXT_PUBLIC_API ?? "") +
          "/api/devenv/" +
          options.uuid +
          "/files/" +
          options.filename,
        value,
        {
          withCredentials: true,
        },
      ),
    onSuccess: (_, value) => {
      devenvStateMutation.mutate("ok");
      queryClient.setQueryData<string>(
        ["devenv", options.uuid, "files", options.filename, "content"],
        () => {
          return value;
        },
      );
    },
  });
}
