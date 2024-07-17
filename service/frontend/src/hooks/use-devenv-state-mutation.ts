"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";

export type DevenvStateMutationOptions = {
  uuid?: string;
};

export function useDevenvStateMutation(options: DevenvStateMutationOptions) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (state: "ok" | "dirty") => {
      queryClient.setQueryData(["devenv", options.uuid, "state"], state);
    },
  });
}
