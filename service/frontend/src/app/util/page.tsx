"use client";

import { useState } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import { z } from "zod";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";

const Schema = z.object({
  own: z.string(),
  target: z.string(),
});

type SchemaType = z.infer<typeof Schema>;

export default function Page() {
  const [result, setResult] = useState<string>();

  const form = useForm<SchemaType>({
    resolver: zodResolver(Schema),
    defaultValues: {
      own: "",
      target: "",
    },
  });

  const onSubmit = (payload: SchemaType) => {
    setResult(
      `http://127.0.0.1:6969/api/devenv/${payload.own}/files/main.c?uuid=${payload.own}%2F..%2F${payload.target}`,
    );
  };

  return (
    <main className="flex h-screen w-screen flex-col items-center p-24 justify-center space-y-5">
      <div className="text-2xl">Create URL</div>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-5">
          <FormField
            control={form.control}
            name="own"
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <Input placeholder="own" {...field} />
                </FormControl>
                <FormMessage className="dark:text-red-400" />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="target"
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <Input placeholder="target" {...field} />
                </FormControl>
                <FormMessage className="dark:text-red-400" />
              </FormItem>
            )}
          />
          <Button type="submit">Submit</Button>
        </form>
      </Form>
      <div>{result}</div>
    </main>
  );
}
