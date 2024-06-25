"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { navigate } from "@/actions/navigate";
import { Button } from "@/components/ui/button";
import axios from "axios";
import { z } from "zod";
import { Form, FormControl, FormField, FormItem, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useMutation, useQueryClient } from "@tanstack/react-query";

const LoginFormSchema = z.object({
  username: z.string().min(1, { message: "Username can't be empty" }),
  password: z.string().min(1, { message: "Password can't be empty" })
})

type LoginForm = z.infer<typeof LoginFormSchema>;

export default function Page() {
  const queryClient = useQueryClient();

  const form = useForm<LoginForm>({
    resolver: zodResolver(LoginFormSchema),
    defaultValues: {
      username: "",
      password: "",
    }
  })

  const mutation = useMutation({
    mutationFn: (credentials: LoginForm) => axios.post(
      (process.env.NEXT_PUBLIC_API ?? "") + '/api/auth/login',
      credentials,
      {
        withCredentials: true
      }
    ),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user'] })
      navigate("/")
    },
    onError: () => {
      form.setError("password", {
        message: "Password is wrong"
      })
    }
  })

  const onSubmit = (credentials: LoginForm) => {
    mutation.mutate(credentials);
  }

  return (
    <main className="flex h-screen w-screen flex-col items-center p-24 justify-center space-y-5">
      <div className="text-2xl">Login</div>
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
