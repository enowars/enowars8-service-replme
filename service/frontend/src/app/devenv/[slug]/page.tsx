"use client";

import dynamic from "next/dynamic";
import FileTree from "@/components/file-tree";
import { ResizableHandle, ResizablePanel, ResizablePanelGroup } from "@/components/ui/resizable";
import { useState } from "react";
import ExecTerm from "@/components/exec-term";
import { Button } from "@/components/ui/button";

const Editor = dynamic(() => import('@/components/editor'), {
  ssr: false
})

export default function Page({ params }: { params: { slug: string } }) {
  const [showTerminal, setShowTerminal] = useState<number>(0)

  const [currentFile, setCurrentFile] = useState<string>();

  return (
    <main className="h-screen w-screen pt-16">
      <ResizablePanelGroup direction="horizontal">
        <ResizablePanel defaultSize={15}>
          <div className="flex flex-col w-full h-full p-4 space-y-5">
            <Button className="" onClick={() => setShowTerminal(showTerminal + 1)}>
              Run
            </Button>
            <FileTree className="flex flex-col w-full h-full space-y-5" devenvUuid={params.slug} selectedFile={currentFile} setSelectedFile={setCurrentFile} />
          </div>
        </ResizablePanel>
        <ResizableHandle />
        <ResizablePanel defaultSize={85}>
          <ResizablePanelGroup direction="vertical">
            <ResizablePanel>
              <Editor className="w-full h-full" devenvUuid={params.slug} filename={currentFile} />
            </ResizablePanel>
            {Boolean(showTerminal) && <>
              <ResizableHandle />
              <ResizablePanel>
                <ExecTerm className="w-full h-full" id={String(showTerminal)} path={"/api/devenv/" + params.slug + "/exec"} />
              </ResizablePanel>
            </>}
          </ResizablePanelGroup>
        </ResizablePanel>

      </ResizablePanelGroup>
    </main>
  );
}
