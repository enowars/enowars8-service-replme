"use client";

import { Devenv, PatchDevenvRequest } from "@/lib/types";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import axios from "axios";

export type PatchDevenvMutationOptions = {
  uuid: string;
  onSuccess?: (data: Devenv) => void;
};

export function usePatchDevenvMutation(options: PatchDevenvMutationOptions) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (patch: PatchDevenvRequest) =>
      axios.patch<void>(
        (process.env.NEXT_PUBLIC_API ?? "") + "/api/devenv/" + options.uuid,
        patch,
        {
          withCredentials: true,
        },
      ),

    onSuccess: (_, data) => {
      queryClient.invalidateQueries({ queryKey: ["devenvs"] });
      queryClient.setQueryData<Devenv>(["devenv", options.uuid], (oldData) => {
        if (!oldData) return oldData;
        if (data.name && data.name !== "") oldData.name = data.name;
        if (data.buildCmd && data.buildCmd !== "")
          oldData.buildCmd = data.buildCmd;
        if (data.runCmd && data.runCmd !== "") oldData.runCmd = data.runCmd;
        if (options.onSuccess) options.onSuccess(oldData);
        return { ...oldData };
      });
    },
  });
}
