"use client"

import { ExitIcon } from "@radix-ui/react-icons";
import { Button } from "./ui/button";
import { navigate } from "@/actions/navigate";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import axios from "axios";

export function LogoutButton() {

  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: () => axios.post(
      process.env.NEXT_PUBLIC_API + '/api/auth/logout',
      undefined,
      {
        withCredentials: true
      }
    ),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user'] })
      navigate("/")
    }
  });

  return (
    <Button variant="outline" size="icon" onClick={() => mutation.mutate()}>
      <ExitIcon className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100" />
    </Button>
  )
}

