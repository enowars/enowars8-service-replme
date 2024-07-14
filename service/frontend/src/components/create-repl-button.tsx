"use client";

import { navigate } from "@/actions/navigate";
import { randomString } from "@/lib/utils";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import axios from "axios";
import { Button, ButtonProps } from "./ui/button";
import { ReloadIcon } from "@radix-ui/react-icons";
import React from "react";
import { PiTerminalWindowLight } from "react-icons/pi";

type Credentials = {
  username: string,
  password: string,
}

const CreateReplButton = React.forwardRef<HTMLButtonElement, ButtonProps>(
  (props, ref) => {
    const client = useQueryClient();

    const mutation = useMutation({
      mutationFn: (credentials: Credentials) => axios.post<{ id: string }>(
        (process.env.NEXT_PUBLIC_API ?? "") + '/api/repl',
        credentials,
        {
          withCredentials: true
        }
      ),
      onSuccess: (response) => {
        client.invalidateQueries({ queryKey: ['repl-sessions'] })
        navigate("/repl/" + response.data.id)
      },
    })

    const handleCreateRepl = (event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
      const username = randomString(60);
      const password = randomString(60);
      mutation.mutate({ username, password });
      if (props.onClick) props.onClick(event)
    }

    return (
      <Button ref={ref} {...props} disabled={mutation.isPending} onClick={handleCreateRepl}>
        {mutation.isPending ? <ReloadIcon className="mr-2 h-4 w-4 animate-spin" /> :
          <PiTerminalWindowLight className="mr-2 h-4 w-4" />} REPLME!
      </Button>
    )
  }
)

CreateReplButton.displayName = "CreateReplButton";

export default CreateReplButton;
