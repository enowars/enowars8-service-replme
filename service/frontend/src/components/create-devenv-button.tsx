"use client";

import { Button, ButtonProps } from "./ui/button";
import { CodeIcon, ReloadIcon } from "@radix-ui/react-icons";
import React from "react";
import { randomString } from "@/lib/utils";
import { useCreateDevenvMutation } from "@/hooks/use-create-devenv-mutation";

const CreateDevenvButton = React.forwardRef<HTMLButtonElement, ButtonProps>(
  (props, ref) => {
    const createDevenvMutation = useCreateDevenvMutation();

    const handleCreateDevenv = (
      event: React.MouseEvent<HTMLButtonElement, MouseEvent>,
    ) => {
      const name = randomString(10);
      createDevenvMutation.mutate({
        name,
        buildCmd: "gcc -o main main.c",
        runCmd: "./main",
      });
      if (props.onClick) props.onClick(event);
    };

    return (
      <Button
        ref={ref}
        {...props}
        disabled={createDevenvMutation.isPending}
        onClick={handleCreateDevenv}
      >
        {createDevenvMutation.isPending ? (
          <ReloadIcon className="mr-2 h-4 w-4 animate-spin" />
        ) : (
          <CodeIcon className="mr-2 h-4 w-4" />
        )}{" "}
        DEVENVME!
      </Button>
    );
  },
);

CreateDevenvButton.displayName = "CreateDevenvButton";

export default CreateDevenvButton;
