"use client";

import { navigate } from "@/actions/navigate";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import axios from "axios";

export type CreateReplMutationPayload = {
  username: string;
  password: string;
};

export function useCreateReplMutation() {
  const client = useQueryClient();

  return useMutation({
    mutationFn: (credentials: CreateReplMutationPayload) =>
      axios.post<{ id: string }>(
        (process.env.NEXT_PUBLIC_API ?? "") + "/api/repl",
        credentials,
        {
          withCredentials: true,
        },
      ),
    onSuccess: (response) => {
      client.invalidateQueries({ queryKey: ["repl-sessions"] });
      navigate("/repl/" + response.data.id);
    },
  });
}
