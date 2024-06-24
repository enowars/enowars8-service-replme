"use client";

import { navigate } from "@/actions/navigate";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import axios from "axios";
import { Button, ButtonProps } from "./ui/button";
import { CodeIcon, ReloadIcon } from "@radix-ui/react-icons";
import React from "react";
import { CreateDevenvRequest, CreateDevenvResponse } from "@/lib/types";
import { randomString } from "@/lib/utils";

const CreateDevenvButton = React.forwardRef<HTMLButtonElement, ButtonProps>(
  (props, ref) => {
    const client = useQueryClient();

    const mutation = useMutation({
      mutationFn: (credentials: CreateDevenvRequest) => axios.post<CreateDevenvResponse>(
        process.env.NEXT_PUBLIC_API + '/api/devenv',
        credentials,
        {
          withCredentials: true
        }
      ),
      onSuccess: (response) => {
        client.invalidateQueries({ queryKey: ['devenvs'] })
        navigate("/devenv/" + response.data.devenvUuid)
      },
    })

    const handleCreateRepl = (event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
      const name = randomString(10);
      mutation.mutate({
        name,
        buildCmd: "gcc -o main main.c",
        runCmd: "./main"
      });
      if (props.onClick) props.onClick(event)
    }

    return (
      <Button ref={ref} {...props} disabled={mutation.isPending} onClick={handleCreateRepl}>
        {mutation.isPending ? <ReloadIcon className="mr-2 h-4 w-4 animate-spin" /> :
          <CodeIcon className="mr-2 h-4 w-4" />} DEVENVME!
      </Button>
    )
  }
)

CreateDevenvButton.displayName = "CreateDevenvButton";

export default CreateDevenvButton;
