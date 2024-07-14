"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query"
import axios from "axios"

export type DevenvFileContentMutationOptions = {
  uuid: string;
  filename?: string
}

export function useDevenvFileContentMutation(options: DevenvFileContentMutationOptions) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (value: string) => axios.post(
      (process.env.NEXT_PUBLIC_API ?? "") + "/api/devenv/" + options.uuid + "/files/" + options.filename,
      value,
      {
        withCredentials: true
      }
    ),
    onSuccess: (_, value) => {
      queryClient.setQueryData<string>(
        ['devenv', options.uuid, 'files', options.filename, 'content'],
        () => {
          return value
        }
      )
    },
  })
}
