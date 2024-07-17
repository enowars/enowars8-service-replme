"use client";

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
import { useRegisterMutation } from "@/hooks/use-register-mutation";

const RegisterFormSchema = z.object({
  username: z
    .string()
    .min(4, { message: "Minimum length 4" })
    .max(64, { message: "Maximum length 4" })
    .regex(/^[a-zA-Z0-9]*$/, { message: "Only alphanumeric" }),
  password: z
    .string()
    .min(4, { message: "Minimum length 4" })
    .max(64, { message: "Maximum length 4" }),
});

type RegisterForm = z.infer<typeof RegisterFormSchema>;

export default function Page() {
  const form = useForm<RegisterForm>({
    resolver: zodResolver(RegisterFormSchema),
    defaultValues: {
      username: "",
      password: "",
    },
  });

  const mutation = useRegisterMutation({
    onError: () => {
      form.setError("username", {
        message: "That did not work, user exists?",
      });
    },
  });

  const onSubmit = (credentials: RegisterForm) => {
    mutation.mutate(credentials);
  };

  return (
    <main className="flex h-screen w-screen flex-col items-center p-24 justify-center space-y-5">
      <div className="text-2xl">Register</div>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-5">
          <FormField
            control={form.control}
            name="username"
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <Input placeholder="username" {...field} />
                </FormControl>
                <FormMessage className="dark:text-red-400" />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="password"
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <Input type="password" placeholder="password" {...field} />
                </FormControl>
                <FormMessage className="dark:text-red-400" />
              </FormItem>
            )}
          />
          <Button type="submit">Submit</Button>
        </form>
      </Form>
    </main>
  );
}
