"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import React, { useEffect } from "react";
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "./ui/drawer";
import { Button } from "./ui/button";
import { GearIcon } from "@radix-ui/react-icons";
import { z } from "zod";
import { useDevenvQuery } from "@/hooks/use-devenv-query";
import { usePatchDevenvMutation } from "@/hooks/use-patch-devenv-mutation";
import { useForm } from "react-hook-form";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "./ui/input";

export type DevenvSettingsMenuProps = {
  uuid: string;
};

const DevenvFormSchema = z.object({
  name: z.string().min(1, { message: "Name command can't be empty" }),
  buildCmd: z.string().min(1, { message: "Build command can't be empty" }),
  runCmd: z.string().min(1, { message: "Run command can't be empty" }),
});

type DevenvForm = z.infer<typeof DevenvFormSchema>;

const DevenvSettingsMenu: React.FC<DevenvSettingsMenuProps> = (props) => {
  const devenvQuery = useDevenvQuery({
    uuid: props.uuid,
  });

  const form = useForm<DevenvForm>({
    resolver: zodResolver(DevenvFormSchema),
    defaultValues: {
      buildCmd: "",
      runCmd: "",
    },
  });

  const patchDevenvMutation = usePatchDevenvMutation({
    uuid: props.uuid,
  });

  useEffect(() => {
    if (devenvQuery.isSuccess) {
      form.setValue("name", devenvQuery.data.name);
      form.setValue("buildCmd", devenvQuery.data.buildCmd);
      form.setValue("runCmd", devenvQuery.data.runCmd);
    }
  }, [devenvQuery.isSuccess, devenvQuery.data]);

  const onSubmit = (devenvPatch: DevenvForm) => {
    patchDevenvMutation.mutate(devenvPatch);
  };

  return (
    <Drawer>
      <DrawerTrigger asChild>
        <Button className="relative" variant="outline">
          <GearIcon className="mr-2 h-4 w-4" /> Settings
        </Button>
      </DrawerTrigger>
      <DrawerContent className="w-full">
        <div className="w-full flex flex-col items-center pb-5">
          <DrawerHeader>
            <DrawerTitle>Devenv Settings</DrawerTitle>
            <DrawerDescription>
              Set the build command and run command
            </DrawerDescription>
          </DrawerHeader>
          <div className="flex flex-row justify-center items-center w-full space-x-5 overflow-auto py-5">
            <Form {...form}>
              <form
                onSubmit={form.handleSubmit(onSubmit)}
                className="space-y-5 w-full max-w-lg px-5"
              >
                <FormField
                  control={form.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Name</FormLabel>
                      <FormControl>
                        <Input placeholder="name" {...field} />
                      </FormControl>
                      <FormMessage className="dark:text-red-400" />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="buildCmd"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Build command</FormLabel>
                      <FormControl>
                        <Input placeholder="build command" {...field} />
                      </FormControl>
                      <FormMessage className="dark:text-red-400" />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="runCmd"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Run command</FormLabel>
                      <FormControl>
                        <Input placeholder="run command" {...field} />
                      </FormControl>
                      <FormMessage className="dark:text-red-400" />
                    </FormItem>
                  )}
                />
                <DrawerClose asChild>
                  <Button
                    className="w-full"
                    disabled={!form.formState.isDirty}
                    type="submit"
                  >
                    Save
                  </Button>
                </DrawerClose>
              </form>
            </Form>
          </div>
        </div>
      </DrawerContent>
    </Drawer>
  );
};

export default DevenvSettingsMenu;
