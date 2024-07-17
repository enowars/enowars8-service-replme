"use client";

import { navigate } from "@/actions/navigate";
import { CreateDevenvRequest, CreateDevenvResponse } from "@/lib/types";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import axios from "axios";

export function useCreateDevenvMutation() {
  const client = useQueryClient();

  return useMutation({
    mutationFn: (credentials: CreateDevenvRequest) =>
      axios.post<CreateDevenvResponse>(
        (process.env.NEXT_PUBLIC_API ?? "") + "/api/devenv",
        credentials,
        {
          withCredentials: true,
        },
      ),
    onSuccess: (response) => {
      client.invalidateQueries({ queryKey: ["devenvs"] });
      navigate("/devenv/" + response.data.devenvUuid);
    },
  });
}
