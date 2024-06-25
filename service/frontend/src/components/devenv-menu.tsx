"use client";

import { Button } from "@/components/ui/button"
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "@/components/ui/drawer"
import { useQuery } from "@tanstack/react-query";
import axios from "axios";
import { CodeIcon } from "@radix-ui/react-icons";
import { navigate } from "@/actions/navigate";
import { Devenv } from "@/lib/types";
import CreateDevenvButton from "./create-devenv-button";

const DevenvMenu = () => {
  const query = useQuery(
    {
      queryKey: ["devenvs"],
      queryFn: () => axios.get<Devenv[]>((process.env.NEXT_PUBLIC_API ?? "") + '/api/devenv', { withCredentials: true }),
    }
  )

  const numSessions = query.data?.data?.length ?? 0;

  return (
    <Drawer>
      <DrawerTrigger asChild>
        <Button className="relative" variant="outline" size="icon">
          <CodeIcon className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all" />
          {Boolean(numSessions) && <div className="absolute w-5 h-5 bg-red-400 -top-2 -right-2 rounded-full text-white">{numSessions}</div>}
        </Button>
      </DrawerTrigger>
      <DrawerContent className="w-full">
        <div className="w-full flex flex-col items-center pb-5">
          <DrawerHeader>
            <DrawerTitle>Your Devenvs</DrawerTitle>
            <DrawerDescription>Open a devenv by clicking it</DrawerDescription>
          </DrawerHeader>

          <div className="flex flex-row justify-center items-center w-full space-x-5 overflow-auto">
            {query.data?.data?.map((devenv) => (
              <DrawerClose asChild key={devenv.id}>
                <Button variant="outline" onClick={() => navigate("/devenv/" + devenv.id)}>
                  {devenv.name}
                </Button>
              </DrawerClose>
            ))}
            <DrawerClose asChild>
              <CreateDevenvButton id="close-button" />
            </DrawerClose>
          </div>
        </div>
      </DrawerContent>
    </Drawer>
  )
}

export default DevenvMenu;
