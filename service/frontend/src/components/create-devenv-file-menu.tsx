"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "@/components/ui/drawer";
import { Button } from "./ui/button";
import { PlusIcon } from "@radix-ui/react-icons";
import { Form, FormControl, FormField, FormItem } from "./ui/form";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { useDevenvCreateFileMutation } from "@/hooks/use-devenv-create-file-mutation";
import { Input } from "./ui/input";

const CreateFileFormSchema = z.object({
  name: z.string().min(1).max(30),
});

type CreateFileForm = z.infer<typeof CreateFileFormSchema>;

export type CreateDevenvFileMenuProps = {
  uuid: string;
};

const CreateDevenvFileMenu: React.FC<CreateDevenvFileMenuProps> = (props) => {
  const createFileMutation = useDevenvCreateFileMutation({
    uuid: props.uuid,
  });

  const form = useForm<CreateFileForm>({
    resolver: zodResolver(CreateFileFormSchema),
    defaultValues: {
      name: "",
    },
  });

  const onCreateFileSubmit = (file: CreateFileForm) => {
    createFileMutation.mutate(file);
  };

  return (
    <Drawer>
      <DrawerTrigger asChild>
        <Button className="relative" variant="outline">
          <PlusIcon className="mr-2 h-4 w-4" /> New File
        </Button>
      </DrawerTrigger>
      <DrawerContent className="w-full">
        <div className="w-full flex flex-col items-center pb-5">
          <DrawerHeader>
            <DrawerTitle>New File</DrawerTitle>
            <DrawerDescription>Add a new file to your devenv</DrawerDescription>
          </DrawerHeader>
          <div className="flex flex-row justify-center items-center w-full space-x-5 overflow-auto py-5">
            <Form {...form}>
              <form
                onSubmit={form.handleSubmit(onCreateFileSubmit)}
                className="space-y-5 w-full max-w-lg "
              >
                <FormField
                  control={form.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem>
                      <FormControl>
                        <Input placeholder="filename" {...field} />
                      </FormControl>
                    </FormItem>
                  )}
                />
                <DrawerClose asChild>
                  <Button type="submit" className="w-full mb-5">
                    Add file
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

export default CreateDevenvFileMenu;
