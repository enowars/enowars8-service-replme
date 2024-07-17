"use client";

import { randomString } from "@/lib/utils";
import { Button, ButtonProps } from "./ui/button";
import { ReloadIcon } from "@radix-ui/react-icons";
import React from "react";
import { PiTerminalWindowLight } from "react-icons/pi";
import { useCreateReplMutation } from "@/hooks/use-create-repl-mutation";

const CreateReplButton = React.forwardRef<HTMLButtonElement, ButtonProps>(
  (props, ref) => {
    const createReplMutation = useCreateReplMutation();
    const handleCreateRepl = (
      event: React.MouseEvent<HTMLButtonElement, MouseEvent>,
    ) => {
      const username = randomString(60);
      const password = randomString(60);
      createReplMutation.mutate({ username, password });
      if (props.onClick) props.onClick(event);
    };
    return (
      <Button
        ref={ref}
        {...props}
        disabled={createReplMutation.isPending}
        onClick={handleCreateRepl}
      >
        {createReplMutation.isPending ? (
          <ReloadIcon className="mr-2 h-4 w-4 animate-spin" />
        ) : (
          <PiTerminalWindowLight className="mr-2 h-4 w-4" />
        )}{" "}
        REPLME!
      </Button>
    );
  },
);

CreateReplButton.displayName = "CreateReplButton";
export default CreateReplButton;
